package controllers_test

import (
	"context"
	"encoding/json"
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
	faker "github.com/go-faker/faker/v4"
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

		handler := &controllers.ApiHandlerV1{
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
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), e.Error()))
}

func TestAddFileLinks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testPipline := func(tUid string, mReq *domain.AddLinksRequest,
		mr *domain.AddLinksResponse, e error) *httptest.ResponseRecorder {

		mockUCase := new(mocks.ArchiverUC)
		mockUCase.On("AddLinks", tUid, mReq).Return(mr, e)

		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)

		j, e := json.Marshal(mReq)
		assert.NoError(t, e)

		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, config.ApiV1GroupName+"/:taskUuid/add-file-links",
			strings.NewReader(string(j)))

		req.SetPathValue("taskUuid", tUid)
		ctx.Request = req

		assert.NoError(t, err)

		handler := &controllers.ApiHandlerV1{
			Ias: mockUCase,
		}
		handler.AddFileLinks(ctx)

		mockUCase.AssertExpectations(t)
		return rec
	}

	mReq := new(domain.AddLinksRequest)
	assert.NoError(t, faker.FakeData(mReq))
	sl := make([]string, 2)
	copy(sl, mReq.Links[:2])
	mResp := &domain.AddLinksResponse{AddedLinks: sl, HasReachedLimit: false}
	tUid := uuid.NewString()
	rec := testPipline(tUid, mReq, mResp, nil)
	assert.Equal(t, http.StatusAccepted, rec.Code)

	rec = testPipline(tUid, mReq, nil, domain.ErrAddingImpossible)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), domain.ErrAddingImpossible.Error()))

	rec = testPipline(tUid, mReq, nil, domain.ErrTaskNotFound)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), domain.ErrTaskNotFound.Error()))
}

func TestGetStatus(t *testing.T) {

}
