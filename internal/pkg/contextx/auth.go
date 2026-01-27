package contextx

import (
	"context"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

// WithIdentity returns a new context with the Identity.
func WithIdentity(ctx context.Context, identity *domains.Identity) context.Context {
	return context.WithValue(ctx, constants.JWTClaims, &domains.Claims{
		Identity: *identity,
	})
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
