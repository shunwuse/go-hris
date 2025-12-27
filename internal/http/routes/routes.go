package routes

import (
	"github.com/go-chi/chi/v5"
)

type Routes []IRoute

type IRoute interface {
	Setup(router chi.Router)
}

func NewRoutes(
	healthRoute *HealthRoute,
	userRoute *UserRoute,
	approvalRoute *ApprovalRoute,
	authRoute *AuthRoute,
	metricsRoute *MetricsRoute,
) Routes {
	return Routes{
		healthRoute,
		userRoute,
		approvalRoute,
		authRoute,
		metricsRoute,
	}
}

func (r Routes) Setup(router chi.Router) {
	for _, route := range r {
		route.Setup(router)
	}
}
