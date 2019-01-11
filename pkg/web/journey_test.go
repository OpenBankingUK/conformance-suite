package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	gmocks "bitbucket.org/openbankingteam/conformance-suite/pkg/generation/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJourneySetDiscoveryModelValidatesModel(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	generator := &gmocks.Generator{}
	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(t, discovery.NoValidationFailures, failures)
	validator.AssertExpectations(t)
	generator.AssertExpectations(t)
}

func TestJourneySetDiscoveryModelHandlesErrorFromValidator(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, errors.New("validator error"))
	generator := &gmocks.Generator{}
	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.Error(t, err)
	assert.Equal(t, "error setting discovery model: validator error", err.Error())
	assert.Nil(t, failures)
}

func TestJourneySetDiscoveryModelReturnsFailuresFromValidator(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	failure := discovery.ValidationFailure{
		Key:   "DiscoveryModel.Name",
		Error: "Field 'Name' is required",
	}
	expectedFailures := discovery.ValidationFailures{failure}
	validator.On("Validate", discoveryModel).Return(expectedFailures, nil)
	generator := &gmocks.Generator{}
	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(t, expectedFailures, failures)
}

func TestJourneyTestCasesCantGenerateIfDiscoveryNotSet(t *testing.T) {
	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewWebJourney(generator, validator)

	testCases, err := journey.TestCases()

	assert.Error(t, err)
	assert.Nil(t, testCases)
}

func TestJourneyTestCasesGenerate(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	expectedTestCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).Return(expectedTestCases)
	journey := NewWebJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	testCases, err := journey.TestCases()

	assert.NoError(t, err)
	assert.Equal(t, expectedTestCases, testCases)
}

func TestJourneyTestCasesDoesntREGenerate(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	expectedTestCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).
		Return(expectedTestCases).Times(1)

	journey := NewWebJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)
	firstRunTestCases, err := journey.TestCases()
	require.NoError(t, err)

	testCases, err := journey.TestCases()

	assert.NoError(t, err)
	assert.Equal(t, expectedTestCases, testCases)
	assert.Equal(t, firstRunTestCases, testCases)
	generator.AssertExpectations(t)
}

func TestJourneyRunTestCasesCantRunIfNoTestCases(t *testing.T) {
	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewWebJourney(generator, validator)

	result, err := journey.RunTests()

	assert.EqualError(t, err, "error running test cases, test cases not set")
	assert.Equal(t, reporting.Result{}, result)
}

func TestJourneyRunTestCases(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	testCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).
		Return(testCases).Times(1)

	journey := NewWebJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	_, err = journey.TestCases()
	require.NoError(t, err)

	result, err := journey.RunTests()

	assert.NoError(t, err)
	noResult := []reporting.Specification([]reporting.Specification{})
	assert.Equal(t, noResult, result.Specifications)
	generator.AssertExpectations(t)
}
