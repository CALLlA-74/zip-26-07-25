package repository

import (
	"github.com/CALLlA-74/zip-26-07-25/domain"
	localstorage "github.com/CALLlA-74/zip-26-07-25/internal/repository/local-storage"
)

type TaskRepo struct {
	storage *localstorage.TaskStorage
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		storage: localstorage.NewTaskStorage(),
	}
}

func (tr *TaskRepo) CreateTask(newTask *domain.Task) error {
	if newTask == nil {
		return domain.ErrInternalServerError
	}
	return tr.storage.CreateTask(&localstorage.TaskDTO{
		Task: newTask,
	})
}

func (tr *TaskRepo) GetByUid(uid string) (*domain.Task, error) {
	dto, err := tr.storage.GetByUid(uid)
	if err != nil {
		return nil, err
	}
	return dto.Task, nil
}

func (tr *TaskRepo) Update(updTask *domain.Task) error {
	if updTask == nil {
		return domain.ErrInternalServerError
	}
	return tr.storage.Update(&localstorage.TaskDTO{
		Task: updTask,
	})
}
