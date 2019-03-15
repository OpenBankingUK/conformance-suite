// Enum to json.UnmarshalJson and json.MarshalJSON taken from: https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87
// We might be able to use https://github.com/campoy/jsonenums - might do this after the package is tested.
package report

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ReportStatus int

const (
	ReportStatusPending ReportStatus = iota + 1
	ReportStatusComplete
	ReportStatusError
)

var ReportStatusPendingToString = map[ReportStatus]string{
	ReportStatusPending:  "Pending",
	ReportStatusComplete: "Complete",
	ReportStatusError:    "Error",
}

var ReportStatusPendingToID = map[string]ReportStatus{
	"Pending":  ReportStatusPending,
	"Complete": ReportStatusComplete,
	"Error":    ReportStatusError,
}

func (r ReportStatus) String() string {
	return ReportStatusPendingToString[r]
}

// MarshalJSON marshals the enum as a quoted json string
func (r ReportStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	value, ok := ReportStatusPendingToString[r]
	if !ok {
		return nil, fmt.Errorf("%d is an invalid enum for ReportStatus", r)
	}
	buffer.WriteString(value)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *ReportStatus) UnmarshalJSON(data []byte) error {
	var status string
	err := json.Unmarshal(data, &status)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	value, ok := ReportStatusPendingToID[status]
	if !ok {
		return fmt.Errorf("%q is an invalid enum for ReportStatus", status)
	}
	*r = value
	return nil
}
