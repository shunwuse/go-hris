package dtos

import (
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

type UserGetList struct {
	domains.OffsetQuery

	Name string         `schema:"name"`
	Role constants.Role `schema:"role"`
}

type UserResponse struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	CreatedTime     string `json:"created_time"`
	LastUpdatedTime string `json:"last_updated_time"`
}

type UserPathParams struct {
	ID uint `schema:"id"`
}

type UserCreate struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Name     string         `json:"name"`
	Role     constants.Role `json:"role"`
}

type UserUpdate struct {
	Name string `json:"name"` // name is optional
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Username     string           `json:"username"`
	Roles        []constants.Role `json:"roles"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
}
