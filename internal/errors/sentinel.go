package errors

// Sentinel errors - predefined common errors
var (
	// Common errors
	ErrNotFound      = &Error{Code: CodeNotFound, Message: "resource not found"}
	ErrAlreadyExists = &Error{Code: CodeAlreadyExists, Message: "resource already exists"}
	ErrInvalidInput  = &Error{Code: CodeInvalidInput, Message: "invalid input"}
	ErrInternalError = &Error{Code: CodeInternalError, Message: "internal error"}

	// Authentication & Authorization errors
	ErrUnauthorized       = &Error{Code: CodeUnauthorized, Message: "unauthorized"}
	ErrForbidden          = &Error{Code: CodeForbidden, Message: "forbidden"}
	ErrInvalidCredentials = &Error{Code: CodeInvalidCredentials, Message: "invalid credentials"}
	ErrTokenExpired       = &Error{Code: CodeTokenExpired, Message: "token expired"}
	ErrTokenInvalid       = &Error{Code: CodeTokenInvalid, Message: "token invalid"}

	// Validation errors
	ErrValidationFailed = &Error{Code: CodeValidationFailed, Message: "validation failed"}

	// Business logic errors
	ErrInsufficientPermissions = &Error{Code: CodeInsufficientPermissions, Message: "insufficient permissions"}
	ErrOperationNotAllowed     = &Error{Code: CodeOperationNotAllowed, Message: "operation not allowed"}
	ErrConflict                = &Error{Code: CodeConflict, Message: "conflict"}

	// Infrastructure errors
	ErrDatabaseError = &Error{Code: CodeDatabaseError, Message: "database error"}
)
