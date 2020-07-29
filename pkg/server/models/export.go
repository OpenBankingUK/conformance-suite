package models

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ExportRequest - Request to `/api/export`.
type ExportRequest struct {
	Environment         string   `json:"environment"`           // Environment used for testing
	Implementer         string   `json:"implementer"`           // Implementer/Brand Name
	AuthorisedBy        string   `json:"authorised_by"`         // Authorised by
	JobTitle            string   `json:"job_title"`             // Job Title
	Products            []string `json:"products"`              // Products tested, e.g., "Business, Personal, Cards"
	HasAgreed           bool     `json:"has_agreed"`            // I agree
	AddDigitalSignature bool     `json:"add_digital_signature"` // Sign this report
}

func (e *ExportRequest) requiresTCAgreement() bool {
	if e.Environment == "production" {
		return true
	}
	return false
}

func (e ExportRequest) Validate() error {
	rules := []*validation.FieldRules{
		validation.Field(&e.Environment, validation.Required),
		validation.Field(&e.Implementer, validation.Required),
		validation.Field(&e.AuthorisedBy, validation.Required),
		validation.Field(&e.JobTitle, validation.Required),
		validation.Field(&e.Products, validation.Required, validation.By(productsValuesValidator)),
	}

	if e.requiresTCAgreement() {
		rules = append(rules, validation.Field(&e.HasAgreed, validation.Required, validation.In(true)))
	}

	return validation.ValidateStruct(&e, rules...)
}

func productsValuesValidator(value interface{}) error {
	values, ok := value.([]string)
	if !ok {
		return fmt.Errorf("pkg/server/models.ExportRequest: 'products' (%+q) is not []string", value)
	}

	supportedValues := []string{
		"Business",
		"Personal",
		"Cards",
	}
	if len(values) > len(supportedValues) {
		return fmt.Errorf("pkg/server/models.ExportRequest: 'products' (%d) contains more than supported values (%d)", len(values), len(supportedValues))
	}

	for _, value := range values {
		if countOccurrences(values, value) >= 2 {
			return fmt.Errorf("pkg/server/models.ExportRequest: 'products' (%+q) contains duplicate value (%+q)", values, value)
		}
		if countOccurrences(supportedValues, value) == 0 {
			return fmt.Errorf("pkg/server/models.ExportRequest: 'products' (%+q) invalid value provided (%+q)", values, value)
		}
	}

	return nil
}

func countOccurrences(slice []string, str string) int {
	count := 0
	for _, s := range slice {
		if s == str {
			count += 1
		}
	}
	return count
}

// ExportResults - Contains `ExportRequest` and results of test run.
type ExportResults struct {
	ExportRequest    ExportRequest                             `json:"export_request"`
	HasPassed        bool                                      `json:"has_passed"`
	Results          map[results.ResultKey][]results.TestCase  `json:"results"`
	Tokens           []events.AcquiredAccessToken              `json:"tokens"`
	DiscoveryModel   discovery.Model                           `json:"discovery_model"`
	ResponseFields   string                                    `json:"-"`
	TLSVersionResult map[string]*discovery.TLSValidationResult `json:"-"`
	JWSStatus        string                                    `json:"jws_status"`
}
