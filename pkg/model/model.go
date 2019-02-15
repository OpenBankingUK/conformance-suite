package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"gopkg.in/resty.v1"
)

// Manifest is the high level container for test suite definition
// It contains a list of all the rules required to be passed for conformance testing
// Each rule can have multiple testcases which contribute to testing that particular rule
// So essentially Manifest is a container
type Manifest struct {
	Context     string    `json:"@context"`         // JSONLD contest reference
	ID          string    `json:"@id"`              // JSONLD ID reference
	Type        string    `json:"@type"`            // JSONLD Type reference
	Name        string    `json:"name"`             // Name of the manifest
	Description string    `json:"description"`      // Description of the Manifest and what it contains
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
	ID         string         `json:"@id,omitempty"`     // JSONLD ID Reference
	Type       []string       `json:"@type,omitempty"`   // JSONLD type array
	Name       string         `json:"name,omitempty"`    // Name
	Purpose    string         `json:"purpose,omitempty"` // Purpose of the testcase in simple words
	Input      Input          `json:"input,omitempty"`   // Input Object
	Context    Context        `json:"context,omitempty"` // Local Context Object
	Expect     Expect         `json:"expect,omitempty"`  // Expected object
	ParentRule *Rule          `json:"-"`                 // Allows accessing parent Rule
	Request    *resty.Request `json:"-"`                 // The request that's been generated in order to call the endpoint
	Header     http.Header    `json:"-"`                 // ResponseHeader
	Body       string         `json:"-"`                 // ResponseBody
	Bearer     string         `json:"bearer,omitempty"`  // Bear token if presented
}

// Prepare a Testcase for execution at and endpoint,
// results in a standard http request that encapsulates the testcase request
// as defined in the test case object with any context inputs/replacements etc applied
func (t *TestCase) Prepare(ctx *Context) (*resty.Request, error) {
	t.AppEntry("Prepare Entry")
	defer t.AppExit("Prepare Exit")

	// Apply Context at end of creating request - get/put values into contexts
	t.ApplyContext(ctx)

	return t.ApplyInput(ctx)
}

// Validate takes the http response that results as a consequence of sending the testcase http
// request to the endpoint implementation. Validate is responsible for checking the http status
// code and running the set of 'Matches' within the 'Expect' object, to determine if all the
// match conditions are met - which would mean the validation passed.
// The context object is passed as part of the validation as its allows the match clauses to
// examine the request object and 'push' response variables into the context object for use
// in downstream test cases which are potentially part of this testcase sequence
// returns true - validation successful
//         false - validation unsuccessful
//         error - adds detail to validation failure
//         NOTE: Validate will only return false if a check fails - no checks = true
func (t *TestCase) Validate(resp *resty.Response, rulectx *Context) (bool, error) {
	if rulectx == nil {
		return false, t.AppErr("error Valdate:rulectx == nil")
	}
	t.Body = resp.String()
	if len(t.Body) == 0 { // The response body can only be read once from the raw response
		// so we cache it in the testcase
		// Check that there is a value body in the raw response of the resty response object
		// Also - if the response is a redirect (302/StatusFound) then we continue as we'll have no body.
		if resp != nil && resp.StatusCode() != http.StatusFound {
			if resp != nil && (resp.RawResponse != nil) && (resp.RawResponse.Body != nil) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(resp.RawResponse.Body)
				if err != nil {
					return false, t.AppErr("Validate: " + err.Error())
				}
				t.Body = buf.String()
			}
		}
	}
	t.Header = resp.Header()
	return t.ApplyExpects(resp, rulectx)
}

