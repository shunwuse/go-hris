package repository

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type ApprovalRepository interface {
	FindAllWithRelations(ctx context.Context) ([]*entgen.Approval, error)
	FindAllWithCursor(ctx context.Context, query domains.CursorQuery, filter domains.ApprovalFilter) (*domains.CursorResult[*entgen.Approval], error)
	FindByID(ctx context.Context, id uint) (*entgen.Approval, error)
	Create(ctx context.Context, approval *domains.ApprovalCreate) (*entgen.Approval, error)
	UpdateStatusByID(ctx context.Context, id uint, status constants.ApprovalStatus, approverID uint) error
}
