//go:generate mockery -name Validator
package discovery

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// Validator defines a generic validation engine
type Validator interface {
	Validate(*Model) (ValidationFailures, error)
}

// NewFuncValidator creates a Validator engine using functional implementation
func NewFuncValidator(checker model.ConditionalityChecker) Validator {
	return funcWrapperValidator{
		validatorFunc:         Validate,
		conditionalityChecker: checker,
	}
}

// funcWrapperValidator is a wrapper for functional validator
type funcWrapperValidator struct {
	validatorFunc         func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error)
	conditionalityChecker model.ConditionalityChecker
}

func (v funcWrapperValidator) Validate(model *Model) (ValidationFailures, error) {
	_, failures, err := v.validatorFunc(v.conditionalityChecker, model)
	if err != nil {
		return nil, err
	}

	if len(failures) == 0 {
		return NoValidationFailures(), nil
	}

	return ValidationFailures(failures), nil
}

// ValidationFailure - Records validation failure key and error.
// e.g. ValidationFailure{
//        Key:   "DiscoveryModel.Name",
//        Error: "Field 'Name' is required",
//      }
type ValidationFailure struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// ValidationFailures represents a list of failures from a Validator engine
type ValidationFailures []ValidationFailure

// NoValidationFailures represents a return of no validation errors found
func NoValidationFailures() ValidationFailures {
	return ValidationFailures{}
}

// Empty returns if there are validation failures present
func (v ValidationFailures) Empty() bool {
	return len(v) == 0
}
