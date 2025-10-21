package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
)

type ApprovalRoute struct {
	jwtMiddleware      middlewares.JWTMiddleware
	approvalController controllers.ApprovalController
}

func NewApprovalRoute(
	jwtMiddleware middlewares.JWTMiddleware,
	approvalController controllers.ApprovalController,
) ApprovalRoute {
	return ApprovalRoute{
		jwtMiddleware:      jwtMiddleware,
		approvalController: approvalController,
	}
}

func (r ApprovalRoute) Setup(router chi.Router) {
	slog.Info("Setting up approval routes")

	router.Route("/approvals", func(approvalRouter chi.Router) {
		approvalRouter.Use(r.jwtMiddleware.Handler())
		approvalRouter.Get("/", r.approvalController.GetApprovals)
		approvalRouter.Post("/", r.approvalController.AddApproval)
		approvalRouter.Put("/action", r.approvalController.ActionApproval)
	})
}
