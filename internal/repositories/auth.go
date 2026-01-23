package repositories

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/ent/entgen/refreshtoken"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"go.uber.org/zap"
)

type AuthRepository struct {
	logger *logger.Logger
	*database.Database
}

func NewAuthRepository(
	log *logger.Logger,
	db *database.Database,
) repository.AuthRepository {
	return &AuthRepository{
		logger:   log,
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

func (r *AuthRepository) FindRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*entgen.RefreshToken, error) {
	token, err := r.GetClient(ctx).RefreshToken.Query().
		Where(refreshtoken.TokenHash(tokenHash)).
		Only(ctx)
	if err != nil {
		if entgen.IsNotFound(err) {
			return nil, errors.ErrTokenInvalid
		}
		r.logger.WithContext(ctx).Error("failed to find refresh token", zap.Error(err))
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

func (r *AuthRepository) DeleteExpiredTokens(ctx context.Context) (int, error) {
	affected, err := r.GetClient(ctx).RefreshToken.Delete().
		Where(refreshtoken.ExpiresAtLT(time.Now())).
		Exec(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to delete expired tokens", zap.Error(err))
		return 0, errors.ErrDatabaseError
	}

	return affected, nil
}
