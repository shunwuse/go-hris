package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type ApprovalService interface {
	GetApprovals(ctx context.Context) ([]domains.Approval, error)
	GetApprovalsWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[domains.Approval], error)
	GetApprovalByID(ctx context.Context, id uint) (*domains.Approval, error)
	AddApproval(ctx context.Context, approval *domains.ApprovalCreate) error
	ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus) error
}
