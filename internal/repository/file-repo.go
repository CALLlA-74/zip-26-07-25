package repository

import (
	"github.com/CALLlA-74/zip-26-07-25/domain"
	localstorage "github.com/CALLlA-74/zip-26-07-25/internal/repository/local-storage"
)

type FileRepo struct {
	storage *localstorage.FileStorage
}

func NewFileRepo() *FileRepo {
	return &FileRepo{
		storage: localstorage.NewFileStorage(),
	}
}

func (fr *FileRepo) Store(f *domain.File) error {
	if f == nil {
		return domain.ErrInternalServerError
	}
	return fr.storage.Store(&localstorage.FileDTO{
		File: f,
	})
}

func (fr *FileRepo) GetByTaskUuid(taskUid string) []*domain.File {
	files := fr.storage.GetByTaskUuid(taskUid)
	res := make([]*domain.File, 0, len(files))
	for _, f := range files {
		res = append(res, f.File)
	}

	return res
}
