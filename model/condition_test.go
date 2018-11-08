package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test we can read the conditionality file - performed in package init()
// and that the endPointConditionality structure - which holds all the endpoint conditions
// in a package local manner, can be read
func TestConditionality(t *testing.T) {
	count := len(endpointConditionality)
	result := count > 10 // check we have more than an arbitary number of conditions
	assert.Equal(t, result, true)
}

// Test all Mandatory endpoints are correctly configured in model
func TestMandatoryEndpoint(t *testing.T) {
	result, err := IsMandatory("POST", "/account-access-consents")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("GET", "/account-access-consents/{ConsentId}")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("DELETE", "/account-access-consents/{ConsentId}")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("GET", "/accounts")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("GET", "/accounts/{AccountId}")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("GET", "/accounts/{AccountId}/balances")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsMandatory("GET", "/accounts/{AccountId}/transactions")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/scheduled-payments")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/scheduled-payments")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/party")
	assert.Nil(t, err)
	assert.True(t, result)
}

// Test all Conditional endpoints are correctly configured in model
func TestConditionalEndpoint(t *testing.T) {
	result, err := IsConditional("GET", "/accounts/{AccountId}/direct-debits")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/standing-orders")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/product")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/offers")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/statements")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/statements/{StatementId}")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/party")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/beneficiaries")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsConditional("GET", "/accounts/{AccountId}/statements/{StatementId}/transactions")
	assert.Nil(t, err)
	assert.True(t, result)
}

// Test all Optional  endpoints are correctly configured in model
func TestOptionalEndpoint(t *testing.T) {
	result, err := IsOptional("GET", "/offers")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/party")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/standing-orders")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/direct-debits")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/beneficiaries")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/transactions")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/products")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/accounts/{AccountId}/statements/{StatementId}/file")
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = IsOptional("GET", "/statements")
	assert.Nil(t, err)
	assert.True(t, result)
}
