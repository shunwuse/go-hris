package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/dtos"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/services"
)

type ApprovalController struct {
	logger          lib.Logger
	approvalService services.ApprovalService
}

func NewApprovalController(
	logger lib.Logger,
	approvalService services.ApprovalService,
) ApprovalController {
	return ApprovalController{
		logger:          logger,
		approvalService: approvalService,
	}
}

// GetApprovals godoc
//
// @Summary Get approvals
// @Description Get all approvals
// @Tags approvals
// @security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ApprovalResponse
// @Router /approvals [get]
func (c ApprovalController) GetApprovals(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.Errorf("Error user not authorized to get approvals")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to get approvals",
		})
		return
	}

	approvals, err := c.approvalService.GetApprovals(ctx)
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

// AddApproval godoc
//
// @Summary Add approval
// @Description Add a new approval
// @Tags approvals
// @security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {string} message "Approval added successfully"
// @Router /approvals [post]
func (c ApprovalController) AddApproval(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
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

	err := c.approvalService.AddApproval(ctx, approval)
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

// ActionApproval godoc
//
// @Summary Action approval
// @Description Action an approval
// @Tags approvals
// @security BasicAuth
// @Accept json
// @Produce json
// @Param action body dtos.ApprovalAction true "Approval action object"
// @Success 200 {string} message "Approval actioned successfully"
// @Router /approvals/action [put]
func (c ApprovalController) ActionApproval(ctx *gin.Context) {
	token := ctx.MustGet(constants.JWTClaims).(domains.TokenPayload)
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

	if !isActionValid(action) {
		c.logger.Errorf("Error invalid action: %v", action)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid action",
		})
		return
	}

	err = c.approvalService.ActionApproval(ctx, approvalID, action, userID)
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

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
