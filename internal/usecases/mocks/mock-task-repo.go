package mocks

import (
	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type TaskRepo struct {
	mock.Mock
}

func (tr *TaskRepo) CreateTask(newTask *domain.Task) error {
	ret := tr.Called(newTask)
	newTask.Uuid = uuid.NewString()

	if rf, ok := ret.Get(0).(func(*domain.Task) error); ok {
		return rf(newTask)
	}
	return ret.Error(0)
}

func (tr *TaskRepo) GetByUid(uid string) (*domain.Task, error) {
	ret := tr.Called(uid)

	var (
		d *domain.Task
		e error
	)
	if rf, ok := ret.Get(0).(func(string) (*domain.Task, error)); ok {
		return rf(uid)
	}

	if rf, ok := ret.Get(0).(func(string) *domain.Task); ok {
		d = rf(uid)
	} else if ret.Get(0) != nil {
		d = ret.Get(0).(*domain.Task)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		e = rf(uid)
	} else if ret.Get(1) != nil {
		e = ret.Error(1)
	}

	return d, e
}

func (tr *TaskRepo) Update(updTask *domain.Task) error {
	ret := tr.Called(updTask)

	if rf, ok := ret.Get(0).(func(*domain.Task) error); ok {
		return rf(updTask)
	}
	return ret.Error(0)
}
