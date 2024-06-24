package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExampleRoute struct
type ExampleRoute struct {
}

func (r ExampleRoute) Setup(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func NewExampleRoute() ExampleRoute {
	return ExampleRoute{}
}
