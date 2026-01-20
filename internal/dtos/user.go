package dtos

import (
	"slices"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
)

type UserGetList struct {
	domains.OffsetQuery

	Name string         `schema:"name"`
	Role constants.Role `schema:"role"`
}

func (d *UserGetList) Validate() error {
	if d.Role != "" {
		if !slices.Contains(
			[]constants.Role{
				constants.Admin,
				constants.Manager,
				constants.Staff},
			d.Role,
		) {
			return errors.ErrValidationFailed
		}
	}

	return nil
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

func (d *UserPathParams) Validate() error {
	if d.ID == 0 {
		return errors.ErrValidationFailed.WithDetails(map[string]string{
			"id": "id must be greater than 0",
		})
	}

	return nil
}

type UserCreate struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	Name     string         `json:"name"`
	Role     constants.Role `json:"role"`
}

func (d *UserCreate) Validate() error {
	details := make(map[string]string)

	if d.Username == "" {
		details["username"] = "username is required"
	}

	if d.Password == "" {
		details["password"] = "password is required"
	} else if len(d.Password) < 6 {
		details["password"] = "password must be at least 6 characters"
	}

	if d.Name == "" {
		details["name"] = "name is required"
	}

	if d.Role == "" {
		details["role"] = "role is required"
	} else if !slices.Contains(
		[]constants.Role{
			constants.Manager,
			constants.Staff},
		d.Role,
	) {
		details["role"] = "invalid role"
	}

	if len(details) > 0 {
		return errors.ErrValidationFailed.WithDetails(details)
	}

	return nil
}

type UserUpdate struct {
	Name string `json:"name"` // name is optional
}

func (d *UserUpdate) Validate() error {
	if d.Name == "" {
		return errors.ErrValidationFailed.WithDetails(map[string]string{
			"name": "name is required",
		})
	}

	return nil
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d *UserLogin) Validate() error {
	details := make(map[string]string)

	if d.Username == "" {
		details["username"] = "username is required"
	}

	if d.Password == "" {
		details["password"] = "password is required"
	}

	if len(details) > 0 {
		return errors.ErrValidationFailed.WithDetails(details)
	}

	return nil
}

type UserLoginResponse struct {
	Username     string           `json:"username"`
	Roles        []constants.Role `json:"roles"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
}
