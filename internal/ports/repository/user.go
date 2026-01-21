package repository

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/domains"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*entgen.User, error)
	FindAllWithOffset(ctx context.Context, query domains.OffsetQuery, filter domains.UserFilter) (*domains.OffsetResult[*entgen.User], error)
	FindByUsername(ctx context.Context, username string) (*entgen.User, error)
	FindByID(ctx context.Context, id uint) (*entgen.User, error)
	Create(ctx context.Context, username string, name string) (*entgen.User, error)
	Update(ctx context.Context, id uint, name string) error
	CreatePassword(ctx context.Context, hash string, owner *entgen.User) (*entgen.Password, error)
	AssignRole(ctx context.Context, userID uint, roleID uint) (*entgen.UserRole, error)
}
