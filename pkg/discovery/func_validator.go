package discovery

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	validation "gopkg.in/go-playground/validator.v9"
)

// Version returns the current version of the Discovery Model parser
func Version() string {
	version := "v0.4.0"
	return version
}

func MinVersion() string {
	minversion := "v0.3.0"
	return minversion
}

// SupportedVersions - returns map of supported versions
func SupportedVersions() map[string]bool {
	return map[string]bool{
		Version():    true,
		MinVersion(): true,
	}
}

// SupportedTokenAcquisitions returns a collection of supported token acquisition methods
func SupportedTokenAcquisitions() []string {
	return []string{"psu", "headless", "store"}
}

const (
	fieldErrMsgFormat            = "Field validation for '%s' failed on the '%s' tag"
	versionErrMsgFormat          = "DiscoveryVersion '%s' not in list of supported versions"
	tokenAcquisitionErrMsgFormat = "TokenAcquisition '%s' not in list of supported methods"
	requiredErrorFormat          = "Field '%s' is required"
	emptyArrayErrorFormat        = "Field '%s' cannot be empty"
	fileOrHttpsErrorFormat       = "Field '%s' must be 'file://' or 'https://'"
)

// Validate - validates a discovery model, returns true when valid,
// returns false and array of ValidationFailure structs when not valid.
func Validate(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
	failures := make([]ValidationFailure, 0)

	v := validation.New()
	httpsValidate := func(f validation.FieldLevel) bool {
		value := f.Field().String()

		if value == "" {
			return true
		}

		u, err := url.Parse(value)
		return err == nil && (u.Scheme == "file" || u.Scheme == "https")
	}
	if err := v.RegisterValidation("fileorhttps", httpsValidate); err != nil {
		return false, nil, errors.Wrap(err, "register `fileorhttps` validation")
	}

	if err := v.Struct(discovery); err != nil {
		failures = appendStructValidationErrors(err.(validation.ValidationErrors), failures)
		return false, failures, nil
	}
	failures = appendOtherValidationErrors(failures, checker, discovery, hasValidDiscoveryVersion)
	failures = appendOtherValidationErrors(failures, checker, discovery, hasValidTokenAcquisitionMethod)
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
		fieldError := msg
		key := strings.Replace(fieldError.Namespace(), "Model.DiscoveryModel", "DiscoveryModel", 1)
		var message string
		switch fieldError.Tag() {
		default:
			message = fmt.Sprintf(fieldErrMsgFormat, fieldError.Field(), fieldError.Tag())
		case "required":
			message = fmt.Sprintf(requiredErrorFormat, key)
		case "gt":
			message = fmt.Sprintf(emptyArrayErrorFormat, key)
		case "fileorhttps":
			message = fmt.Sprintf(fileOrHttpsErrorFormat, key)
		}
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
		failures = append(failures, newFailures...)
	}
	return failures
}

// checker passed to match function definition expectation in appendOtherValidationErrors function.
func hasValidDiscoveryVersion(_ model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure) {
	var failures []ValidationFailure
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
func hasValidTokenAcquisitionMethod(_ model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure) {
	var failures []ValidationFailure

	for _, method := range SupportedTokenAcquisitions() {
		if method == discovery.DiscoveryModel.TokenAcquisition {
			return true, failures
		}
	}

	failure := ValidationFailure{
		Key:   "DiscoveryModel.TokenAcquisition",
		Error: fmt.Sprintf(tokenAcquisitionErrMsgFormat, discovery.DiscoveryModel.TokenAcquisition),
	}
	failures = append(failures, failure)
	return false, failures
}

// checker passed to match function definition expectation in appendOtherValidationErrors function.
func hasValidAPISpecifications(_ model.ConditionalityChecker, discoveryConfig *Model) (bool, []ValidationFailure) {
	var failures []ValidationFailure
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
		if specification.URL.String() != discoveryItem.APISpecification.URL {
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
	var failures []ValidationFailure

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationFromSchemaVersion(schemaVersion)
		if err != nil {
			continue // err already added to failures in hasValidAPISpecifications
		}

		for endpointIndex, endpoint := range discoveryItem.Endpoints {
			isPresent, err := checker.IsPresent(endpoint.Method, endpoint.Path, specification.Identifier)
			if err != nil {
				failure := ValidationFailure{
					Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints[%d]", discoveryItemIndex, endpointIndex),
					Error: err.Error(),
				}
				logrus.WithFields(logrus.Fields{
					"function": "HasValidEndpoints",
					"module":   "func_validator",
					"package":  "discovery",
				}).Debugf("failure=%#v", failure)

				failures = append(failures, failure)
				continue
			}
			if !isPresent {
				failure := ValidationFailure{
					Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].Endpoints[%d]", discoveryItemIndex, endpointIndex),
					Error: fmt.Sprintf("Invalid endpoint Method='%s', Path='%s'", endpoint.Method, endpoint.Path),
				}

				logrus.WithFields(logrus.Fields{
					"function": "HasValidEndpoints",
					"module":   "func_validator",
					"package":  "discovery",
				}).Debugf("failure=%#v", failure)

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
	var failures []ValidationFailure

	for discoveryItemIndex, discoveryItem := range discoveryConfig.DiscoveryModel.DiscoveryItems {
		schemaVersion := discoveryItem.APISpecification.SchemaVersion
		specification, err := model.SpecificationFromSchemaVersion(schemaVersion)
		if err != nil {
			continue // err already added to failures in hasValidAPISpecifications
		}

		var discoveryEndpoints []model.Input
		for _, endpoint := range discoveryItem.Endpoints {
			discoveryEndpoints = append(discoveryEndpoints, model.Input{Endpoint: endpoint.Path, Method: endpoint.Method})
		}
		var missingMandatory []model.Input // disable mandatory endpoint checking for payment apis - to allow single token test runs
		if strings.HasPrefix(specification.Identifier, "payment") {
			logrus.Trace("Skipping payment spec mandatory endpoint check")
		} else {
			missingMandatory, err = checker.MissingMandatory(discoveryEndpoints, specification.Identifier)
		}
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
			logrus.Warnf("Missing mandatory endpoint Method='%s', Path='%s'", mandatoryEndpoint.Method, mandatoryEndpoint.Endpoint)
			failures = append(failures, failure)
		}
	}

	if len(failures) > 0 {
		return false, failures
	}

	return true, failures
}
