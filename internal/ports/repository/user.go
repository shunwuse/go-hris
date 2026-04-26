package repository

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]domains.User, error)
	FindAllWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[domains.User], error)
	FindByUsername(ctx context.Context, username string) (*domains.User, error)
	FindByID(ctx context.Context, id uint) (*domains.UserWithRoles, error)
	Create(ctx context.Context, username string, name string) (*domains.User, error)
	Update(ctx context.Context, id uint, name string) error
	CreatePassword(ctx context.Context, hash string, ownerID uint) error
	AssignRole(ctx context.Context, userID uint, roleID uint) error
}
