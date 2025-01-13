package schema

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

type queryParamValidator struct {
	finder finder
}

func newQueryParamValidator(f finder) Validator {
	return &queryParamValidator{
		finder: f,
	}
}

func (v *queryParamValidator) Validate(r HTTPResponse) ([]Failure, error) {
	failures := []Failure{}

	// Get operation and handle potential error
	op, err := v.finder.Operation(r.Method, r.Path)
	if err != nil {
		if err == ErrNotFound {
			// If the path/method combination isn't found, we might want to skip validation
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to find operation in spec")
	}

	// Parse URL and extract query parameters
	requestURL, err := url.Parse(r.Path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse request URL")
	}
	queryValues := requestURL.Query()

	// Validate each parameter
	for _, param := range op.Parameters {
		// Only process query parameters
		if param.ParamProps.In != "query" {
			continue
		}

		paramName := param.ParamProps.Name
		queryValue := queryValues.Get(paramName)

		// Check required parameters
		if param.ParamProps.Required && queryValue == "" {
			failures = append(failures, newFailure(fmt.Sprintf(
				"required query parameter '%s' is missing", paramName)))
			continue
		}

		// Skip validation if parameter is not provided and not required
		if queryValue == "" {
			continue
		}

		// Validate against schema
		if err := validateQueryParamValue(param, queryValue); err != nil {
			failures = append(failures, newFailure(fmt.Sprintf(
				"query parameter '%s' validation failed: %s", paramName, err.Error())))
		}
	}

	return failures, nil
}

// Implement the IsRequestProperty interface method
func (v *queryParamValidator) IsRequestProperty(method, path, propertyPath string) (bool, string, error) {
	// Query parameters don't need this implementation as it's mainly for body validation
	return false, "", nil
}

func validateQueryParamValue(param spec.Parameter, value string) error {
	// Handle items directly if it's an array
	if param.Type == "array" {
		return validateArrayQueryParam(param, value)
	}

	switch param.Type {
	case "string":
		return validateStringQueryParam(param, value)
	case "integer":
		return validateIntegerQueryParam(param, value)
	case "number":
		return validateNumberQueryParam(param, value)
	case "boolean":
		return validateBooleanQueryParam(param, value)
	default:
		return fmt.Errorf("unsupported parameter type: %s", param.Type)
	}
}

func validateStringQueryParam(param spec.Parameter, value string) error {
	// Check enum if specified
	if len(param.Enum) > 0 {
		valid := false
		for _, enum := range param.Enum {
			if fmt.Sprint(enum) == value {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("value must be one of: %v", param.Enum)
		}
	}

	// Check minLength
	if param.MinLength != nil {
		if len(value) < int(*param.MinLength) {
			return fmt.Errorf("value length must be >= %d", *param.MinLength)
		}
	}

	// Check maxLength
	if param.MaxLength != nil {
		if len(value) > int(*param.MaxLength) {
			return fmt.Errorf("value length must be <= %d", *param.MaxLength)
		}
	}

	// Check pattern
	if param.Pattern != "" {
		matched, err := regexp.MatchString(param.Pattern, value)
		if err != nil {
			return fmt.Errorf("invalid pattern: %v", err)
		}
		if !matched {
			return fmt.Errorf("value does not match pattern: %s", param.Pattern)
		}
	}

	return nil
}

func validateIntegerQueryParam(param spec.Parameter, value string) error {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer value")
	}

	// Check minimum
	if param.Minimum != nil {
		min := int64(*param.Minimum)
		if param.ExclusiveMinimum {
			if num <= min {
				return fmt.Errorf("value must be > %d", min)
			}
		} else if num < min {
			return fmt.Errorf("value must be >= %d", min)
		}
	}

	// Check maximum
	if param.Maximum != nil {
		max := int64(*param.Maximum)
		if param.ExclusiveMaximum {
			if num >= max {
				return fmt.Errorf("value must be < %d", max)
			}
		} else if num > max {
			return fmt.Errorf("value must be <= %d", max)
		}
	}

	// Check multipleOf
	if param.MultipleOf != nil {
		multiple := int64(*param.MultipleOf)
		if multiple != 0 && num%multiple != 0 {
			return fmt.Errorf("value must be multiple of %d", multiple)
		}
	}

	return nil
}

func validateNumberQueryParam(param spec.Parameter, value string) error {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid number value")
	}

	// Check minimum
	if param.Minimum != nil {
		if param.ExclusiveMinimum {
			if num <= *param.Minimum {
				return fmt.Errorf("value must be > %v", *param.Minimum)
			}
		} else if num < *param.Minimum {
			return fmt.Errorf("value must be >= %v", *param.Minimum)
		}
	}

	// Check maximum
	if param.Maximum != nil {
		if param.ExclusiveMaximum {
			if num >= *param.Maximum {
				return fmt.Errorf("value must be < %v", *param.Maximum)
			}
		} else if num > *param.Maximum {
			return fmt.Errorf("value must be <= %v", *param.Maximum)
		}
	}

	// Check multipleOf
	if param.MultipleOf != nil && *param.MultipleOf != 0 {
		multiple := *param.MultipleOf
		quotient := num / multiple
		if math.Round(quotient) != quotient {
			return fmt.Errorf("value must be multiple of %v", multiple)
		}
	}

	return nil
}

func validateBooleanQueryParam(_ spec.Parameter, value string) error {
	validValues := map[string]bool{
		"true":  true,
		"false": true,
		"1":     true,
		"0":     true,
	}

	if !validValues[strings.ToLower(value)] {
		return fmt.Errorf("invalid boolean value: must be true/false or 1/0")
	}

	return nil
}

func validateArrayQueryParam(param spec.Parameter, value string) error {
	if param.Items == nil {
		return fmt.Errorf("array parameter must define 'items' schema")
	}

	// Determine the delimiter based on collectionFormat
	var delimiter string
	switch param.CollectionFormat {
	case "csv", "": // empty defaults to csv
		delimiter = ","
	case "ssv":
		delimiter = " "
	case "tsv":
		delimiter = "\t"
	case "pipes":
		delimiter = "|"
	default:
		return fmt.Errorf("unsupported collection format: %s", param.CollectionFormat)
	}

	// Split the value into items
	items := strings.Split(value, delimiter)

	// Validate array length constraints
	if param.MinItems != nil && len(items) < int(*param.MinItems) {
		return fmt.Errorf("array must contain at least %d items", *param.MinItems)
	}
	if param.MaxItems != nil && len(items) > int(*param.MaxItems) {
		return fmt.Errorf("array must contain at most %d items", *param.MaxItems)
	}

	// Validate each item
	for i, item := range items {
		// Create a temporary parameter for the item
		itemParam := spec.Parameter{
			SimpleSchema: spec.SimpleSchema{
				Type:   param.Items.Type,
				Format: param.Items.Format,
			},
			CommonValidations: spec.CommonValidations{
				Enum:      param.Items.Enum,
				Minimum:   param.Items.Minimum,
				Maximum:   param.Items.Maximum,
				MinLength: param.Items.MinLength,
				MaxLength: param.Items.MaxLength,
			},
		}

		if err := validateQueryParamValue(itemParam, item); err != nil {
			return fmt.Errorf("invalid item at index %d: %v", i, err)
		}
	}

	// Check uniqueItems if required
	if param.UniqueItems {
		seen := make(map[string]struct{})
		for _, item := range items {
			if _, exists := seen[item]; exists {
				return fmt.Errorf("array items must be unique")
			}
			seen[item] = struct{}{}
		}
	}

	return nil
}
