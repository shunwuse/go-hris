package repositories

import "github.com/shunwuse/go-hris/lib"

type ApprovalRepository struct {
	logger lib.Logger
	lib.Database
}

func NewApprovalRepository(
	logger lib.Logger,
	db lib.Database,
) ApprovalRepository {
	return ApprovalRepository{
		logger:   logger,
		Database: db,
	}
}
