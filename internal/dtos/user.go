package dtos

import (
	"github.com/shunwuse/go-hris/internal/constants"
)

type GetUserResponse struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	CreatedTime     string `json:"created_time"`
	LastUpdatedTime string `json:"last_updated_time"`
}

type UserCreate struct {
	Username string         `json:"username" binding:"required" validate:"alphanum"`
	Password string         `json:"password" binding:"required"`
	Name     string         `json:"name" binding:"required"`
	Role     constants.Role `json:"role" binding:"required"`
}

type UserUpdate struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"omitempty"` // name is optional
}

type UserLogin struct {
	Username string `json:"username" binding:"required" validate:"alphanum"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Token    string   `json:"token"`
}
