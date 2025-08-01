package controllers

import (
	"net/http"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/gin-gonic/gin"
)

type iArchiverService interface {
	CreateTask() (*domain.CreateTaskResponse, error)
	AddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error)
	GetStatus(taskUid string) (*domain.TaskStatusResponse, error)
}

type ApiHandlerV1 struct {
	Ias iArchiverService
}

func NewArchiverHandler(ias iArchiverService, g gin.IRouter) *ApiHandlerV1 {
	g.POST("/create-task")
	g.PATCH("/:taskUuid/add-file-links")
	g.GET("/:taskUuid/status")
	return &ApiHandlerV1{
		Ias: ias,
	}
}

func (ah ApiHandlerV1) CreateTask(context *gin.Context) {
	task, e := ah.Ias.CreateTask()
	if e != nil {
		context.JSON(errToStatusCode(e), &domain.ErrorResponse{Message: e.Error()})
		return
	}
	context.JSON(http.StatusCreated, task)
}

func (ah ApiHandlerV1) AddFileLinks(context *gin.Context) {
	taskUuid := context.Param("taskUuid")

	addLinksReq := new(domain.AddLinksRequest)
	if e := context.BindJSON(addLinksReq); e != nil {
		context.JSON(errToStatusCode(e), &domain.ErrorResponse{Message: e.Error()})
		return
	}

	resp, e := ah.Ias.AddLinks(taskUuid, addLinksReq)
	if e != nil {
		context.JSON(errToStatusCode(e), &domain.ErrorResponse{Message: e.Error()})
		return
	}
	context.JSON(http.StatusAccepted, resp)
}

func (ah ApiHandlerV1) GetStatus(context *gin.Context) {
	taskUuid := context.Param("taskUuid")
	resp, e := ah.Ias.GetStatus(taskUuid)
	if e != nil {
		context.JSON(errToStatusCode(e), &domain.ErrorResponse{Message: e.Error()})
		return
	}
	context.JSON(http.StatusOK, resp)
}

func errToStatusCode(e error) int {
	if e == nil {
		return http.StatusOK
	}

	switch e {
	case domain.ErrBusyServer:
		return http.StatusServiceUnavailable
	case domain.ErrAddingImpossible:
		return http.StatusForbidden
	case domain.ErrTaskNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
