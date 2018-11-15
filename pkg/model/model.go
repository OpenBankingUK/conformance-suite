package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

// Manifest is the high level container for test suite definition
// It contains a list of all the rules required to be passed for conformance testing
// Each rule can have multiple testcases which contribute to testing that particular rule
// So essentially Manifest is a container
type Manifest struct {
	Context     string    `json:"@context"`         // JSONLD contest reference
	ID          string    `json:"@id"`              // JSONLD ID reference
	Type        string    `json:"@type"`            // JSONLD Type reference
	Name        string    `json:"name"`             // Name of the manifiest
	Description string    `json:"description"`      // Description of the Mainfest and what it contains
	BaseIri     string    `json:"baseIri"`          // Base Iri
	Sections    []Context `json:"section_contexts"` // Section specific contexts
	Rules       []Rule    `json:"rules"`            // All the rules in the Manifest
}

// Rule - Define a specific location within a specification that is being tested
// Rule also identifies all the tests that must be passed in order to show that the rule
// implementation in conformant with the specific section in the referenced specification
type Rule struct {
	ID           string       `json:"@id"`             // JSONLD ID reference
	Type         []string     `json:"@type,omitempty"` // JSONLD type reference
	Name         string       `json:"name"`            // A short meaningful name for this rule
	Purpose      string       `json:"purpose"`         // The purpose of this rule
	Specref      string       `json:"specref"`         // Description of area of spec/name/version/section under test
	Speclocation string       `json:"speclocation"`    // specific http reference to location in spec under test covered by this rule
	Tests        [][]TestCase `json:"tests"`           // Tests - allows for many testcases - array of arrays - to be associated with this rule
}

// TestCase defines a test that will be run and needs to be passed as part of the conformance suite
// in order to determine implementation conformance to a specification.
// Testcase have three major sections
// Input:
//     Defines the inputs that are required by the testcase. This effectively involves preparing the http request object
// Context:
//     Provides a link between Discovery information and the testcase
// Expects:
//     Examines the http response to the testcase Input in order to determine if the expected conditions existing in the response
//     and therefore the testcase has passed
//
type TestCase struct {
	ID         string        `json:"@id,omitempty"`     // JSONLD ID Reference
	Type       []string      `json:"@type,omitempty"`   // JSONLD type array
	Name       string        `json:"name,omitempty"`    // Name
	Purpose    string        `json:"purpose,omitempty"` // Purpose of the testcase in simple words
	Input      Input         `json:"input,omitempty"`   // Input Object
	Context    Context       `json:"context,omitempty"` // Context Object
	Expect     Expect        `json:"expect,omitempty"`  // Expected object
	ParentRule *Rule         // Allows accessing parent Rule
	Request    *http.Request // The request that's been generated in order to call the endpoint
}

// Input defines the content of the http request object used to execute the test case
// Input is built up typically from the openapi/swagger definition of the method/endpoint for a particualar
// specification. Additional properties/fields/headers can be added or change in order to setup the http
// request object of the specific test case. Once setup correctly,the testcase gives the http request object
// to the parent Rule which determine how to execute the requestion object. On execution an http response object
// is received and passed back to the testcase for validation using the Expects object.
type Input struct {
	Method   string `json:"method"`   // http Method that this test case uses
	Endpoint string `json:"endpoint"` // resource endpoint where the http object needs to be sent to get a response
}

// Context is intended to handle two types of object and make them available to various parts of the suite including
// testcases. The first set are objects created as a result of the discovery phase, which capture discovery model
// information like endpoints and conditional implementation indicators. The other set of data is information passed
// between a sequeuence of test cases, for example AccountId - extracted from the output of one testcase (/Accounts) and fed in
// as part of the input of another testcase for example (/Accounts/{AccountId}/transactions}
type Context struct {
	ID      string `json:"@id,omitempty"`
	BaseURL string `json:"baseurl,omitempty"`
}

// Expect defines a structure for expressing testcase result expectations.
type Expect struct {
	StatusCode       int  `json:"status-code,omitempty"`       // Http response code
	SchemaValidation bool `json:"schema-validation,omitempty"` // Flag to indicate if we need schema validation -
	// provides the ability to switch off schema validation
	Matches []Match `json:"matches,omitempty"` // An array of zero or more match items which must be 'passed' for the testcase to succeed
}

// Match encapsulates a conditional statement that must 'match' in order to succeed.
// Matches should -
// - match using a specific JSON field and a value
// - match using a Regex expression
// - match a specific header field to a value
// - match using a Regex expression on a header field
type Match struct {
	Description string `json:"description,omitempty"`
	JSON        string `json:"json,omitempty"`
	Value       string `json:"value,omitempty"`
	Regex       string `json:"regex,omitempty"`
	Header      string `json:"header,omitempty"`
}

