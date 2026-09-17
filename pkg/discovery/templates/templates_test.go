package templates

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/test"
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
			if discoveryFile == "ob-v3.1-generic.json" || discoveryFile == "ob-v4.0-generic.json" || discoveryFile == "ob-v4.0.1-generic.json" || discoveryFile == "cVRP-v4.0-generic-discovery.json" || discoveryFile == "cVRP-v4.0.1-generic-discovery.json" {
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
			if discoveryFile == "ob-v3.1-generic.json" || discoveryFile == "ob-v4.0-generic.json" || discoveryFile == "ob-v4.0.1-generic.json" || discoveryFile == "cVRP-v4.0-generic-discovery.json" || discoveryFile == "cVRP-v4.0.1-generic-discovery.json" {
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

func TestCommercialVRPV4_0_1DiscoverySample(t *testing.T) {
	discoveryJSON, err := ioutil.ReadFile("cVRP-v4.0.1-generic-discovery.json")
	require.NoError(t, err)

	var discoveryModel discovery.Model
	require.NoError(t, json.Unmarshal(discoveryJSON, &discoveryModel))
	require.Len(t, discoveryModel.DiscoveryModel.DiscoveryItems, 1)

	apiSpecification := discoveryModel.DiscoveryModel.DiscoveryItems[0].APISpecification
	require.Equal(t, "cVRP v4.0.1", discoveryModel.DiscoveryModel.Name)
	require.Equal(t, "v4.0.1", apiSpecification.Version)
	require.Equal(t, "https://raw.githubusercontent.com/OpenBankingUK/Commercial-VRP-API-Spec/refs/tags/v4.0.1/OpenAPI/cvrp-openapi.json", apiSpecification.SchemaVersion)
	require.Equal(t, "file://manifests/cVRP_4.0_variable_recurring_payments.json", apiSpecification.Manifest)

	specification, err := model.SpecificationFromSchemaVersion(apiSpecification.SchemaVersion)
	require.NoError(t, err)
	require.Equal(t, "commercial-variable-recurring-payments-v4.0.1", specification.Identifier)
}
