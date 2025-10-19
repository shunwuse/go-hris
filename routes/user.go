package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
)

type UserRoute struct {
	logger         lib.Logger
	jwtMiddleware  middlewares.JWTMiddleware
	userController controllers.UserController
}

func NewUserRoute(
	logger lib.Logger,
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
	r.logger.Info("Setting up user routes")

	router.Route("/users", func(userRouter chi.Router) {
		userRouter.Use(r.jwtMiddleware.Handler())
		userRouter.Get("/", r.userController.GetUsers)
		userRouter.Post("/", r.userController.CreateUser)
		userRouter.Put("/", r.userController.UpdateUser)
	})

	router.Post("/login", r.userController.Login)
}