// ApplyInput - creates an HTTP request for this test case
// The reason why we're doing this is that a testcase behaves like an http object
// It produces an http.Request - which can be sent to a server
// It consumes and http.Response - which it uses to validate the response against "Expects"
// TestCase lifecycle:
//     Create a Testcase Object
//     Create / retrieve the http request object
//     Apply context information to the request object
//     Rule - manages passing the request object from the testcase to an appropriate endpoint handler (like the proxy)
//     Rule - receives http response from endpoint and provides it back to testcase
//     Testcase evaluates the http response object using its 'Expects' clause
//     Testcase passes or fails depending on the 'Expects' outcome
func (t *TestCase) ApplyInput() (*http.Request, error) {
	// NOTE: This is an initial implementation to get things moving - expect a lot of change in this function
	var err error
	if (t.Input == Input{}) {
		return nil, errors.New("Testcase Input empty")
	}
	if t.Input.Method != "GET" { // Only get Supported Initially
		return nil, errors.New("Testcase Method Only support GET currently")
	}
	req, err := http.NewRequest(t.Input.Method, t.Input.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	t.Request = req // store the request in the testcase

	req, err = t.ApplyContext() // Apply Context at end of creating request
	return req, err
}

// ApplyContext - at the end of ApplyInputs on the testcase - we have an initial http request object
// ApplyContext, applys context parameters to the http object.
// Context parameter typically involve variables that originaled in discovery
// The functionality of ApplyContext will grow significantly over time.
func (t *TestCase) ApplyContext() (*http.Request, error) {
	req := t.Request
	if (t.Context != Context{}) && (t.Context.BaseURL != "") {
		u, err := url.Parse(t.Context.BaseURL + t.Input.Endpoint) // expand url in request to be full pathname including Discovery endpoint info from context
		if err != nil {
			return nil, errors.New("Error parsing context baseURL: (" + t.Context.BaseURL + ")")
		}
		req.URL = u
	}
	return req, nil
}

// ApplyExpects runs the Expects section of the testcase to evaluate if the response from the system under test passes or fails
// The Expects section of a testcase can contain multiple conditions that need to be met to pass a testcase
// When a test failes, it the ApplyExpects that is responsible for reporting back information about the failure, why it occured, where it occured etc.
func (t *TestCase) ApplyExpects(res *http.Response) (bool, error) {
	if t.Expect.StatusCode != res.StatusCode { // Status codes don't match
		return false, fmt.Errorf("(%s):%s: HTTP Status code does not match: expected %d got %d", t.ID, t.Name, t.Expect.StatusCode, res.StatusCode)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close() // standard tidying
	body := string(bodyBytes)
	for _, match := range t.Expect.Matches {
		jsonMatch := match.JSON // check if there is a JSON match to be satisifed
		if len(jsonMatch) > 0 {
			matched := gjson.Get(body, jsonMatch)
			if matched.String() != match.Value { // check the value of the JSON match - is equal to the 'value' parameter of the match section within the testcase 'Expects' area
				return false, fmt.Errorf("(%s):%s: Json Match: expected %s got %s", t.ID, t.Name, match.Value, matched.String())
			}
		}
	}
	return true, nil
}

// Various helpers - main to dump struct contents to console

// Dump Manifest helper
func (m *Manifest) Dump() {
	fmt.Printf(m.String())
}

func (m *Manifest) String() string {
	return fmt.Sprintf("MANIFEST\nName: %s\nDescription: %s\nRules: %d\n", m.Name, m.Description, len(m.Rules))
}

// Dump Rule helper
func (r *Rule) Dump() {
	fmt.Printf(r.String())
}

func (r *Rule) String() string {
	return fmt.Sprintf("RULE\nName: %s\nPurpose: %s\nSpecRef: %s\nSpec Location: %s\nTests: %d\n",
		r.Name, r.Purpose, r.Specref, r.Speclocation, len(r.Tests))
}

// Dump - TestCase helper
func (t *TestCase) Dump() {
	fmt.Printf("TESTCASE\nID: %s\nName: %s\nPurpose: %s\n", t.ID, t.Name, t.Purpose)
	fmt.Printf("Input: ")
	pkgutils.DumpJSON(t.Input)
	fmt.Printf("Context: ")
	pkgutils.DumpJSON(t.Context)
	fmt.Printf("Expect: ")
	pkgutils.DumpJSON(t.Expect)
}

// RunTests - runs all the tests for aTestRule
func (r *Rule) RunTests() {
	for _, testSequence := range r.Tests {
		for _, tester := range testSequence {
			// testcase.ApplyInput
			// testcase.ApplyContext
			// testcase.ApplyExpects
			_ = tester // placeholder
			fmt.Println("Test Result: ", true)
		}
	}
}
