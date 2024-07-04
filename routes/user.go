package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
)

type UserRoute struct {
	logger         lib.Logger
	userController controllers.UserController
}

func NewUserRoute() UserRoute {
	logger := lib.GetLogger()

	userController := controllers.NewUserController()

	return UserRoute{
		logger:         logger,
		userController: userController,
	}
}

func (r UserRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up user routes")

	userRouter := router.Group("/users", middlewares.NewJWTMiddleware().Handler())
	userRouter.GET("", r.userController.GetUsers)
	userRouter.POST("", r.userController.CreateUser)
	userRouter.PUT("", r.userController.UpdateUser)

	router.POST("/login", r.userController.Login)
}
