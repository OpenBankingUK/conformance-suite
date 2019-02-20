package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dataHolder struct {
	Method   string
	Endpoint string
}

var mandatoryData = []dataHolder{
	{"POST", "/account-access-consents"},
	{"GET", "/account-access-consents/{ConsentId}"},
	{"DELETE", "/account-access-consents/{ConsentId}"},
	{"GET", "/accounts"},
	{"GET", "/accounts/{AccountId}"},
	{"GET", "/accounts/{AccountId}/balances"},
	{"GET", "/accounts/{AccountId}/transactions"},
}

var conditionalData = []dataHolder{
	{"GET", "/accounts/{AccountId}/scheduled-payments"},
	{"GET", "/accounts/{AccountId}/direct-debits"},
	{"GET", "/accounts/{AccountId}/standing-orders"},
	{"GET", "/accounts/{AccountId}/product"},
	{"GET", "/accounts/{AccountId}/offers"},
	{"GET", "/accounts/{AccountId}/statements"},
	{"GET", "/accounts/{AccountId}/statements/{StatementId}"},
	{"GET", "/accounts/{AccountId}/party"},
	{"GET", "/accounts/{AccountId}/beneficiaries"},
	{"GET", "/accounts/{AccountId}/statements/{StatementId}/transactions"},
}

var optionalData = []dataHolder{
	{"GET", "/scheduled-payments"},
	{"GET", "/party"},
	{"GET", "/offers"},
	{"GET", "/standing-orders"},
	{"GET", "/direct-debits"},
	{"GET", "/beneficiaries"},
	{"GET", "/transactions"},
	{"GET", "/products"},
	{"GET", "/accounts/{AccountId}/statements/{StatementId}/file"},
	{"GET", "/statements"},
}

const ACCOUNT_SPEC_ID = "account-transaction-v3.1"

// Test we can read the conditionality file - performed in package init()
// and that the endPointConditionality structure - which holds all the endpoint conditions
// in a package local manner, can be read
func TestConditionality(t *testing.T) {
	count := len(endpointConditionality[ACCOUNT_SPEC_ID])
	result := count > 10 // check we have more than an arbitrary number of conditions
	assert.Equal(t, result, true)
}

// Test that GetEndpointConditionality returns a clone of endpointConditionality.
func TestGetEndpointConditionality(t *testing.T) {
	assert := assert.New(t)
	specification := ACCOUNT_SPEC_ID
	assert.Len(GetEndpointConditionality(specification), len(endpointConditionality[specification]))
	assert.EqualValues(endpointConditionality[specification], GetEndpointConditionality(specification))

	// modify returned clone and ensure the original wasn't modified - just being pedantic
	clone := GetEndpointConditionality(specification)
	clone[0].Endpoint = ""

	assert.EqualValues(endpointConditionality[specification], GetEndpointConditionality(specification))
}

func TestConditionalityChecker(t *testing.T) {
	checker := NewConditionalityChecker()
	specification := ACCOUNT_SPEC_ID

	t.Run("IsPresent true for endpoint method mix in specification", func(t *testing.T) {
		result, err := checker.IsPresent("POST", "/account-access-consents", specification)
		require.Nil(t, err)
		require.True(t, result)
	})

	t.Run("IsPresent false for endpoint method mix not in specification", func(t *testing.T) {
		result, err := checker.IsPresent("PUT", "/account-access-consents", specification)
		require.Nil(t, err)
		require.False(t, result)
	})

	t.Run("MissingMandatory returns array of missing mandatory endpoints", func(t *testing.T) {
		result, err := checker.MissingMandatory([]Input{}, specification)
		require.Nil(t, err)

		expectedAllMissing := []Input{}
		for _, tt := range mandatoryData {
			expectedAllMissing = append(expectedAllMissing, Input{Method: tt.Method, Endpoint: tt.Endpoint})
		}
		assert.Equal(t, result, expectedAllMissing)
	})

	t.Run("MissingMandatory returns empty array when no missing mandatory endpoints", func(t *testing.T) {
		allMandatory := []Input{}
		for _, tt := range mandatoryData {
			allMandatory = append(allMandatory, Input{Method: tt.Method, Endpoint: tt.Endpoint})
		}
		result, err := checker.MissingMandatory(allMandatory, specification)
		assert.Nil(t, err)

		expectedEmptyMissing := []Input{}
		assert.Equal(t, result, expectedEmptyMissing)
	})
}

// Test all Mandatory endpoints are correct and configured in the model
func TestMandatoryData(t *testing.T) {
	checker := NewConditionalityChecker()
	for _, tt := range mandatoryData {
		result, err := checker.IsMandatory(tt.Method, tt.Endpoint, ACCOUNT_SPEC_ID)
		require.Nil(t, err)
		require.True(t, result)
	}
}

// Test all Conditional endpoints are correctly configured in model
func TestConditionalData(t *testing.T) {
	checker := NewConditionalityChecker()
	for _, tt := range conditionalData {
		result, err := checker.IsConditional(tt.Method, tt.Endpoint, ACCOUNT_SPEC_ID)
		require.Nil(t, err)
		require.True(t, result)
	}
}

// Test all Optional  endpoints are correctly configured in model
func TestOptionalData(t *testing.T) {
	checker := NewConditionalityChecker()
	for _, tt := range optionalData {
		result, err := checker.IsOptional(tt.Method, tt.Endpoint, ACCOUNT_SPEC_ID)
		require.Nil(t, err)
		require.True(t, result)
	}
}
