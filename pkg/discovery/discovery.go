package discovery

import (
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
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
	Method                string                `json:"method" validate:"required"`
	Path                  string                `json:"path" validate:"required,uri"`
	ConditionalProperties []ConditionalProperty `json:"conditionalProperties,omitempty" validate:"dive"`
}

// ConditionalProperty - Conditional schema property that has been implemented
type ConditionalProperty struct {
	Schema             string `json:"schema,omitempty" validate:"required"`
	Name               string `json:"name,omitempty" validate:"-"`     // transitional - will be required in a future version
	PropertyDeprecated string `json:"property,omitempty" validate:"-"` // property to be deprecated in favour of 'name'
	Path               string `json:"path,omitempty" validate:"required"`
	Required           bool   `json:"required,omitempty" validate:"-"`
	Request            bool   `json:"request,omitempty" validate:"-"` // indicates a request property that can be entered by the use
	Value              string `json:"value,omitempty" validate:"-"`
}

type ConditionalAPIProperties struct {
	Name      string          `json:"name,omitempty"`
	Endpoints []ModelEndpoint `json:"endpoints,omitempty"`
}

type CBPIIDebtorAccount struct {
	Identification string `json:"identification"`
	SchemeName     string `json:"scheme_name"`
	Name           string `json:"name,omitempty"`
}

func (c CBPIIDebtorAccount) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.SchemeName, validation.Required, validation.Length(1, 40)),
		validation.Field(&c.Identification, validation.Required, validation.Length(1, 256)),
		validation.Field(&c.Name, validation.Length(1, 70)),
	)
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

// validator, err := schema.NewSwaggerOBSpecValidator(item.APISpecification.Name, item.APISpecification.Version)
func GetConditionalProperties(disco *Model) ([]ConditionalAPIProperties, bool, error) {
	var haveProperties bool
	conditionalProps := make([]ConditionalAPIProperties, 0, len(disco.DiscoveryModel.DiscoveryItems))
	for k, discoitem := range disco.DiscoveryModel.DiscoveryItems {
		validator, err := schema.NewSwaggerOBSpecValidator(discoitem.APISpecification.Name, discoitem.APISpecification.Version)
		if err != nil {
			logrus.Error(err)
			return nil, false, err
		}
		conditionalProp := ConditionalAPIProperties{Name: discoitem.APISpecification.Name}
		for _, endpoint := range discoitem.Endpoints {
			if len(endpoint.ConditionalProperties) > 0 {

				for _, prop := range endpoint.ConditionalProperties {
					isRequest, _, err := validator.IsRequestProperty(endpoint.Method, endpoint.Path, prop.Path)
					if err != nil {
						logrus.Error(err)
						return nil, false, err
					}
					if isRequest {
						endpoint.ConditionalProperties[k].Request = true
					}

				}
				conditionalProp.Endpoints = append(conditionalProp.Endpoints, endpoint)
				haveProperties = true
			}
		}
		if haveProperties {
			conditionalProps = append(conditionalProps, conditionalProp)
		}
	}

	return conditionalProps, haveProperties, nil
}

func GetDiscoveryItemConditionalProperties(item ModelDiscoveryItem) []ModelEndpoint {
	endpoints := make([]ModelEndpoint, 0)
	for _, endpoint := range item.Endpoints {
		if len(endpoint.ConditionalProperties) > 0 {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}
