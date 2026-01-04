package dtos

import "github.com/shunwuse/go-hris/internal/constants"

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
