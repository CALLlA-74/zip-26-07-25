package usecases

import (
	"container/list"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type iTaskRepo interface {
	CreateTask() (*domain.Task, error)
	GetByUid(uid string) (*domain.Task, error)
	Update(updTask *domain.Task) error
}

type iFileRepo interface {
	Store(f *domain.File) error
	GetByTaskUuid(uid string) []*domain.File
}

type AchiverUC struct {
	itr     iTaskRepo
	ifl     iFileRepo
	taskSem *semaphore.Weighted
}

func NewArchiverUC(itr iTaskRepo, ifl iFileRepo) *AchiverUC {
	return &AchiverUC{
		itr:     itr,
		ifl:     ifl,
		taskSem: semaphore.NewWeighted(config.TASK_LIMIT),
	}
}

func (auc *AchiverUC) CreateTask() (*domain.CreateTaskResponse, error) {
	if !auc.taskSem.TryAcquire(1) {
		return nil, domain.ErrBusyServer
	}

	task, err := auc.itr.CreateTask()
	if err != nil {
		auc.taskSem.Release(1)
		return nil, err
	}

	return &domain.CreateTaskResponse{TaskUuid: task.Uuid}, nil
}

func (auc *AchiverUC) AddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error) {
	resp, err := auc.tryToAddLinks(taskUid, addLinksReq)
	if err != nil {
		return nil, err
	}

	auc.loadFiles(taskUid, resp.AddedLinks)

	return resp, nil
}

func (auc *AchiverUC) tryToAddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error) {
	isOk := false
	delta := int64(0)
	hasReachedLimit := false
	for !isOk {
		task, err := auc.itr.GetByUid(taskUid)
		if err != nil {
			return nil, err
		}

		if task.TaskStatus != domain.WAITING_LINKS || task.NumOfLinks >= config.LINKS_LIMIT {
			return nil, domain.ErrAddingImpossible
		}

		delta = config.LINKS_LIMIT - task.NumOfLinks

		task.NumOfLinks = min(config.LINKS_LIMIT, task.NumOfLinks+int64(len(addLinksReq.Links)))
		task.Version++

		if task.NumOfLinks == config.LINKS_LIMIT {
			task.TaskStatus = domain.PROCESSING
			hasReachedLimit = true
		}

		if err := auc.itr.Update(task); err != nil {
			if err != domain.ErrVersionConflict {
				return nil, err
			} else {
				isOk = true
			}
		}
	}

	return &domain.AddLinksResponse{
		AddedLinks:      addLinksReq.Links[:delta],
		HasReachedLimit: hasReachedLimit,
	}, nil
}

func (auc *AchiverUC) loadFiles(taskUid string, links []string) {
	exceptFiles := func(files *list.List, err error) {
		for iter := files.Front(); iter != nil; iter = iter.Next() {
			f := iter.Value.(*domain.File)
			f.ErrorMessage = err.Error()
			f.Path = ""
			if e := auc.ifl.Store(f); e != nil {
				logrus.Error(e)
			}
		}
	}

	loadAndValidate := func(wg *sync.WaitGroup, inp chan *list.List) {
		defer wg.Done()
		files := <-inp
		f := files.Front().Value.(*domain.File)

		if e := domain.DownloadFile(f.Link, f.Path); e != nil {
			logrus.Errorf("Load err file by link: %s; msg: %s", f.Link, e.Error())
			exceptFiles(files, domain.ErrLoadFile)
			return
		}

		fType := domain.ValidateFile(f.Path)
		if fType == domain.UNKNOWN_TYPE {
			os.Remove(f.Path)
			logrus.Infof("File \"%s\" has unsupportable type", f.Path)
			exceptFiles(files, domain.ErrLoadFile)
			return
		}

		comm := fmt.Sprintf("%s.%s", time.Now().Format("20060102150405"), fType)
		path := f.Path
		for iter := files.Front(); iter != nil; iter = iter.Next() {
			f := iter.Value.(*domain.File)
			f.Path = fmt.Sprintf("%s_%s", f.Path, comm)
			if e := domain.CopyFile(path, f.Path); e != nil {
				f.Path = ""
				f.ErrorMessage = domain.ErrInternalServerError.Error()
			}
			auc.ifl.Store(f)
		}
		os.Remove(path)
	}

	go func() {
		defer auc.taskSem.Release(1)

		linksMap := make(map[string]*list.List)
		for idx, v := range links {
			val, ok := linksMap[v]
			if !ok {
				val = list.New().Init()
				linksMap[v] = val
			}
			f := &domain.File{
				TaskUuid: taskUid,
				Link:     v,
				Path:     fmt.Sprintf("%s/%d_%s_", config.DownloadPath, idx+1, taskUid),
			}
			val.PushBack(f)
		}

		wg := new(sync.WaitGroup)
		wg.Add(len(linksMap))

		for _, v := range linksMap {
			ch := make(chan *list.List)
			go loadAndValidate(wg, ch)
			ch <- v
		}
		wg.Wait()

		files := auc.ifl.GetByTaskUuid(taskUid)
		archName := fmt.Sprintf("%s_%s.zip", taskUid, time.Now().Format("20060102150405"))
		domain.PackToArchiver(fmt.Sprintf("%s/%s", config.DownloadPath, archName), files)
		task, err := auc.itr.GetByUid(taskUid)
		if err != nil {
			logrus.Errorf("Get task %s error: %s", taskUid, err)
			return
		}

		task.ArchiveName = archName
		task.ArchiveLink = fmt.Sprintf("%s/%s", config.LoadArchGroupName, archName)
		task.TaskStatus = domain.FINISHED

		if e := auc.itr.Update(task); e != nil {
			logrus.Errorf("Update task %s error: %s", taskUid, e)
			return
		}
	}()
}

func (auc *AchiverUC) GetStatus(taskUid string) (*domain.TaskStatusResponse, error) {
	task, err := auc.itr.GetByUid(taskUid)
	if err != nil {
		return nil, err
	}

	files := auc.ifl.GetByTaskUuid(taskUid)

	fLinks := make([]*domain.FailedLink, 0, len(files))
	for _, f := range files {
		if len(f.ErrorMessage) > 0 {
			fLinks = append(fLinks, &domain.FailedLink{
				Link:    f.Link,
				Message: f.ErrorMessage,
			})
		}
	}

	return &domain.TaskStatusResponse{
		Status:      task.TaskStatus,
		FailedLinks: fLinks,
		ArchiveLink: task.ArchiveLink,
	}, nil
}
