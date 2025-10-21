package repositories

import (
	"github.com/shunwuse/go-hris/internal/infra"
)

type UserRepository struct {
	infra.Database
}

func NewUserRepository(
	db infra.Database,
) UserRepository {
	return UserRepository{
		Database: db,
	}
}
