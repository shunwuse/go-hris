package domains

import (
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
)

type User struct {
	ID        uint
	Username  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserWithRoles struct {
	ID           uint
	Username     string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PasswordHash string
	Roles        []constants.Role
}

type UserWithPermissions struct {
	ID           uint
	Username     string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PasswordHash string
	Roles        []constants.Role
	Permissions  constants.Permissions
}

type Role struct {
	ID   uint
	Name constants.Role
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
