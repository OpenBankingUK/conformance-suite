package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/OpenBankingUK/conformance-suite/pkg/authentication"
	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
	"github.com/OpenBankingUK/conformance-suite/pkg/schemaprops"

	"github.com/sirupsen/logrus"

	"github.com/OpenBankingUK/conformance-suite/pkg/tracer"
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
	ID                  string           `json:"@id,omitempty"`                  // JSONLD ID Reference
	Type                []string         `json:"@type,omitempty"`                // JSONLD type array
	Name                string           `json:"name,omitempty"`                 // Name
	Detail              string           `json:"detail,omitempty"`               // Detailed description of the test case
	RefURI              string           `json:"refURI,omitempty"`               // Reference URI for the test case
	Purpose             string           `json:"purpose,omitempty"`              // Purpose of the testcase in simple words
	Input               Input            `json:"input,omitempty"`                // Input Object
	Context             Context          `json:"context,omitempty"`              // Local Context Object
	Expect              Expect           `json:"expect,omitempty"`               // Expected object
	ExpectOneOf         []Expect         `json:"expect_one_of,omitempty"`        // Slice of possible expected objects
	ExpectLastIfAll     []Expect         `json:"expect_last_if_all,omitempty"`   // Slice of expected objects if all before last one passed the last one needs too
	ParentRule          *Rule            `json:"-"`                              // Allows accessing parent Rule
	Request             *resty.Request   `json:"-"`                              // The request that's been generated in order to call the endpoint
	Header              http.Header      `json:"-"`                              // ResponseHeader
	Body                string           `json:"-"`                              // ResponseBody
	Bearer              string           `json:"bearer,omitempty"`               // Bear token if presented
	DoNotCallEndpoint   bool             `json:"do_not_call_endpoint,omitempty"` // If we should not call the endpoint, see `components/PSUConsentProviderComponent.json`
	ExpectArrayResults  bool             `json:"expect_array_results,omitempty"` // Compare response body lengths between each expect (currently used by ExpectLastIfAll)
	APIName             string           `json:"apiName"`
	APIVersion          string           `json:"apiVersion"`
	Validator           schema.Validator `json:"-"` // Swagger schema validator
	ValidateSignature   bool             `json:"validateSignature,omitempty"`
	StatusCode          string           `json:"statusCode,omitempty"`
	ResultArray         []string         `json:"-"` // represents Result array
	ResultPresenceArray []bool           `json:"-"` // represents Result bool array with information if the fields were found based on JSON query in the Results
}

// MakeTestCase builds an empty testcase
func MakeTestCase() TestCase {
	i := Input{}
	i.FormData = make(map[string]string)
	i.QueryParameters = make(map[string]string)
	i.Generation = make(map[string]string)
	i.Headers = make(map[string]string)
	i.RemoveClaims = []string{}
	i.Claims = make(map[string]string)

	tc := TestCase{Input: i, Validator: schema.NewNullValidator()}
	return tc
}

