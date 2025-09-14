package routes

import (
	"github.com/gin-gonic/gin"
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

func (r UserRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up user routes")

	userRouter := router.Group("/users", r.jwtMiddleware.Handler())
	userRouter.GET("", r.userController.GetUsers)
	userRouter.POST("", r.userController.CreateUser)
	userRouter.PUT("", r.userController.UpdateUser)

	router.POST("/login", r.userController.Login)
}
