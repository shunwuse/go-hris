package routes

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/shunwuse/go-hris/docs/swagger"
	"github.com/shunwuse/go-hris/lib"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerRoute struct {
	logger lib.Logger
}

func NewSwaggerRoute(
	logger lib.Logger,
) SwaggerRoute {
	return SwaggerRoute{
		logger: logger,
	}
}

func (r SwaggerRoute) Setup(router chi.Router) {
	r.logger.Info("Setting up swagger routes")

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
