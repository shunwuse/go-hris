package dtos

import "github.com/shunwuse/go-hris/constants"

type ApprovalResponse struct {
	ID           uint    `json:"id"`
	CreatorName  string  `json:"creator_name"`
	ApproverName *string `json:"approver_name"`
	Status       string  `json:"status"`
}

type ApprovalAction struct {
	ID     uint                     `json:"id" binding:"required"`
	Action constants.ApprovalStatus `json:"action" binding:"required"`
}
