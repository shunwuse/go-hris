package services

import (
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/repositories"
)

type ApprovalService struct {
	logger             lib.Logger
	approvalRepository repositories.ApprovalRepository
}

func NewApprovalService() ApprovalService {
	logger := lib.GetLogger()

	// Initialize repositories
	approvalRepository := repositories.NewApprovalRepository()

	return ApprovalService{
		logger:             logger,
		approvalRepository: approvalRepository,
	}
}

func (s ApprovalService) GetApprovals() ([]models.Approval, error) {
	var approvals []models.Approval

	result := s.approvalRepository.Preload("Creator").Preload("Approver").Find(&approvals)
	if result.Error != nil {
		s.logger.Errorf("Error getting approvals: %v", result.Error)
		return nil, result.Error
	}

	return approvals, nil
}
