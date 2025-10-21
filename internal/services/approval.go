package services

import (
	"context"
	"log/slog"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/approval"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
)

type approvalService struct {
	approvalRepository repositories.ApprovalRepository
}

func NewApprovalService(
	approvalRepository repositories.ApprovalRepository,
) service.ApprovalService {
	return approvalService{
		approvalRepository: approvalRepository,
	}
}

func (s approvalService) GetApprovals(ctx context.Context) ([]*entgen.Approval, error) {
	// var approvals []*entgen.Approval

	approvals, err := s.approvalRepository.Client.Approval.
		Query().
		WithCreator().
		WithApprover().
		All(ctx)
	if err != nil {
		slog.Error("Error getting approvals", "error", err)
		return nil, err
	}

	return approvals, nil
}

func (s approvalService) AddApproval(ctx context.Context, approval domains.ApprovalCreate) error {
	_, err := s.approvalRepository.Client.Approval.
		Create().
		SetStatus(approval.Status).
		SetCreatorID(approval.CreatorID).
		Save(ctx)
	if err != nil {
		slog.Error("Error adding approval", "error", err)
		return err
	}

	return nil
}

func (s approvalService) ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error {
	err := s.approvalRepository.Client.Approval.
		Update().
		Where(
			approval.IDEQ(approvalID),
			approval.StatusEQ(constants.ApprovalStatusPending),
		).
		SetStatus(action).
		SetApproverID(approverID).
		Exec(ctx)
	if err != nil {
		slog.Error("Error updating approval", "error", err)
		return err
	}

	// if result.RowsAffected == 0 {
	// 	slog.Error("Error updating approval: approval not found or already actioned")
	// 	return errors.New("approval not found or already actioned")
	// }

	return nil
}
