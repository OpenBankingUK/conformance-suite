package model

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccountConsent(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	assert.Nil(t, err)
	assert.Equal(t, "accounts.TokenProvider", c.Name)
}

func TestUseAComponent(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	assert.Nil(t, err)
	assert.Equal(t, "accounts.TokenProvider", c.Name)

	log.Println("Input Parameters ---------------------")
	for k, v := range c.InputParameters {
		log.Printf("%s, %s\n", k, v)
	}

	log.Println("Output Parameters ---------------------")
	for k, v := range c.OutputParameters {
		log.Printf("%s, %s\n", k, v)
	}
}

func TestRequiredParametersPresent(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	require.Nil(t, err)
	ctx := Context{
		"client_id":              "myid",
		"fapi_financial_id":      "finid",
		"basic_authentication":   "basicauth",
		"token_endpoint":         "tokend",
		"authorisation_endpoint": "authend",
		"resource_server":        "resend",
		"redirect_url":           "redirurl",
		"permission_payload":     "permpay",
		"result_token":           "mytoken",
	}
	err = c.ValidateParameters(&ctx)
	assert.Nil(t, err)
}

func TestInputParametersMissing(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	require.Nil(t, err)
	ctx := Context{
		"client_id":              "myid",
		"fapi_financial_id":      "finid",
		"basic_authentication":   "basicauth",
		"authorisation_endpoint": "authend",
		"resource_server":        "resend",
		"redirect_url":           "redirurl",
		"permission_payload":     "permpay",
		"result_token":           "mytoken",
	}
	err = c.ValidateParameters(&ctx)
	assert.NotNil(t, err)
	t.Logf(err.Error())
}

func TestOutputParametersMissing(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	require.Nil(t, err)
	ctx := Context{
		"client_id":              "myid",
		"fapi_financial_id":      "finid",
		"basic_authentication":   "basicauth",
		"token_endpoint":         "tokend",
		"authorisation_endpoint": "authend",
		"resource_server":        "resend",
		"redirect_url":           "redirurl",
		"permission_payload":     "permpay",
		"result_token1":          "mytoken",
	}
	err = c.ValidateParameters(&ctx)
	assert.NotNil(t, err)
	t.Logf(err.Error())
}

func TestComponentGetTests(t *testing.T) {
	c, err := LoadComponent("testdata/tokenProviderComponent.json")
	require.Nil(t, err)
	tests := c.GetTests()
	assert.Equal(t, 4, len(tests))
	assert.Equal(t, "Code Exchange", tests[3].Name)
}

func TestComponentHeadlessLoad(t *testing.T) {
	c, err := LoadComponent("../../components/headlessTokenProviderComponent.json")
	require.Nil(t, err)
	ctx := Context{
		"client_id":              "myid",
		"x-fapi-financial-id":    "finid",
		"basic_authentication":   "basicauth",
		"token_endpoint":         "tokend",
		"authorisation_endpoint": "authend",
		"resource_server":        "resend",
		"redirect_url":           "redirurl",
		"permission_payload":     "permpay",
		"result_token":           "mytoken",
	}
	err = c.ValidateParameters(&ctx)
	assert.Nil(t, err)

}
