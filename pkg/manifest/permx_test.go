package manifest

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestPermx(t *testing.T) {
	tests, err := GenerateTestCases("Account and Transaction API Specification", "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)
	testcasePermissions, err := getTestCasePermissions(tests)
	assert.Nil(t, err)
	requiredTokens, err := getRequiredTokens(testcasePermissions)
	dumpJSON(requiredTokens)
}

func TestGetScriptConsentTokens(t *testing.T) {
	tests, err := GenerateTestCases("Account and Transaction API Specification", "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)
	testcasePermissions, err := getTestCasePermissions(tests)
	assert.Nil(t, err)
	requiredTokens, err := getRequiredTokens(testcasePermissions)
	populateTokens(requiredTokens)
	dumpJSON(requiredTokens)
}

func populateTokens(gatherer []RequiredTokens) error {
	fmt.Printf("%d entries found\n", len(gatherer))
	for k, tokenGatherer := range gatherer {
		if len(tokenGatherer.Perms) == 0 {
			continue
		}
		token, err := getToken(tokenGatherer.Perms)
		if err != nil {
			return err
		}
		tokenGatherer.Token = token
		gatherer[k] = tokenGatherer

	}
	return nil
}

func getToken(perms []string) (string, error) {
	// for headless - get the okens for the permissions here

	return "abigfattoken", nil
}
