package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
	"github.com/OpenBankingUK/conformance-suite/pkg/test"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"
)

var (
	// standingOrder JSON queries
	numberOfPaymentsJSON = "Data.StandingOrder.#.NumberOfPayments"
	finalPaymentDateTime = "Data.StandingOrder.#.FinalPaymentDateTime"
	finalPaymentAmount   = "Data.StandingOrder.#.FinalPaymentAmount"

	matchNumberOfPayments   = Match{Description: "test", JSON: numberOfPaymentsJSON}
	matchNoNumberOfPayments = Match{Description: "test", JSONNotPresent: numberOfPaymentsJSON}

	matchFinalPaymentDateTime   = Match{Description: "test", JSON: finalPaymentDateTime}
	matchNoFinalPaymentDateTime = Match{Description: "test", JSONNotPresent: finalPaymentDateTime}

	matchFinalPaymentAmount = Match{Description: "test", JSON: finalPaymentAmount}

	tc1 = TestCase{ExpectLastIfAll: []Expect{
		{Matches: []Match{matchNumberOfPayments}},
		{Matches: []Match{matchFinalPaymentDateTime}},
		{StatusCode: 403},
	},
		Validator:          schema.NewNullValidator(),
		ExpectArrayResults: true,
	}

	tc2 = TestCase{ExpectLastIfAll: []Expect{
		{Matches: []Match{matchNoNumberOfPayments}},
		{Matches: []Match{matchNoFinalPaymentDateTime}},
		{Matches: []Match{matchFinalPaymentAmount}},
		{StatusCode: 403},
	},
		Validator:          schema.NewNullValidator(),
		ExpectArrayResults: true,
	}
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
// TestChainedTestCases - tests passing test scenarios which do match the responses.
func TestChainedTestCases(t *testing.T) {
	manifest, err := loadManifest("testdata/passAccountId.json")
	require.NoError(t, err)

	expectedEndpoints := []string{
		"http://myaspsp/accounts/",
		"http://myaspsp/accounts/500000000000000000000007",
	}

	executor := &executor{} // Allows rule testcase execution strategies to be dynamically added to rules
	rulectx := Context{}    // create a context to hold the passed parameters

	for i, tc := range manifest.Rules[0].Tests[0] { // get testcases of the first rule
		tc.Validator = schema.NewNullValidator()
		req, err := tc.Prepare(&rulectx) // Prepare calls ApplyInput and ApplyContext on testcase
		require.NoError(t, err)
		assert.Equal(t, expectedEndpoints[i], tc.Input.Endpoint)
		resp, err := executor.ExecuteTestCase(req, &tc, &rulectx) // send the request to be executed resulting in a response
		assert.NoError(t, err)
		result, errs := tc.Validate(resp, &rulectx) // Validate checks response against the match rules and processes any contextPuts present
		assert.Nil(t, errs)
		assert.True(t, result)
	}

	acctid, exist := rulectx.Get("AccountId")
	assert.True(t, exist)
	assert.Equal(t, "500000000000000000000007", acctid)
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

func TestPassingExpectsLastIfAll(t *testing.T) {
	for i, body := range okTestCasesStandingOrders {
		t.Run(fmt.Sprintf("Standing Orders - OK test case %d", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(200, "OK", body)
			result, errs := tc1.Validate(resp, emptyContext)
			assert.True(t, result)
			assert.Empty(t, errs)
		})
	}

	for i, body := range badTestCasesStandingOrders {
		t.Run(fmt.Sprintf("Standing Orders - Bad test case %d with the expected status code", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(403, "Forbidden", body)
			result, errs := tc1.Validate(resp, emptyContext)
			assert.True(t, result)
			assert.Empty(t, errs)
		})
	}

	for i, body := range okTestCasesStandingOrders2 {
		t.Run(fmt.Sprintf("Standing Orders - OK test case %d", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(200, "OK", body)
			result, errs := tc2.Validate(resp, emptyContext)
			assert.True(t, result)
			assert.Empty(t, errs)
		})
	}

	for i, body := range badTestCasesStandingOrders2 {
		t.Run(fmt.Sprintf("Standing Orders - Bad test case %d with the expected status code", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(403, "Forbidden", body)
			result, errs := tc2.Validate(resp, emptyContext)
			assert.True(t, result)
			assert.Empty(t, errs)
		})
	}
}

func TestFailingExpectsLastIfAll(t *testing.T) {
	for i, body := range badTestCasesStandingOrders {
		t.Run(fmt.Sprintf("Standing Orders - Bad test case %d with the unexpected status code", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(200, "OK", body)
			result, errs := tc1.Validate(resp, emptyContext)
			assert.False(t, result)
			assert.NotEmpty(t, errs)
		})
	}

	for i, body := range badTestCasesStandingOrders2 {
		t.Run(fmt.Sprintf("Standing Orders - Bad test case %d with the unexpected status code", i), func(t *testing.T) {
			resp := test.CreateHTTPResponse(200, "OK", body)
			result, errs := tc2.Validate(resp, emptyContext)
			assert.False(t, result)
			assert.NotEmpty(t, errs)
		})
	}
}
