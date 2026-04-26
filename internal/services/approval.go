package services

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
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

func (s *approvalService) GetApprovals(ctx context.Context) ([]domains.Approval, error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadApproval) {
		s.logger.WithContext(ctx).Error("not authorized to get approvals")
		return nil, errors.ErrForbidden
	}

	return s.approvalRepository.FindAllWithRelations(ctx)
}

func (s *approvalService) GetApprovalsWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[domains.Approval], error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadApproval) {
		s.logger.WithContext(ctx).Error("not authorized to get approvals")
		return nil, errors.ErrForbidden
	}

	return s.approvalRepository.FindAllWithCursor(ctx, query, filter)
}

func (s *approvalService) GetApprovalByID(ctx context.Context, id uint) (*domains.Approval, error) {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return nil, errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionReadApproval) {
		s.logger.WithContext(ctx).Error("not authorized to get approval")
		return nil, errors.ErrForbidden
	}

	return s.approvalRepository.FindByID(ctx, id)
}

func (s *approvalService) AddApproval(ctx context.Context, approval *domains.ApprovalCreate) error {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return errors.ErrUnauthorized
	}

	if !identity.Can(constants.PermissionCreateApproval) {
		s.logger.WithContext(ctx).Error("not authorized to add approval")
		return errors.ErrForbidden
	}

	approval.CreatorID = identity.UserID

	_, err := s.approvalRepository.Create(ctx, approval)

	return err
}

func (s *approvalService) ActionApproval(ctx context.Context, approvalID uint, action constants.ApprovalStatus) error {
	identity, ok := contextx.GetIdentity(ctx)
	if !ok {
		s.logger.WithContext(ctx).Error("failed to get identity from context")
		return errors.ErrUnauthorized
	}

	if !identity.CanAll(constants.PermissionReadApproval, constants.PermissionActionApproval) {
		s.logger.WithContext(ctx).Error("not authorized to action approval")
		return errors.ErrForbidden
	}

	if !isActionValid(action) {
		s.logger.WithContext(ctx).Error("invalid approval action", zap.String("action", string(action)))
		return errors.ErrValidationFailed
	}

	return s.approvalRepository.UpdateStatusByID(ctx, approvalID, action, identity.UserID)
}

func isActionValid(action constants.ApprovalStatus) bool {
	return action == constants.ApprovalStatusApproved || action == constants.ApprovalStatusRejected
}
