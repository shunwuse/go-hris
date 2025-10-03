package service

import (
	"context"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/models"
)

type ApprovalService interface {
	GetApprovals(ctx context.Context) ([]models.Approval, error)
	AddApproval(ctx context.Context, approval models.Approval) error
	ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error
}
