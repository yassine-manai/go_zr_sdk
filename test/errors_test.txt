package errors

import (
	"errors"
	"testing"
)

func TestSDKError(t *testing.T) {
	err := NewSDKError(ErrorTypeValidation, "invalid input", nil)

	if err.Type != ErrorTypeValidation {
		t.Errorf("expected type %s, got %s", ErrorTypeValidation, err.Type)
	}

	if err.Message != "invalid input" {
		t.Errorf("expected message 'invalid input', got %s", err.Message)
	}
}

func TestSDKErrorWithStatusCode(t *testing.T) {
	err := NewSDKError(ErrorTypeNotFound, "resource not found", nil).WithStatusCode(404)

	if err.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", err.StatusCode)
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("email", "invalid format", "bad-email")

	if err.Field != "email" {
		t.Errorf("expected field 'email', got %s", err.Field)
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}

func TestMultiValidationError(t *testing.T) {
	multiErr := &MultiValidationError{}
	multiErr.Add("email", "required", nil)
	multiErr.Add("age", "must be positive", -5)

	if !multiErr.HasErrors() {
		t.Error("should have errors")
	}

	if len(multiErr.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(multiErr.Errors))
	}

	if multiErr.Return() == nil {
		t.Error("should return error when there are validation errors")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "network error is retryable",
			err:      NewNetworkError("connection refused", nil),
			expected: true,
		},
		{
			name:     "service unavailable is retryable",
			err:      NewServiceUnavailableError("maintenance", "ui-service", 60),
			expected: true,
		},
		{
			name:     "rate limit is retryable",
			err:      NewRateLimitError("too many requests", 30, 100, 0),
			expected: true,
		},
		{
			name:     "validation error is not retryable",
			err:      NewValidationError("email", "invalid", "bad"),
			expected: false,
		},
		{
			name:     "not found is not retryable",
			err:      NewNotFoundError("user not found", "user", "123"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.expected {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorUnwrapping(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := NewNetworkError("network failed", baseErr)

	if !errors.Is(wrappedErr, baseErr) {
		t.Error("should be able to unwrap to base error")
	}
}

func TestGetRetryAfter(t *testing.T) {
	err := NewRateLimitError("rate limit", 60, 100, 0)
	retryAfter := GetRetryAfter(err)

	if retryAfter != 60 {
		t.Errorf("expected retry after 60, got %d", retryAfter)
	}
}
