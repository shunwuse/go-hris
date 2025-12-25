package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/infra"
)

// HealthRoute struct
type HealthRoute struct {
	logger           *infra.Logger
	healthController *controllers.HealthController
}

func NewHealthRoute(
	logger *infra.Logger,
	healthController *controllers.HealthController,
) *HealthRoute {
	return &HealthRoute{
		logger:           logger,
		healthController: healthController,
	}
}

func (r *HealthRoute) Setup(router chi.Router) {
	r.logger.Info("setting up health routes")

	router.Get("/health", r.healthController.Check)
}
