package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"github.com/pkg/errors"
)

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() []generation.SpecificationTestCases
}

type journey struct {
	generator           generation.Generator
	testCases           []generation.SpecificationTestCases
	validator           discovery.Validator
	validDiscoveryModel *discovery.Model
}

var journeyInstance Journey

// NewWebJourney creates an instance for a user journey, assumes one user only no concurrency
// so a singleton is returned
func NewWebJourney(generator generation.Generator, validator discovery.Validator) Journey {
	if journeyInstance == nil {
		journeyInstance = &journey{
			generator: generator,
			validator: validator,
		}
	}
	return journeyInstance
}

func (wj *journey) SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error) {
	failures, err := wj.validator.Validate(discoveryModel)
	if err != nil {
		return nil, errors.Wrap(err, "error setting discovery model")
	}

	if !failures.Empty() {
		return failures, nil
	}

	wj.validDiscoveryModel = discoveryModel
	wj.testCases = wj.generator.GenerateSpecificationTestCases(discoveryModel.DiscoveryModel)
	return discovery.NoValidationFailures, nil
}

func (wj *journey) TestCases() []generation.SpecificationTestCases {
	return wj.testCases
}
