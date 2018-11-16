package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
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

	t.Run("rule has a RunTests() function", func(t *testing.T) {
		rule.RunTests() // Run Tests for a Rule
	})

	testcase := rule.Tests[0][0]
	t.Run("testcase has Dump() function", func(t *testing.T) {
		testcase.Dump()
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
			fmt.Printf("Register %s %s\n", meth, newPath)
		}
	}
}

// Interate over swagger file and generate all testcases
func TestGenerateSwaggerTestCases(t *testing.T) {
	doc, err := loadOpenAPI(false)
	require.NoError(t, err)
	var testcases []TestCase
	testNo := 1000
	for path, props := range doc.Spec().Paths.Paths {
		for meth, op := range getOperations(&props) {
			testNo++
			successStatus := 0
			for i := range op.OperationProps.Responses.ResponsesProps.StatusCodeResponses {
				if i > 199 && i < 300 {
					successStatus = i
				}
			}
			input := Input{Method: meth, Endpoint: path}
			expect := Expect{StatusCode: successStatus, SchemaValidation: true}
			testcase := TestCase{ID: fmt.Sprintf("#t%4.4d", testNo), Input: input, Expect: expect, Name: op.Description}
			testcases = append(testcases, testcase)
		}
	}
	dumpTestCases(testcases)
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadManifest(filename string) (Manifest, error) {
	plan, _ := ioutil.ReadFile(filename)
	var i Manifest
	err := json.Unmarshal(plan, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadModel() (Manifest, error) {
	plan, _ := ioutil.ReadFile("testdata/testmanifest.json")
	var m Manifest
	err := json.Unmarshal(plan, &m)
	if err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// Utility to load the 2.0 swagger spec for testing purposes
func loadOpenAPI(print bool) (*loads.Document, error) {
	doc, err := loads.Spec("testdata/rwspec2-0.json")
	if print {
		var jsondoc []byte
		jsondoc, _ = json.MarshalIndent(doc.Spec(), "", "    ")
		fmt.Println(string(jsondoc))
	}
	return doc, err
}

// Utility to Dump out an array of test cases in JSON formaT
func dumpTestCases(testcases []TestCase) {
	var model []byte
	model, _ = json.MarshalIndent(testcases, "", "    ")
	//fmt.Println(string(model))
	_ = model

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
I want the testcases to communicate paramaters between themselves using a context
I want the results of one test case being used as input to the other testcase
I want to use json pattern matching to extract the first returned AccountId from the first testcase and
use that value as the accountid parameter for the second testcase
*/

func TestChainedTestCases(t *testing.T) {
	manifest, err := loadManifest("testdata/passAccountId.json")
	require.NoError(t, err)
	assert.Equal(t, manifest.Name, "Basic Swagger 2.0 test run")

	defer gock.Off()
	gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

	for _, rule := range manifest.Rules { // Iterate over Rules
		rule.Executor = &executor{}
		for _, testcases := range rule.Tests {
			ctx := NewContext()
			ctx.Put("{AccountId}", "1231231")
			for _, testcase := range testcases {
				fmt.Println("\n==============Dumping testcase =-->")
				testcase.Dump()
				ctx.Dump()

				myReq, _ := testcase.Prepare(ctx, nil)      // Apply inputs, context - results on http object and context
				resp, err := rule.Execute(myReq, &testcase) // execute the testcase
				testcase.Validate(resp, ctx)

				_, _, _ = resp, err, myReq
			}
			_ = ctx
		}
		//rule.ProcessTestCases() // what does this even mean?
		// take the first testcase array,
		// is it a single testcase or multiple?
		// if single ... process/run
		// if multiple ... check that all requirements to run are met
		//- like available context variables
		//- like if an input variable is specificed - its used
		//- like all input variables are available
		//- like if an expects variable is specificed its used
	}
	fmt.Println("\n==============END =-->")
}

type executor struct {
}

func (e *executor) ExecuteTestCase(r *http.Request, t *TestCase, ctx *Context) (*http.Response, error) {
	return nil, nil
}

/*
As a developer I'd like many types of matches in my toolkit
- some matches simple match the response for success for failure
- some matches match a value and put it in the context - for access by other testcases



As a developer I'd like my testcase which uses accountid between testcases to be run against ozone bank

    			  "expect": {
                        "status-code": 200,
                        "matches": [{
                            "description": "A json match on response body",
                            "json": "Data.Account.Accountid",
                            "value": "@AccountId"   // store result in context - what if its not there? Match fails!!! with appropriate message
                        }]
                    }
                    "input": {
                        "method": "GET",
                        "endpoint": "/accounts/{AccountId}",
                        "contextGet": [{ // get a variable from the context and put it somewhere in the request
                            "name": "AccountId", // index into context
                            "purpose": "supplies the account number to be queried",
							"replaceInEndpoint": "{AccountId}" // replace strategy, also replace with jsonexpression - SJSON friend
							"replaceInBody": "{AccountId}" // replace strategy, also replace with jsonexpression - SJSON friend - after body constructed!!! (not get)
							"replaceInHeader/put header?"
							"replaceWithJson":"Data.OBresponse1.Field[3]" // build a response object and but this value in there
							"rawJsonBody":"{}" // raw json to but in the request body - simply use this as the body rather than building
                        }]
                    },
*/

// Context Tests
//
// Put things into a context
// Put and get strings from a context
// Put and get numbers from a context
// Put ang get structures from a context (deep vs shallow)
// Read context from a configuration file
//

/*

Test cases to handle looping through paginated output ... how?

Expects {
	pagecount indicator ... 1 of 6
	// current page
	// total pages
	// logic while current page < total page and not error
	// must be able to report currently page and total pages in any errors
	// pagereader:Testcase
	// firstpage:Testcase
	// loopconstruct:Rule? ... so testcases dont know and run unmodified
	// output feeds to input of same testcase until condition met
	// status - match, fail, loop

	Or ... changing parameter in context referred to by input section
	update context parameter in expects section
	repeat keyword in section
	repeatUntil {
		maxtimes: 10
		Match[] - condition
	}

}


///////////// TESTCASE GENERATION
// iterate over swagger
Get Endpoint/permission combinations
.. so figure out all positive permission permutations
end up with a list of structs that contains a permission + and method/endpoint tied to a testcase

for each permissioned endpoint, identify which permission set it can be satisfied from.
Figure out a set of rules for declarative testcase permission annotations + specifiy default if makes it clearer


OpenAPI - what does it do, what does is provide, is it upto date and reliable, can it be used to enrich our
our core model?

*/
