package model

import (
	"encoding/json"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"github.com/stretchr/testify/require"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"github.com/stretchr/testify/assert"
)

const simplejson = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

// create a context
// create a json match object
// create an http response
// send response to match object
// have match object parse response and extract json field into context parameter
var emptyContext = &Context{}
var emptyTestCase = &TestCase{}

func TestContextPutFromMatch(t *testing.T) {
	ctx := Context{}
	m := Match{JSON: "name.last", Description: "simple match test", ContextName: "LastName"}
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	tc := TestCase{Body: resp.String()}
	assert.True(t, m.PutValue(&tc, &ctx))
	ctxvalue, exists := ctx.Get(m.ContextName)
	assert.True(t, exists)
	assert.Equal(t, "Prichard", ctxvalue)
}

// Create a testcase that defines the basic matchers
// json matcher
func TestJSONBodyValue(t *testing.T) {
	match := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	tc := TestCase{Expect: Expect{Matches: []Match{match}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, success)
}

// check
func TestJSONBodyValueMismatch(t *testing.T) {
	match := Match{JSON: "name.first", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	tc := TestCase{Expect: Expect{Matches: []Match{match}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, success)
}

var status200 = []byte(`{"expect": {"status-code": 200}}`)

// check match on status code is detected
func TestMatchOnStatusCode(t *testing.T) {
	var tc TestCase
	err := json.Unmarshal(status200, &tc)
	require.NoError(t, err)
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	result, errs := tc.ApplyExpects(resp, nil)
	assert.Nil(t, errs)
	assert.True(t, result)
}

// check status code mismatch is detected
func TestNoMatchOnStatusCode(t *testing.T) {
	var tc TestCase
	err := json.Unmarshal(status200, &tc)
	require.NoError(t, err)
	resp := test.CreateHTTPResponse(404, "File Not Found", simplejson)
	result, errs := tc.ApplyExpects(resp, nil)
	assert.Equal(t, "():: HTTP Status code does not match: expected 200 got 404", errs[0].Error())
	assert.False(t, result)
}

const statusok = `{"status":"isReplacement"}`

// check header value match is detected
func TestMatchResponseHeaderValue(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/borg"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/borg")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header value match is detected
func TestMatchResponseHeaderValueCaseInsensitive(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/borg"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "content-type", "application/borg")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header value mismatch is detected
func TestNoMatchResponseHeaderValue(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/klingon"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/borg")
	result, errs := tc.Validate(resp, emptyContext)
	assert.Contains(t, errs[0].Error(), "expected (application/klingon) got (application/borg)")
	assert.False(t, result)
}

// detect invalid match type
func TestInvalidMatchType(t *testing.T) {
	m := Match{Description: "type test"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/json")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check getType returns known type
func TestGetMatchType(t *testing.T) {
	m := Match{Description: "type test", MatchType: BodyJSONCount}
	gettype := m.GetType()
	assert.Equal(t, m.MatchType, gettype)
}

// check Header Regex match is detected
func TestCheckHeaderRegexMatch(t *testing.T) {
	m := Match{Description: "test", Header: "Authorization", Regex: "^Basic\\s.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckHeaderRegexMatchCaseInsensitive(t *testing.T) {
	m := Match{Description: "test", Header: "authorization", Regex: "^Basic\\s.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header regex mismatch is detected
func TestCheckHeaderRegexMismatch(t *testing.T) {
	m := Match{Description: "test", Header: "authorization", Regex: "^Basic\\s.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basics YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckHeaderRegexCompileFail(t *testing.T) {
	m := Match{Description: "test", Header: "Authorization", Regex: `[ ]\K(?<!\d`}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check header present is detected
func TestCheckHeaderPresentMatch(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckHeaderPresentMatchCaseInsensitive(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "AuthoriZation", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header present fail is detected
func TestCheckHeaderPresentMismatch(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", statusok, "Security_token", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check body regex match is detected
func TestCheckBodyRegexMatch(t *testing.T) {
	m := Match{Description: "test", Regex: ".*London Bridge.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"London Bridge\"}")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check body regex mismatch is detect
func TestCheckBodyRegexMismatch(t *testing.T) {
	m := Match{Description: "test", Regex: ".*London Bridge.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"London !! Bridge\"}")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyRegexWithContextPut(t *testing.T) {
	m := Match{Description: "test", Regex: ".*", ContextName: "mybody"}
	ctx := Context{}
	ca := ContextAccessor{Context: &Context{}, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{StatusCode: 200, ContextPut: ca}, Context: Context{}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"London !! Bridge\"}")
	result, err := tc.Validate(resp, &ctx)
	assert.Nil(t, err)
	assert.True(t, result)
	value, exists := ctx.Get("mybody")
	assert.True(t, exists)
	assert.Equal(t, value, "{\"status\":\"London !! Bridge\"}")
}

func TestCheckBodyRegexWithContextPutBadRegex(t *testing.T) {
	m := Match{Description: "test", Regex: "x.*", ContextName: "mybody"}
	ctx := Context{}
	ca := ContextAccessor{Context: &Context{}, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{StatusCode: 200, ContextPut: ca}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"London !! Bridge\"}")
	result, err := tc.Validate(resp, &ctx)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyRegexCompileError(t *testing.T) {
	m := Match{Description: "test", Regex: ".*\\KLondon Bridge.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"London Bridge\"}")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check BodyJSONPresent match is detected
func TestCheckBodyJsonMatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridge\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check BodyJSONPresent mismatch is detected
func TestCheckBodyJsonMisMatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridges\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

const testbankjson = `{"banks": [{"name": "Barclays" },{"name": "HSBC"},{"name": "Lloyds" }]}`

func TestCheckJsonBodyCount(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.#", Count: 3}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckJsonBodyCountMismatch(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.#", Count: 2}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyJSONRegex(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.name", Regex: "^L.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckBodyJSONRegexMismatch(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.name", Regex: "^B.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyJSONRegexMismatchJSONPattern(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.names", Regex: "^B.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyJSONRegexCompileFail(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.name", Regex: "^L.\\K*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyLength(t *testing.T) {
	var len int64 = 35
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpainFallsMainlyOnThePlain")
	result, err := tc.Validate(resp, emptyContext)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckBodyLengthMismatch(t *testing.T) {
	var len int64 = 11
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpainFallsMainlyOnThePlain")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyLengthMismatch2(t *testing.T) {
	var len int64 = 35
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "")
	result, err := tc.Validate(resp, emptyContext)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestMatchStringOutput(t *testing.T) {
	var len int64 = 77
	var num int64 = 88
	var count int64 = 88
	m := Match{Description: "test",
		ContextName: "context", Header: "header", HeaderPresent: "presentheader",
		Regex: "myregex", JSON: "myjson", Value: "myvalue", Numeric: num, Count: count, BodyLength: &len,
		ReplaceEndpoint: "myreplace", Authorisation: "myauthorisation", Result: "myresult"}
	assert.Equal(t, "{\"match_type\":11,\"description\":\"test\",\"name\":\"context\",\"header\":\"header\",\"header-present\":\"presentheader\",\"regex\":\"myregex\",\"json\":\"myjson\",\"value\":\"myvalue\",\"numeric\":88,\"count\":88,\"body-length\":77,\"replaceInEndpoint\":\"myreplace\",\"authorisation\":\"myauthorisation\",\"result\":\"myresult\"}",
		m.String())

}

func TestCheckAuthorisation(t *testing.T) {
	m := Match{Description: "AuthTest", Authorisation: "Bearer"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Authorization", "Bearer 1010110101010101")
	result, err := tc.Validate(resp, emptyContext)
	assert.Equal(t, "1010110101010101", tc.Expect.Matches[0].Result)
	assert.Nil(t, err)
	assert.True(t, result)
}
func TestCheckAuthorisationNotPresent(t *testing.T) {
	m := Match{Description: "AuthTest", Authorisation: "Bearer"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain")
	result, err := tc.Validate(resp, emptyContext)
	assert.Equal(t, "", tc.Expect.Matches[0].Result)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckAuthorisationIncorrectValue(t *testing.T) {
	m := Match{Description: "AuthTest", Authorisation: "Bearer"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Authorisation", "Beardy 12312312")
	result, err := tc.Validate(resp, emptyContext)
	assert.Equal(t, "", tc.Expect.Matches[0].Result)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestContextPutHeaderRegexContextSubFieldCapture(t *testing.T) {
	m := Match{Description: "AuthCode xChange", Header: "Location", Regex: "code=(.*)&+.*", ContextName: "mycode"}
	ctx := Context{}
	ctxPut := ContextAccessor{Context: &ctx, Matches: []Match{m}}
	tc := TestCase{Context: Context{}, Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()}

	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Location", "https://mysite/auth?code=1234&redir=here")
	result, err := tc.Validate(resp, &ctx)
	assert.True(t, result)
	ctxCode, exists := ctx.Get("mycode")
	assert.True(t, exists)
	assert.True(t, exists)
	assert.Equal(t, "1234", ctxCode) // code from location header now accessible in context
	assert.Nil(t, err)
}

func TestContextPutHeaderRegexContextSubFieldCaptureFail(t *testing.T) {
	m := Match{Description: "AuthCode xChange", Header: "Location", Regex: "xcode=(.*)&+.*", ContextName: "mycode"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Location", "https://mysite/auth?code=1234&redir=here")
	result, err := tc.Validate(resp, &c)
	assert.False(t, result)
	assert.NotNil(t, err)
}

func TestContextPutHeaderRegexContextSubFieldCompileFaile(t *testing.T) {
	m := Match{Description: "AuthCode xChange", Header: "Location", Regex: "code=\\K(.*)&+.*", ContextName: "mycode"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Location", "https://mysite/auth?code=1234&redir=here")
	result, err := tc.Validate(resp, &c)
	assert.False(t, result)
	assert.NotNil(t, err)
}

// check contextPut BodyJSONPresent match is detected
func TestContextPutCheckBodyJsonMatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge", ContextName: "attractions"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj

	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridge\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, &c)
	assert.Nil(t, err)
	assert.True(t, result)
	ctxAttract, exists := c.Get("attractions")
	assert.True(t, exists)
	assert.Equal(t, "Tower Bridge", ctxAttract) // check body value now in context
}

func TestContextPutCheckBodyJsonMatchMismatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge1", ContextName: "attractions"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridge\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, &c)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// Create a testcase that defines the basic matchers
// json matcher
func TestContextPutJSONBodyValue(t *testing.T) {
	m := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, &c)
	assert.Nil(t, err)
	assert.True(t, success)
}

func TestContextPutJSONBodyValueFail(t *testing.T) {
	m := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard1"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, &c)
	assert.NotNil(t, err)
	assert.False(t, success)
}

// Create a testcase that defines the basic matchers
// json matcher
func TestValidateNilContext(t *testing.T) {
	tc := TestCase{}
	resp := test.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, success)
}

// json matcher
func TestContextPutAuthorisation(t *testing.T) {
	m := Match{Description: "test", ContextName: "token", Authorisation: "bearer"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", simplejson, "Authorisation", "Bearer 10101011010")
	success, err := tc.Validate(resp, &c)
	assert.Nil(t, err)
	assert.True(t, success)
	ctxToken, exist := c.Get("token")
	assert.True(t, exist)
	assert.Equal(t, "10101011010", ctxToken)
}

func TestContextPutAuthorisationFail(t *testing.T) {
	m := Match{Description: "test", ContextName: "token", Authorisation: "bearer"}
	c := Context{}
	ctxPut := ContextAccessor{Context: &c, Matches: []Match{m}}
	tc := TestCase{Expect: Expect{ContextPut: ctxPut, StatusCode: 200}, Validator: schema.NewNullValidator()} // note: the "Match" lives in the contextPut obj
	resp := test.CreateHTTPResponse(200, "OK", simplejson, "Authorisation", "Barat 10101011010")
	success, err := tc.Validate(resp, &c)
	assert.NotNil(t, err)
	assert.False(t, success)
	ctxToken, exist := c.Get("token")
	assert.False(t, exist)
	assert.Equal(t, nil, ctxToken)
}
