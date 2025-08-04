package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/gin-contrib/static"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func NewLoadArchiverHandler(g gin.IRouter) {
	g.Use(func(ctx *gin.Context) {
		logrus.Infof("Load by url: %s", ctx.Request.URL.Path)
		ctx.Next()
		if ctx.Writer == nil || ctx.Writer.Status() >= http.StatusBadRequest {
			logrus.Infoln("Writer: ", ctx.Writer)
			return
		}

		splt := strings.Split(ctx.Request.URL.Path, "/")
		ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", splt[len(splt)-1]))
	}).StaticFS("/", static.LocalFile(config.DownloadPath, false))
}
