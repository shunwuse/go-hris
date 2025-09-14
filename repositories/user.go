package repositories

import "github.com/shunwuse/go-hris/lib"

type UserRepository struct {
	logger lib.Logger
	lib.Database
}

func NewUserRepository(
	logger lib.Logger,
	db lib.Database,
) UserRepository {
	return UserRepository{
		logger:   logger,
		Database: db,
	}
}
