package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/dtos"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/ports/service"
)

type ApprovalController struct {
	logger          lib.Logger
	approvalService service.ApprovalService
}

func NewApprovalController(
	logger lib.Logger,
	approvalService service.ApprovalService,
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
func (c ApprovalController) GetApprovals(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.Errorf("Error user not authorized to get approvals")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to get approvals",
		})
		return
	}

	approvals, err := c.approvalService.GetApprovals(r.Context())
	if err != nil {
		c.logger.Errorf("Error getting approvals: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error getting approvals",
		})
		return
	}

	approvalsResponse := make([]dtos.ApprovalResponse, 0)
	for _, approval := range approvals {
		approvalResponse := dtos.ApprovalResponse{
			ID:          approval.ID,
			CreatorName: approval.Edges.Creator.Name,
			Status:      string(approval.Status),
		}

		if approval.Edges.Approver != nil {
			approvalResponse.ApproverName = &approval.Edges.Approver.Name
		}

		approvalsResponse = append(approvalsResponse, approvalResponse)

	}

	render.JSON(w, r, map[string]any{
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
func (c ApprovalController) AddApproval(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.Contains(constants.PermissionCreateApproval); !hasPermission {
		c.logger.Errorf("Error user not authorized to add approval")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to add approval",
		})
		return
	}

	userID := token.UserID

	approval := domains.ApprovalCreate{
		CreatorID: userID,
		Status:    constants.ApprovalStatusPending,
	}

	err := c.approvalService.AddApproval(r.Context(), approval)
	if err != nil {
		c.logger.Errorf("Error adding approval: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "Error adding approval",
		})
		return
	}

	render.JSON(w, r, map[string]string{
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
func (c ApprovalController) ActionApproval(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// check all permissions
	if hasPermission := permissions.ContainsAll(constants.Permissions{
		constants.PermissionReadApproval,
		constants.PermissionActionApproval,
	}); !hasPermission {
		c.logger.Errorf("Error user not authorized to action approval")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "User not authorized to action approval",
		})
		return
	}

	userID := token.UserID

	var actionRequest dtos.ApprovalAction
	err := render.DecodeJSON(r.Body, &actionRequest)
	if err != nil {
		c.logger.Errorf("Error binding action request: %v", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid request",
		})
		return
	}

	approvalID := actionRequest.ID
	action := actionRequest.Action

	if !isActionValid(action) {
		c.logger.Errorf("Error invalid action: %v", action)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "Invalid action",
		})
		return
	}

	err = c.approvalService.ActionApproval(r.Context(), approvalID, action, userID)
	if err != nil {
		c.logger.Errorf("Error actioning approval: %v", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "Approval actioned successfully",
	})
}

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
