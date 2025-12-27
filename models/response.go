package models

// Response is a generic API response wrapper
type Response[T any] struct {
	Success  bool      `json:"success"`
	Data     T         `json:"data,omitempty"`
	Error    *APIError `json:"error,omitempty"`
	Metadata Metadata  `json:"metadata"`
}

// ListResponse is a generic paginated list response
type ListResponse[T any] struct {
	Success    bool               `json:"success"`
	Data       []T                `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
	Error      *APIError          `json:"error,omitempty"`
	Metadata   Metadata           `json:"metadata"`
}

// APIError represents an error from the API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewResponse creates a successful response
func NewResponse[T any](data T, requestID string) Response[T] {
	return Response[T]{
		Success:  true,
		Data:     data,
		Metadata: NewMetadata(requestID),
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse[T any](code, message, details string) Response[T] {
	return Response[T]{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Metadata: NewMetadata(""),
	}
}

// NewListResponse creates a paginated list response
func NewListResponse[T any](data []T, pagination PaginationResponse, requestID string) ListResponse[T] {
	return ListResponse[T]{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Metadata:   NewMetadata(requestID),
	}
}
