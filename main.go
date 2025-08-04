package main

import (
	"fmt"
	"net/http"

	"github.com/CALLlA-74/zip-26-07-25/config"
	"github.com/CALLlA-74/zip-26-07-25/internal/controllers"
	"github.com/CALLlA-74/zip-26-07-25/internal/repository"
	"github.com/CALLlA-74/zip-26-07-25/internal/usecases"
	"github.com/gin-gonic/gin"
)

func main() {
	tr := repository.NewTaskRepo()
	fr := repository.NewFileRepo()
	auc := usecases.NewArchiverUC(tr, fr)
	router := controllers.NewRouter(auc)

	router.Router.GET("/manage/health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	router.Router.Run(fmt.Sprintf(":%d", config.HostPort))
}
