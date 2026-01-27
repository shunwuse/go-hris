package response

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/errors"
)

// domainCodeToHTTPStatus maps domain error codes to HTTP status codes.
// This mapping is HTTP-specific, not in domain layer.
func domainCodeToHTTPStatus(code string) int {
	if status, ok := codeToStatus[code]; ok {
		return status
	}

	return http.StatusInternalServerError
}

var codeToStatus = map[string]int{
	// 400 [Bad Request]
	errors.CodeInvalidInput: http.StatusBadRequest,
	// 401 [Unauthorized]
	errors.CodeUnauthorized:       http.StatusUnauthorized,
	errors.CodeInvalidCredentials: http.StatusUnauthorized,
	errors.CodeTokenExpired:       http.StatusUnauthorized,
	errors.CodeTokenInvalid:       http.StatusUnauthorized,
	// 403 [Forbidden]
	errors.CodeForbidden:           http.StatusForbidden,
	errors.CodeOperationNotAllowed: http.StatusForbidden,
	// 404 [Not Found]
	errors.CodeNotFound: http.StatusNotFound,
	// 409 [Conflict]
	errors.CodeAlreadyExists: http.StatusConflict,
	errors.CodeConflict:      http.StatusConflict,
	// 422 [Unprocessable Entity]
	errors.CodeValidationFailed: http.StatusUnprocessableEntity,
	// 500 [Internal Server Error]
	errors.CodeDatabaseError: http.StatusInternalServerError,
	errors.CodeInternalError: http.StatusInternalServerError,
}
