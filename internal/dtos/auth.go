package dtos

import (
	"github.com/shunwuse/go-hris/internal/errors"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (d *RefreshRequest) Validate() error {
	if d.RefreshToken == "" {
		return errors.ErrValidationFailed.WithDetails(map[string]string{
			"refresh_token": "refresh token is required",
		})
	}

	return nil
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (d *LogoutRequest) Validate() error {
	return nil
}
