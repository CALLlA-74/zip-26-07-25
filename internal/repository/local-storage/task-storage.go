package localstorage

import (
	"sync"
	"sync/atomic"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/google/uuid"
)

type TaskModel struct {
	task domain.Task

	Version int64
	rSm     atomic.Int64
	wMx     sync.Mutex
}

type TaskDTO struct {
	Task *domain.Task
}

type TaskStorage struct {
	tasks map[string]*TaskModel
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]*TaskModel),
	}
}

func (ts *TaskStorage) CreateTask(dto *TaskDTO) error {
	uid := uuid.NewString()
	for _, ok := ts.tasks[uid]; ok; _, ok = ts.tasks[uid] {
		uid = uuid.NewString()
	}
	dto.Task.Uuid = uid

	newModel := &TaskModel{
		task: *dto.Task,
	}

	ts.tasks[uid] = newModel
	return nil
}

func (ts *TaskStorage) GetByUid(uid string) (*TaskDTO, error) {
	model, ok := ts.tasks[uid]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}

	model.wMx.Lock()
	model.rSm.Add(1)
	model.wMx.Unlock()
	defer model.rSm.Add(-1)

	t := model.task
	res := &TaskDTO{
		Task: &t,
	}
	return res, nil
}

func (ts *TaskStorage) Update(updDTO *TaskDTO) error {
	model, ok := ts.tasks[updDTO.Task.Uuid]
	if !ok {
		return domain.ErrTaskNotFound
	}

	model.wMx.Lock()
	defer model.wMx.Unlock()
	for model.rSm.Load() > 0 {
	}

	if model.Version+1 != updDTO.Task.Version {
		return domain.ErrVersionConflict
	}

	model.Version = updDTO.Task.Version
	model.task = domain.Task{
		TaskStatus:  updDTO.Task.TaskStatus,
		NumOfLinks:  updDTO.Task.NumOfLinks,
		ArchiveName: updDTO.Task.ArchiveName,
		ArchiveLink: updDTO.Task.ArchiveLink,
		Version:     updDTO.Task.Version,
	}

	return nil
}
