package service

import (
	"context"

	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/models"
)

type AuthService interface {
	GenerateToken(ctx context.Context, user *models.User) (string, error)
	AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error)
}
