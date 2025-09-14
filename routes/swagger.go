package routes

import (
	"github.com/gin-gonic/gin"

	_ "github.com/shunwuse/go-hris/docs/swagger"
	"github.com/shunwuse/go-hris/lib"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRoute struct {
	logger lib.Logger
}

func NewSwaggerRoute(
	logger lib.Logger,
) SwaggerRoute {
	return SwaggerRoute{
		logger: logger,
	}
}

func (r SwaggerRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up swagger routes")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
