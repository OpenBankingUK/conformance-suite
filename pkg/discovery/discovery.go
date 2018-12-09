package discovery

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	validation "gopkg.in/go-playground/validator.v9"
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
	DiscoveryItems   []ModelDiscoveryItem `json:"discoveryItems" validate:"required,gt=0,dive"`
}

// ModelDiscoveryItem - Each discovery item contains information related to a particular specification version.
type ModelDiscoveryItem struct {
	APISpecification       ModelAPISpecification `json:"apiSpecification,omitempty" validate:"required"`
	OpenidConfigurationURI string                `json:"openidConfigurationUri,omitempty" validate:"required,url"`
	ResourceBaseURI        string                `json:"resourceBaseUri,omitempty" validate:"required,url"`
	ResourceIds            ResourceIds           `json:"resourceIds,omitempty"`
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
}

// ModelEndpoint - Endpoint and methods that have been implemented by implementer.
type ModelEndpoint struct {
	Method                string                       `json:"method" validate:"required"`
	Path                  string                       `json:"path" validate:"required,uri"`
	ConditionalProperties []ModelConditionalProperties `json:"conditionalProperties,omitempty" validate:"dive"`
}

// ModelConditionalProperties - Conditional schema properties implemented by implementer.
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
	version := "v0.1.0"
	return version
}

// SupportedVersions - returns map of supported versions
func SupportedVersions() map[string]bool {
	return map[string]bool{
		Version(): true,
	}
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
	if !SupportedVersions()[discovery.DiscoveryModel.DiscoveryVersion] {
		failures = append(failures, `Key: 'Model.DiscoveryModel.DiscoveryVersion' Error:DiscoveryVersion `+
			discovery.DiscoveryModel.DiscoveryVersion+` not in list of supported versions`)
	}
	pass, messages, _ := hasValidAPISpecifications(discovery)
	if !pass {
		for _, message := range messages {
			failures = append(failures, message)
		}
	}
	pass, messages, _ = HasValidEndpoints(checker, discovery)
	if !pass {
		for _, message := range messages {
			failures = append(failures, message)
		}
	}

	pass, messages, _ = HasMandatoryEndpoints(checker, discovery)
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

// unmarshalDiscoveryJSON - used for testing to get discovery model from JSON.
// In production, we use echo.Context Bind to load configuration from JSON in HTTP POST.
func unmarshalDiscoveryJSON(discoveryJSON string) (*Model, error) {
	discovery := &Model{}
	err := json.Unmarshal([]byte(discoveryJSON), &discovery)
	return discovery, err
}

func hasValidAPISpecifications(discoveryConfig *Model) (bool, []string, error) {
	errs := []string{}
	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {

		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationFromSchemaVersion(schemaVersion)
		if err != nil {
			failure := fmt.Sprintf("Key: 'Model.DiscoveryModel.DiscoveryItems[%d].APISpecification.SchemaVersion' Error:'SchemaVersion' not supported by suite '%s'",
				discoveryItemIndex, schemaVersion)
			errs = append(errs, failure)
			continue
		}
		if specification.Name != discoveryItem.APISpecification.Name {
			failure := fmt.Sprintf("Key: 'Model.DiscoveryModel.DiscoveryItems[%d].APISpecification.Name' Error:'Name' should be '%s' when schemaVersion is '%s'",
				discoveryItemIndex, specification.Name, schemaVersion)
			errs = append(errs, failure)
		}
		if specification.Version != discoveryItem.APISpecification.Version {
			failure := fmt.Sprintf("Key: 'Model.DiscoveryModel.DiscoveryItems[%d].APISpecification.Version' Error:'Version' should be '%s' when schemaVersion is '%s'",
				discoveryItemIndex, specification.Version, schemaVersion)
			errs = append(errs, failure)
		}
		if specification.URL != discoveryItem.APISpecification.URL {
			failure := fmt.Sprintf("Key: 'Model.DiscoveryModel.DiscoveryItems[%d].APISpecification.URL' Error:'URL' should be '%s' when schemaVersion is '%s'",
				discoveryItemIndex, specification.URL, schemaVersion)
			errs = append(errs, failure)
		}

	}
	if len(errs) > 0 {
		return false, errs, nil
	}
	return true, errs, nil
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
			continue // err already added to failures in hasValidAPISpecifications
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
			continue // err already added to failures in hasValidAPISpecifications
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
