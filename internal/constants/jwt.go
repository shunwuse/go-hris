package constants

// JWT custom claim keys.
const (
	ClaimUsername    = "username"
	ClaimCreatedAt   = "created_at"
	ClaimRoles       = "roles"
	ClaimPermissions = "permissions"
	ClaimType        = "type"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)
