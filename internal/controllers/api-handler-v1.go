package controllers

import (
	"fmt"
	"net/http"

	"github.com/CALLlA-74/zip-26-07-25/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type iArchiverUsecase interface {
	CreateTask() (*domain.CreateTaskResponse, error)
	AddLinks(taskUid string, addLinksReq *domain.AddLinksRequest) (*domain.AddLinksResponse, error)
	GetStatus(taskUid string) (*domain.TaskStatusResponse, error)
}

type ApiHandlerV1 struct {
	Ias iArchiverUsecase
}

func NewArchiverHandler(ias iArchiverUsecase, g gin.IRouter) *ApiHandlerV1 {
	ahv1 := &ApiHandlerV1{
		Ias: ias,
	}

	g.POST("/create-task", ahv1.CreateTask)
	g.PATCH("/:taskUuid/add-file-links", ahv1.AddFileLinks)
	g.GET("/:taskUuid/status", ahv1.GetStatus)
	return ahv1
}

func (ah ApiHandlerV1) CreateTask(context *gin.Context) {
	task, e := ah.Ias.CreateTask()
	fmt.Println(task, e)
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
		logrus.Error(e)
		context.JSON(http.StatusBadRequest, &domain.ErrorResponse{Message: e.Error()})
		return
	}

	resp, e := ah.Ias.AddLinks(taskUuid, addLinksReq)
	fmt.Println(resp, e)
	if e != nil {
		context.JSON(errToStatusCode(e), &domain.ErrorResponse{Message: e.Error()})
		return
	}
	context.JSON(http.StatusAccepted, resp)
}

func (ah ApiHandlerV1) GetStatus(context *gin.Context) {
	taskUuid := context.Param("taskUuid")
	resp, e := ah.Ias.GetStatus(taskUuid)
	fmt.Println(resp, e)
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

	logrus.Error(e)
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
