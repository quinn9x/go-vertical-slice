// Package validation provides utility functions for validating data.
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	apperrors "github.com/quinn9x/go-vertical-slice/internal/shared/errors"
)

// Validator is a wrapper around the go-playground/validator library that provides validation functionality.
type Validator struct {
	validate *validator.Validate
}

// New creates a new instance of Validator.
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Struct validates the provided struct value and returns a ValidationError if validation fails.
func (v *Validator) Struct(value any) error {
	err := v.validate.Struct(value)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	fields := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := fieldName(value, fieldErr.Field())

		fields[field] = messageFor(fieldErr)
	}

	return apperrors.NewValidationError(fields)
}

func messageFor(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"

	case "min":
		return fmt.Sprintf("must be at least %s characters", err.Param())

	case "max":
		return fmt.Sprintf("must be at most %s characters", err.Param())

	case "gt":
		return fmt.Sprintf("must be greater than %s", err.Param())

	default:
		return "is invalid"
	}
}

func fieldName(value any, field string) string {
	valueType := reflect.TypeOf(value)

	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	fieldStruct, ok := valueType.FieldByName(field)
	if !ok {
		return strings.ToLower(field)
	}

	jsonTag := fieldStruct.Tag.Get("json")
	if jsonTag != "" {
		name := strings.Split(jsonTag, ",")[0]

		if name != "" && name != "-" {
			return name
		}
	}

	return strings.ToLower(field)
}
