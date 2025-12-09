package constants

// JWT custom claim keys.
const (
	ClaimUserID      = "user_id"
	ClaimUsername    = "username"
	ClaimCreatedAt   = "created_at"
	ClaimRoles       = "roles"
	ClaimPermissions = "permissions"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)
