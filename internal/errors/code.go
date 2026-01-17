package errors

const (
	// Common errors
	CodeNotFound      = "NOT_FOUND"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeInvalidInput  = "INVALID_INPUT"
	CodeInternalError = "INTERNAL_ERROR"

	// Authentication & Authorization errors
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeTokenExpired       = "TOKEN_EXPIRED"
	CodeTokenInvalid       = "TOKEN_INVALID"

	// Validation errors
	CodeValidationFailed = "VALIDATION_FAILED"

	// Business logic errors
	CodeForbidden           = "FORBIDDEN"
	CodeOperationNotAllowed = "OPERATION_NOT_ALLOWED"
	CodeConflict            = "CONFLICT"

	// Infrastructure errors
	CodeDatabaseError = "DATABASE_ERROR"
)
