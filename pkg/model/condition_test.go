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

// Test we can read the conditionality file - performed in package init()
// and that the endPointConditionality structure - which holds all the endpoint conditions
// in a package local manner, can be read
func TestConditionality(t *testing.T) {
	count := len(endpointConditionality)
	result := count > 10 // check we have more than an arbitary number of conditions
	assert.Equal(t, result, true)
}

// Test that GetEndpointConditionality returns a clone of endpointConditionality.
func TestGetEndpointConditionality(t *testing.T) {
	assert := assert.New(t)

	assert.Len(GetEndpointConditionality(), len(endpointConditionality))
	assert.EqualValues(endpointConditionality, GetEndpointConditionality())

	// modify returned clone and ensure the original wasn't modified - just being pedantic
	clone := GetEndpointConditionality()
	clone[0].Endpoint = ""

	assert.EqualValues(endpointConditionality, GetEndpointConditionality())
}

// Test all Mandatory endpoints are correct and configured in the model
func TestMandatoryData(t *testing.T) {
	checker := ConditionalityChecker{}
	for _, tt := range mandatoryData {
		result, err := checker.IsMandatory(tt.Method, tt.Endpoint)
		require.Nil(t, err)
		require.True(t, result)
	}
}

// Test all Conditional endpoints are correctly configured in model
func TestConditionalData(t *testing.T) {
	checker := ConditionalityChecker{}
	for _, tt := range conditionalData {
		result, err := checker.IsConditional(tt.Method, tt.Endpoint)
		require.Nil(t, err)
		require.True(t, result)
	}
}

// Test all Optional  endpoints are correctly configured in model
func TestOptionalData(t *testing.T) {
	checker := ConditionalityChecker{}
	for _, tt := range optionalData {
		result, err := checker.IsOptional(tt.Method, tt.Endpoint)
		require.Nil(t, err)
		require.True(t, result)
	}
}
