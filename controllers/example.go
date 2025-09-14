package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/services"
)

// ExampleController struct
type ExampleController struct {
	logger         lib.Logger
	exampleService services.ExampleService
}

func NewExampleController(
	logger lib.Logger,
	exampleService services.ExampleService,
) ExampleController {
	return ExampleController{
		logger:         logger,
		exampleService: exampleService,
	}
}

func (c ExampleController) Ping(ctx *gin.Context) {
	c.logger.Info("Ping controller")

	message := c.exampleService.Ping()

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
