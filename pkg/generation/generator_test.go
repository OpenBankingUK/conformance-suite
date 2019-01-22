package generation_test

import (
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
)

func testLoadDiscoveryModel(t *testing.T) *discovery.ModelDiscovery {
	t.Helper()
	template, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-generic.json")
	require.NoError(t, err)
	require.NotNil(t, template)
	json := string(template)
	model, err := discovery.UnmarshalDiscoveryJSON(json)
	require.NoError(t, err)
	return &model.DiscoveryModel
}

func TestGenerateSpecificationTestCases(t *testing.T) {
	discovery := *testLoadDiscoveryModel(t)
	generator := generation.NewGenerator()
	cases := generator.GenerateSpecificationTestCases(discovery)

	t.Run("returns slice of SpecificationTestCases, one per discovery item", func(t *testing.T) {
	 	require.NotNil(t, cases)
	 	assert.Equal(t, len(discovery.DiscoveryItems), len(cases))
	})

	t.Run("returns each SpecificationTestCases with a Specification matching discovery item", func(t *testing.T) {
	 	require.Equal(t, len(discovery.DiscoveryItems), len(cases))

	 	for i, specificationCases := range cases {
	 		if specificationCases.Specification.Name == "CustomTest-GetOzoneToken" {
	 			continue
	 		}
	 		expectedSpec := discovery.DiscoveryItems[i].APISpecification
	 		assert.Equal(t, expectedSpec, specificationCases.Specification)
	 	}
	})

}
