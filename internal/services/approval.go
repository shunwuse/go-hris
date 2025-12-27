package services

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
	"go.uber.org/zap"
)

type approvalService struct {
	logger             *infra.Logger
	approvalRepository *repositories.ApprovalRepository
}

func NewApprovalService(
	logger *infra.Logger,
	approvalRepository *repositories.ApprovalRepository,
) service.ApprovalService {
	return &approvalService{
		logger:             logger,
		approvalRepository: approvalRepository,
	}
}

func (s *approvalService) GetApprovals(ctx context.Context) ([]*entgen.Approval, error) {
	return s.approvalRepository.FindAllWithRelations(ctx)
}

func (s *approvalService) AddApproval(ctx context.Context, approval *domains.ApprovalCreate) error {
	_, err := s.approvalRepository.Create(ctx, approval)

	return err
}

func (s *approvalService) ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus, approverID uint) error {
	if !isActionValid(action) {
		s.logger.WithContext(ctx).Error("invalid approval action", zap.String("action", string(action)))
		return errors.ErrValidationFailed
	}

	return s.approvalRepository.UpdateStatusByID(ctx, approvalID, action, approverID)
}

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
