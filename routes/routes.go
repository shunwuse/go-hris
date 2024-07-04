package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/middlewares"
)

type Routes []IRoute

type IRoute interface {
	Setup(router *gin.Engine)
}

func NewRoutes() Routes {
	exampleRoute := NewExampleRoute()
	UserRoute := NewUserRoute()
	ApprovalRoute := NewApprovalRoute()

	// swagger route
	swaggerRoute := NewSwaggerRoute()

	return Routes{
		exampleRoute,
		UserRoute,
		ApprovalRoute,
		swaggerRoute,
	}
}

func (r Routes) Setup(router *gin.Engine) {
	// register global middleware
	router.Use(middlewares.NewDBTrxMiddleware().Handler())

	for _, route := range r {
		route.Setup(router)
	}
}
