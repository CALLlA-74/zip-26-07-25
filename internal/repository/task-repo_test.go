package repository_test

import (
	"sync"
	"testing"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/CALLlA-74/zip-26-07-25/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	tr := repository.NewTaskRepo()

	e := tr.CreateTask(nil)
	assert.Error(t, e)
	assert.Equal(t, domain.ErrInternalServerError.Error(), e.Error())

	task := &domain.Task{
		Uuid:       "",
		TaskStatus: domain.WAITING_LINKS,
	}
	e = tr.CreateTask(task)
	assert.NoError(t, e)
	assert.Equal(t, true, len(task.Uuid) > 0)
}

func TestGetByUid(t *testing.T) {
	tr := repository.NewTaskRepo()
	task := &domain.Task{
		TaskStatus: domain.WAITING_LINKS,
	}
	e := tr.CreateTask(task)
	assert.NoError(t, e)
	assert.Equal(t, true, len(task.Uuid) > 0)

	task2, e2 := tr.GetByUid("")
	assert.Error(t, e2)
	assert.Equal(t, domain.ErrTaskNotFound.Error(), e2.Error())
	assert.Empty(t, task2)

	task3, e3 := tr.GetByUid(task.Uuid)
	assert.NoError(t, e3)
	assert.NotEmpty(t, task3)
	assert.Equal(t, task.Uuid, task3.Uuid)
	assert.Equal(t, task.TaskStatus, task3.TaskStatus)
}

func TestGetAndUpdate(t *testing.T) {
	tr := repository.NewTaskRepo()

	task := &domain.Task{
		TaskStatus: domain.WAITING_LINKS,
	}
	e := tr.CreateTask(task)
	assert.NoError(t, e)
	assert.Equal(t, true, len(task.Uuid) > 0)

	e2 := tr.Update(nil)
	assert.Error(t, e2)
	assert.Equal(t, domain.ErrInternalServerError.Error(), e2.Error())

	task3 := &domain.Task{
		Uuid:       uuid.NewString(),
		TaskStatus: domain.PROCESSING,
	}
	e3 := tr.Update(task3)
	assert.Error(t, e3)
	assert.Equal(t, domain.ErrTaskNotFound.Error(), e3.Error())

	numOfGorutines := int64(1e3)
	wg := new(sync.WaitGroup)
	wg.Add(int(numOfGorutines))
	for i := int64(0); i < numOfGorutines; i++ {
		go func() {
			defer wg.Done()
			var task4 *domain.Task
			var e4 error
			for {
				task4, e4 = tr.GetByUid(task.Uuid)
				assert.NotEmpty(t, task4)
				assert.NoError(t, e4)
				task4.Uuid = task.Uuid
				task4.NumOfLinks++
				task4.Version++
				e42 := tr.Update(task4)
				if e42 == nil || (e42 != nil && e42.Error() != domain.ErrVersionConflict.Error()) {
					break
				}
			}
		}()
	}
	wg.Wait()
	task5, e5 := tr.GetByUid(task.Uuid)
	assert.NotEmpty(t, task5)
	assert.NoError(t, e5)
	assert.Equal(t, numOfGorutines, task5.NumOfLinks)
}
