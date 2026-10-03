// Package apperrors provides structured error handling for the application.
package apperrors

// Error represents a structured error with a code, message, and optional fields.
type Error struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// NewNotFound creates a new Error instance representing a "not found" error with the provided message.
func NewNotFound(message string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: message,
	}
}

// NewBadRequest creates a new Error instance representing a "bad request" error with the provided message.
func NewBadRequest(message string) *Error {
	return &Error{
		Code:    CodeBadRequest,
		Message: message,
	}
}

// NewConflict creates a new Error instance representing a "conflict" error with the provided message.
func NewConflict(message string) error {
	return &Error{
		Code:    CodeConflict,
		Message: message,
	}
}
