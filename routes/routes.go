package routes

import (
	"github.com/gin-gonic/gin"
)

type Routes []IRoute

type IRoute interface {
	Setup(router *gin.Engine)
}

func NewRoutes(
	exampleRoute ExampleRoute,
	userRoute UserRoute,
	approvalRoute ApprovalRoute,
	swaggerRoute SwaggerRoute,
) Routes {
	return Routes{
		exampleRoute,
		userRoute,
		approvalRoute,
		swaggerRoute,
	}
}

func (r Routes) Setup(router *gin.Engine) {
	for _, route := range r {
		route.Setup(router)
	}
}
