package repositories

import (
	"github.com/shunwuse/go-hris/internal/infra"
)

type UserRepository struct {
	logger infra.Logger
	infra.Database
}

func NewUserRepository(
	logger infra.Logger,
	db infra.Database,
) UserRepository {
	return UserRepository{
		logger:   logger,
		Database: db,
	}
}
