package model

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const simplejson = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

// create a context
// create a json match object
// create an http response
// send response to match object
// have match object parse response and extract json field into context parameter
func TestContextPutFromMatch(t *testing.T) {
	ctx := Context{}
	m := Match{JSON: "name.last", Description: "simple match test", ContextName: "LastName"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, m.PutValue(string(bodyBytes), &ctx))
	assert.Equal(t, ctx.Get(m.ContextName), "Prichard")
}

// JSON field match on response string, and return field value + context variable name for context insertion
func TestContextGetFromContext(t *testing.T) {
	ctx := Context{}
	m := Match{JSON: "name.first", Description: "simple match test", ContextName: "FirstName"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	value, variable := m.GetValue(string(bodyBytes))
	assert.Equal(t, "Janet", value)
	assert.Equal(t, ctx.Get(m.ContextName), ctx.Get(variable))
}

// Create a testcase that defines the basic matchers
// json matcher
func TestJSONMatcher(t *testing.T) {
	match := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	tc := TestCase{Expect: Expect{Matches: []Match{match}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, success)
}

// check
func TestJSONMatcherMismatch(t *testing.T) {
	match := Match{JSON: "name.first", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	tc := TestCase{Expect: Expect{Matches: []Match{match}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	success, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, success)
}

var status200 = []byte(`{"expect": {"status-code": 200}}`)

// check match on status code is detected
func TestMatchOnStatusCode(t *testing.T) {
	var tc TestCase
	json.Unmarshal(status200, &tc)
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	result, err := tc.ApplyExpects(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check status code mismatch is detected
func TestNoMatchOnStatusCode(t *testing.T) {
	var tc TestCase
	json.Unmarshal(status200, &tc)
	resp := pkgutils.CreateHTTPResponse(404, "File Not Found", simplejson)
	result, err := tc.ApplyExpects(resp, nil)
	assert.Equal(t, "():: HTTP Status code does not match: expected 200 got 404", err.Error())
	assert.False(t, result)
}

const statusok = `{"status":"ok"}`

// check header value match is detected
func TestMatchResponseHeaderValue(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/borg"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/borg")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header value match is detected
func TestMatchResponseHeaderValueCaseInsensitive(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/borg"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "content-type", "application/borg")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header value mismatch is detected
func TestNoMatchResponseHeaderValue(t *testing.T) {
	m := Match{Description: "header test", Header: "Content-Type", Value: "application/klingon"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/borg")
	result, err := tc.Validate(resp, nil)
	assert.Contains(t, err.Error(), "expected (application/klingon) got (application/borg)")
	assert.False(t, result)
}

// detect invalid match type
func TestInvalidMatchType(t *testing.T) {
	m := Match{Description: "type test"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Content-Type", "application/json")
	result, err := tc.Validate(resp, nil)
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
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckHeaderRegexMatchCaseInsensitive(t *testing.T) {
	m := Match{Description: "test", Header: "authorization", Regex: "^Basic\\s.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header regex mismatch is detected
func TestCheckHeaderRegexMismatch(t *testing.T) {
	m := Match{Description: "test", Header: "authorization", Regex: "^Basic\\s.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Authorization", "Basics YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check header present is detected
func TestCheckHeaderPresentMatch(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "authorization", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckHeaderPresentMatchCaseInsensitive(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "AuthoriZation", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check header present fail is detected
func TestCheckHeaderPresentMismat1ch(t *testing.T) {
	m := Match{Description: "test", HeaderPresent: "authorization"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", statusok, "Security_token", "Basic YjMzODg4ZGMtYzg==")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check body regex match is detected
func TestCheckBodyRegexMatch(t *testing.T) {
	m := Match{Description: "test", Regex: ".*London Bridge.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "{\"status\":\"London Bridge\"}")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check body regex mismatch is detect
func TestCheckBodyRegexMismatch(t *testing.T) {
	m := Match{Description: "test", Regex: ".*London Bridge.*"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "{\"status\":\"London !! Bridge\"}")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

// check BodyJSONPresent match is detected
func TestCheckBodyJsonMatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridge\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

// check BodyJSONPresent mismatch is detected
func TestCheckBodyJsonMisMatch(t *testing.T) {
	m := Match{Description: "test", JSON: "tourist-attractions.bridge"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "{\"status\":\"OK\",\"tourist-attractions\":{\"bridges\":\"Tower Bridge\"}}")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

const testbankjson = `{"banks": [{"name": "Barclays" },{"name": "HSBC"},{"name": "Lloyds" }]}`

func TestCheckJsonBodyCount(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.#", Count: 3}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckJsonBodyCountMismatch(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.#", Count: 2}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyJSONRegex(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.name", Regex: "^L.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckBodyJSONRegexMismatch(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.name", Regex: "^B.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyJSONRegexMismatchJSONPattern(t *testing.T) {
	m := Match{Description: "test", JSON: "banks.2.names", Regex: "^B.*s$"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", testbankjson)
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyLength(t *testing.T) {
	var len int64 = 35
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "TheRainInSpainFallsMainlyOnThePlain")
	result, err := tc.Validate(resp, nil)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckBodyLengthMismatch(t *testing.T) {
	var len int64 = 11
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "TheRainInSpainFallsMainlyOnThePlain")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestCheckBodyLengthMismatch2(t *testing.T) {
	var len int64 = 35
	m := Match{Description: "test", BodyLength: &len}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := pkgutils.CreateHTTPResponse(200, "OK", "")
	result, err := tc.Validate(resp, nil)
	assert.NotNil(t, err)
	assert.False(t, result)
}
