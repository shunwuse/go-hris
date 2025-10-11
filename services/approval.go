package services

import (
	"context"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/approval"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/ports/service"
	"github.com/shunwuse/go-hris/repositories"
)

type approvalService struct {
	logger             lib.Logger
	approvalRepository repositories.ApprovalRepository
}

func NewApprovalService(
	logger lib.Logger,
	approvalRepository repositories.ApprovalRepository,
) service.ApprovalService {
	return approvalService{
		logger:             logger,
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
		s.logger.Errorf("Error getting approvals: %v", err)
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
		s.logger.Errorf("Error adding approval: %v", err)
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
		s.logger.Errorf("Error updating approval: %v", err)
		return err
	}

	// if result.RowsAffected == 0 {
	// 	s.logger.Errorf("Error updating approval: approval not found or already actioned")
	// 	return errors.New("approval not found or already actioned")
	// }

	return nil
}
