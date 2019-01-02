package model

import (
	"os"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tracer.Silent = true
	os.Exit(m.Run())
}

func TestCreateRequestEmptyEndpointOrMethod(t *testing.T) {
	i := &Input{}
	req, err := i.CreateRequest(emptyTestCase, emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)

	i = &Input{Endpoint: "http://google.com"}
	req, err = i.CreateRequest(emptyTestCase, emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)

	i = &Input{Method: "GET"}
	req, err = i.CreateRequest(emptyTestCase, emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestInputGetValuesMissingContextVariable(t *testing.T) {
	match := Match{Description: "simple match test", ContextName: "GetValueToFind"}
	accessor := ContextAccessor{Matches: []Match{match}}
	i := &Input{Method: "GET", Endpoint: "http://google.com", ContextGet: accessor}
	req, err := i.CreateRequest(emptyTestCase, emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestCreateRequestionNilContext(t *testing.T) {
	i := &Input{Method: "GET", Endpoint: "http://google.com"}
	req, err := i.CreateRequest(emptyTestCase, nil)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestCreateRequestionNilTestcase(t *testing.T) {
	i := &Input{Method: "GET", Endpoint: "http://google.com"}
	req, err := i.CreateRequest(nil, emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestCreateRequestNilHeaderContext(t *testing.T) {
	headers := map[string]string{
		"Myheader": "myValue",
	}
	i := &Input{Method: "GET", Endpoint: "http://google.com", Headers: headers}
	req, err := i.CreateRequest(emptyTestCase, emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	for k, v := range req.Header {
		assert.Equal(t, "Myheader", k)
		assert.Equal(t, "myValue", v[0])
	}
}

func TestCreateRequestHeaderContext(t *testing.T) {
	headers := map[string]string{
		"Myheader": "$replacement",
	}
	ctx := Context{
		"replacement": "myNewValue",
	}
	i := &Input{Method: "GET", Endpoint: "http://google.com", Headers: headers}
	req, err := i.CreateRequest(emptyTestCase, &ctx)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	for k, v := range req.Header {
		assert.Equal(t, "Myheader", k)
		assert.Equal(t, "myNewValue", v[0])
	}
}

func TestCreateRequestHeaderContextFails(t *testing.T) {
	headers := map[string]string{
		"Myheader": "$replacement",
	}
	ctx := Context{
		"nomatch": "myNewValue",
	}
	i := &Input{Method: "GET", Endpoint: "http://google.com", Headers: headers}
	req, err := i.CreateRequest(emptyTestCase, &ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestCheckAuthorizationTokenProcessed(t *testing.T) {
	m := Match{Description: "TokenProcessing", Authorisation: "Bearer"}
	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Authorization", "Bearer 1010110101010101")
	result, err := tc.Validate(resp, emptyContext)
	assert.Equal(t, "1010110101010101", tc.Expect.Matches[0].Result)
	assert.Nil(t, err)
	assert.True(t, result)

	ctx := &Context{
		"access_token": "1010101010101010",
	}
	match := Match{Description: "test", ContextName: "access_token", Authorisation: "bearer"}
	accessor := ContextAccessor{Matches: []Match{match}}
	i := &Input{Method: "GET", Endpoint: "http://google.com", ContextGet: accessor}
	req, err := i.CreateRequest(emptyTestCase, ctx)
	assert.Nil(t, err)
	assert.NotNil(t, req)
}

func TestFormData(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST", FormData: map[string]string{
		"grant_type": "client_credentials",
		"scope":      "accounts openid"}}
	ctx := Context{"baseurl": "http://mybaseurl"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, 2, len(req.FormData))
}

func TestFormDataMissingContextVariable(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST", FormData: map[string]string{
		"grant_type": "$client_credentials",
		"scope":      "accounts openid"}}
	ctx := Context{"baseurl": "http://mybaseurl"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestInputBody(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST", RequestBody: "The Rain in Spain Falls Mainly on the Plain"}
	ctx := Context{"baseurl": "http://mybaseurl"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "The Rain in Spain Falls Mainly on the Plain", req.Body.(string))
}
