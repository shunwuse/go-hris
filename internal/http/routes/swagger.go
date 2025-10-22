package routes

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/shunwuse/go-hris/docs/swagger"
	"github.com/shunwuse/go-hris/internal/infra"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerRoute struct {
	logger infra.Logger
}

func NewSwaggerRoute(
	logger infra.Logger,
) SwaggerRoute {
	return SwaggerRoute{
		logger: logger,
	}
}

func (r SwaggerRoute) Setup(router chi.Router) {
	r.logger.Info("setting up swagger routes")

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
