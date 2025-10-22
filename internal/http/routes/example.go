package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/infra"
)

// ExampleRoute struct
type ExampleRoute struct {
	logger            infra.Logger
	exampleController controllers.ExampleController
}

func (r ExampleRoute) Setup(router chi.Router) {
	r.logger.Info("setting up example routes")

	router.Get("/ping", r.exampleController.Ping)
}

func NewExampleRoute(
	logger infra.Logger,
	exampleController controllers.ExampleController,
) ExampleRoute {
	return ExampleRoute{
		logger:            logger,
		exampleController: exampleController,
	}
}
