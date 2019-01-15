package generation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This Example runs and verifies example code. See: https://golang.org/pkg/testing/#hdr-Examples
// Its purpose is to exercise the discovery to test case mapping
func ExampleGetImplementedTestCases() {
	results := []model.TestCase{}
	disco, err := loadModelOBv3Ozone()
	if err != nil {
		// This Example function fails when output does not match expectation below
		fmt.Println(err.Error())
	}
	testNo := 1000

	replacements := make(map[string]string)
	for _, v := range disco.DiscoveryModel.DiscoveryItems {
		result := GetImplementedTestCases(&v, testNo, replacements)
		results = append(results, result...)
		testNo += 1000
	}

	data, err := json.MarshalIndent(results[0], "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(data))
	// Output:
	// {
	//     "@id": "#t1000",
	//     "name": "Create Account Access Consents",
	//     "input": {
	//         "method": "POST",
	//         "endpoint": "/account-access-consents"
	//     },
	//     "expect": {
	//         "status-code": 201,
	//         "schema-validation": true,
	//         "contextPut": {}
	//     }
	// }
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

func TestGetConsentTempatedTestCasesPOST(t *testing.T) {
	newpath := "POST /open-banking/v3.0/aisp/account-access-consents"
	customTestCases, err := getTemplatedTestCases(newpath)
	require.NoError(t, err)
	dumpTestCases(customTestCases)
	assert.Equal(t, 3, len(customTestCases))
}

func TestGetConsentTempatedTestCasesGETDELETE(t *testing.T) {
	newpath := "GET /open-banking/v3.0/aisp/account-access-consents"
	customTestCases, err := getTemplatedTestCases(newpath)
	require.NoError(t, err)
	assert.Equal(t, 0, len(customTestCases))
}
