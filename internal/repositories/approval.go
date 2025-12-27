package repositories

import (
	"context"
	"strconv"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/approval"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/utils"
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
	approvals, err := r.GetClient(ctx).Approval.Query().
		WithCreator().
		WithApprover().
		All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find all approvals", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return approvals, nil
}

func (r *ApprovalRepository) FindAllWithCursor(ctx context.Context, query domains.CursorQuery) (*domains.CursorResult[*entgen.Approval], error) {
	dbQuery := r.GetClient(ctx).Approval.Query().
		WithCreator().
		WithApprover().
		Order(entgen.Desc(approval.FieldID)).
		Limit(query.Limit + 1)

	if query.Cursor != "" {
		parts, err := utils.DecodeCursor(query.Cursor)
		if err != nil || len(parts) == 0 {
			r.logger.WithContext(ctx).Error("failed to decode cursor", zap.Error(err))
			return nil, errors.ErrInvalidInput
		}

		decodedID, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			r.logger.WithContext(ctx).Error("failed to parse cursor ID", zap.Error(err))
			return nil, errors.ErrInvalidInput
		}
		dbQuery = dbQuery.Where(approval.IDLT(uint(decodedID)))
	}

	approvals, err := dbQuery.All(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to find approvals with cursor", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	hasMore := len(approvals) > query.Limit
	var nextCursor string
	if hasMore {
		nextCursor = utils.EncodeCursor(approvals[query.Limit-1].ID)
		approvals = approvals[:query.Limit]
	}

	return &domains.CursorResult[*entgen.Approval]{
		Items:      approvals,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *ApprovalRepository) Create(ctx context.Context, approval *domains.ApprovalCreate) (*entgen.Approval, error) {
	appr, err := r.GetClient(ctx).Approval.Create().
		SetStatus(approval.Status).
		SetCreatorID(approval.CreatorID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create approval", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return appr, nil
}

func (r *ApprovalRepository) UpdateStatusByID(ctx context.Context, id uint, status constants.ApprovalStatus, approverID uint) error {
	affected, err := r.GetClient(ctx).Approval.Update().
		Where(
			approval.IDEQ(id),
			approval.StatusEQ(constants.ApprovalStatusPending),
		).
		SetStatus(status).
		SetApproverID(approverID).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to update approval status", zap.Error(err))
		return errors.ErrDatabaseError
	}

	if affected == 0 {
		return errors.ErrNotFound
	}

	return nil
}
