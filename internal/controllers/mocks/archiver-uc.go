package mocks

import (
	"github.com/CALLlA-74/zip-26-07-25/domain"

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

	var (
		p0 *domain.CreateTaskResponse
		p1 error
	)
	if rf, ok := ret.Get(0).(func() *domain.CreateTaskResponse); ok {
		p0 = rf()
	} else if rf, ok := ret.Get(0).(*domain.CreateTaskResponse); ok {
		p0 = rf
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		p1 = rf()
	} else if rf, ok := ret.Get(1).(error); ok {
		p1 = rf
	}

	return p0, p1
}

func (_m *ArchiverUC) AddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error) {
	ret := _m.Called(taskUid, addLinksReq)

	if len(ret) <= 0 {
		panic("No return value specified for AddLinks")
	}

	if rf, ok := ret.Get(0).(func(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error)); ok {
		return rf(taskUid, addLinksReq)
	}

	var (
		p0 *domain.AddLinksResponse
		p1 error
	)
	if rf, ok := ret.Get(0).(func(taskUid string, addLinksReq *domain.AddLinksRequest) *domain.AddLinksResponse); ok {
		p0 = rf(taskUid, addLinksReq)
	} else if rf, ok := ret.Get(0).(*domain.AddLinksResponse); ok {
		p0 = rf
	}

	if rf, ok := ret.Get(1).(func(taskUid string, addLinksReq *domain.AddLinksRequest) error); ok {
		p1 = rf(taskUid, addLinksReq)
	} else if rf, ok := ret.Get(1).(error); ok {
		p1 = rf
	}

	return p0, p1
}

func (_m *ArchiverUC) GetStatus(taskUid string) (*domain.TaskStatusResponse, error) {
	ret := _m.Called(taskUid)

	if rf, ok := ret.Get(0).(func(taskUid string) (*domain.TaskStatusResponse, error)); ok {
		return rf(taskUid)
	}

	var (
		p0 *domain.TaskStatusResponse
		p1 error
	)
	if rf, ok := ret.Get(0).(func(taskUid string) *domain.TaskStatusResponse); ok {
		p0 = rf(taskUid)
	} else if rf, ok := ret.Get(0).(*domain.TaskStatusResponse); ok {
		p0 = rf
	}

	if rf, ok := ret.Get(1).(func(taskUid string) error); ok {
		p1 = rf(taskUid)
	} else if rf, ok := ret.Get(1).(error); ok {
		p1 = rf
	}

	return p0, p1
}
