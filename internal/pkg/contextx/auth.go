package contextx

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

// SystemIdentity is a predefined identity for internal system operations.
var SystemIdentity = &domains.Identity{
	UserID:      0,
	Username:    "system",
	Roles:       []constants.Role{constants.Admin},
	Permissions: constants.AllPermissions,
}

// WithIdentity returns a new context with the Identity.
func WithIdentity(ctx context.Context, identity *domains.Identity) context.Context {
	return context.WithValue(ctx, constants.JWTClaims, &domains.Claims{
		Identity: *identity,
	})
}

// WithSystemIdentity returns a new context with the SystemIdentity.
func WithSystemIdentity(ctx context.Context) context.Context {
	return WithIdentity(ctx, SystemIdentity)
}

// GetIdentity extracts the Identity from the context.
func GetIdentity(ctx context.Context) (*domains.Identity, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return nil, false
	}

	return &claims.Identity, true
}

// WithClaims returns a new context with the Claims.
func WithClaims(ctx context.Context, claims *domains.Claims) context.Context {
	return context.WithValue(ctx, constants.JWTClaims, claims)
}

// GetClaims extracts the Claims from the context.
func GetClaims(ctx context.Context) (*domains.Claims, bool) {
	val := ctx.Value(constants.JWTClaims)
	if val == nil {
		return nil, false
	}

	claims, ok := val.(*domains.Claims)
	return claims, ok
}
