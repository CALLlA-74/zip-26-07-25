package controllers

import (
	"github.com/CALLlA-74/zip-service/config"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewLoadArchiverHandler(g gin.IRouter) {
	g.Use(static.Serve(config.LoadArchGroupName, static.LocalFile(config.DownloadPath, false)))
}
