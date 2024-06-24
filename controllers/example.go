package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExampleController struct
type ExampleController struct {
}

func NewExampleController() ExampleController {
	return ExampleController{}
}

func (c ExampleController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
