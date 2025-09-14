package routes

import (
	"github.com/gin-gonic/gin"
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

func (r ApprovalRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up approval routes")

	approvalRouter := router.Group("/approvals", r.jwtMiddleware.Handler())
	approvalRouter.GET("", r.approvalController.GetApprovals)
	approvalRouter.POST("", r.approvalController.AddApproval)
	approvalRouter.PUT("/action", r.approvalController.ActionApproval)
}
