package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/request"
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
	identity, ok := request.GetIdentity(r)
	if !ok {
		c.logger.WithContext(r.Context()).Error("failed to get identity from context")
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	// Check if user has permission to read approvals.
	if hasPermission := identity.Permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to get approvals")
		response.Error(w, errors.ErrForbidden)
		return
	}

	// Get pagination parameters.
	var query dtos.ApprovalGetList
	if err := request.DecodeQuery(r, &query); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode query params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}
	query.Normalize()

	if err := query.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate query params", zap.Error(err))
		response.Error(w, err)
		return
	}

	filter := domains.ApprovalFilter{
		Status: query.Status,
	}

	result, err := c.approvalService.GetApprovalsWithCursor(r.Context(), query.CursorQuery, filter)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get approvals", zap.Error(err))
		response.Error(w, err)
		return
	}

	approvalsResponse := make([]dtos.ApprovalResponse, len(result.Items))
	for idx, approval := range result.Items {
		approvalResponse := dtos.ApprovalResponse{
			ID:          approval.ID,
			CreatorName: approval.Edges.Creator.Name,
			Status:      approval.Status,
		}

		if approval.Edges.Approver != nil {
			approvalResponse.ApproverName = &approval.Edges.Approver.Name
		}

		approvalsResponse[idx] = approvalResponse
	}

	response.CursorList(w, approvalsResponse, response.CursorPaginationMeta{
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	})
}

func (c *ApprovalController) GetApproval(w http.ResponseWriter, r *http.Request) {
	identity, ok := request.GetIdentity(r)
	if !ok {
		c.logger.WithContext(r.Context()).Error("failed to get identity from context")
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	// Check if user has permission to read approvals.
	if hasPermission := identity.Permissions.Contains(constants.PermissionReadApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to get approval")
		response.Error(w, errors.ErrForbidden)
		return
	}

	var pathParams dtos.ApprovalPathParams
	if err := request.DecodePath(r, &pathParams); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode path params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := pathParams.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate path params", zap.Error(err))
		response.Error(w, err)
		return
	}

	approval, err := c.approvalService.GetApprovalByID(r.Context(), pathParams.ID)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to get approval", zap.Error(err), zap.Uint("approval_id", pathParams.ID))
		response.Error(w, err)
		return
	}

	resp := dtos.ApprovalResponse{
		ID:          approval.ID,
		CreatorName: approval.Edges.Creator.Name,
		Status:      approval.Status,
	}

	if approval.Edges.Approver != nil {
		resp.ApproverName = &approval.Edges.Approver.Name
	}

	response.OK(w, resp)
}

func (c *ApprovalController) AddApproval(w http.ResponseWriter, r *http.Request) {
	identity, ok := request.GetIdentity(r)
	if !ok {
		c.logger.WithContext(r.Context()).Error("failed to get identity from context")
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	// Check if user has permission to create approvals.
	if hasPermission := identity.Permissions.Contains(constants.PermissionCreateApproval); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to add approval")
		response.Error(w, errors.ErrForbidden)
		return
	}

	approval := domains.ApprovalCreate{
		CreatorID: identity.UserID,
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
	identity, ok := request.GetIdentity(r)
	if !ok {
		c.logger.WithContext(r.Context()).Error("failed to get identity from context")
		response.Error(w, errors.ErrUnauthorized)
		return
	}

	// Check if user has permission to action approvals.
	if hasPermission := identity.Permissions.ContainsAll(constants.Permissions{
		constants.PermissionReadApproval,
		constants.PermissionActionApproval,
	}); !hasPermission {
		c.logger.WithContext(r.Context()).Error("user not authorized to action approval")
		response.Error(w, errors.ErrForbidden)
		return
	}

	var pathParams dtos.ApprovalPathParams
	if err := request.DecodePath(r, &pathParams); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode path params", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := pathParams.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate path params", zap.Error(err))
		response.Error(w, err)
		return
	}

	var actionRequest dtos.ApprovalAction
	if err := request.DecodeJSON(r, &actionRequest); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to decode action request", zap.Error(err))
		response.Error(w, errors.ErrInvalidInput)
		return
	}

	if err := actionRequest.Validate(); err != nil {
		c.logger.WithContext(r.Context()).Error("failed to validate action request", zap.Error(err))
		response.Error(w, err)
		return
	}

	err := c.approvalService.ActionApproval(r.Context(), pathParams.ID, actionRequest.Action, identity.UserID)
	if err != nil {
		c.logger.WithContext(r.Context()).Error("failed to action approval", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.OK(w, "approval actioned successfully")
}
