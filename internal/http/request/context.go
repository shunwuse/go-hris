package request

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
)

// GetClaims extracts the Claims from the request context.
// It returns the claims and true if found and valid, otherwise nil and false.
func GetClaims(r *http.Request) (*domains.Claims, bool) {
	val := r.Context().Value(constants.JWTClaims)
	if val == nil {
		return nil, false
	}

	claims, ok := val.(*domains.Claims)
	return claims, ok
}

// GetIdentity extracts the Identity from the request context.
// It returns the identity and true if found and valid, otherwise nil and false.
func GetIdentity(r *http.Request) (*domains.Identity, bool) {
	claims, ok := GetClaims(r)
	if !ok {
		return nil, false
	}

	return &claims.Identity, true
}
