package generation

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/stretchr/testify/require"
)

// This test case intentionally doesn't assert anything
// Its purpose is to exercise the discovery to test case mapping
func TestEnumerateOpenApiTestcases(t *testing.T) {
	dmodel, err := loadModelOBv3Ozone()
	require.NoError(t, err)
	for _, dItem := range dmodel.DiscoveryModel.DiscoveryItems {
		fmt.Printf("\n=========================================\n%s\n=========================================", dItem.APISpecification.Name)
		fmt.Printf("\n%s\n--------------\n", dItem.APISpecification.Version)
		doc, err := loadSpec(dItem.APISpecification.SchemaVersion, false)
		require.NoError(t, err)
		printSpec(doc, dItem.ResourceBaseURI, dItem.APISpecification.Version) // print the endpoints in the spec
		fmt.Printf("\nResourceIds\n-----------\n")
		printResourceIds(&dItem)
		fmt.Printf("\nImplemented\n--------------\n")
		printImplemented(dItem, dItem.Endpoints, dItem.APISpecification.Version) // print what this org has implemeneted
	}
}

// This test cases intentionally doesn't assert anything
// Its purpose is to exercise the discovery to test case mapping
func TestGenerateTestCases(t *testing.T) {
	results := []model.TestCase{}
	disco, _ := loadModelOBv3Ozone()

	for _, v := range disco.DiscoveryModel.DiscoveryItems {
		result := GetImplementedTestCases(&v)
		results = append(results, result...)
	}

	for _, tc := range results {
		fmt.Println(string(pkgutils.DumpJSON(&tc)))
	}
}
