package routes

import (
	"github.com/go-chi/chi/v5"
)

type Routes []IRoute

type IRoute interface {
	Setup(router chi.Router)
}

func NewRoutes(
	monitorRoute *MonitorRoute,
	userRoute *UserRoute,
	approvalRoute *ApprovalRoute,
	authRoute *AuthRoute,
) Routes {
	return Routes{
		monitorRoute,
		userRoute,
		approvalRoute,
		authRoute,
	}
}

func (r Routes) Setup(router chi.Router) {
	for _, route := range r {
		route.Setup(router)
	}
}
