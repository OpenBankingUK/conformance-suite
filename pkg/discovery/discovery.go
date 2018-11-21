package discovery

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	validator "gopkg.in/go-playground/validator.v9"
)

// Model ... TODO: Document.
type Model struct {
	DiscoveryModel ModelDiscovery `json:"discoveryModel" validate:"required,dive"`
}

// ModelDiscovery ... TODO: Document.
type ModelDiscovery struct {
	Version        string               `json:"version" validate:"required"`
	DiscoveryItems []ModelDiscoveryItem `json:"discoveryItems" validate:"required,gt=0,dive"`
}

// ModelDiscoveryItem ... TODO: Document.
type ModelDiscoveryItem struct {
	APISpecification       ModelAPISpecification `json:"apiSpecification" validate:"required"`
	OpenidConfigurationURI string                `json:"openidConfigurationUri" validate:"required,url"`
	ResourceBaseURI        string                `json:"resourceBaseUri" validate:"required,url"`
	Endpoints              []ModelEndpoint       `json:"endpoints" validate:"required,gt=0,dive"`
}

// ModelAPISpecification ... TODO: Document.
type ModelAPISpecification struct {
	Name          string `json:"name" validate:"required"`
	URL           string `json:"url" validate:"required,url"`
	Version       string `json:"version" validate:"required"`
	SchemaVersion string `json:"schemaVersion" validate:"required,url"`
}

// ModelEndpoint ... TODO: Document.
type ModelEndpoint struct {
	Method                string                       `json:"method" validate:"required"`
	Path                  string                       `json:"path" validate:"required,uri"`
	ConditionalProperties []ModelConditionalProperties `json:"conditionalProperties,omitempty" validate:"dive"`
}

// ModelConditionalProperties ... TODO: Document.
type ModelConditionalProperties struct {
	Schema   string `json:"schema" validate:"required"`
	Property string `json:"property" validate:"required"`
	Path     string `json:"path" validate:"required"`
}

var (
	// use a single instance of Validate, it caches struct info
	validate = validator.New()
)

// Version returns the current version of the Discovery Model parser
func Version() string {
	version := "v0.0.1"
	return version
}

// FromJSONString - used for testing.
// In production, we use echo.Context Bind to load configuration from JSON in HTTP POST.
func FromJSONString(checker model.ConditionalityChecker, configStr string) (*Model, error) {
	discoveryConfig := &Model{}

	err := json.Unmarshal([]byte(configStr), &discoveryConfig)
	if err != nil {
		return nil, err
	}

	// returns nil or ValidationErrors ( []FieldError )
	if err := validate.Struct(discoveryConfig); err != nil {
		// // translate all error at once
		// errs := err.(validator.ValidationErrors)
		// errsMap := errs.Translate(nil)
		return nil, err
	}
	if _, err := HasValidEndpoints(checker, discoveryConfig); err != nil {
		return nil, err
	}

	if _, err := HasMandatoryEndpoints(checker, discoveryConfig); err != nil {
		return nil, err
	}

	return discoveryConfig, nil
}

// HasValidEndpoints - checks that all the endpoints defined in the discovery
// model are either mandatory, conditional or optional.
// Return false and errors indicating which endpoints are not valid.
func HasValidEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, error) {
	errs := []string{}

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		// ignore if it isn't accounts as we don't have the definitions for payments just yet
		if discoveryItem.APISpecification.Name != "Account and Transaction API Specification" {
			continue
		}
		specification := "account-transaction-v3.0"

		for _, endpoint := range discoveryItem.Endpoints {
			isPresent, err := checker.IsPresent(endpoint.Method, endpoint.Path, specification)
			if (err != nil) {
				return false, err
			}
			if !isPresent {
				err := fmt.Sprintf(
					"discoveryItemIndex=%d, invalid endpoint Method=%s, Path=%s",
					discoveryItemIndex,
					endpoint.Method,
					endpoint.Path,
				)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return false, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return true, nil
}

// HasMandatoryEndpoints - checks that all the mandatory endpoints have been defined in each
// discovery model, otherwise it returns a error with all the missing mandatory endpoints separated
// by a newline.
func HasMandatoryEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, error) {
	errs := []string{}

	// filter out non-mandatory endpoints, i.e., just store the mandatory endpoints.
	// this is just for accounts at the moment, conditionality for payments has not be defined just yet.
	mandatoryEndpoints := []model.Conditionality{}
	specification := "account-transaction-v3.0"
	endpoints := model.GetEndpointConditionality(specification)
	for _, endpoint := range endpoints {
		isMandatory, err := checker.IsMandatory(endpoint.Method, endpoint.Endpoint, specification)
		if err != nil {
			continue
		}

		if isMandatory {
			mandatoryEndpoints = append(mandatoryEndpoints, endpoint)
		}
	}

	// check that each mandatory endpoint is included in the discovery model
	for _, mandatoryEndpoint := range mandatoryEndpoints {
		for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
			// ignore if it isn't accounts as we don't have the definitions for payments just yet
			if discoveryItem.APISpecification.Name != "Account and Transaction API Specification" {
				continue
			}

			isPresent := false
			for _, endpoint := range discoveryItem.Endpoints {
				isPresent = endpoint.Method == mandatoryEndpoint.Method && endpoint.Path == mandatoryEndpoint.Endpoint
				if isPresent {
					break
				}
			}

			if !isPresent {
				err := fmt.Sprintf(
					"discoveryItemIndex=%d, missing mandatory endpoint Method=%s, Path=%s",
					discoveryItemIndex,
					mandatoryEndpoint.Method,
					mandatoryEndpoint.Endpoint,
				)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return false, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return true, nil
}
