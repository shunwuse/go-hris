package services

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type approvalService struct {
	logger             *logger.Logger
	approvalRepository repository.ApprovalRepository
}

func NewApprovalService(
	log *logger.Logger,
	approvalRepository repository.ApprovalRepository,
) service.ApprovalService {
	return &approvalService{
		logger:             log,
		approvalRepository: approvalRepository,
	}
}

func (s *approvalService) GetApprovals(ctx context.Context) ([]*entgen.Approval, error) {
	return s.approvalRepository.FindAllWithRelations(ctx)
}

func (s *approvalService) GetApprovalsWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[*entgen.Approval], error) {
	return s.approvalRepository.FindAllWithCursor(ctx, query, filter)
}

func (s *approvalService) GetApprovalByID(ctx context.Context, id uint) (*entgen.Approval, error) {
	return s.approvalRepository.FindByID(ctx, id)
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
