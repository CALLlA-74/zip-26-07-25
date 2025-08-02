package controllers

import (
	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router       *gin.Engine
	apiV1Handler *ApiHandlerV1
}

func NewRouter(ias iArchiverUsecase) *Router {
	router := gin.Default()
	r := &Router{
		router:       router,
		apiV1Handler: NewArchiverHandler(ias, router.Group(config.ApiV1GroupName)),
	}
	NewLoadArchiverHandler(router.Group(config.LoadArchGroupName))

	return r
}
