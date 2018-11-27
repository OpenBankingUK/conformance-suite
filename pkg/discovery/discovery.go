package discovery

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	validation "gopkg.in/go-playground/validator.v9"
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
	validator = validation.New()
)

// Version returns the current version of the Discovery Model parser
func Version() string {
	version := "v0.0.1"
	return version
}

const fieldErrMsg = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"

// Validate - validates a discovery model, returns true when valid,
// returns false and validation failure messages when not valid.
func Validate(checker model.ConditionalityChecker, discovery *Model) (bool, []string, error) {
	failures := make([]string, 0)

	if err := validator.Struct(discovery); err != nil {
		errs := err.(validation.ValidationErrors)
		for _, msg := range errs {
			failure := validation.FieldError(msg)
			message := fmt.Sprintf(fieldErrMsg, failure.Namespace(), failure.Field(), failure.Tag())
			failures = append(failures, message)
		}
		return false, failures, nil
	}

	pass, messages, _ := HasValidEndpoints(checker, discovery)
	if !pass {
		for _, message := range messages {
			failures = append(failures, message)
		}
	}

	if len(failures) > 0 {
		return false, failures, nil
	}
	return true, failures, nil
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
	if err := validator.Struct(discoveryConfig); err != nil {
		// // translate all error at once
		// errs := err.(validator.ValidationErrors)
		// errsMap := errs.Translate(nil)
		return nil, err
	}
	if _, _, err := HasValidEndpoints(checker, discoveryConfig); err != nil {
		return nil, err
	}

	if _, _, err := HasMandatoryEndpoints(checker, discoveryConfig); err != nil {
		return nil, err
	}

	return discoveryConfig, nil
}

// HasValidEndpoints - checks that all the endpoints defined in the discovery
// model are either mandatory, conditional or optional.
// Return false and errors indicating which endpoints are not valid.
func HasValidEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []string, error) {
	errs := []string{}

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationIdentifierFromSchemaVersion(schemaVersion)
		if err != nil {
			warning := fmt.Sprintf("discoveryItemIndex=%d, "+err.Error(), discoveryItemIndex)
			errs = append(errs, warning)
			continue
		}

		for _, endpoint := range discoveryItem.Endpoints {
			isPresent, err := checker.IsPresent(endpoint.Method, endpoint.Path, specification)
			if err != nil {
				warning := fmt.Sprintf("discoveryItemIndex=%d, "+err.Error(), discoveryItemIndex)
				errs = append(errs, warning)
				continue
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
		return false, errs, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return true, errs, nil
}

// HasMandatoryEndpoints - checks that all the mandatory endpoints have been defined in each
// discovery model, otherwise it returns a error with all the missing mandatory endpoints separated
// by a newline.
func HasMandatoryEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []string, error) {
	errs := []string{}

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationIdentifierFromSchemaVersion(schemaVersion)
		if err != nil {
			warning := fmt.Sprintf("discoveryItemIndex=%d, "+err.Error(), discoveryItemIndex)
			errs = append(errs, warning)
			continue
		}

		discoveryEndpoints := []model.Input{}
		for _, endpoint := range discoveryItem.Endpoints {
			discoveryEndpoints = append(discoveryEndpoints, model.Input{Endpoint: endpoint.Path, Method: endpoint.Method})
		}
		missingMandatory, err := checker.MissingMandatory(discoveryEndpoints, specification)
		if err != nil {
			warning := fmt.Sprintf("discoveryItemIndex=%d, "+err.Error(), discoveryItemIndex)
			errs = append(errs, warning)
			continue
		}
		for _, mandatoryEndpoint := range missingMandatory {
			err := fmt.Sprintf(
				"discoveryItemIndex=%d, missing mandatory endpoint Method=%s, Path=%s",
				discoveryItemIndex,
				mandatoryEndpoint.Method,
				mandatoryEndpoint.Endpoint,
			)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return false, errs, fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return true, errs, nil
}
