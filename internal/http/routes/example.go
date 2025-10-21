package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
)

// ExampleRoute struct
type ExampleRoute struct {
	exampleController controllers.ExampleController
}

func (r ExampleRoute) Setup(router chi.Router) {
	slog.Info("Setting up example routes")

	router.Get("/ping", r.exampleController.Ping)
}

func NewExampleRoute(
	exampleController controllers.ExampleController,
) ExampleRoute {
	return ExampleRoute{
		exampleController: exampleController,
	}
}
