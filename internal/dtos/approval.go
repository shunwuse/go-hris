package dtos

import (
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type GetApprovalsRequest struct {
	domains.CursorQuery

	Status constants.ApprovalStatus `schema:"status"`
}

type ApprovalResponse struct {
	ID           uint                     `json:"id"`
	CreatorName  string                   `json:"creator_name"`
	ApproverName *string                  `json:"approver_name"`
	Status       constants.ApprovalStatus `json:"status"`
}

type ApprovalPathParams struct {
	ID uint `schema:"id"`
}

type ApprovalAction struct {
	Action constants.ApprovalStatus `json:"action" binding:"required"`
}