// Prepare a Testcase for execution at and endpoint,
// results in a standard http request that encapsulates the testcase request
// as defined in the test case object with any context inputs/replacements etc applied
func (t *TestCase) Prepare(ctx *Context) (*resty.Request, error) {
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
func (t *TestCase) Validate(resp *resty.Response, ctx *Context) (bool, []error) {
	if ctx == nil {
		return false, []error{t.AppErr("error Valdate:rulectx == nil")}
	}
	t.Body = resp.String()
	if len(t.Body) == 0 {
		logrus.WithField("testcase", t.String()).Debug("Validate: resty.body is empty")
	}
	t.Header = resp.Header()
	pass, errs := t.ApplyExpects(resp, ctx)

	var failures []schema.Failure
	if t.Expect.SchemaValidation {
		if t.Validator == nil {
			return false, []error{t.AppErr("Validate: schema validator is nil")}
		}

		var err error
		fmt.Printf("\n\n validator: %+v \n\n", t.Validator)
		failures, err = t.Validator.Validate(schema.HTTPResponse{
			Method:     t.Input.Method,
			Path:       t.Input.Endpoint,
			Header:     resp.Header(),
			Body:       strings.NewReader(t.Body),
			StatusCode: resp.StatusCode(),
		})
		if err != nil {
			return false, []error{t.AppErr("Validate: " + err.Error())}
		}
		for _, failure := range failures {
			errs = append(errs, errors.New(failure.Message))
		}
	} else {
		logSchemaValidationOffWarning(t)
	}

	// Apply Signature Validator
	if t.ValidateSignature && !disableJws {
		xJwsSignature := resp.Header().Get("x-jws-signature")
		logrus.Warn("Validating Signature: " + xJwsSignature)
		logrus.Warn("body: ", t.Body)
		valid, err := validateSignature(xJwsSignature, t.Body, ctx)
		if err != nil {
			return false, []error{t.AppErr("Signature validation failed: " + err.Error())}
		}
		if !valid {
			errs = append(errs, errors.New("Invalid x-jws-signature found - unable to validate"))
		} else {
			logrus.Infoln("x-jws-signature validation succeded")
		}
	}

	// Gather fields within json response - for reporting
	collector := schemaprops.GetPropertyCollector()
	collector.CollectProperties(t.Input.Method, t.Input.Endpoint, t.Body, resp.StatusCode())

	return pass, errs
}

func validateSignature(signature, body string, ctx *Context) (bool, error) {
	var pass bool
	if signature != "" {
		jwksURI, err := ctx.GetString("jwks_uri")
		if err != nil {
			return false, errors.New("ValidateSignature - JWKS_URI not present ")
		}

		b64encoding, err := authentication.GetB64Encoding(ctx)
		if err != nil {
			return false, errors.New("ValidationSignature cannot get B64Encoding: " + err.Error())
		}

		pass, err = authentication.ValidateSignature(signature, body, jwksURI, b64encoding)
		if err != nil {
			return false, errors.New("Invalid x-jws-signature found - unable to validate: " + err.Error())
		}
		if !pass {
			return false, errors.New("Invalid x-jws-signature - fails validation")
		}
		logrus.Infoln("x-jws-signature validation succeded")
	} else {
		return false, errors.New("x-jws-signature header not found for Validation")
	}
	logrus.Tracef("Signature validation succeeded")
	return pass, nil
}

func logSchemaValidationOffWarning(testCase *TestCase) {
	// Don't log warning if it is one of the TestCases in the ignore list.
	type IgnoreTestCase struct {
		ID   string
		Name string
	}
	ignoredTestCases := []IgnoreTestCase{
		{
			ID:   "#compPsuConsent01",
			Name: "ClientCredential Grant",
		},
		{
			ID:   "#ct0001",
			Name: "ClientCredential Grant",
		},
		{
			ID:   "#ccg0001",
			Name: "ClientCredential Grant",
		},
		{
			ID:   "#compPsuConsent03",
			Name: "Ozone Headless Consent Flow",
		},
		{
			ID:   "#compPsuConsent03",
			Name: "PSU Consent Token Exchange",
		},
		{
			ID:   "#ct0003",
			Name: "Ozone Headless Consent Flow",
		},
		{
			ID:   "#ct0004",
			Name: "Code Exchange",
		},
	}

	// Check if TestCase is in the ignored list.
	found := false
	for _, ignoredTestCase := range ignoredTestCases {
		if ignoredTestCase.ID == testCase.ID && ignoredTestCase.Name == testCase.Name {
			found = true
			break
		}
	}

	// Only log warning if it is not in the ignore list.
	if !found {
		logrus.WithFields(logrus.Fields{
			"module":   "TestCase",
			"function": "Validate",
			"package":  "model",
			"TestCase": testCase.ID,
		}).Warn(`TestCase.Expect.SchemaValidation is false`)
	}
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
// It consumes an http.Response - which it uses to validate the response against "Expects"
// TestCase lifecycle:
//     Create a Testcase Object
//     Create / retrieve the http request object
//     Apply context information to the request object
//     Rule - manages passing the request object from the testcase to an appropriate endpoint handler (like the proxy)
//     Rule - receives http response from endpoint and provides it back to testcase
//     Testcase evaluates the http response object using its 'Expects' clause
//     Testcase passes or fails depending on the 'Expects' outcome
func (t *TestCase) ApplyInput(rulectx *Context) (*resty.Request, error) {
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
	if !strings.HasPrefix(t.Input.Endpoint, baseURL) {
		t.Input.Endpoint = baseURL + t.Input.Endpoint
	}
}

// ApplyExpects runs the Expects section of the testcase to evaluate if the response from the system under test passes or fails
// The Expects section of a testcase can contain multiple conditions that need to be met to pass a testcase
// When a test fails, ApplyExpects is responsible for reporting back information about the failure, why it occurred, where it occurred etc.
//
// The ApplyExpect section is also responsible for running and contextPut clauses.
// contextPuts are responsible for updated context variables with values selected from the test case response
// contextPuts will only be executed if the ApplyExpects standards match tests pass
// if any of the ApplyExpects match tests fail - ApplyExpects returns false and contextPuts aren't executed
func (t *TestCase) ApplyExpects(res *resty.Response, rulectx *Context) (bool, []error) {
	if res == nil { // if we've not got a response object to check, always return false
		return false, []error{t.AppErr("nil http.Response - cannot process ApplyExpects")}
	}
	ok, err := t.validateExpect(t.Expect, res)
	if !ok {
		return ok, []error{err}
	}

	if result, failedExpects := t.validateExpectsOneOf(res); !result && failedExpects != nil {
		return result, failedExpects
	}

	if result, failedExpects := t.validateExpectsLastIfAll(res); !result && failedExpects != nil {
		return result, failedExpects
	}

	if err := t.Expect.ContextPut.PutValues(t, rulectx); err != nil {
		return false, []error{t.AppErr("ApplyExpects Returns FALSE " + err.Error())}
	}
	return true, nil
}

func (t *TestCase) validateExpect(expect Expect, res *resty.Response) (bool, error) {
	// Status code `-1` is specified in test cases if we want to ignore the HTTP status code.
	if expect.StatusCode > 0 && expect.StatusCode != res.StatusCode() {
		return false, t.AppErr(fmt.Sprintf("(%s):%s: HTTP Status code does not match: expected %d got %d", t.ID, t.Name, expect.StatusCode, res.StatusCode()))
	}

	t.AppMsg(fmt.Sprintf("Status check isReplacement: expected [%d] got [%d]", expect.StatusCode, res.StatusCode()))
	for k, match := range expect.Matches {
		match.ExpectResults = t.ExpectArrayResults
		checkResult, got := match.Check(t)
		if !checkResult {
			return false, t.AppErr(fmt.Sprintf("ApplyExpects Returns False on match %s : %s", match.String(), got.Error()))
		}

		expect.Matches[k].Result = match.Result
		t.ResultPresenceArray = match.ResultPresenceArray
		t.ResultArray = match.ResultArray

		t.AppMsg(fmt.Sprintf("Checked Match: %s: result: %s", match.Description, expect.Matches[k].Result))
	}

	return true, nil
}

// validateExpectsOneOf - validates the slice of expects. It is OK when at least one of has passed.
func (t *TestCase) validateExpectsOneOf(res *resty.Response) (bool, []error) {
	failedExpects := make([]error, 0, len(t.ExpectOneOf))
	for _, expect := range t.ExpectOneOf {
		ok, err := t.validateExpect(expect, res)
		if !ok {
			failedExpects = append(failedExpects, err)
		}
	}

	// t.ExpectOneOf represents an optional list of []Expect one of which must be met
	// since the usage of t.ExpectOneOf is optional, t.ExpectOneOf can be empty
	// in this case the validation is skipped.
	// When t.ExpectOneOf is not empty, at least one of the Expect must be successful
	if len(t.ExpectOneOf) > 0 && len(failedExpects) == len(t.ExpectOneOf) {
		return false, failedExpects
	}

	return true, nil
}

// validateExpectsLastIfAll - validates the last expect when all beofre have passed.
func (t *TestCase) validateExpectsLastIfAll(res *resty.Response) (bool, []error) {
	if t.ExpectArrayResults {
		return t.validateExpectsLastIfAllArrayResults(res)
	}

	failedExpects := make([]error, 0, len(t.ExpectLastIfAll))

	for i, expect := range t.ExpectLastIfAll {
		ok, err := t.validateExpect(expect, res)

		switch i {
		case len(t.ExpectLastIfAll) - 1:
			if expect.StatusCode <= 0 {
				failedExpects = append(failedExpects, t.AppErr(fmt.Sprintf("(%s):%s: Missing Status code in the last expect", t.ID, t.Name)))
				return false, failedExpects
			}

			if !ok && len(failedExpects) == 0 {
				failedExpects = append(failedExpects, err)
				return false, failedExpects
			}
		default:
			if !ok {
				failedExpects = append(failedExpects, err)
			}
		}
	}

	return true, nil
}

func (t *TestCase) validateExpectsLastIfAllArrayResults(res *resty.Response) (bool, []error) {
	failedExpects := make([]error, 0, len(t.ExpectLastIfAll))

	var previoustMatch Match
	var currentMatch Match
	var conditionChecks []bool

	for i, expect := range t.ExpectLastIfAll {
		//
		if len(expect.Matches) > 1 {
			failedExpects = append(failedExpects, fmt.Errorf("ApplyExpects Returns False on expect %d : %d unsupported amount of matches", i, len(expect.Matches)))
			return false, failedExpects
		}

		ok, err := t.validateExpect(expect, res)

		if len(expect.Matches) == 1 {
			currentMatch = expect.Matches[0]
			currentMatch.ResultPresenceArray = t.ResultPresenceArray
		}

		switch i {
		case 0:
			conditionChecks = currentMatch.ResultPresenceArray
		case len(t.ExpectLastIfAll) - 1:
			// the last expect should include Status Code
			if expect.StatusCode <= 0 {
				failedExpects = append(failedExpects, t.AppErr(fmt.Sprintf("(%s):%s: Missing Status code in the last expect", t.ID, t.Name)))
				return false, failedExpects
			}

			finds := make([]string, 0, len(conditionChecks))
			for i, conditionCheck := range conditionChecks {
				if !conditionCheck {
					finds = append(finds, strconv.Itoa(i))
				}
			}

			if !ok && len(finds) != len(conditionChecks) {
				failedExpects = append(failedExpects, err)
				if len(finds) > 0 {
					failedExpects = append(failedExpects, t.AppErr(fmt.Sprintf("Matches returned False on array elements: %s", strings.Join(finds, ", "))))
				}
				return false, failedExpects
			}
		default:
			if len(previoustMatch.ResultPresenceArray) != len(currentMatch.ResultPresenceArray) {
				return true, nil
			}

			for j, presence := range currentMatch.ResultPresenceArray {
				conditionChecks[j] = presence && previoustMatch.ResultPresenceArray[j]
			}

		}

		previoustMatch = currentMatch
	}

	return true, nil
}

// InjectBearerToken injects a bear token header into the testcase, token can either be the actual bearer token or a parameter starting with '$'
func (t *TestCase) InjectBearerToken(token string) {
	if t.Input.Headers == nil {
		t.Input.Headers = map[string]string{}
	}
	t.Input.Headers["Authorization"] = "Bearer " + token
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

// Clone a testcase
func (t *TestCase) Clone() TestCase {
	tc := TestCase{}

	tc.ID = t.ID
	tc.Type = t.Type
	tc.Name = t.Name
	tc.Purpose = t.Purpose
	tc.Bearer = t.Bearer
	tc.Input = t.Input.Clone()
	tc.Context = Context{}
	tc.Context.PutContext(&t.Context)
	tc.Expect = t.Expect.Clone()

	logrus.Debugf("cloned test -\n before: %#v\nafter : %#v\n ", t, tc)
	return tc
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
	var ignoreErrors bool
	phase, exist := ctx.Get("phase")
	if !exist || phase == "generation" {
		ignoreErrors = true
	}

	field, isReplacement := getReplacementField(source)
	if !isReplacement {
		return source, nil
	}
	if len(field) == 0 {
		if ignoreErrors {
			return source, nil
		}
		return source, errors.New("field not found in context " + field)
	}
	replacement, exist := ctx.Get(field)
	if !exist {
		if ignoreErrors {
			return source, nil
		}
		return source, errors.New("replacement not found in context: " + source)
	}
	contextField, ok := replacement.(string)
	if !ok {
		if ignoreErrors {
			return source, nil
		}
		return source, errors.New("replacement is not of type string: " + source)
	}
	result := strings.Replace(source, "$"+field, contextField, 1)
	return result, nil
}

var singleDollarRegex = regexp.MustCompile(`[^\$]?\$([\w|\-|_]*)`)

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
func (t *TestCase) ProcessReplacementFields(ctx *Context, showReplacementErrors bool) {
	var err error
	logger := logrus.StandardLogger()

	t.Input.Endpoint, err = replaceContextField(t.Input.Endpoint, ctx) // errors if field not present in context - which is isReplacement for this function
	if err != nil {
		t.logReplaceError("Endpoint", err, logger, showReplacementErrors)
	}

	t.Input.RequestBody, err = replaceContextField(t.Input.RequestBody, ctx)
	if err != nil {
		t.logReplaceError("RequestBody", err, logger, showReplacementErrors)
	}

	t.processReplacementFormData(ctx)
	t.processReplacementHeaders(ctx, logger, showReplacementErrors)
	t.processReplacementClaims(ctx)

	// If customer ip value is not set, remove it from headers
	customerIP, err := ctx.GetString("x-fapi-customer-ip-address")
	if err == nil && customerIP == "" {
		delete(t.Input.Headers, "x-fapi-customer-ip-address")
	}

	for k := range t.Context {
		param, ok := t.Context[k].(string)
		if !ok {
			continue
		}
		t.Context[k], err = replaceContextField(param, ctx)
		if err != nil {
			t.logReplaceError("param", err, logger, showReplacementErrors)
		}
	}

	for k, v := range t.Expect.ContextPut.Matches {
		t.Expect.ContextPut.Matches[k].ContextName, err = replaceContextField(v.ContextName, ctx)
		if err != nil {
			t.logReplaceError("ContextName", err, logger, showReplacementErrors)
		}
	}

	for idx, match := range t.Expect.Matches {
		match.ProcessReplacementFields(ctx)
		t.Expect.Matches[idx] = match
	}
}

func (t *TestCase) logReplaceError(field string, err error, logger *logrus.Logger, showReplacementErrors bool) {
	if showReplacementErrors {
		logger.WithError(err).Errorf("processing %s replacement fields", field)
	} else {
		logger.WithError(err).Debugf("processing %s replacement fields", field)
	}
}

func (t *TestCase) processReplacementFormData(ctx *Context) {
	var err error
	for k := range t.Input.FormData {
		t.Input.FormData[k], err = replaceContextField(t.Input.FormData[k], ctx)
		if err != nil {
			logrus.StandardLogger().WithError(err).Error("processing replacement fields")
		}
	}
}

func (t *TestCase) processReplacementHeaders(ctx *Context, logger *logrus.Logger, showReplacementErrors bool) {
	var err error
	for k := range t.Input.Headers {
		t.Input.Headers[k], err = replaceContextField(t.Input.Headers[k], ctx)
		if err != nil {
			field := fmt.Sprintf("Headers[%s]", k)
			t.logReplaceError(field, err, logger, showReplacementErrors)
		}
	}
}

func (t *TestCase) processReplacementClaims(ctx *Context) {
	var err error
	for k := range t.Input.Claims {
		t.Input.Claims[k], err = replaceContextField(t.Input.Claims[k], ctx)
		if err != nil {
			logrus.StandardLogger().WithError(err).Error("processing replacement fields")
		}
	}
}

// Clone - preforms deep copy of expect object
func (e *Expect) Clone() Expect {
	ex := Expect{}
	ex.StatusCode = e.StatusCode
	ex.SchemaValidation = e.SchemaValidation
	for _, match := range e.Matches {
		m := match.Clone()
		ex.Matches = append(ex.Matches, m)
	}
	return ex
}

// LoadTestCaseFromJSONFile a single testcase from a json file
func LoadTestCaseFromJSONFile(filename string) (TestCase, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return TestCase{}, err
	}
	var m TestCase
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return TestCase{}, err
	}
	return m, nil
}
