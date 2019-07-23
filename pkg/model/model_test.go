package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"
)

func TestLoadModel(t *testing.T) {
	model, err := loadModel()
	require.NoError(t, err)

	t.Run("model String() returns string representation", func(t *testing.T) {
		expected := "MANIFEST\nName: Basic Swagger 2.0 test run\nDescription: Tests appropriate behaviour of the Open Banking Limited 2.0 Read/Write APIs\nRules: 3\n"
		assert.Equal(t, expected, model.String())
	})

	rule := model.Rules[0]
	t.Run("rule String() returns string representation", func(t *testing.T) {
		expected := "RULE\nName: Get Accounts Basic Rule\nPurpose: Accesses the Accounts endpoint and retrives a list of PSU accounts\nSpecRef: Read Write 2.0 section subsection 1 point 1a\nSpec Location: https://openbanking.org.uk/rw2.0spec/errata#1.1a\nTests: 1\n"
		assert.Equal(t, expected, rule.String())
	})

	t.Run("rule has a Name", func(t *testing.T) {
		assert.Equal(t, rule.Name, "Get Accounts Basic Rule")
	})
}

// Enumerates all OpenAPI calls from swagger file
func TestEnumerateOpenApiTestcases(t *testing.T) {
	doc, err := loadOpenAPI(false)
	require.NoError(t, err)
	base := "https://myaspsp.resourceserver:443/"

	for path, props := range doc.Spec().Paths.Paths {
		for meth := range getOperations(&props) {
			newPath := base + path
			assert.NotNil(t, meth)
			assert.NotNil(t, newPath)
		}
	}
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadManifest(filename string) (Manifest, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return Manifest{}, err
	}
	var i Manifest
	err = json.Unmarshal(plan, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadModel() (Manifest, error) {
	plan, err := ioutil.ReadFile("testdata/testmanifest.json")
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// Utility to load the 2.0 swagger spec for testing purposes
func loadOpenAPI(print bool) (*loads.Document, error) {
	doc, err := loads.Spec("testdata/rwspec2-0.json")
	if err != nil {
		return nil, err
	}
	if print {
		var jsondoc []byte
		jsondoc, err = json.MarshalIndent(doc.Spec(), "", "    ")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(jsondoc))
	}
	return doc, err
}

// Utilities to walk the swagger tree
// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOperations(props *spec.PathItem) map[string]*spec.Operation {
	ops := map[string]*spec.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}

	// Keep those != nil
	for key, op := range ops {
		if op == nil {
			delete(ops, key)
		}
	}
	return ops
}

// Use cases

/*
As a developer I want to perform a test where I load some json which defines a manifest, rule and testcases
I want the rule to manage the execution of the test that includes two test cases
I want the testcases to communicate parameters between themselves using a context
I want the results of one test case being used as input to the other testcase
I want to use json pattern matching to extract the first returned AccountId from the first testcase and
use that value as the accountid parameter for the second testcase
*/
func TestChainedTestCases(t *testing.T) {
	manifest, err := loadManifest("testdata/passAccountId.json")
	require.NoError(t, err)

	rule := manifest.Rules[0] // get the first rule
	executor := &executor{}   // Allows rule testcase execution strategies to be dynamically added to rules
	rulectx := Context{}      // create a context to hold the passed parameters
	tc01 := rule.Tests[0][0]  // get the first testcase of the first rule
	tc01.Validator = schema.NewNullValidator()
	req, err := tc01.Prepare(&rulectx) // Prepare calls ApplyInput and ApplyContext on testcase
	require.NoError(t, err)
	resp, err := executor.ExecuteTestCase(req, &tc01, &rulectx) // send the request to be executed resulting in a response
	assert.NoError(t, err)
	result, errs := tc01.Validate(resp, &rulectx) // Validate checks response against the match rules and processes any contextPuts present
	assert.Nil(t, errs)
	assert.True(t, result)

	tc02 := rule.Tests[0][1] // get the second testcase of the first rule
	tc02.Validator = schema.NewNullValidator()
	req, err = tc02.Prepare(&rulectx) // Prepare
	assert.NoError(t, err)
	resp, err = executor.ExecuteTestCase(req, &tc02, &rulectx) // Execute
	assert.NoError(t, err)
	result, errs = tc02.Validate(resp, &rulectx) // Validate checks the match rules and processes any contextPuts present
	assert.Nil(t, errs)
	assert.True(t, result)
	acctid, exist := rulectx.Get("AccountId")
	assert.True(t, exist)
	assert.Equal(t, "500000000000000000000007", acctid)
	assert.Equal(t, "http://myaspsp/accounts/500000000000000000000007", tc02.Input.Endpoint)
}

type executor struct {
}

func (e *executor) ExecuteTestCase(r *resty.Request, t *TestCase, ctx *Context) (*resty.Response, error) {
	responseKey := t.Input.Method + " " + t.Input.Endpoint
	return chainTest[responseKey](), nil
}

var chainTest = map[string]func() *resty.Response{
	"GET http://myaspsp/accounts/":                         httpAccountCall(),
	"GET http://myaspsp/accounts/{AccountId}":              httpAccountIDCall(),
	"GET http://myaspsp/accounts/500000000000000000000007": httpAccountID007Call(),
}

func httpAccountCall() func() *resty.Response {
	return func() *resty.Response {
		return test.CreateHTTPResponse(200, "OK", string(getAccountResponse))
	}
}

func httpAccountIDCall() func() *resty.Response {
	return func() *resty.Response {
		return test.CreateHTTPResponse(200, "OK", string(getAccountResponse), "content-type", "klingon/text")
	}
}

func httpAccountID007Call() func() *resty.Response {
	return func() *resty.Response {
		return test.CreateHTTPResponse(200, "OK", string(account0007), "content-type", "klingon/text")
	}
}

func TestGetReplacementField(t *testing.T) {
	testCases := []struct {
		stringToCheck string
		value         string
		isFn          bool
		isReplacement bool
		err           error
	}{
		{
			stringToCheck: "$hello",
			value:         "hello",
			isReplacement: true,
			isFn:          false,
			err:           nil,
		},
		{
			stringToCheck: "hello",
			value:         "hello",
			isReplacement: false,
			isFn:          false,
			err:           nil,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			value, isReplacement := getReplacementField(tc.stringToCheck)
			assert.Equal(t, tc.value, value)
			assert.Equal(t, tc.isReplacement, isReplacement)
		})
	}

}
