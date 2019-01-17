package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// func TestInputGetValuesMissingContextVariable(t *testing.T) {
// 	match := Match{Description: "simple match test", ContextName: "GetValueToFind"}
// 	accessor := ContextAccessor{Matches: []Match{match}}
// 	i := &Input{Method: "GET", Endpoint: "http://google.com", ContextGet: accessor}
// 	req, err := i.CreateRequest(emptyTestCase, emptyContext)
// 	assert.NotNil(t, err)
// 	assert.Nil(t, req)
// }

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

func TestSetBearerAuthTokenFromContext(t *testing.T) {
	headers := map[string]string{
		"authorization": "Bearer $access_token",
	}
	ctx := Context{
		"access_token": "myShineyNewAccessTokenHotOffThePress",
	}
	i := &Input{Method: "GET", Endpoint: "http://google.com", Headers: headers}
	req, err := i.CreateRequest(emptyTestCase, &ctx)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	for k, v := range req.Header {
		assert.Equal(t, "Authorization", k)
		assert.Equal(t, "Bearer myShineyNewAccessTokenHotOffThePress", v[0])
	}
}

func TestCreateHeaderContextMissingForReplacement(t *testing.T) {
	ctx := Context{
		"nomatch": "myNewValue",
	}
	result, err := ReplaceContextField("$replacement", &ctx)
	assert.NotNil(t, err)
	assert.Equal(t, "$replacement", result)

}

// func TestCheckAuthorizationTokenProcessed(t *testing.T) {
// 	m := Match{Description: "TokenProcessing", Authorisation: "Bearer"}
// 	tc := TestCase{Expect: Expect{Matches: []Match{m}, StatusCode: 200}}
// 	resp := test.CreateHTTPResponse(200, "OK", "TheRainInSpain", "Authorization", "Bearer 1010110101010101")
// 	result, err := tc.Validate(resp, emptyContext)
// 	assert.Equal(t, "1010110101010101", tc.Expect.Matches[0].Result)
// 	assert.Nil(t, err)
// 	assert.True(t, result)

// 	ctx := &Context{
// 		"access_token": "1010101010101010",
// 	}
// 	match := Match{Description: "test", ContextName: "access_token", Authorisation: "bearer"}
// 	accessor := ContextAccessor{Matches: []Match{match}}
// 	i := &Input{Method: "GET", Endpoint: "http://google.com", ContextGet: accessor}
// 	req, err := i.CreateRequest(emptyTestCase, ctx)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, req)
// }

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

func TestInputClaims(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST",
		Generation: map[string]string{
			"strategy": "consenturl",
		},
		Claims: map[string]string{
			"iss":          "8672384e-9a33-439f-8924-67bb14340d71",
			"scope":        "openid accounts",
			"redirect_url": "https://test.example.co.uk/redir",
			"responseType": "code",
		}}
	ctx := Context{"baseurl": "http://mybaseurl"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	m, _ := url.ParseQuery(req.URL)
	assert.Equal(t, m["request"][0], "eyJhbGciOiJub25lIn0.eyJhdWQiOiIiLCJjbGFpbXMiOnsiaWRfdG9rZW4iOnsib3BlbmJhbmtpbmdfaW50ZW50X2lkIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIifX19LCJpc3MiOiI4NjcyMzg0ZS05YTMzLTQzOWYtODkyNC02N2JiMTQzNDBkNzEiLCJyZWRpcmVjdF91cmkiOiJodHRwczovL3Rlc3QuZXhhbXBsZS5jby51ay9yZWRpciIsInNjb3BlIjoib3BlbmlkIGFjY291bnRzIn0.")
}

func TestInputClaimsWithContextReplacementParameters(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST",
		Generation: map[string]string{
			"strategy": "consenturl",
		},
		Claims: map[string]string{
			"aud":          "$baseurl",
			"iss":          "8672384e-9a33-439f-8924-67bb14340d71",
			"scope":        "openid accounts",
			"redirect_url": "https://test.example.co.uk/redir",
			"consentId":    "$consent_id",
			"responseType": "code",
		}}
	ctx := Context{"baseurl": "http://mybaseurl", "consent_id": "myconsentid"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	m, _ := url.ParseQuery(req.URL)
	assert.Equal(t, m["request"][0], "eyJhbGciOiJub25lIn0.eyJhdWQiOiJodHRwOi8vbXliYXNldXJsIiwiY2xhaW1zIjp7ImlkX3Rva2VuIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJlc3NlbnRpYWwiOnRydWUsInZhbHVlIjoibXljb25zZW50aWQifX19LCJpc3MiOiI4NjcyMzg0ZS05YTMzLTQzOWYtODkyNC02N2JiMTQzNDBkNzEiLCJyZWRpcmVjdF91cmkiOiJodHRwczovL3Rlc3QuZXhhbXBsZS5jby51ay9yZWRpciIsInNjb3BlIjoib3BlbmlkIGFjY291bnRzIn0.")

}

func TestInputClaimsConsentId(t *testing.T) {
	ctx := Context{"consent_id": "aac-fee2b8eb-ce1b-48f1-af7f-dc8f576d53dc", "xchange_code": "10e9d80b-10d4-4abd-9fe0-15789cc512b5", "baseurl": "https://modelobankauth2018.o3bank.co.uk:4101", "access_token": "18d5a754-0b76-4a8f-9c68-dc5caaf812e2"}
	_ = ctx
	i := Input{Endpoint: "/accounts", Method: "POST",
		Generation: map[string]string{
			"strategy": "consenturl",
		},
		Claims: map[string]string{
			"aud":          "$baseurl",
			"iss":          "8672384e-9a33-439f-8924-67bb14340d71",
			"scope":        "openid accounts",
			"redirect_url": "https://test.example.co.uk/redir",
			"consentId":    "$consent_id",
			"responseType": "code",
		}}
	tc := TestCase{Input: i, Context: ctx}
	res, err := i.CreateRequest(&tc, &ctx)
	assert.NoError(t, err, "create request should succeed")
	assert.NotNil(t, res)
}
