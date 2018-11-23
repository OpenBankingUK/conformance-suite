package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	// Get /accounts example json response from ozone
	getAccountResponse = []byte(`
	{
	    "Data": {
	        "Account": [{
	            "AccountId": "500000000000000000000001",
	            "Currency": "GBP",
	            "Nickname": "xxxx0101",
	            "AccountType": "Personal",
	            "AccountSubType": "CurrentAccount",
	            "Account": [{
	                "SchemeName": "SortCodeAccountNumber",
	                "Identification": "10000119820101",
	                "Name": "Mr. Roberto Rastapopoulos & Mr. Mitsuhirato",
	                "SecondaryIdentification": "Roll No. 001"
	            }]
	        }, {
	            "AccountId": "500000000000000000000007",
	            "Currency": "GBP",
	            "Nickname": "xxxx0001",
	            "AccountType": "Business",
	            "AccountSubType": "CurrentAccount",
	            "Account": [{
	                "SchemeName": "SortCodeAccountNumber",
	                "Identification": "10000190210001",
	                "Name": "Marios Amazing Carpentry Supplies Limited"
	            }]
	        }]
	    },
	    "Links": {
	        "Self": "http://modelobank2018.o3bank.co.uk/open-banking/v2.0/accounts/"
	    },
	    "Meta": {
	        "TotalPages": 1
	    }
    }    
	`)
	basicTestCase = []byte(`
	{
        "@id": "#t1008",
        "name": "Get a list of accounts",
        "input": {
            "method": "GET",
            "endpoint": "/accounts"
        },
        "context": {
			"baseurl":"http://myaspsp"
		},
        "expect": {
            "status-code": 200,
            "schema-validation": true
        }
    }
	`)

	jsonTestCase = []byte(`
	{
        "@id": "#t1000",
        "name": "Get a list of accounts",
        "input": {
            "method": "GET",
            "endpoint": "/accounts"
        },
        "context": {
			"baseurl":"http://myaspsp"
		},
        "expect": {
            "status-code": 200,
			"schema-validation": true,
			"matches": [{
				"description": "A json match on response body",
				"json": "Data.Account.0.AccountId",
				"value": "500000000000000000000001"
			}]			
        }
    }
	`)
)

// Reads a single testcase from json bytes
func TestReadSingleTestCaseFromJsonBytes(t *testing.T) {
	// testcase in json format
	var testcase TestCase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)
	assert.Equal(t, "#t1008", testcase.ID)
	assert.Equal(t, "GET", testcase.Input.Method)
	assert.Equal(t, "/accounts", testcase.Input.Endpoint)
	assert.Equal(t, true, testcase.Expect.SchemaValidation)
	pkgutils.DumpJSON(testcase)
}

// Simple gock test
func TestRunGock(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	res, err := http.Get("http://foo.com/bar")
	assert.NoError(t, err)
	_ = res
	assert.Equal(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, string(body)[:13], `{"foo":"bar"}`)
	fmt.Printf("Body: %s\n", string(body))

	// Verify that we don't have pending mocks
	assert.Equal(t, gock.IsDone(), true)

}

func TestSimpleMock(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))
	req, _ := http.NewRequest("GET", "http://myaspsp/accounts", nil)
	res, err := (&http.Client{}).Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}

// TestMockedTestCase - creates http request and response objects, sends them to a mocked
// service which uses "gock" as a mock http server.
// Create a test case
// Runs the testcase against the mock server
// checks the resultcode
// Noted - the mocked service sends 'getAccountResponse' as the response body
func TestMockedTestCase(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

	var testcase TestCase // get the testcase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(nil)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	res, err := (&http.Client{}).Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	result, err := testcase.ApplyExpects(res)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Nil(t, err)
	assert.Equal(t, result, true)

	// Verify that we don't have pending mocks
	assert.Equal(t, gock.IsDone(), true)
}

// test a testcase against mock service which supplies incorrect http status code
func TestResponseStatusCodeMismatch(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(201).BodyString(string(getAccountResponse))

	var testcase TestCase // get the testcase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(nil)
	assert.Nil(t, err)

	res, err := (&http.Client{}).Do(req)
	assert.Nil(t, err)

	result, err := testcase.ApplyExpects(res)
	assert.NotNil(t, err)
	fmt.Println(err)
	assert.Equal(t, result, false)

	// Verify that we don't have pending mocks (gock specific)
	assert.Equal(t, gock.IsDone(), true)
}

