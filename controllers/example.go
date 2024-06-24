package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
)

// ExampleController struct
type ExampleController struct {
	logger lib.Logger
}

func NewExampleController() ExampleController {
	logger := lib.GetLogger()

	return ExampleController{
		logger: logger,
	}
}

func (c ExampleController) Ping(ctx *gin.Context) {
	c.logger.Info("Ping controller")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
