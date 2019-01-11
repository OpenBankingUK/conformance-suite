package reporting

import "github.com/google/uuid"

// Result of a full test run
type Result struct {
	Id             uuid.UUID       `json:"id"`
	Specifications []Specification `json:"specifications"`
}

// Specification of a full test run
type Specification struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	URL           string `json:"url"`
	SchemaVersion string `json:"schemaVersion"`
	Pass          bool   `json:"pass"`
	Tests         []Test `json:"tests"`
}

// Test result for a run
type Test struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Pass     bool   `json:"pass"`
}
