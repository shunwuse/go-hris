package response

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"net/http"

	"github.com/shunwuse/go-hris/internal/errors"
)

// JSON encodes data as JSON response with status code
func JSON(w http.ResponseWriter, code int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(buf.Bytes()) //nolint:errcheck
}

// OK sends successful response with data (200 OK)
func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

// Created sends resource created response (201 Created)
func Created(w http.ResponseWriter, data any) {
	if data == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}

	JSON(w, http.StatusCreated, data)
}

// NoContent sends no content response (204 No Content)
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ServiceUnavailable sends service unavailable response (503 Service Unavailable)
func ServiceUnavailable(w http.ResponseWriter, data any) {
	JSON(w, http.StatusServiceUnavailable, data)
}

// List sends simple list response
func List(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

// OffsetList sends offset-based paginated list response
func OffsetList(w http.ResponseWriter, data any, meta OffsetPaginationMeta) {
	JSON(w, http.StatusOK, OffsetListResponse{
		Data: data,
		Meta: meta,
	})
}

// CursorList sends cursor-based paginated list response
func CursorList(w http.ResponseWriter, data any, meta CursorPaginationMeta) {
	JSON(w, http.StatusOK, CursorListResponse{
		Data: data,
		Meta: meta,
	})
}

// Error encodes domain error to HTTP error response
func Error(w http.ResponseWriter, err error) {
	// Extract custom Error (including sentinel errors)
	var appErr *errors.Error
	if stderrors.As(err, &appErr) {
		status := domainCodeToHTTPStatus(appErr.Code())
		JSON(w, status, ErrorResponse{
			Error: ErrorDetail{
				Code:    appErr.Code(),
				Message: appErr.Message(),
				Details: appErr.Details(),
			},
		})

		return
	}

	// Fallback for unexpected errors
	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code:    errors.CodeInternalError,
			Message: err.Error(),
		},
	})
}

// domainCodeToHTTPStatus maps domain error codes to HTTP status codes
// This mapping is HTTP-specific, not in domain layer
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
