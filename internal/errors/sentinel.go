package errors

// Sentinel errors - predefined common errors
var (
	// Common errors
	ErrNotFound      = New(CodeNotFound, "resource not found")
	ErrAlreadyExists = New(CodeAlreadyExists, "resource already exists")
	ErrInvalidInput  = New(CodeInvalidInput, "invalid input")
	ErrInternalError = New(CodeInternalError, "internal error")

	// Authentication & Authorization errors
	ErrUnauthorized       = New(CodeUnauthorized, "unauthorized")
	ErrInvalidCredentials = New(CodeInvalidCredentials, "invalid credentials")
	ErrTokenExpired       = New(CodeTokenExpired, "token expired")
	ErrTokenInvalid       = New(CodeTokenInvalid, "token invalid")

	// Validation errors
	ErrValidationFailed = New(CodeValidationFailed, "validation failed")

	// Business logic errors
	ErrForbidden           = New(CodeForbidden, "forbidden")
	ErrOperationNotAllowed = New(CodeOperationNotAllowed, "operation not allowed")
	ErrConflict            = New(CodeConflict, "conflict")

	// Infrastructure errors
	ErrDatabaseError = New(CodeDatabaseError, "database error")
)
