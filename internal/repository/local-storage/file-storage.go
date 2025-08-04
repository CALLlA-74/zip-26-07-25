package localstorage

import (
	"sync"
	"sync/atomic"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/google/uuid"
)

type FileModel struct {
	file domain.File

	Version int64
	rSm     atomic.Int64
	wMx     sync.Mutex
}

type FileDTO struct {
	File *domain.File
}

type FileStorage struct {
	files       map[string]*FileModel
	filesByTask map[string][]*FileModel
}

func NewFileStorage() *FileStorage {
	return &FileStorage{
		files:       make(map[string]*FileModel),
		filesByTask: make(map[string][]*FileModel),
	}
}

func (fs *FileStorage) Store(dto *FileDTO) error {
	uid := uuid.NewString()
	for _, ok := fs.filesByTask[uid]; ok; _, ok = fs.filesByTask[uid] {
		uid = uuid.NewString()
	}
	dto.File.FileUid = uid

	newModel := &FileModel{
		file: *dto.File,
	}

	fs.files[uid] = newModel
	fs.filesByTask[dto.File.TaskUuid] = append(fs.filesByTask[dto.File.TaskUuid], newModel)
	return nil
}

func (fs *FileStorage) GetByTaskUuid(taskUid string) []*FileDTO {
	files, ok := fs.filesByTask[taskUid]
	if !ok {
		return []*FileDTO{}
	}

	res := make([]*FileDTO, 0, len(files))
	for _, f := range files {
		t := f.file
		res = append(res, &FileDTO{
			File: &t,
		})
	}
	return res
}
