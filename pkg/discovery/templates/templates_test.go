package templates

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/stretchr/testify/assert"
)

// Test that all the *.json discovery files parse correctly.
func TestDiscoverySamples_Examples_Parse_Correctly(t *testing.T) {
	// discoveryFiles, err := filepath.Glob("*.json")
	// assert.NoError(t, err)

	discoveryFiles := []string{
		// "custom.json",
		"ob-v3.0-generic.json",
		"ob-v3.0-ozone.json",
	}
	for _, discoveryFile := range discoveryFiles {
		t.Run("Parses_Without_Error_"+discoveryFile, func(t *testing.T) {
			assert := assert.New(t)

			discoveryJSON, err := ioutil.ReadFile(discoveryFile)
			assert.NoError(err)
			assert.NotNil(discoveryJSON)

			discoveryModel := &discovery.Model{}
			assert.NoError(json.Unmarshal(discoveryJSON, &discoveryModel))

			checker := model.NewConditionalityChecker()
			result, failures, err := discovery.Validate(checker, discoveryModel)
			assert.True(result)
			assert.NoError(err)
			assert.Equal([]string{}, failures)
		})
	}
}
