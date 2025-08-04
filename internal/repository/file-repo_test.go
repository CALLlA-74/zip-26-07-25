package repository_test

import (
	"testing"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/CALLlA-74/zip-26-07-25/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	tf := repository.NewFileRepo()

	e := tf.Store(nil)
	assert.Error(t, e)
	assert.Equal(t, domain.ErrInternalServerError.Error(), e.Error())

	file := &domain.File{
		FileUid:  "",
		TaskUuid: uuid.NewString(),
	}
	e = tf.Store(file)
	assert.NoError(t, e)
	assert.Equal(t, true, len(file.FileUid) > 0)
}

func TestGetByTaskUuid(t *testing.T) {
	tf := repository.NewFileRepo()
	taskUid := uuid.NewString()

	file := &domain.File{
		FileUid:  "",
		TaskUuid: taskUid,
	}
	e := tf.Store(file)
	assert.NoError(t, e)
	assert.Equal(t, true, len(file.FileUid) > 0)

	files := tf.GetByTaskUuid("")
	assert.Equal(t, true, len(files) == 0)

	files2 := tf.GetByTaskUuid(taskUid)
	assert.Equal(t, 1, len(files2))
}
