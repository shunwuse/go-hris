package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
)

type UserRoute struct {
	jwtMiddleware  middlewares.JWTMiddleware
	userController controllers.UserController
}

func NewUserRoute(
	jwtMiddleware middlewares.JWTMiddleware,
	userController controllers.UserController,
) UserRoute {
	return UserRoute{
		jwtMiddleware:  jwtMiddleware,
		userController: userController,
	}
}

func (r UserRoute) Setup(router chi.Router) {
	slog.Info("Setting up user routes")

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.Use(r.jwtMiddleware.Handler())
		userRouter.Get("/", r.userController.GetUsers)
		userRouter.Post("/", r.userController.CreateUser)
		userRouter.Put("/", r.userController.UpdateUser)
	})

	router.Post("/login", r.userController.Login)
}
