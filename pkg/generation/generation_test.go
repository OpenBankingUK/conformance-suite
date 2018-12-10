package generation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/go-openapi/loads"
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
	testNo := 1000

	for _, v := range disco.DiscoveryModel.DiscoveryItems {
		result := GetImplementedTestCases(&v, false, testNo)
		results = append(results, result...)
		testNo += 1000
	}

	for _, tc := range results {
		fmt.Println(string(pkgutils.DumpJSON(&tc)))
	}
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadModelOBv3Ozone() (discovery.Model, error) {
	filedata, _ := ioutil.ReadFile("testdata/disco.json")
	var d discovery.Model
	err := json.Unmarshal(filedata, &d)
	if err != nil {
		return discovery.Model{}, err
	}
	return d, nil
}

func printSpec(doc *loads.Document, base, spec string) {
	for path, props := range doc.Spec().Paths.Paths {
		for method := range getOperations(&props) {
			newPath := base + path
			condition := getConditionality(method, path, spec)
			fmt.Printf("[%s] %s %s\n", condition, method, newPath) // give to testcase along with any conditionality?
		}
	}
}
