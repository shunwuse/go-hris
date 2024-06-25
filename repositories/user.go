package repositories

import "github.com/shunwuse/go-hris/lib"

type UserRepository struct {
	logger lib.Logger
	lib.Database
}

func NewUserRepository() UserRepository {
	logger := lib.GetLogger()
	db := lib.GetDatabase()

	return UserRepository{
		logger:   logger,
		Database: db,
	}
}
