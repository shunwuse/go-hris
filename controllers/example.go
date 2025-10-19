package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/ports/service"
)

// ExampleController struct
type ExampleController struct {
	logger         lib.Logger
	exampleService service.ExampleService
}

func NewExampleController(
	logger lib.Logger,
	exampleService service.ExampleService,
) ExampleController {
	return ExampleController{
		logger:         logger,
		exampleService: exampleService,
	}
}

func (c ExampleController) Ping(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Ping controller")

	message := c.exampleService.Ping(r.Context())

	render.JSON(w, r, map[string]any{
		"message": message,
	})
}
