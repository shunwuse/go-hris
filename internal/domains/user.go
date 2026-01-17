package domains

import (
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
)

type UserWithPermissions struct {
	*entgen.User

	Permissions constants.Permissions
}

type UserFilter struct {
	Name string
	Role constants.Role
}

type UserCreate struct {
	Username string
	Name     string
	Password string
}

type UserUpdate struct {
	ID   uint
	Name string
}
