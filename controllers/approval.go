package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/services"
)

type ApprovalController struct {
	logger          lib.Logger
	approvalService services.ApprovalService
}

func NewApprovalController() ApprovalController {
	logger := lib.GetLogger()

	// Initialize services
	approvalService := services.NewApprovalService()

	return ApprovalController{
		logger:          logger,
		approvalService: approvalService,
	}
}

func (c ApprovalController) GetApprovals(ctx *gin.Context) {
	approvals, err := c.approvalService.GetApprovals()
	if err != nil {
		c.logger.Errorf("Error getting approvals: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting approvals",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": approvals,
	})
}
