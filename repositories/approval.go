package repositories

import "github.com/shunwuse/go-hris/lib"

type ApprovalRepository struct {
	logger lib.Logger
	lib.Database
}

func NewApprovalRepository() ApprovalRepository {
	logger := lib.GetLogger()
	db := lib.GetDatabase()

	return ApprovalRepository{
		logger:   logger,
		Database: db,
	}
}
