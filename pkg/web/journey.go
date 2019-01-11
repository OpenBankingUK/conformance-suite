package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/pkg/errors"
)

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() ([]generation.SpecificationTestCases, error)
	RunTests() (reporting.Result, error)
}

type journey struct {
	generator           generation.Generator
	testCases           []generation.SpecificationTestCases
	validator           discovery.Validator
	validDiscoveryModel *discovery.Model
	reportService       reporting.Service
}

// NewWebJourney creates an instance for a user journey
func NewWebJourney(generator generation.Generator, validator discovery.Validator) Journey {
	return &journey{
		generator:     generator,
		validator:     validator,
		reportService: reporting.NewMockedService(),
	}
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
	wj.testCases = nil

	return discovery.NoValidationFailures, nil
}

var errDiscoveryModelNotSet = errors.New("error generation test cases, discovery model not set")

func (wj *journey) TestCases() ([]generation.SpecificationTestCases, error) {
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	if wj.testCases == nil {
		wj.testCases = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
	}
	return wj.testCases, nil
}

var errTestCasesNotSet = errors.New("error running test cases, test cases not set")

func (wj *journey) RunTests() (reporting.Result, error) {
	if wj.testCases == nil {
		return reporting.Result{}, errTestCasesNotSet
	}
	return wj.reportService.Run(wj.testCases)
}
