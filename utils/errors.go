package utils

import (
	"fmt"
)

type ValidationError struct {
	Fields map[string]string
}

func (v *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{Fields: fields}
}

func IsValidationError(err error) (*ValidationError, bool) {
	ve, ok := err.(*ValidationError)
	return ve, ok
}

func WrapValidationError(err error) error {
	if ve, ok := err.(*ValidationError); ok {
		return fmt.Errorf("validation: %w", ve)
	}
	return err
}
