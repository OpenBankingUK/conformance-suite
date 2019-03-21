package manifest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermx(t *testing.T) {
	tests, err := GenerateTestCases("TestSpec", "http://mybaseurl")
	assert.Nil(t, err)
	testcasePermissions, err := GetTestCasePermissions(tests)
	assert.Nil(t, err)
	requiredTokens, err := GatherTokens(testcasePermissions)
	dumpJSON(requiredTokens)
}

func TestGetScriptConsentTokens(t *testing.T) {
	tests, err := GenerateTestCases("TestSpec", "http://mybaseurl")
	assert.Nil(t, err)
	testcasePermissions, err := GetTestCasePermissions(tests)
	assert.Nil(t, err)
	requiredTokens, err := GatherTokens(testcasePermissions)
	populateTokens(requiredTokens)
	dumpJSON(requiredTokens)
}

func populateTokens(gatherer []TokenGatherer) error {
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
