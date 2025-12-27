package service

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]*entgen.User, error)
	CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error
	GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error)
	GetUserByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error)
	UpdateUser(ctx context.Context, user *domains.UserUpdate) error
}
