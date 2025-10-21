package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	_ "github.com/shunwuse/go-hris/docs/swagger"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerRoute struct {
}

func NewSwaggerRoute() SwaggerRoute {
	return SwaggerRoute{}
}

func (r SwaggerRoute) Setup(router chi.Router) {
	slog.Info("Setting up swagger routes")

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
