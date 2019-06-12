// Enum to json.UnmarshalJson and json.MarshalJSON taken from: https://gist.github.com/lummie/7f5c237a17853c031a57277371528e87
// We might be able to use https://github.com/campoy/jsonenums - might do this after the package is tested.
package report

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CertifiedByEnvironment - the environment the `Report` was generated in.
type CertifiedByEnvironment int

const (
	// ReportCertifiedByEnvironmentTesting - testing environment.
	CertifiedByEnvironmentTesting CertifiedByEnvironment = iota + 1
	// ReportCertifiedByEnvironmentProduction - production environment.
	CertifiedByEnvironmentProduction
	CertifiedByEnvironmentSandbox
)

func certifiedByEnvironmentToString() map[CertifiedByEnvironment]string {
	return map[CertifiedByEnvironment]string{
		CertifiedByEnvironmentTesting:    "testing",
		CertifiedByEnvironmentProduction: "production",
		CertifiedByEnvironmentSandbox:    "sandbox",
	}
}

func certifiedByEnvironmentToID() map[string]CertifiedByEnvironment {
	return map[string]CertifiedByEnvironment{
		"testing":    CertifiedByEnvironmentTesting,
		"production": CertifiedByEnvironmentProduction,
		"sandbox":    CertifiedByEnvironmentSandbox,
	}
}

func (r CertifiedByEnvironment) String() string {
	return certifiedByEnvironmentToString()[r]
}

// MarshalJSON marshals the enum as a quoted json string
func (r CertifiedByEnvironment) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	value, ok := certifiedByEnvironmentToString()[r]
	if !ok {
		return nil, fmt.Errorf("%d is an invalid enum for CertifiedByEnvironment", r)
	}
	buffer.WriteString(value)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *CertifiedByEnvironment) UnmarshalJSON(data []byte) error {
	var environment string
	err := json.Unmarshal(data, &environment)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	value, ok := certifiedByEnvironmentToID()[environment]
	if !ok {
		return fmt.Errorf("%q is an invalid enum for CertifiedByEnvironment", environment)
	}
	*r = value
	return nil
}
