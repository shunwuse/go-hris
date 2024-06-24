package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/controllers"
)

// ExampleRoute struct
type ExampleRoute struct {
	exampleController controllers.ExampleController
}

func (r ExampleRoute) Setup(router *gin.Engine) {
	router.GET("/ping", r.exampleController.Ping)
}

func NewExampleRoute() ExampleRoute {
	exampleController := controllers.NewExampleController()

	return ExampleRoute{
		exampleController: exampleController,
	}
}
