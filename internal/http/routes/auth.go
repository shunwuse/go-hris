package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/infra/logger"
)

type AuthRoute struct {
	logger         *logger.Logger
	jwtMiddleware  *middlewares.JWTMiddleware
	authController *controllers.AuthController
}

func NewAuthRoute(
	log *logger.Logger,
	jwtMiddleware *middlewares.JWTMiddleware,
	authController *controllers.AuthController,
) *AuthRoute {
	return &AuthRoute{
		logger:         log,
		jwtMiddleware:  jwtMiddleware,
		authController: authController,
	}
}

func (r *AuthRoute) Setup(router chi.Router) {
	r.logger.Info("setting up auth routes")

	router.Post("/login", r.authController.Login)

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/refresh", r.authController.RefreshToken)

		authRouter.Group(func(protectedRouter chi.Router) {
			protectedRouter.Use(r.jwtMiddleware.Handler())
			protectedRouter.Post("/logout", r.authController.Logout)
			protectedRouter.Post("/logout-all", r.authController.LogoutAll)
		})
	})
}
