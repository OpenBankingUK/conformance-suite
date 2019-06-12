// Enum to json.UnmarshalJson and json.MarshalJSON taken from: https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87
// We might be able to use https://github.com/campoy/jsonenums - might do this after the package is tested.
package report

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Status - the status of the `Report`.
type Status int

const (
	// StatusPending - The `Report` is pending.
	StatusPending Status = iota + 1
	// StatusComplete - The `Report` is complete.
	StatusComplete
	// StatusError - The `Report` is in error.
	StatusError
)

func reportStatusPendingToString() map[Status]string {
	return map[Status]string{
		StatusPending:  "Pending",
		StatusComplete: "Complete",
		StatusError:    "Error",
	}
}

func reportStatusPendingToID() map[string]Status {
	return map[string]Status{
		"Pending":  StatusPending,
		"Complete": StatusComplete,
		"Error":    StatusError,
	}
}

func (r Status) String() string {
	return reportStatusPendingToString()[r]
}

// MarshalJSON - marshals the enum as a quoted json string
func (r Status) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	value, ok := reportStatusPendingToString()[r]
	if !ok {
		return nil, fmt.Errorf("%d is an invalid enum for Status", r)
	}
	buffer.WriteString(value)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON - unmashals a quoted json string to the enum value
func (r *Status) UnmarshalJSON(data []byte) error {
	var status string
	err := json.Unmarshal(data, &status)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	value, ok := reportStatusPendingToID()[status]
	if !ok {
		return fmt.Errorf("%q is an invalid enum for Status", status)
	}
	*r = value
	return nil
}
