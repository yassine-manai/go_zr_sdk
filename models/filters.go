package models

import (
	"time"
)

// SortOrder represents sort direction
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SortBy represents sorting parameters
type SortBy struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

// NewSortBy creates a sort parameter
func NewSortBy(field string, order SortOrder) SortBy {
	if order == "" {
		order = SortOrderAsc
	}
	return SortBy{
		Field: field,
		Order: order,
	}
}

// DateRange represents a date range filter
type DateRange struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

// NewDateRange creates a date range filter
func NewDateRange(from, to *time.Time) DateRange {
	return DateRange{
		From: from,
		To:   to,
	}
}

// IsValid checks if the date range is valid
func (d DateRange) IsValid() bool {
	if d.From != nil && d.To != nil {
		return d.To.After(*d.From) || d.To.Equal(*d.From)
	}
	return true
}
