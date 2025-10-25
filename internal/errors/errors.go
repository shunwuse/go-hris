package errors

import "fmt"

type Error struct {
	Code    string
	Message string
	Details map[string]string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return e.Code
}

func New(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) WithDetails(details map[string]string) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}
