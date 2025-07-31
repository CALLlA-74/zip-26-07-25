package mocks

import (
	"errors"

	"github.com/CALLlA-74/zip-service/domain"

	"github.com/stretchr/testify/mock"
)

type ArchiverUC struct {
	mock.Mock
}

func (_m *ArchiverUC) CreateTask() (*domain.CreateTaskResponse, error) {
	ret := _m.Called()

	if len(ret) <= 0 {
		panic("No return value specified for CreateTask")
	}

	if rf, ok := ret.Get(0).(func() (*domain.CreateTaskResponse, error)); ok {
		return rf()
	}

	return nil, errors.New("Cast error")
}

func (_m *ArchiverUC) AddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error) {
	ret := _m.Called(taskUid, addLinksReq)

	if len(ret) <= 0 {
		panic("No return value specified for AddLinks")
	}

	if rf, ok := ret.Get(0).(func(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error)); ok {
		return rf(taskUid, addLinksReq)
	}

	return nil, errors.New("Cast error")
}

func (_m *ArchiverUC) GetStatus(taskUid string) (*domain.TaskStatusResponse, error) {
	ret := _m.Called(taskUid)

	if rf, ok := ret.Get(0).(func(taskUid string) (*domain.TaskStatusResponse, error)); ok {
		return rf(taskUid)
	}

	return nil, errors.New("Cast error")
}
