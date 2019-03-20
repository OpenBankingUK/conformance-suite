package manifest

import (
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
