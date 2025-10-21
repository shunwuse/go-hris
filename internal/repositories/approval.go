package repositories

import (
	"github.com/shunwuse/go-hris/internal/infra"
)

type ApprovalRepository struct {
	infra.Database
}

func NewApprovalRepository(
	db infra.Database,
) ApprovalRepository {
	return ApprovalRepository{
		Database: db,
	}
}
