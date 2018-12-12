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

// ValidationFailure - Records validation failure key and error.
// e.g. ValidationFailure{
//        Key:   "DiscoveryModel.Name",
//        Error: "Field validation for 'Name' failed on the 'required' tag",
//      }
type ValidationFailure struct {
	Key   string `json:"key"`
	Error string `json:"error"`
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

const (
	fieldErrMsgFormat   = "Field validation for '%s' failed on the '%s' tag"
	versionErrMsgFormat = "DiscoveryVersion '%s' not in list of supported versions"
)

// Validate - validates a discovery model, returns true when valid,
// returns false and array of ValidationFailure structs when not valid.
func Validate(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
	failures := []ValidationFailure{}

	if err := validator.Struct(discovery); err != nil {
		failures = appendStructValidationErrors(err.(validation.ValidationErrors), failures)
		return false, failures, nil
	}
	failures = appendOtherValidationErrors(failures, checker, discovery, hasValidDiscoveryVersion)
	failures = appendOtherValidationErrors(failures, checker, discovery, hasValidAPISpecifications)
	failures = appendOtherValidationErrors(failures, checker, discovery, HasValidEndpoints)
	failures = appendOtherValidationErrors(failures, checker, discovery, HasMandatoryEndpoints)
	if len(failures) > 0 {
		return false, failures, nil
	}
	return true, failures, nil
}

func appendStructValidationErrors(errs validation.ValidationErrors, failures []ValidationFailure) []ValidationFailure {
	for _, msg := range errs {
		fieldError := validation.FieldError(msg)
		key := strings.Replace(fieldError.Namespace(), "Model.DiscoveryModel", "DiscoveryModel", 1)
		message := fmt.Sprintf(fieldErrMsgFormat, fieldError.Field(), fieldError.Tag())
		failure := ValidationFailure{
			Key:   key,
			Error: message,
		}
		failures = append(failures, failure)
	}
	return failures
}

func appendOtherValidationErrors(failures []ValidationFailure, checker model.ConditionalityChecker, discovery *Model,
	validationFn func(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []ValidationFailure)) []ValidationFailure {
	pass, newFailures := validationFn(checker, discovery)
	if !pass {
		for _, message := range newFailures {
			failures = append(failures, message)
		}
	}
	return failures;
}

// unmarshalDiscoveryJSON - used for testing to get discovery model from JSON.
// In production, we use echo.Context Bind to load configuration from JSON in HTTP POST.
func unmarshalDiscoveryJSON(discoveryJSON string) (*Model, error) {
	discovery := &Model{}
	err := json.Unmarshal([]byte(discoveryJSON), &discovery)
	return discovery, err
}

// checker passed to match function definition expectation in appendOtherValidationErrors function.
func hasValidDiscoveryVersion(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure) {
	failures := []ValidationFailure{}
	if !SupportedVersions()[discovery.DiscoveryModel.DiscoveryVersion] {
		failure := ValidationFailure{
			Key:   "DiscoveryModel.DiscoveryVersion",
			Error: fmt.Sprintf(versionErrMsgFormat, discovery.DiscoveryModel.DiscoveryVersion),
		}
		failures = append(failures, failure)
		return false, failures
	}
	return true, failures
}

// checker passed to match function definition expectation in appendOtherValidationErrors function.
func hasValidAPISpecifications(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []ValidationFailure) {
	failures := []ValidationFailure{}
	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {

		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationFromSchemaVersion(schemaVersion)
		if err != nil {
			failure := ValidationFailure{
				Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].APISpecification.SchemaVersion", discoveryItemIndex),
				Error: fmt.Sprintf("'SchemaVersion' not supported by suite '%s'", schemaVersion),
			}
			failures = append(failures, failure)
			continue
		}
		if specification.Name != discoveryItem.APISpecification.Name {
			failure := ValidationFailure{
				Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].APISpecification.Name", discoveryItemIndex),
				Error: fmt.Sprintf("'Name' should be '%s' when schemaVersion is '%s'", specification.Name, schemaVersion),
			}
			failures = append(failures, failure)
		}
		if specification.Version != discoveryItem.APISpecification.Version {
			failure := ValidationFailure{
				Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].APISpecification.Version", discoveryItemIndex),
				Error: fmt.Sprintf("'Version' should be '%s' when schemaVersion is '%s'", specification.Version, schemaVersion),
			}
			failures = append(failures, failure)
		}
		if specification.URL != discoveryItem.APISpecification.URL {
			failure := ValidationFailure{
				Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].APISpecification.URL", discoveryItemIndex),
				Error: fmt.Sprintf("'URL' should be '%s' when schemaVersion is '%s'", specification.URL, schemaVersion),
			}
			failures = append(failures, failure)
		}

	}
	if len(failures) > 0 {
		return false, failures
	}
	return true, failures
}

// HasValidEndpoints - checks that all the endpoints defined in the discovery
// model are either mandatory, conditional or optional.
// Return false and ValidationFailure structs indicating which endpoints are not valid.
func HasValidEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []ValidationFailure) {
	failures := []ValidationFailure{}

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationIdentifierFromSchemaVersion(schemaVersion)
		if err != nil {
			continue // err already added to failures in hasValidAPISpecifications
		}

		for endpointIndex, endpoint := range discoveryItem.Endpoints {
			isPresent, err := checker.IsPresent(endpoint.Method, endpoint.Path, specification)
			if err != nil {
				failure := ValidationFailure{
					Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints[%d]", discoveryItemIndex, endpointIndex),
					Error: err.Error(),
				}
				failures = append(failures, failure)
				continue
			}
			if !isPresent {
				failure := ValidationFailure{
					Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints[%d]", discoveryItemIndex, endpointIndex),
					Error: fmt.Sprintf("Invalid endpoint Method='%s', Path='%s'", endpoint.Method, endpoint.Path),
				}
				failures = append(failures, failure)
			}
		}
	}

	if len(failures) > 0 {
		return false, failures
	}

	return true, failures
}

// HasMandatoryEndpoints - checks that all the mandatory endpoints have been defined in each
// discovery model, otherwise it returns ValidationFailure structs for each missing mandatory endpoint.
func HasMandatoryEndpoints(checker model.ConditionalityChecker, discoveryConfig *Model) (bool, []ValidationFailure) {
	failures := []ValidationFailure{}

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
			failure := ValidationFailure{
				Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints", discoveryItemIndex),
				Error: err.Error(),
			}
			failures = append(failures, failure)
			continue
		}
		for _, mandatoryEndpoint := range missingMandatory {
			failure := ValidationFailure{
				Key: fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints", discoveryItemIndex),
				Error: fmt.Sprintf("Missing mandatory endpoint Method='%s', Path='%s'", mandatoryEndpoint.Method,
					mandatoryEndpoint.Endpoint),
			}
			failures = append(failures, failure)
		}
	}

	if len(failures) > 0 {
		return false, failures
	}

	return true, failures
}
