package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	gmocks "bitbucket.org/openbankingteam/conformance-suite/pkg/generation/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJourneySetDiscoveryModelValidatesModel(t *testing.T) {
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)

	discoveryModelInternal := discovery.ModelDiscovery{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModelInternal).Return([]generation.SpecificationTestCases{})

	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	assert.Equal(t, discovery.NoValidationFailures, failures)
	validator.AssertExpectations(t)
}

func TestJourneySetDiscoveryModelHandlesErrorFromValidator(t *testing.T) {
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, errors.New("validator error"))

	discoveryModelInternal := discovery.ModelDiscovery{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModelInternal).Return([]generation.SpecificationTestCases{})

	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.Error(t, err)

	assert.Equal(t, "error setting discovery model: validator error", err.Error())
	assert.Nil(t, failures)
}

func TestJourneySetDiscoveryModelReturnsFailuresFromValidator(t *testing.T) {
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	failure := discovery.ValidationFailure{
		Key:   "DiscoveryModel.Name",
		Error: "Field 'Name' is required",
	}
	expectedFailures := discovery.ValidationFailures{failure}
	validator.On("Validate", discoveryModel).Return(expectedFailures, nil)

	discoveryModelInternal := discovery.ModelDiscovery{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModelInternal).Return([]generation.SpecificationTestCases{})

	journey := NewWebJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	assert.Equal(t, expectedFailures, failures)
}
