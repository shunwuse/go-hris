package routes

import (
	"github.com/go-chi/chi/v5"
)

type Routes []IRoute

type IRoute interface {
	Setup(router chi.Router)
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

func (r Routes) Setup(router chi.Router) {
	for _, route := range r {
		route.Setup(router)
	}
}
