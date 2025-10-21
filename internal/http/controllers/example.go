package controllers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

// ExampleController struct
type ExampleController struct {
	exampleService service.ExampleService
}

func NewExampleController(
	exampleService service.ExampleService,
) ExampleController {
	return ExampleController{
		exampleService: exampleService,
	}
}

func (c ExampleController) Ping(w http.ResponseWriter, r *http.Request) {
	slog.Info("Ping controller")

	message := c.exampleService.Ping(r.Context())

	render.JSON(w, r, map[string]any{
		"message": message,
	})
}
