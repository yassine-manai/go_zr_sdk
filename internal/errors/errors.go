package errors

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrConnectionFailed = errors.New("connection failed")
	ErrNotFound         = errors.New("resource not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrTimeout          = errors.New("operation timed out")
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation         ErrorType = "validation"
	ErrorTypeAuthentication     ErrorType = "authentication"
	ErrorTypeAuthorization      ErrorType = "authorization"
	ErrorTypeNotFound           ErrorType = "not_found"
	ErrorTypeNetwork            ErrorType = "network"
	ErrorTypeRateLimit          ErrorType = "rate_limit"
	ErrorTypeServiceUnavailable ErrorType = "service_unavailable"
	ErrorTypeInternal           ErrorType = "internal"
	ErrorTypeDatabase           ErrorType = "database"
)

// SDKError is the base error type for all SDK errors
type SDKError struct {
	Type       ErrorType // Error category
	Message    string    // Human-readable message
	StatusCode int       // HTTP status code (if applicable)
	Err        error     // Original/wrapped error
}

// Error implements the error interface
func (e *SDKError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *SDKError) Unwrap() error {
	return e.Err
}

// Is checks if the error matches the target
func (e *SDKError) Is(target error) bool {
	t, ok := target.(*SDKError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// NewSDKError creates a new SDK error
func NewSDKError(errType ErrorType, message string, err error) *SDKError {
	return &SDKError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// WithStatusCode adds HTTP status code to the error
func (e *SDKError) WithStatusCode(code int) *SDKError {
	e.StatusCode = code
	return e
}
