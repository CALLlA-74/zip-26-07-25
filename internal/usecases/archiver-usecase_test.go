package usecases_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/CALLlA-74/zip-26-07-25/internal/usecases"
	"github.com/CALLlA-74/zip-26-07-25/internal/usecases/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	t.Run("stress-test", func(t *testing.T) {
		mocksTaskRepo := new(mocks.TaskRepo)
		mockFileRepo := new(mocks.FileRepo)

		mocksTaskRepo.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(nil).Times(config.TASK_LIMIT)

		uc := usecases.NewArchiverUC(mocksTaskRepo, mockFileRepo)

		numOfGOrutines := int(1e3)
		wg := new(sync.WaitGroup)
		wg.Add(numOfGOrutines)

		start := false
		var numOfTasks atomic.Int64
		for i := 0; i < numOfGOrutines; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				for !start {
				}
				resp, e := uc.CreateTask()

				if e != nil {
					assert.Empty(t, resp)
					assert.Equal(t, domain.ErrBusyServer.Error(), e.Error())
				} else {
					numOfTasks.Add(1)
					assert.NotEmpty(t, resp)
				}
			}(wg)
		}
		start = true
		wg.Wait()
		assert.Equal(t, config.TASK_LIMIT, int(numOfTasks.Load()))
	})

	t.Run("first-fail-second-success", func(t *testing.T) {
		req := mock.AnythingOfType("*domain.Task")
		mocksTaskRepo := new(mocks.TaskRepo)
		mockFileRepo := new(mocks.FileRepo)
		mocksTaskRepo.On("CreateTask", req).Return(domain.ErrInternalServerError).Once()
		mocksTaskRepo.On("CreateTask", req).Return(nil).Times(config.TASK_LIMIT)

		uc := usecases.NewArchiverUC(mocksTaskRepo, mockFileRepo)

		resp, e := uc.CreateTask()
		assert.Empty(t, resp)
		assert.Error(t, e)
		assert.Equal(t, domain.ErrInternalServerError.Error(), e.Error())

		for i := 0; i < config.TASK_LIMIT; i++ {
			r, err := uc.CreateTask()
			assert.Empty(t, err)
			assert.NotEmpty(t, r)
		}
	})
}

func TestAddLinks(t *testing.T) {
	init := func() (*mocks.TaskRepo, *mocks.FileRepo, string, *usecases.AchiverUC) {
		mocksTaskRepo := new(mocks.TaskRepo)
		mockFileRepo := new(mocks.FileRepo)
		uc := usecases.NewArchiverUC(mocksTaskRepo, mockFileRepo)
		uid := uuid.NewString()

		for i := 0; i < config.TASK_LIMIT; i++ {
			mocksTaskRepo.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(nil).Once()
			resp, e := uc.CreateTask()
			assert.NoError(t, e)
			assert.NotEmpty(t, resp)
		}
		return mocksTaskRepo, mockFileRepo, uid, uc
	}

	req := &domain.AddLinksRequest{
		Links: []string{
			"https://content-27.foto.my.mail.ru/community/inache/_groupsphoto/h-10206.jpg",
			"http://wiki.sunfounder.cc/images/f/f8/Bluetooth_4.0_BLE_module_datasheet.pdf",
			"https://content-27.foto.my.mail.ru/community/inache/_groupsphoto/h-10206.jpg",
		},
	}

	t.Run("test-adding-impossible", func(t *testing.T) {
		mtr, _, uid, uc := init()

		task := &domain.Task{
			Uuid:       uid,
			TaskStatus: domain.PROCESSING,
		}

		mtr.On("GetByUid", uid).Return(task, nil).Once()

		r, e := uc.AddLinks(uid, req)
		assert.Empty(t, r)
		assert.Equal(t, domain.ErrAddingImpossible.Error(), e.Error())
	})

	t.Run("test-version-conflict-and-internal-serv-err", func(t *testing.T) {
		mtr, _, uid, uc := init()

		task := domain.Task{
			Uuid:       uid,
			TaskStatus: domain.WAITING_LINKS,
		}
		t1 := task
		t2 := task

		mtr.On("GetByUid", uid).Return(&t1, nil).Once()
		mtr.On("Update", &t1).Return(domain.ErrVersionConflict).Once()

		mtr.On("GetByUid", uid).Return(&t2, nil).Once()
		mtr.On("Update", &t2).Return(domain.ErrInternalServerError).Once()

		r, e := uc.AddLinks(uid, req)
		assert.Empty(t, r)
		assert.Equal(t, domain.ErrInternalServerError.Error(), e.Error())
	})

	t.Run("test-success", func(t *testing.T) {
		mtr, mfr, uid, uc := init()

		task := &domain.Task{
			Uuid:       uid,
			TaskStatus: domain.WAITING_LINKS,
		}

		mtr.On("GetByUid", uid).Return(task, nil).Once()
		mtr.On("Update", task).Return(nil).Once()

		for i := 0; i < len(req.Links); i++ {
			mfr.On("Store", mock.AnythingOfType("*domain.File")).Return(nil).Once()
		}

		r, e := uc.AddLinks(uid, req)
		assert.Empty(t, e)
		assert.NotEmpty(t, r)
		assert.Equal(t, len(req.Links), len(r.AddedLinks))
		assert.Equal(t, len(r.AddedLinks) == config.LINKS_LIMIT, r.HasReachedLimit == true)

		/*mtr.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(nil).Once()
		resp, err := uc.CreateTask()
		assert.NoError(t, err)
		assert.NotEmpty(t, resp)*/
	})
}

