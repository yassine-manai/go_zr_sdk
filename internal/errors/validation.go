package errors

import "fmt"

// ValidationError represents validation failures
type ValidationError struct {
	Field   string // Field that failed validation
	Message string // Validation error message
	Value   any    // The invalid value
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("validation failed for field '%s': %s (value: %v)",
			e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// MultiValidationError holds multiple validation errors
type MultiValidationError struct {
	Errors []*ValidationError
}

// Error implements the error interface
func (e *MultiValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("validation failed with %d errors: %s",
		len(e.Errors), e.Errors[0].Message)
}

// Add adds a validation error to the collection
func (e *MultiValidationError) Add(field, message string, value interface{}) {
	e.Errors = append(e.Errors, NewValidationError(field, message, value))
}

// HasErrors returns true if there are validation errors
func (e *MultiValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Return returns the error if there are validation errors, nil otherwise
func (e *MultiValidationError) Return() error {
	if e.HasErrors() {
		return e
	}
	return nil
}
