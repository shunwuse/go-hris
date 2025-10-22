package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/infra"
)

type UserRoute struct {
	logger         infra.Logger
	jwtMiddleware  middlewares.JWTMiddleware
	userController controllers.UserController
}

func NewUserRoute(
	logger infra.Logger,
	jwtMiddleware middlewares.JWTMiddleware,
	userController controllers.UserController,
) UserRoute {
	return UserRoute{
		logger:         logger,
		jwtMiddleware:  jwtMiddleware,
		userController: userController,
	}
}

func (r UserRoute) Setup(router chi.Router) {
	r.logger.Info("setting up user routes")

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.Use(r.jwtMiddleware.Handler())
		userRouter.Get("/", r.userController.GetUsers)
		userRouter.Post("/", r.userController.CreateUser)
		userRouter.Put("/", r.userController.UpdateUser)
	})

	router.Post("/login", r.userController.Login)
}