// Expect defines a structure for expressing testcase result expectations.
type Expect struct {
	StatusCode       int  `json:"status-code,omitempty"`       // Http response code
	SchemaValidation bool `json:"schema-validation,omitempty"` // Flag to indicate if we need schema validation -
	// provides the ability to switch off schema validation
	Matches    []Match         `json:"matches,omitempty"`    // An array of zero or more match items which must be 'passed' for the testcase to succeed
	ContextPut ContextAccessor `json:"contextPut,omitempty"` // allows storing of test response fragments in context variables
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
func (t *TestCase) ApplyInput(rulectx *Context) (*resty.Request, error) {
	t.AppEntry("ApplyInput entry")
	defer t.AppExit("ApplyInput exit")

	if t.Input.Method == "" {
		return nil, t.AppErr("error: TestCase input cannot have empty input.Method")
	}
	req, err := t.Input.CreateRequest(t, rulectx)
	if err != nil {
		return nil, t.AppErr("createRequest: " + err.Error())
	}

	t.Request = req // store the request in the testcase
	return req, err
}

// ApplyContext - at the end of ApplyInputs on the testcase - we have an initial http request object
// ApplyContext, applies context parameters to the http object.
// Context parameter typically involve variables that originated in discovery
// The functionality of ApplyContext will grow significantly over time.
func (t *TestCase) ApplyContext(rulectx *Context) {
	t.AppEntry("ApplyContext entry")
	defer t.AppExit("ApplyContext exit")

	if rulectx != nil {
		for k, v := range t.Context { // put testcase context values into rule context ...
			rulectx.Put(k, v)
		}
	}

	baseURL, err := t.Context.GetString("baseurl")
	if err == ErrNotFound {
		t.AppMsg("no base url - using only input.Endpoint")
		return
	} else if err != nil {
		t.AppMsg("error getting baseUrl from ctx using only default in input.Endpoint")
		return
	}

	// "convention" puts baseurl as prefix to endpoint in testcase"
	t.Input.Endpoint = baseURL + t.Input.Endpoint
}

// ApplyExpects runs the Expects section of the testcase to evaluate if the response from the system under test passes or fails
// The Expects section of a testcase can contain multiple conditions that need to be met to pass a testcase
// When a test fails, ApplyExpects is responsible for reporting back information about the failure, why it occurred, where it occurred etc.
//
// The ApplyExpect section is also responsible for running and contextPut clauses.
// contextPuts are responsible for updated context variables with values selected from the test case response
// contextPuts will only be executed if the ApplyExpects standards match tests pass
// if any of the ApplyExpects match tests fail - ApplyExpects returns false and contextPuts aren't executed
func (t *TestCase) ApplyExpects(res *resty.Response, rulectx *Context) (bool, error) {
	t.AppEntry("ApplyExpects entry")
	defer t.AppExit("ApplyExpects exit")

	if res == nil { // if we've not got a response object to check, always return false
		return false, t.AppErr("nil http.Response - cannot process ApplyExpects")
	}

	if t.Expect.StatusCode != res.StatusCode() { // Status codes don't match
		return false, t.AppErr(fmt.Sprintf("(%s):%s: HTTP Status code does not match: expected %d got %d", t.ID, t.Name, t.Expect.StatusCode, res.StatusCode()))
	}

	t.AppMsg(fmt.Sprintf("Status check isReplacement: expected [%d] got [%d]", t.Expect.StatusCode, res.StatusCode()))
	for k, match := range t.Expect.Matches {
		checkResult, got := match.Check(t)
		if checkResult == false {
			return false, t.AppErr(fmt.Sprintf("ApplyExpects Returns False on match %s : %s", match.String(), got.Error()))
		}

		t.Expect.Matches[k].Result = match.Result
		t.AppMsg(fmt.Sprintf("Checked Match: %s: result: %s", match.Description, t.Expect.Matches[k].Result))
	}

	if err := t.Expect.ContextPut.PutValues(t, rulectx); err != nil {
		return false, t.AppErr("ApplyExpects Returns FALSE " + err.Error())
	}

	return true, nil
}

// AppMsg - application level trace
func (t *TestCase) AppMsg(msg string) string {
	tracer.AppMsg("TestCase", msg, "")
	return msg
}

// AppErr - application level trace error msg
func (t *TestCase) AppErr(msg string) error {
	tracer.AppErr("TestCase", msg, "")
	return errors.New(msg)
}

// AppEntry - application level trace error msg
func (t *TestCase) AppEntry(msg string) string {
	tracer.AppEntry("TestCase", msg)
	return msg
}

// AppExit - application level trace error msg
func (t *TestCase) AppExit(msg string) string {
	tracer.AppExit("TestCase", msg)
	return msg
}

func (t *TestCase) String() string {
	bites, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return t.AppErr(fmt.Sprintf("error stringifying TestCase %s %s %s", t.ID, t.Name, err.Error())).Error()
	}
	return string(bites)
}

// Various helpers - main to dump struct contents to console

func (m *Manifest) String() string {
	return fmt.Sprintf("MANIFEST\nName: %s\nDescription: %s\nRules: %d\n", m.Name, m.Description, len(m.Rules))
}

func (r *Rule) String() string {
	return fmt.Sprintf("RULE\nName: %s\nPurpose: %s\nSpecRef: %s\nSpec Location: %s\nTests: %d\n",
		r.Name, r.Purpose, r.Specref, r.Speclocation, len(r.Tests))
}

// replaceContextField
func replaceContextField(source string, ctx *Context) (string, error) {
	field, isReplacement := getReplacementField(source)
	if !isReplacement {
		return source, nil
	}
	if len(field) == 0 {
		return source, errors.New("field not found in context " + field)
	}
	replacement, exist := ctx.Get(field)
	if !exist {
		return source, errors.New("replacement not found in context: " + source)
	}
	contextField, ok := replacement.(string)
	if !ok {
		return source, errors.New("replacement is not of type string: " + source)
	}
	result := strings.Replace(source, "$"+field, contextField, 1)
	return result, nil
}

var singleDollarRegex = regexp.MustCompile(`[^\$]?\$(\w*)`)

// GetReplacementField examines the input string and returns the first character
// sequence beginning with '$' and ending with whitespace. '$$' sequence acts as an escape value
// A zero length string is return if now Replacement Fields are found
// returns a boolean to indicate if the field contains a field beginning with a $
func getReplacementField(value string) (string, bool) {
	isReplacement := isReplacementField(value)
	if !isReplacement {
		return value, false
	}
	result := singleDollarRegex.FindStringSubmatch(value)
	if result == nil {
		return "", false
	}
	return result[len(result)-1], true
}

func isReplacementField(value string) bool {
	index := strings.Index(value, "$")
	return index != -1
}

// ProcessReplacementFields prefixed by '$' in the testcase Input and Context sections
// Call to pre-process custom test cases from discovery model
func (t *TestCase) ProcessReplacementFields(ctx *Context) {

	t.Input.Endpoint, _ = replaceContextField(t.Input.Endpoint, ctx) // errors if field not present in context - which is isReplacement for this function
	t.Input.RequestBody, _ = replaceContextField(t.Input.RequestBody, ctx)

	for k := range t.Input.FormData {
		t.Input.FormData[k], _ = replaceContextField(t.Input.FormData[k], ctx)
	}
	for k := range t.Input.Headers {
		t.Input.Headers[k], _ = replaceContextField(t.Input.Headers[k], ctx)
	}
	for k := range t.Input.Claims {
		t.Input.Claims[k], _ = replaceContextField(t.Input.Claims[k], ctx)
	}
	for k := range t.Context {
		param, ok := t.Context[k].(string)
		if !ok {
			continue
		}
		t.Context[k], _ = replaceContextField(param, ctx)
	}
}
