package models

import (
	"encoding/xml"
	"time"
)

// Timestamp represents a standardized timestamp
type Timestamp struct {
	time.Time
}

// NewTimestamp creates a new timestamp from time.Time
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// Now returns current timestamp
func Now() Timestamp {
	return Timestamp{Time: time.Now()}
}

// Metadata contains common metadata for responses
type Metadata struct {
	RequestID string    `json:"request_id,omitempty"`
	Timestamp Timestamp `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// NewMetadata creates metadata with current timestamp
func NewMetadata(requestID string) Metadata {
	return Metadata{
		RequestID: requestID,
		Timestamp: Now(),
		Version:   "v1",
	}
}

// ErrorResponse represents XML error response from API
type ErrorResponse struct {
	XMLName xml.Name `xml:"http://gsph.sub.com/cust/types errorResponse"`
	Error   struct {
		ErrCode      string `xml:"errCode"`
		ShortMsg     string `xml:"shortMsg"`
		Message      string `xml:"message"`
		CauseMessage string `xml:"causeMessage"`
	} `xml:"error"`
}
