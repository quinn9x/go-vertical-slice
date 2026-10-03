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
