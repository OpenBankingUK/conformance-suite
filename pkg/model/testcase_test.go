package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"github.com/stretchr/testify/assert"
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

	account0007 = []byte(`
	{
		"Data": {
			"Account": [
				{
					"AccountId": "500000000000000000000007",
					"Currency": "GBP",
					"Nickname": "xxxx0001",
					"AccountType": "Business",
					"AccountSubType": "CurrentAccount",
					"Account": [
						{
							"SchemeName": "IBAN",
							"Identification": "GB29PAPA20000390210099",
							"Name": "Mario Carpentry"
						}
					],
					"Servicer": {
						"SchemeName": "BICFI",
						"Identification": "PAPAUK00XXX"
					}
				}
			]
		},
		"Links": {
			"Self": "http://modelobank2018.o3bank.co.uk/open-banking/v2.0/accounts/500000000000000000000007"
		},
		"Meta": {
			"TotalPages": 1
		}
	}`)

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

	expectOneOfTestCase = []byte(`
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
            "status-code": 0,
            "schema-validation": true
	},
	"expect_one_of": [
	    {
	        "status-code": 400,
	        "schema-validation": true
	    },
	    {
		"status-code": 200,
		"schema-validation": true
	    }	
	]
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

	data, err := json.MarshalIndent(testcase, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(data))
	// Output:
	// {
	// 	"@id": "#t1008",
	// 	"name": "Get a list of accounts",
	//	"input": {
	//	"method": "GET",
	//		"endpoint": "/accounts",
	//		"contextGet": {}
	// },
	//	"context": {
	//	"baseurl": "http://myaspsp"
	// },
	//	"expect": {
	//	"status-code": 200,
	//		"schema-validation": true,
	//		"contextPut": {}
	// }
	// }
}

// TestMockedTestCase - creates http request and response objects, sends them to a mocked
// service which uses "gock" as a mock http server.
// Create a test case
// Runs the testcase against the mock server
// checks the resultcode
// Noted - the mocked service sends 'getAccountResponse' as the response body
func TestMockedTestCase(t *testing.T) {
	var testcase TestCase // get the testcase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(&Context{})
	assert.Nil(t, err)
	assert.NotNil(t, req)

	res := test.CreateHTTPResponse(200, "OK", string(getAccountResponse))
	result, errs := testcase.ApplyExpects(res, nil)
	assert.Nil(t, errs)
	assert.Equal(t, res.StatusCode(), 200)
	assert.Nil(t, err)
	assert.Equal(t, result, true)
}

func TestMockedTestCaseExpectOneOfSucceeds(t *testing.T) {
	var testcase TestCase // get the testcase
	err := json.Unmarshal(expectOneOfTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(&Context{})
	assert.Nil(t, err)
	assert.NotNil(t, req)

	res := test.CreateHTTPResponse(200, "OK", string(getAccountResponse))
	result, errs := testcase.ApplyExpects(res, nil)
	assert.Nil(t, errs)
	assert.Equal(t, res.StatusCode(), 200)
	assert.Nil(t, err)
	assert.Equal(t, result, true)
}

func TestMockedTestCaseExpectOneOfFails(t *testing.T) {
	var testcase TestCase // get the testcase
	err := json.Unmarshal(expectOneOfTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(&Context{})
	assert.Nil(t, err)
	assert.NotNil(t, req)

	res := test.CreateHTTPResponse(404, "OK", string(getAccountResponse))
	result, errs := testcase.ApplyExpects(res, nil)
	assert.NotNil(t, errs)
	assert.Equal(t, res.StatusCode(), 404)
	assert.Nil(t, err)
	assert.Equal(t, result, false)
}

// test a testcase against mock service which supplies incorrect http status code
func TestResponseStatusCodeMismatch(t *testing.T) {

	var testcase TestCase // get the testcase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)

	res := test.CreateHTTPResponse(201, "OK", string(getAccountResponse))

	result, errs := testcase.ApplyExpects(res, nil)
	assert.NotNil(t, errs)
	assert.Equal(t, result, false)

}

// Check that a json response field can be checked using Expects
func TestJsonExpectMatch(t *testing.T) {
	var testcase TestCase // get the testcase
	err := json.Unmarshal(jsonTestCase, &testcase)
	testcase.Validator = schema.NewNullValidator()
	assert.NoError(t, err)

	res := test.CreateHTTPResponse(200, "OK", string(getAccountResponse))

	result, errs := testcase.Validate(res, emptyContext)
	assert.Nil(t, errs)
	assert.Equal(t, result, true)
}

func TestApplyInputNoGetMethod(t *testing.T) {
	tc := TestCase{}
	req, err := tc.Prepare(emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}
