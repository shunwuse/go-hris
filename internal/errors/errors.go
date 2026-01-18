package errors

import (
	"fmt"
)

// Error represents a domain-specific error.
type Error struct {
	code    string
	message string
	details map[string]string
}

// New creates a new Error instance.
func New(code string, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.message != "" {
		return fmt.Sprintf("[%s] %s", e.code, e.message)
	}
	return e.code
}

// Is supports errors.Is by comparing the internal error code.
// This allows comparing cloned errors (with details) against their sentinel counterparts.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.code == t.code
}

// ========================================
// Fluent API (Immutability)
// ========================================

// WithDetails returns a new Error instance with the provided details.
func (e *Error) WithDetails(details map[string]string) *Error {
	return &Error{
		code:    e.code,
		message: e.message,
		details: details,
	}
}

// ========================================
// Getter methods
// ========================================

func (e *Error) Code() string               { return e.code }
func (e *Error) Message() string            { return e.message }
func (e *Error) Details() map[string]string { return e.details }
