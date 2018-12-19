package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)


func TestJourneySetDiscoveryModelValidatesModel(t *testing.T) {
	// this is a global var/singleton so we need to reset it's state between tests
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	journey := NewWebJourney(validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	assert.Equal(t, discovery.NoValidationFailures, failures)
	validator.AssertExpectations(t)
}

func TestJourneySetDiscoveryModelHandlesErrorFromValidator(t *testing.T) {
	// this is a global var/singleton so we need to reset it's state between tests
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, errors.New("validator error"))
	journey := NewWebJourney(validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.Error(t, err)

	assert.Equal(t, "error setting discovery model: validator error", err.Error())
	assert.Nil(t, failures)
}

func TestJourneySetDiscoveryModelReturnsFailuresFromValidator(t *testing.T) {
	// this is a global var/singleton so we need to reset it's state between tests
	journeyInstance = nil
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, nil)
	journey := NewWebJourney(validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	assert.Equal(t, expectedFailures, failures)
}
