package discovery

import (
	"encoding/json"
)

// Model - Top level struct holding discovery model.
type Model struct {
	DiscoveryModel ModelDiscovery `json:"discoveryModel" validate:"required,dive"`
}

// ModelDiscovery - Holds fields describing model, and array of discovery items.
// For detailed documentation see ./doc/permissions.md file.
type ModelDiscovery struct {
	Name             string               `json:"name" validate:"required"`
	Description      string               `json:"description" validate:"required"`
	DiscoveryVersion string               `json:"discoveryVersion" validate:"required"`
	TokenAcquisition string               `json:"tokenAcquisition" validate:"required"`
	DiscoveryItems   []ModelDiscoveryItem `json:"discoveryItems" validate:"required,dive"`
	CustomTests      []CustomTest         `json:"customTests" validate:"-"`
}

// ModelDiscoveryItem - Each discovery item contains information related to a particular specification version.
type ModelDiscoveryItem struct {
	APISpecification       ModelAPISpecification `json:"apiSpecification,omitempty" validate:"required"`
	OpenidConfigurationURI string                `json:"openidConfigurationUri,omitempty" validate:"required,url"`
	ResourceBaseURI        string                `json:"resourceBaseUri,omitempty" validate:"required,url"`
	ResourceIds            ResourceIds           `json:"resourceIds,omitempty" validate:"-"`
	Endpoints              []ModelEndpoint       `json:"endpoints,omitempty" validate:"required,gt=0,dive"`
}

// ResourceIds section allows the replacement of endpoint resourceid values with real data parameters like accountid
type ResourceIds map[string]string

// ModelAPISpecification ... TODO: Document.
type ModelAPISpecification struct {
	Name          string `json:"name" validate:"required"`
	URL           string `json:"url" validate:"required,url"`
	Version       string `json:"version" validate:"required"`
	SchemaVersion string `json:"schemaVersion" validate:"required,url"`
	Manifest      string `json:"manifest" validate:"required,fileorhttps"`
	SpecType      string `json:"-"`
}

// ModelEndpoint - Endpoint and methods that have been implemented by implementer.
type ModelEndpoint struct {
	Method                string                  `json:"method" validate:"required"`
	Path                  string                  `json:"path" validate:"required,uri"`
	ConditionalProperties []ConditionalProperties `json:"conditionalProperties,omitempty" validate:"dive"`
}

// ModelConditionalProperties - Conditional schema properties implemented by implementer.
type ConditionalProperties struct {
	Schema   string `json:"schema" validate:"required"`
	Property string `json:"property" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Required string `json:"required,omitempty"`
}

// UnmarshalDiscoveryJSON - Used for testing in multiple packages to get discovery
// model from JSON. We tried moving this function to a _test file, but we get
// `go vet` error as it is used from multiple packages.
// In production, we use echo.Context Bind to load configuration from JSON in HTTP POST.
func UnmarshalDiscoveryJSON(discoveryJSON string) (*Model, error) {
	discovery := &Model{}
	err := json.Unmarshal([]byte(discoveryJSON), &discovery)
	return discovery, err
}
