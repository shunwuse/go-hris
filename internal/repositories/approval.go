package repositories

import (
	"github.com/shunwuse/go-hris/internal/infra"
)

type ApprovalRepository struct {
	logger infra.Logger
	infra.Database
}

func NewApprovalRepository(
	logger infra.Logger,
	db infra.Database,
) ApprovalRepository {
	return ApprovalRepository{
		logger:   logger,
		Database: db,
	}
}
