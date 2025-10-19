package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
)

// ExampleRoute struct
type ExampleRoute struct {
	logger            lib.Logger
	exampleController controllers.ExampleController
}

func (r ExampleRoute) Setup(router chi.Router) {
	r.logger.Info("Setting up example routes")

	router.Get("/ping", r.exampleController.Ping)
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
