package models

// CustomerMedia represents a customer media record
type CustomerMedia struct {
	ID          int64     `json:"id"`
	CustomerID  int64     `json:"customer_id"`
	MediaType   string    `json:"media_type"` // image, video, document
	MediaURL    string    `json:"media_url"`
	Filename    string    `json:"filename"`
	FileSize    int64     `json:"file_size"` // bytes
	MimeType    string    `json:"mime_type"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
	CreatedBy   string    `json:"created_by,omitempty"`
}

// CreateCustomerMediaRequest represents a request to create media
type CreateCustomerMediaRequest struct {
	CustomerID  int64    `json:"customer_id" validate:"required"`
	MediaType   string   `json:"media_type" validate:"required,oneof=image video document"`
	MediaURL    string   `json:"media_url" validate:"required,url"`
	Filename    string   `json:"filename" validate:"required"`
	FileSize    int64    `json:"file_size" validate:"required,min=1"`
	MimeType    string   `json:"mime_type" validate:"required"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// UpdateCustomerMediaRequest represents a request to update media
type UpdateCustomerMediaRequest struct {
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// CustomerMediaFilters represents filters for listing media
type CustomerMediaFilters struct {
	CustomerID *int64     `json:"customer_id,omitempty"`
	MediaType  *string    `json:"media_type,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	DateRange  *DateRange `json:"date_range,omitempty"`
}

// CustomerMediaListRequest combines filters, pagination, and sorting
type CustomerMediaListRequest struct {
	Filters    CustomerMediaFilters `json:"filters"`
	Pagination PaginationRequest    `json:"pagination"`
	Sort       *SortBy              `json:"sort,omitempty"`
}

// NewCustomerMediaListRequest creates a list request with defaults
func NewCustomerMediaListRequest() CustomerMediaListRequest {
	return CustomerMediaListRequest{
		Pagination: NewPaginationRequest(1, 20),
		Sort:       &SortBy{Field: "created_at", Order: SortOrderDesc},
	}
}
