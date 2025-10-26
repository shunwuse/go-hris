package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

// ExampleController struct
type ExampleController struct {
	logger         *infra.Logger
	exampleService service.ExampleService
}

func NewExampleController(
	logger *infra.Logger,
	exampleService service.ExampleService,
) *ExampleController {
	return &ExampleController{
		logger:         logger,
		exampleService: exampleService,
	}
}

func (c *ExampleController) Ping(w http.ResponseWriter, r *http.Request) {
	c.logger.WithContext(r.Context()).Info("ping controller invoked")

	message := c.exampleService.Ping(r.Context())

	response.OK(w, message)
}
