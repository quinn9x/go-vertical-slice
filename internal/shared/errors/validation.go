package apperrors

// ValidationError represents an error that occurs when request validation fails.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "request validation failed"
}

// NewValidationError creates a new ValidationError with the provided fields.
func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{
		Fields: fields,
	}
}
