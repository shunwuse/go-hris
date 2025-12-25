package controllers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/response"
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

func (c *ApprovalController) GetApprovals(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to read approvals.
	if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to get approvals")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	approvals, err := c.approvalService.GetApprovals(r.Context())
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get approvals", zap.Error(err))
		response.Error(w, err)
		return
	}

	approvalsResponse := make([]dtos.ApprovalResponse, 0)
	for _, approval := range approvals {
		approvalResponse := dtos.ApprovalResponse{
			ID:          approval.ID,
			CreatorName: approval.Edges.Creator.Name,
			Status:      approval.Status,
		}

		if approval.Edges.Approver != nil {
			approvalResponse.ApproverName = &approval.Edges.Approver.Name
		}

		approvalsResponse = append(approvalsResponse, approvalResponse)

	}

	response.List(w, approvalsResponse)
}

func (c *ApprovalController) AddApproval(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to create approvals.
	if hasPermission := permissions.Contains(constants.PermissionCreateApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to add approval")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	userID := token.UserID

	approval := domains.ApprovalCreate{
		CreatorID: userID,
		Status:    constants.ApprovalStatusPending,
	}

	err := c.approvalService.AddApproval(r.Context(), &approval)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to add approval", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Created(w, "approval added successfully")
}

func (c *ApprovalController) ActionApproval(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(constants.JWTClaims).(domains.TokenPayload)
	permissions := token.Permissions

	// Check if user has permission to action approvals.
	if hasPermission := permissions.ContainsAll(constants.Permissions{
		constants.PermissionReadApproval,
		constants.PermissionActionApproval,
	}); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to action approval")
		response.Error(w, errors.ErrInsufficientPermissions)
		return
	}

	userID := token.UserID

	var actionRequest dtos.ApprovalAction
	err := render.DecodeJSON(r.Body, &actionRequest)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode action request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	approvalID := actionRequest.ID
	action := actionRequest.Action

	if !isActionValid(action) {
		c.logger.WithContext(r.Context()).Error("invalid approval action", zap.String("action", string(action)))
		response.Error(w, errors.ErrValidationFailed)
		return
	}

	err = c.approvalService.ActionApproval(r.Context(), approvalID, action, userID)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to action approval", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "approval actioned successfully")
}

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
