package errors

import (
	"errors"
)

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	var ve *ValidationError
	var mve *MultiValidationError
	return errors.As(err, &ve) || errors.As(err, &mve)
}

// IsAuthenticationError checks if error is an authentication error
func IsAuthenticationError(err error) bool {
	var ae *AuthenticationError
	return errors.As(err, &ae)
}

// IsAuthorizationError checks if error is an authorization error
func IsAuthorizationError(err error) bool {
	var ae *AuthorizationError
	return errors.As(err, &ae)
}

// IsNetworkError checks if error is a network error
func IsNetworkError(err error) bool {
	var ne *NetworkError
	return errors.As(err, &ne)
}

// IsRateLimitError checks if error is a rate limit error
func IsRateLimitError(err error) bool {
	var re *RateLimitError
	return errors.As(err, &re)
}

// IsServiceUnavailableError checks if error is a service unavailable error
func IsServiceUnavailableError(err error) bool {
	var se *ServiceUnavailableError
	return errors.As(err, &se)
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	var nfe *NotFoundError
	return errors.As(err, &nfe)
}

// IsDatabaseError checks if error is a database error
func IsDatabaseError(err error) bool {
	var de *DatabaseError
	return errors.As(err, &de)
}

// IsRetryable determines if an error should trigger a retry
func IsRetryable(err error) bool {
	// Network errors are retryable
	if IsNetworkError(err) {
		return true
	}

	// Service unavailable is retryable
	if IsServiceUnavailableError(err) {
		return true
	}

	// Rate limit errors are retryable (with backoff)
	if IsRateLimitError(err) {
		return true
	}

	// Check for specific SDK errors
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		switch sdkErr.Type {
		case ErrorTypeNetwork, ErrorTypeServiceUnavailable, ErrorTypeRateLimit:
			return true
		}
	}

	return false
}

// GetRetryAfter extracts retry-after duration from error
func GetRetryAfter(err error) int {
	var re *RateLimitError
	if errors.As(err, &re) {
		return re.RetryAfter
	}

	var se *ServiceUnavailableError
	if errors.As(err, &se) {
		return se.RetryAfter
	}

	return 0
}
