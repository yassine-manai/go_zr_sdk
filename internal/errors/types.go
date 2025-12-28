package errors

// AuthenticationError represents authentication failures
type AuthenticationError struct {
	Message string
	Err     error
}

func (e *AuthenticationError) Error() string {
	if e.Err != nil {
		return "authentication failed: " + e.Message + ": " + e.Err.Error()
	}
	return "authentication failed: " + e.Message
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

func NewAuthenticationError(message string, err error) *AuthenticationError {
	return &AuthenticationError{Message: message, Err: err}
}

// AuthorizationError represents authorization failures
type AuthorizationError struct {
	Message  string
	Resource string // Resource being accessed
}

func (e *AuthorizationError) Error() string {
	if e.Resource != "" {
		return "authorization failed for resource '" + e.Resource + "': " + e.Message
	}
	return "authorization failed: " + e.Message
}

func NewAuthorizationError(message, resource string) *AuthorizationError {
	return &AuthorizationError{Message: message, Resource: resource}
}

// NetworkError represents network-related failures
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return "network error: " + e.Message + ": " + e.Err.Error()
	}
	return "network error: " + e.Message
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

func NewNetworkError(message string, err error) *NetworkError {
	return &NetworkError{Message: message, Err: err}
}

// RateLimitError represents rate limiting errors
type RateLimitError struct {
	Message        string
	RetryAfter     int // Seconds to wait before retrying
	Limit          int // Rate limit
	RemainingCalls int // Remaining calls in current window
}

func (e *RateLimitError) Error() string {
	return "rate limit exceeded: " + e.Message
}

func NewRateLimitError(message string, retryAfter, limit, remaining int) *RateLimitError {
	return &RateLimitError{
		Message:        message,
		RetryAfter:     retryAfter,
		Limit:          limit,
		RemainingCalls: remaining,
	}
}

// ServiceUnavailableError represents service unavailability
type ServiceUnavailableError struct {
	Message    string
	Service    string // Which service is unavailable
	RetryAfter int    // Seconds to wait (if known)
}

func (e *ServiceUnavailableError) Error() string {
	if e.Service != "" {
		return "service unavailable (" + e.Service + "): " + e.Message
	}
	return "service unavailable: " + e.Message
}

func NewServiceUnavailableError(message, service string, retryAfter int) *ServiceUnavailableError {
	return &ServiceUnavailableError{
		Message:    message,
		Service:    service,
		RetryAfter: retryAfter,
	}
}

// NotFoundError represents resource not found errors
type NotFoundError struct {
	Message  string
	Resource string // Resource type
	ID       string // Resource identifier
}

func (e *NotFoundError) Error() string {
	if e.Resource != "" && e.ID != "" {
		return e.Resource + " not found: " + e.ID
	}
	return "not found: " + e.Message
}

func NewNotFoundError(message, resource, id string) *NotFoundError {
	return &NotFoundError{
		Message:  message,
		Resource: resource,
		ID:       id,
	}
}

// DatabaseError represents database-related errors
type DatabaseError struct {
	Message   string
	Query     string // SQL query (sanitized)
	Operation string // Operation type (SELECT, INSERT, etc.)
	Err       error
}

func (e *DatabaseError) Error() string {
	if e.Operation != "" {
		return "database error (" + e.Operation + "): " + e.Message
	}
	return "database error: " + e.Message
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

func NewDatabaseError(message, query, operation string, err error) *DatabaseError {
	return &DatabaseError{
		Message:   message,
		Query:     query,
		Operation: operation,
		Err:       err,
	}
}
