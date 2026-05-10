package token

import (
	"context"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	appconfig "github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestJWTService(t *testing.T) (*JWTService, *cache.Cache, string) {
	t.Helper()

	const secret = "test-jwt-secret"

	log := logger.NewNopLogger()
	c := cache.New(&cache.Config{UseMiniredis: true}, log)
	t.Cleanup(func() {
		_ = c.Close()
	})

	svc, ok := New(&appconfig.AuthConfig{
		JWTSecret:               secret,
		JWTExpireMinutes:        60,
		JWTRefreshExpireMinutes: 60,
	}, log, c).(*JWTService)
	require.True(t, ok)

	return svc, c, secret
}

func signTestAccessToken(t *testing.T, secret string, subject string, jti string, expiration time.Time) string {
	t.Helper()

	key, err := jwk.FromRaw([]byte(secret))
	require.NoError(t, err)
	require.NoError(t, key.Set(jwk.AlgorithmKey, jwa.HS256))
	require.NoError(t, key.Set(jwk.KeyIDKey, "test-kid"))

	token, err := jwt.NewBuilder().
		IssuedAt(time.Now().Add(-time.Minute)).
		Expiration(expiration).
		JwtID(jti).
		Subject(subject).
		Claim(constants.ClaimType, string(constants.TokenTypeAccess)).
		Build()
	require.NoError(t, err)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	require.NoError(t, err)

	return string(signed)
}

func parseTestAccessToken(t *testing.T, secret, tokenString string) jwt.Token {
	t.Helper()

	key, err := jwk.FromRaw([]byte(secret))
	require.NoError(t, err)
	require.NoError(t, key.Set(jwk.AlgorithmKey, jwa.HS256))
	require.NoError(t, key.Set(jwk.KeyIDKey, "test-kid"))

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256, key), jwt.WithValidate(false))
	require.NoError(t, err)

	return token
}

func TestJWTService_GenerateAndValidateAccessToken(t *testing.T) {
	svc, _, _ := newTestJWTService(t)

	ctx := context.Background()
	accessToken, err := svc.GenerateAccessToken(ctx, &domains.UserWithPermissions{ID: 42})
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	claims, err := svc.ValidateAccessToken(ctx, accessToken)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, uint(42), claims.Identity.UserID)
	assert.NotEmpty(t, claims.JTI)
	assert.False(t, claims.ExpiresAt.IsZero())
}

func TestJWTService_ValidateAccessToken_ExpiredToken(t *testing.T) {
	svc, _, secret := newTestJWTService(t)

	tokenString := signTestAccessToken(t, secret, "42", "expired-jti", time.Now().Add(-time.Hour))

	claims, err := svc.ValidateAccessToken(context.Background(), tokenString)
	require.ErrorIs(t, err, errors.ErrTokenExpired)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateAccessToken_BlacklistedToken(t *testing.T) {
	svc, _, secret := newTestJWTService(t)

	tokenString := signTestAccessToken(t, secret, "42", "blacklisted-jti", time.Now().Add(time.Hour))
	token := parseTestAccessToken(t, secret, tokenString)

	require.NoError(t, svc.BlacklistToken(context.Background(), token.JwtID(), time.Hour))

	claims, err := svc.ValidateAccessToken(context.Background(), tokenString)
	require.ErrorIs(t, err, errors.ErrTokenInvalid)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateAccessToken_BlacklistLookupFailure(t *testing.T) {
	svc, c, secret := newTestJWTService(t)

	tokenString := signTestAccessToken(t, secret, "42", "cache-failure-jti", time.Now().Add(time.Hour))

	require.NoError(t, c.Close())

	claims, err := svc.ValidateAccessToken(context.Background(), tokenString)
	require.ErrorIs(t, err, errors.ErrServiceUnavailable)
	assert.Nil(t, claims)
}

func TestJWTService_BlacklistToken_StorageFailure(t *testing.T) {
	svc, c, _ := newTestJWTService(t)

	require.NoError(t, c.Close())

	err := svc.BlacklistToken(context.Background(), "blacklist-jti", time.Hour)
	require.ErrorIs(t, err, errors.ErrServiceUnavailable)
}

func TestJWTService_ValidateAccessToken_InvalidTokenString(t *testing.T) {
	svc, _, _ := newTestJWTService(t)

	claims, err := svc.ValidateAccessToken(context.Background(), "not-a-token")
	require.ErrorIs(t, err, errors.ErrTokenInvalid)
	assert.Nil(t, claims)
}
