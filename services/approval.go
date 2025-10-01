package services

import (
	"context"
	"errors"

	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/repositories"
)

type ApprovalService struct {
	logger             lib.Logger
	approvalRepository repositories.ApprovalRepository
}

func NewApprovalService(
	logger lib.Logger,
	approvalRepository repositories.ApprovalRepository,
) ApprovalService {
	return ApprovalService{
		logger:             logger,
		approvalRepository: approvalRepository,
	}
}

func (s ApprovalService) GetApprovals(ctx context.Context) ([]models.Approval, error) {
	var approvals []models.Approval

	result := s.approvalRepository.Preload("Creator").Preload("Approver").Find(&approvals)
	if result.Error != nil {
		s.logger.Errorf("Error getting approvals: %v", result.Error)
		return nil, result.Error
	}

	return approvals, nil
}

func (s ApprovalService) AddApproval(ctx context.Context, approval models.Approval) error {
	result := s.approvalRepository.Create(&approval)
	if result.Error != nil {
		s.logger.Errorf("Error adding approval: %v", result.Error)
		return result.Error
	}

	return nil
}

func (s ApprovalService) ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error {
	result := s.approvalRepository.Where("id = ?", approvalID).Where("status = ?", constants.ApprovalStatusPending).Updates(models.Approval{
		Status:     action,
		ApproverID: &approverID,
	})

	if result.Error != nil {
		s.logger.Errorf("Error updating approval: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		s.logger.Errorf("Error updating approval: approval not found or already actioned")
		return errors.New("approval not found or already actioned")
	}

	return nil
}
