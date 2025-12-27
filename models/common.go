package models

import (
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

// MarshalJSON implements json.Marshaler
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) < 2 {
		return nil
	}
	str = str[1 : len(str)-1]

	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
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
