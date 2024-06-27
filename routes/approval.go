package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
)

type ApprovalRoute struct {
	logger             lib.Logger
	approvalController controllers.ApprovalController
}

func NewApprovalRoute() ApprovalRoute {
	logger := lib.GetLogger()

	approvalController := controllers.NewApprovalController()

	return ApprovalRoute{
		logger:             logger,
		approvalController: approvalController,
	}
}

func (r ApprovalRoute) Setup(router *gin.Engine) {
	r.logger.Info("Setting up approval routes")

	approvalRouter := router.Group("/approvals", middlewares.NewJWTMiddleware().Handler())
	approvalRouter.GET("", r.approvalController.GetApprovals)
	approvalRouter.POST("", r.approvalController.AddApproval)
}
