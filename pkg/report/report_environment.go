// Enum to json.UnmarshalJson and json.MarshalJSON taken from: https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87
// We might be able to use https://github.com/campoy/jsonenums - might do this after the package is tested.
package report

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ReportCertifiedByEnvironment - the environment the `Report` was generated in.
type ReportCertifiedByEnvironment int

const (
	// ReportCertifiedByEnvironmentTesting - testing environment.
	ReportCertifiedByEnvironmentTesting ReportCertifiedByEnvironment = iota + 1
	// ReportCertifiedByEnvironmentProduction - production environment.
	ReportCertifiedByEnvironmentProduction
)

var reportCertifiedByEnvironmentToString = map[ReportCertifiedByEnvironment]string{
	ReportCertifiedByEnvironmentTesting:    "testing",
	ReportCertifiedByEnvironmentProduction: "production",
}

var reportCertifiedByEnvironmentToID = map[string]ReportCertifiedByEnvironment{
	"testing":    ReportCertifiedByEnvironmentTesting,
	"production": ReportCertifiedByEnvironmentProduction,
}

func (r ReportCertifiedByEnvironment) String() string {
	return reportCertifiedByEnvironmentToString[r]
}

// MarshalJSON marshals the enum as a quoted json string
func (r ReportCertifiedByEnvironment) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	value, ok := reportCertifiedByEnvironmentToString[r]
	if !ok {
		return nil, fmt.Errorf("%d is an invalid enum for ReportCertifiedByEnvironment", r)
	}
	buffer.WriteString(value)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *ReportCertifiedByEnvironment) UnmarshalJSON(data []byte) error {
	var environment string
	err := json.Unmarshal(data, &environment)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	value, ok := reportCertifiedByEnvironmentToID[environment]
	if !ok {
		return fmt.Errorf("%q is an invalid enum for ReportCertifiedByEnvironment", environment)
	}
	*r = value
	return nil
}
