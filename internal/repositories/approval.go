package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/approval"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type ApprovalRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewApprovalRepository(
	logger *infra.Logger,
	db *infra.Database,
) *ApprovalRepository {
	return &ApprovalRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *ApprovalRepository) FindAllWithRelations(ctx context.Context) ([]*entgen.Approval, error) {
	approvals, err := r.Client.Approval.Query().
		WithCreator().
		WithApprover().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all approvals", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return approvals, nil
}

func (r *ApprovalRepository) Create(ctx context.Context, status constants.ApprovalStatus, creatorID uint) (*entgen.Approval, error) {
	approval, err := r.Client.Approval.Create().
		SetStatus(status).
		SetCreatorID(creatorID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create approval", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return approval, nil
}

func (r *ApprovalRepository) UpdateStatusByID(ctx context.Context, id uint, status constants.ApprovalStatus, approverID uint) error {
	err := r.Client.Approval.Update().
		Where(
			approval.IDEQ(id),
			approval.StatusEQ(constants.ApprovalStatusPending),
		).
		SetStatus(status).
		SetApproverID(approverID).
		Exec(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to update approval status", zap.Error(err))
		return errors.ErrDatabaseError
	}

	// TODO: Check if rows were affected

	return nil
}
