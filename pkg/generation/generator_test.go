package generation_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
)

func testLoadDiscoveryModel(t *testing.T) *discovery.ModelDiscovery {
	t.Helper()
	template, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	require.NoError(t, err)
	require.NotNil(t, template)
	json := string(template)
	model, err := discovery.UnmarshalDiscoveryJSON(json)
	require.NoError(t, err)
	return &model.DiscoveryModel
}

func TestNewGenerator(t *testing.T) {
	t.Run("returns Generator", func(t *testing.T) {
		generator := generation.NewGenerator()
		require.NotNil(t, generator)
	})
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
		for i, specificationCases := range cases {
			expectedSpec := discovery.DiscoveryItems[i].APISpecification
			assert.Equal(t, expectedSpec, specificationCases.Specification)
		}
	})

	t.Run("returns each SpecificationTestCases with generated TestCases", func(t *testing.T) {
		expectedCount := []int{8, 31}
		for i, specificationCases := range cases {
			assert.Len(t, specificationCases.TestCases, expectedCount[i])
		}
	})
}

// This Example runs and verifies example code. See: https://golang.org/pkg/testing/#hdr-Examples
// We deliberately check only a couple of test cases in output, as we are just
// checking here that test cases are being generated in the general case.
// Unit testing of test case generation under varying output is to be done elsewhere.
func ExampleGenerator() {
	template, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	if err != nil {
		// This Example function fails when output does not match expectation below
		fmt.Println(err.Error())
	}
	json := string(template)
	model, err := discovery.UnmarshalDiscoveryJSON(json)
	if err != nil {
		fmt.Println(err.Error())
	}
	discovery := model.DiscoveryModel
	generator := generation.NewGenerator()
	cases := generator.GenerateSpecificationTestCases(discovery)

	specOneTestCase := string(pkgutils.DumpJSON(cases[0].TestCases[1]))
	specTwoTestCase := string(pkgutils.DumpJSON(cases[1].TestCases[0]))
	fmt.Println(specOneTestCase)
	fmt.Println(specTwoTestCase)
	// Output:
	// {
	//     "@id": "#t1001",
	//     "name": "Get Account Access Consents",
	//     "input": {
	//         "method": "GET",
	//         "endpoint": "/account-access-consents/{ConsentId}",
	//         "contextGet": {}
	//     },
	//     "expect": {
	//         "status-code": 200,
	//         "schema-validation": true,
	//         "contextPut": {}
	//     }
	// }
	// {
	//     "@id": "#t2000",
	//     "name": "Create Domestic Payment Consents",
	//     "input": {
	//         "method": "POST",
	//         "endpoint": "/domestic-payment-consents",
	//         "contextGet": {}
	//     },
	//     "expect": {
	//         "status-code": 201,
	//         "schema-validation": true,
	//         "contextPut": {}
	//     }
	// }
}
