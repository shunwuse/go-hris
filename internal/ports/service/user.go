package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]domains.User, error)
	GetUsersWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[domains.User], error)
	GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error)
	GetUserByID(ctx context.Context, id uint) (*domains.UserWithPermissions, error)
	CreateUser(ctx context.Context, user *domains.UserCreate, role constants.Role) error
	UpdateUser(ctx context.Context, user *domains.UserUpdate) error
}
