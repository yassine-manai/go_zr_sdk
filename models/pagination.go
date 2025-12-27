package models

// PaginationRequest represents pagination parameters for list requests
type PaginationRequest struct {
	Page     int `json:"page"`      // Page number (1-based)
	PageSize int `json:"page_size"` // Number of items per page
	Offset   int `json:"-"`         // Calculated offset (internal use)
}

// NewPaginationRequest creates pagination request with defaults
func NewPaginationRequest(page, pageSize int) PaginationRequest {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// PaginationResponse represents pagination metadata in responses
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewPaginationResponse creates pagination response
func NewPaginationResponse(page, pageSize int, totalItems int64) PaginationResponse {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
