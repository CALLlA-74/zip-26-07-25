package controllers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/CALLlA-74/zip-26-07-25/internal/controllers"

	"github.com/CALLlA-74/zip-26-07-25/internal/controllers/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPipline := func(mt *domain.CreateTaskResponse, e error) *httptest.ResponseRecorder {
		mockUCase := new(mocks.ArchiverUC)
		mockUCase.On("CreateTask").Return(mt, e)

		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, config.ApiV1GroupName+"/create-task",
			strings.NewReader(""))
		assert.NoError(t, err)
		ctx.Request = req

		handler := controllers.ApiHandlerV1{
			Ias: mockUCase,
		}
		handler.CreateTask(ctx)

		mockUCase.AssertExpectations(t)
		return rec
	}

	mockTask := &domain.CreateTaskResponse{
		TaskUuid: uuid.NewString(),
	}
	rec := testPipline(mockTask, nil)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), mockTask.TaskUuid))

	rec = testPipline(nil, domain.ErrBusyServer)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), domain.ErrBusyServer.Error()))

	e := errors.New("Internal Server error")
	rec = testPipline(nil, e)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), e.Error()))
}

func TestAddFileLinks(t *testing.T) {

}

func TestGetStatus(t *testing.T) {

}
