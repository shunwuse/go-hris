package service

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type ApprovalService interface {
	GetApprovals(ctx context.Context) ([]*entgen.Approval, error)
	AddApproval(ctx context.Context, approval domains.ApprovalCreate) error
	ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error
}