// Check that a json response field can be checked using Expects
func TestJsonExpectMatch(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

	var testcase TestCase // get the testcase
	err := json.Unmarshal(jsonTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(nil) // need a proper http object here
	assert.Nil(t, err)

	res, err := (&http.Client{}).Do(req)
	assert.Nil(t, err)

	result, err := testcase.ApplyExpects(res)
	assert.Nil(t, err)
	assert.Equal(t, result, true)

	// Verify that we don't have pending mocks (gock specific)
	assert.Equal(t, gock.IsDone(), true)
}

// REFAPP-468 testcase chaining
var (
	chan01 = []byte(`
	{
        "@id": "#t1000",
        "name": "Get a list of accounts",
        "input": {
            "method": "GET",
            "endpoint": "/accounts"
        },
        "context": {
			"baseurl":"http://myaspsp"
		},
        "expect": {
            "status-code": 200,
			"schema-validation": true,
			"matches": [{
				"description": "A json match on response body",
				"json": "Data.Account.0.AccountId",
				"value": "500000000000000000000001"
			}]			
        }
	}
	`)

	testCaseChain01 = []byte(`
{
	"@id": "#t1000",
	"name": "Get a list of accounts",
	"input": {
		"method": "GET",
		"endpoint": "/accounts"
	},
	"context": {
		"baseurl":"http://myaspsp"
	},
	"expect": {
		"status-code": 200,
		"matches": [{
			"description": "A json match on response body",
			"json": "Data.Account.Accountid",
			"value": "XYZ1231231231231"
		}],
		"contextPut": [{
			"matches": [{
				"name": "AccountId",                                
				"description": "A json match to extract vaiable to context",                                
				"json": "Data.Account.Accountid"
			}]
		}]
    }
}
`)

	testCaseChain02 = []byte(`
{
	"@id": "#t0002",
	"name": "Get Accounts using AccountId",
	"purpose": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
	"input": {
		"method": "GET",
		"endpoint": "/accounts/{AccountId}",
		"contextGet": [{
			"name": "AccountId",
			"purpose": "supplies the account number to be queried",
			"replaceInEndpoint": "{AccountId}"
		}]
	},
	"context": {
		"baseurl":"http://myaspsp"
	},
	"expect": {
		"status-code": 200,
		"matches": [{
			"description": "A json match on response body",
			"json": "Data.Account.Accountid",
			"value": "XYZ1231231231231"
		}]
	}
}
`)
)

// Define 2 test cases - testcase two is run after testcase one
func TestSequencingPt1(t *testing.T) {
	var tc01 TestCase
	json.Unmarshal(testCaseChain01, &tc01)
	assert.Equal(t, "/accounts", tc01.Input.Endpoint)
	fmt.Println(string(pkgutils.DumpJSON(tc01)))

	var tc02 TestCase
	json.Unmarshal(testCaseChain02, &tc02)
	assert.Equal(t, "/accounts/{AccountId}", tc02.Input.Endpoint)
	fmt.Println(string(pkgutils.DumpJSON(tc02)))

}

func TestSequencingTestCases(t *testing.T) {
	manifest, err := loadManifest("testdata/passAccountId.json")
	require.NoError(t, err)
	assert.Equal(t, manifest.Name, "Basic Swagger 2.0 test run")
	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

	executor := executor{}
	rule := manifest.Rules[0]
	rule.Executor = &executor
	tc := rule.Tests[0][0]
	ctx := Context{}
	req, _ := tc.Prepare(&ctx)

	res, err := rule.Execute(req, &tc)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%s\n", body)
	result, _ := tc.Validate(res, &ctx)
	assert.True(t, result)
}

// type executor struct {
// }

// func (e *executor) ExecuteTestCase(r *http.Request, t *TestCase, ctx *Context) (*http.Response, error) {
// 	res, err := (&http.Client{}).Do(r)
// 	return res, err
// }

// // Utility to load Manifest Data Model containing all Rules, Tests and Conditions
// func loadManifest(filename string) (Manifest, error) {
// 	plan, _ := ioutil.ReadFile(filename)
// 	var i Manifest
// 	err := json.Unmarshal(plan, &i)
// 	if err != nil {
// 		return i, err
// 	}
// 	return i, nil
// }

// for _, rule := range manifest.Rules { // Iterate over Rules
// 	rule.Executor = &executor
// 	for _, testcases := range rule.Tests {
// 		ctx := Context{}
// 		ctx["{AccountId}"] = "1231231"
// 		for _, testcase := range testcases {
// 			myReq, _ := testcase.Prepare(&ctx)          // Apply inputs, context - results on http object and context
// 			resp, err := rule.Execute(myReq, &testcase) // execute the testcase

// 			testcase.Validate(resp, &ctx)

// 			_ = err
// 		}
// 		_ = ctx
// 	}
//rule.ProcessTestCases() // what does this even mean?
// take the first testcase array,
// is it a single testcase or multiple?
// if single ... process/run
// if multiple ... check that all requirements to run are met
//- like available context variables
//- like if an input variable is specificed - its used
//- like all input variables are available
//- like if an expects variable is specificed its used
//}
