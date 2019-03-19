package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermx(t *testing.T) {
	tests, err := GenerateTestCases("TestSpec", "http://mybaseurl")
	assert.Nil(t, err)
	tcp, err := GetTestCasePermissions(tests)
	assert.Nil(t, err)
	dumpJSON(tcp)
	requiredTokens, err := GatherTokens(tcp)
	dumpJSON(requiredTokens)
}
