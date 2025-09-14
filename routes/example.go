package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
)

// ExampleRoute struct
type ExampleRoute struct {
	logger            lib.Logger
	exampleController controllers.ExampleController
}

func (r ExampleRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up example routes")

	router.GET("/ping", r.exampleController.Ping)
}

func NewExampleRoute(
	logger lib.Logger,
	exampleController controllers.ExampleController,
) ExampleRoute {
	return ExampleRoute{
		logger:            logger,
		exampleController: exampleController,
	}
}
