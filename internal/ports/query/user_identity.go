package query

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type UserIdentityReader interface {
	GetUserWithPermissionsByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error)
	GetUserWithPermissionsByID(ctx context.Context, userID uint) (*domains.UserWithPermissions, error)
}
