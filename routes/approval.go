package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
)

type ApprovalRoute struct {
	logger             lib.Logger
	jwtMiddleware      middlewares.JWTMiddleware
	approvalController controllers.ApprovalController
}

func NewApprovalRoute(
	logger lib.Logger,
	jwtMiddleware middlewares.JWTMiddleware,
	approvalController controllers.ApprovalController,
) ApprovalRoute {
	return ApprovalRoute{
		logger:             logger,
		jwtMiddleware:      jwtMiddleware,
		approvalController: approvalController,
	}
}

func (r ApprovalRoute) Setup(router chi.Router) {
	r.logger.Info("Setting up approval routes")

	router.Route("/approvals", func(approvalRouter chi.Router) {
		approvalRouter.Use(r.jwtMiddleware.Handler())
		approvalRouter.Get("/", r.approvalController.GetApprovals)
		approvalRouter.Post("/", r.approvalController.AddApproval)
		approvalRouter.Put("/action", r.approvalController.ActionApproval)
	})
}
