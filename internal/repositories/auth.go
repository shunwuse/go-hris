package repositories

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/refreshtoken"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type AuthRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewAuthRepository(
	logger *infra.Logger,
	db *infra.Database,
) *AuthRepository {
	return &AuthRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *AuthRepository) CreateRefreshToken(
	ctx context.Context,
	tokenHash string,
	userID uint,
	expiresAt time.Time,
) (*entgen.RefreshToken, error) {
	query := r.GetClient(ctx).RefreshToken.Create().
		SetTokenHash(tokenHash).
		SetUserID(userID).
		SetExpiresAt(expiresAt)

	token, err := query.Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create refresh token", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return token, nil
}

func (r *AuthRepository) FindValidRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*entgen.RefreshToken, error) {
	now := time.Now()

	token, err := r.GetClient(ctx).RefreshToken.Query().
		Where(
			refreshtoken.TokenHash(tokenHash),
			refreshtoken.Revoked(false),
			refreshtoken.ExpiresAtGT(now),
		).
		Only(ctx)
	if err != nil {
		if entgen.IsNotFound(err) {
			return nil, errors.ErrTokenInvalid
		}
		r.logger.WithContext(ctx).Error("failed to find valid refresh token", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return token, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	now := time.Now()

	affected, err := r.GetClient(ctx).RefreshToken.Update().
		Where(refreshtoken.TokenHash(tokenHash)).
		SetRevoked(true).
		SetRevokedAt(now).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to revoke refresh token", zap.Error(err))
		return errors.ErrDatabaseError
	}

	if affected == 0 {
		return errors.ErrTokenInvalid
	}

	return nil
}

func (r *AuthRepository) RevokeAllRefreshTokensForUser(ctx context.Context, userID uint) error {
	now := time.Now()

	err := r.GetClient(ctx).RefreshToken.Update().
		Where(
			refreshtoken.UserID(userID),
			refreshtoken.Revoked(false),
		).
		SetRevoked(true).
		SetRevokedAt(now).
		Exec(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to revoke all tokens for user", zap.Uint("user_id", userID), zap.Error(err))
		return errors.ErrDatabaseError
	}

	return nil
}
