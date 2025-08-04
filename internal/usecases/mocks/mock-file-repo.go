package mocks

import (
	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/stretchr/testify/mock"
)

type FileRepo struct {
	mock.Mock
}

func (fr *FileRepo) Store(f *domain.File) error {
	ret := fr.Called(f)

	if rf, ok := ret.Get(0).(func(*domain.File) error); ok {
		return rf(f)
	}
	return ret.Error(0)
}

func (fr *FileRepo) GetByTaskUuid(uid string) []*domain.File {
	ret := fr.Called(uid)

	if rf, ok := ret.Get(0).(func(string) []*domain.File); ok {
		return rf(uid)
	}
	return ret.Get(0).([]*domain.File)
}
