package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/dtos"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
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
	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.Errorf("Error user not authorized to get approvals")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to get approvals",
		})
		return
	}

	approvals, err := c.approvalService.GetApprovals()
	if err != nil {
		c.logger.Errorf("Error getting approvals: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting approvals",
		})
		return
	}

	approvalsResponse := make([]dtos.ApprovalResponse, 0)
	for _, approval := range approvals {
		approvalResponse := dtos.ApprovalResponse{
			ID:          approval.ID,
			CreatorName: approval.Creator.Name,
			Status:      string(approval.Status),
		}

		if approval.Approver != nil {
			approvalResponse.ApproverName = &approval.Approver.Name
		}

		approvalsResponse = append(approvalsResponse, approvalResponse)

	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": approvalsResponse,
	})
}

func (c ApprovalController) AddApproval(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionCreateApproval); !hasPermission {
		c.logger.Errorf("Error user not authorized to add approval")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to add approval",
		})
		return
	}

	userID := token.UserID

	approval := models.Approval{
		CreatorID: userID,
		Status:    constants.ApprovalStatusPending,
	}

	err := c.approvalService.AddApproval(approval)
	if err != nil {
		c.logger.Errorf("Error adding approval: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error adding approval",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Approval added successfully",
	})
}

func (c ApprovalController) ActionApproval(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(services.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.ContainsAll(constants.Permissions{
		constants.PermissionReadApproval,
		constants.PermissionActionApproval,
	}); !hasPermission {
		c.logger.Errorf("Error user not authorized to action approval")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to action approval",
		})
		return
	}

	userID := token.UserID

	var actionRequest dtos.ApprovalAction
	err := ctx.ShouldBindJSON(&actionRequest)
	if err != nil {
		c.logger.Errorf("Error binding action request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	approvalID := actionRequest.ID
	action := actionRequest.Action

	err = c.approvalService.ActionApproval(approvalID, action, userID)
	if err != nil {
		c.logger.Errorf("Error actioning approval: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Approval actioned successfully",
	})
}
