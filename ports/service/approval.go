package service

import (
	"context"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/ent/entgen"
)

type ApprovalService interface {
	GetApprovals(ctx context.Context) ([]*entgen.Approval, error)
	AddApproval(ctx context.Context, approval domains.ApprovalCreate) error
	ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error
}
