package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/infra"
)

type AuthRoute struct {
	logger         *infra.Logger
	authController *controllers.AuthController
}

func NewAuthRoute(
	logger *infra.Logger,
	authController *controllers.AuthController,
) *AuthRoute {
	return &AuthRoute{
		logger:         logger,
		authController: authController,
	}
}

func (r *AuthRoute) Setup(router chi.Router) {
	r.logger.Info("setting up auth routes")

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/refresh", r.authController.RefreshToken)
		authRouter.Post("/logout", r.authController.Logout)
	})
}