func TestGetStatus(t *testing.T) {
	t.Run("test-success", func(t *testing.T) {
		mocksTaskRepo := new(mocks.TaskRepo)
		mockFileRepo := new(mocks.FileRepo)
		uc := usecases.NewArchiverUC(mocksTaskRepo, mockFileRepo)
		uid := uuid.NewString()
		task := &domain.Task{}
		assert.NoError(t, faker.FakeData(task))
		task.TaskStatus = domain.FINISHED

		files := []*domain.File{
			{
				TaskUuid:     uid,
				Link:         "",
				ErrorMessage: "",
			},
			{
				TaskUuid:     uid,
				Link:         "",
				ErrorMessage: domain.ErrLoadFile.Error(),
			},
			{
				TaskUuid:     uid,
				Link:         "",
				ErrorMessage: domain.ErrUnsuppType.Error(),
			},
		}

		mocksTaskRepo.On("GetByUid", uid).Return(task, nil).Once()
		mockFileRepo.On("GetByTaskUuid", uid).Return(files, nil)

		resp, e := uc.GetStatus(uid)
		assert.Empty(t, e)
		assert.NotEmpty(t, resp)
		assert.Equal(t, 2, len(resp.FailedLinks))
		assert.Equal(t, domain.ErrLoadFile.Error(), resp.FailedLinks[0].Message)
		assert.Equal(t, domain.ErrUnsuppType.Error(), resp.FailedLinks[1].Message)
	})

	t.Run("test-notfound-internal-server-error", func(t *testing.T) {
		mocksTaskRepo := new(mocks.TaskRepo)
		mockFileRepo := new(mocks.FileRepo)
		uc := usecases.NewArchiverUC(mocksTaskRepo, mockFileRepo)
		uid := uuid.NewString()

		mocksTaskRepo.On("GetByUid", uid).Return(nil, domain.ErrTaskNotFound).Once()
		mocksTaskRepo.On("GetByUid", uid).Return(nil, domain.ErrInternalServerError).Once()

		r, e := uc.GetStatus(uid)
		assert.Empty(t, r)
		assert.NotEmpty(t, e)
		assert.Equal(t, domain.ErrTaskNotFound.Error(), e.Error())

		r, e = uc.GetStatus(uid)
		assert.Empty(t, r)
		assert.NotEmpty(t, e)
		assert.Equal(t, domain.ErrInternalServerError.Error(), e.Error())
	})
}
