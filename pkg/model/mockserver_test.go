package model

/* Tests



 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	// testcase in json format
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
)

// Reads a single testcase from json bytes
func TestReadSingleTestCaseFromJsonBytes(t *testing.T) {
	var testcase TestCase
	err := json.Unmarshal(basicTestCase, &testcase)
	assert.NoError(t, err)
	assert.Equal(t, "#t1008", testcase.ID)
	assert.Equal(t, "GET", testcase.Input.Method)
	assert.Equal(t, "/accounts", testcase.Input.Endpoint)
	assert.Equal(t, true, testcase.Expect.SchemaValidation)
	pkgutils.DumpJSON(testcase)
}

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

	req, err := testcase.Prepare(nil, nil)
	assert.Nil(t, err)

	res, err := (&http.Client{}).Do(req)
	assert.Nil(t, err)

	result, err := testcase.ApplyExpects(res)
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

	req, err := testcase.Prepare(nil, nil)
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
	jsonTestCase := []byte(`
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

	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

	var testcase TestCase // get the testcase
	err := json.Unmarshal(jsonTestCase, &testcase)
	assert.NoError(t, err)

	req, err := testcase.Prepare(nil, nil) // need a proper http object here
	assert.Nil(t, err)

	res, err := (&http.Client{}).Do(req)
	assert.Nil(t, err)

	result, err := testcase.ApplyExpects(res)
	assert.Nil(t, err)
	assert.Equal(t, result, true)

	// Verify that we don't have pending mocks (gock specific)
	assert.Equal(t, gock.IsDone(), true)

}
