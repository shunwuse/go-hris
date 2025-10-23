package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type ApprovalController struct {
	logger          *infra.Logger
	approvalService service.ApprovalService
}

func NewApprovalController(
	logger *infra.Logger,
	approvalService service.ApprovalService,
) *ApprovalController {
	return &ApprovalController{
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

	// Check if user has permission to read approvals.
	if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to get approvals")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "user not authorized to get approvals",
		})
		return
	}

	approvals, err := c.approvalService.GetApprovals(r.Context())
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get approvals", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "failed to get approvals",
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

	// Check if user has permission to create approvals.
	if hasPermission := permissions.Contains(constants.PermissionCreateApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to add approval")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "user not authorized to add approval",
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
		c.logger.WithContext(r.Context()).Error("failed to add approval", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "failed to add approval",
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "approval added successfully",
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

	// Check if user has permission to action approvals.
	if hasPermission := permissions.ContainsAll(constants.Permissions{
		constants.PermissionReadApproval,
		constants.PermissionActionApproval,
	}); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to action approval")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{
			"error": "user not authorized to action approval",
		})
		return
	}

	userID := token.UserID

	var actionRequest dtos.ApprovalAction
	err := render.DecodeJSON(r.Body, &actionRequest)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode action request", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid request",
		})
		return
	}

	approvalID := actionRequest.ID
	action := actionRequest.Action

	if !isActionValid(action) {
		c.logger.WithContext(r.Context()).Error("invalid approval action", zap.String("action", string(action)))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid action",
		})
		return
	}

	err = c.approvalService.ActionApproval(r.Context(), approvalID, action, userID)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to action approval", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "approval actioned successfully",
	})
}

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
