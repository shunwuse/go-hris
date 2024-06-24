package routes

import "github.com/gin-gonic/gin"

type Routes []IRoute

type IRoute interface {
	Setup(router *gin.Engine)
}

func NewRoutes() Routes {
	exampleRoute := NewExampleRoute()

	return Routes{
		exampleRoute,
	}
}

func (r Routes) Setup(router *gin.Engine) {
	for _, route := range r {
		route.Setup(router)
	}
}
