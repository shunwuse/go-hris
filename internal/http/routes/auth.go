package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/infra"
)

type AuthRoute struct {
	logger         *infra.Logger
	jwtMiddleware  *middlewares.JWTMiddleware
	authController *controllers.AuthController
}

func NewAuthRoute(
	logger *infra.Logger,
	jwtMiddleware *middlewares.JWTMiddleware,
	authController *controllers.AuthController,
) *AuthRoute {
	return &AuthRoute{
		logger:         logger,
		jwtMiddleware:  jwtMiddleware,
		authController: authController,
	}
}

func (r *AuthRoute) Setup(router chi.Router) {
	r.logger.Info("setting up auth routes")

	router.Post("/login", r.authController.Login)

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/refresh", r.authController.RefreshToken)
		authRouter.Post("/logout", r.authController.Logout)

		authRouter.Group(func(protectedRouter chi.Router) {
			protectedRouter.Use(r.jwtMiddleware.Handler())
			protectedRouter.Post("/logout-all", r.authController.LogoutAll)
		})
	})
}
