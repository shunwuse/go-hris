package domains

import (
	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/constants"
)

type UserWithPermissions struct {
	*entgen.User

	Permissions constants.Permissions
}

type UserCreate struct {
	Username string
	Name     string

	Password PasswordCreate
}

type UserUpdate struct {
	ID   uint
	Name string
}
