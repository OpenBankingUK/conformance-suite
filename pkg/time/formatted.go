package time

import (
	"encoding/json"
	"time"
)

const (
	// Layout we use to format the time
	Layout = time.RFC3339
)

// Formatted represents a custom time.Time  for consistency across our application
type Formatted time.Time

// NewUTCTime returns a Formatted object with location set to UTC
func NewUTCTime() Formatted {
	return Formatted(time.Now().UTC())
}

// MarshalJSON implements the json.Marshaler interface.
func (d Formatted) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(Layout))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Formatted) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"`+Layout+`"`, string(data))
	*d = Formatted(parsedTime)
	return err
}

// String returns a string representing the time in RFC3339 format
func (d Formatted) String() string {
	return time.Time(d).Format(Layout)
}
