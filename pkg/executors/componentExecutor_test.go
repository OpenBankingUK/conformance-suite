package executors

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTestCases(t *testing.T) {
	tests, err := manifest.GenerateTestCases("TestSpec", "http://mybaseurl")
	assert.Nil(t, err)

	// perms, err := manifest.GetTestCasePermissions(tests)
	// assert.Nil(t, err)
	// requiredTokens, err := manifest.GetRequiredTokens(perms)
	// assert.Nil(t, err)
	//dumpJSON(requiredTokens)

	ctx := model.Context{
		"client_id":              "******myid",
		"fapi_financial_id":      "*****finid",
		"basic_authentication":   "****basicauth",
		"token_endpoint":         "****tokend",
		"authorisation_endpoint": "****authend",
		"resource_server":        "****resend",
		"redirect_url":           "****redirurl",
		"permission_payload":     "****permpay",
		"result_token":           "****mytoken",
	}

	err = AcquireHeadlessTokens(tests, &ctx, RunDefinition{})
	assert.Nil(t, err)

}
