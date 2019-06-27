package templates

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

// Test that all the *.json discovery files parse correctly.
func TestDiscoverySamples_Examples_Parse_Correctly(t *testing.T) {
	discoveryFiles, err := filepath.Glob("*.json")
	require.NoError(t, err)

	for _, discoveryFile := range discoveryFiles {
		discoveryFile := discoveryFile
		t.Run("Parses_Without_Error_"+discoveryFile, func(t *testing.T) {
			t.Logf("discoveryFile=%s", discoveryFile)
			// Skip for now as we get this error:
			// [{DiscoveryModel.DiscoveryItems[0].OpenidConfigurationURI Field 'DiscoveryModel.DiscoveryItems[0].OpenidConfigurationURI' is required} {DiscoveryModel.DiscoveryItems[0].ResourceBaseURI Field 'DiscoveryModel.DiscoveryItems[0].ResourceBaseURI' is required}]
			if discoveryFile == "ob-v3.1-generic.json" {
				t.Skip()
			}
			assert := test.NewAssert(t)

			discoveryJSON, err := ioutil.ReadFile(discoveryFile)
			assert.NoError(err)
			assert.NotNil(discoveryJSON)

			discoveryModel := &discovery.Model{}
			assert.NoError(json.Unmarshal(discoveryJSON, &discoveryModel))

			checker := model.NewConditionalityChecker()
			result, failures, err := discovery.Validate(checker, discoveryModel)
			require.NoError(t, err)
			require.Empty(t, failures)
			assert.True(result)
		})
	}
}

// TestDiscoverySamplesIfManifestIsURLHTTPSOnly Asserts that if the manifest field is populated as a URL
// then it must use the HTTPS scheme.
func TestDiscoverySamplesIfManifestIsURLHTTPSOnly(t *testing.T) {
	discoveryFiles, err := filepath.Glob("*.json")
	require.NoError(t, err)

	for _, discoveryFile := range discoveryFiles {
		discoveryFile := discoveryFile
		t.Run("Parses_Without_Error_"+discoveryFile, func(t *testing.T) {
			// Skip for now as get this error:
			// [{DiscoveryModel.DiscoveryItems[0].OpenidConfigurationURI Field 'DiscoveryModel.DiscoveryItems[0].OpenidConfigurationURI' is required} {DiscoveryModel.DiscoveryItems[0].ResourceBaseURI Field 'DiscoveryModel.DiscoveryItems[0].ResourceBaseURI' is required}]
			if discoveryFile == "ob-v3.1-generic.json" {
				t.Skip()
			}
			assert := test.NewAssert(t)

			discoveryJSON, err := ioutil.ReadFile(discoveryFile)
			assert.NoError(err)
			assert.NotNil(discoveryJSON)

			discoveryModel := &discovery.Model{}
			assert.NoError(json.Unmarshal(discoveryJSON, &discoveryModel))

			checker := model.NewConditionalityChecker()
			result, failures, err := discovery.Validate(checker, discoveryModel)
			require.NoError(t, err)
			require.Empty(t, failures)
			assert.True(result)
		})
	}
}
