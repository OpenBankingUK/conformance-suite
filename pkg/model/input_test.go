package model

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		"replacement":            "myNewValue",
		"authorisation_endpoint": "https://example.com/authorisation",
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
		"access_token":           "myShineyNewAccessTokenHotOffThePress",
		"authorisation_endpoint": "https://example.com/authorisation",
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
		"phase":                  "run",
		"nomatch":                "myNewValue",
		"authorisation_endpoint": "https://example.com/authorisation",
	}
	result, err := replaceContextField("$replacement", &ctx)
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
	ctx := Context{"baseurl": "http://mybaseurl", "authorisation_endpoint": "https://example.com/authorisation"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, 2, len(req.FormData))
}

func TestFormDataMissingContextVariable(t *testing.T) {
	ctx1 := Context{"phase": "run"}
	i := Input{Endpoint: "/accounts", Method: "POST", FormData: map[string]string{
		"grant_type": "$client_credentials",
		"scope":      "accounts openid"}}
	ctx := Context{"baseurl": "http://mybaseurl", "authorisation_endpoint": "https://example.com/authorisation"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(&ctx1)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestInputBody(t *testing.T) {
	i := Input{Endpoint: "/accounts", Method: "POST", RequestBody: "The Rain in Spain Falls Mainly on the Plain"}
	ctx := Context{"baseurl": "http://mybaseurl", "authorisation_endpoint": "https://example.com/authorisation"}
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
	ctx := Context{"baseurl": "http://mybaseurl", "authorisation_endpoint": "https://example.com/authorisation"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	m, err := url.ParseQuery(req.URL)
	require.NoError(t, err)
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
	ctx := Context{"baseurl": "http://mybaseurl", "consent_id": "myconsentid", "authorisation_endpoint": "https://example.com/authorisation"}
	tc := TestCase{Input: i, Context: ctx}
	req, err := tc.Prepare(emptyContext)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	m, err := url.ParseQuery(req.URL)
	require.NoError(t, err)
	assert.Equal(t, m["request"][0], "eyJhbGciOiJub25lIn0.eyJhdWQiOiJodHRwOi8vbXliYXNldXJsIiwiY2xhaW1zIjp7ImlkX3Rva2VuIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJlc3NlbnRpYWwiOnRydWUsInZhbHVlIjoibXljb25zZW50aWQifX19LCJpc3MiOiI4NjcyMzg0ZS05YTMzLTQzOWYtODkyNC02N2JiMTQzNDBkNzEiLCJyZWRpcmVjdF91cmkiOiJodHRwczovL3Rlc3QuZXhhbXBsZS5jby51ay9yZWRpciIsInNjb3BlIjoib3BlbmlkIGFjY291bnRzIn0.")

}

func TestInputClaimsConsentId(t *testing.T) {
	ctx := Context{"consent_id": "aac-fee2b8eb-ce1b-48f1-af7f-dc8f576d53dc", "xchange_code": "10e9d80b-10d4-4abd-9fe0-15789cc512b5", "baseurl": "https://modelobankauth2018.o3bank.co.uk:4101", "access_token": "18d5a754-0b76-4a8f-9c68-dc5caaf812e2", "authorisation_endpoint": "https://example.com/authorisation"}
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

func TestClaimsJWTBearer(t *testing.T) {
	cert, err := authentication.NewCertificate(selfsignedDummypub, selfsignedDummykey)
	require.NoError(t, err)
	ctx := Context{
		"consent_id":             "aac-fee2b8eb-ce1b-48f1-af7f-dc8f576d53dc",
		"xchange_code":           "10e9d80b-10d4-4abd-9fe0-15789cc512b5",
		"baseurl":                "https://matls-sso.openbankingtest.org.uk",
		"access_token":           "18d5a754-0b76-4a8f-9c68-dc5caaf812e2",
		"client_id":              "12312",
		"scope":                  "AuthoritiesReadAccess ASPSPReadAccess TPPReadAll",
		"SigningCert":            cert,
		"signingPrivate":         selfsignedDummykey,
		"signingPublic":          selfsignedDummypub,
		"authorisation_endpoint": "https://example.com/authorisation",
		"nonOBDirectory":         false,
	}

	i := Input{Endpoint: "/as/token.oauth2", Method: "POST",
		Generation: map[string]string{
			"strategy": "jwt-bearer",
		},
		Claims: map[string]string{
			"aud":   "$baseurl",
			"iss":   "$client_id",
			"scope": "AuthoritiesReadAccess ASPSPReadAccess TPPReadAll",
		},
		FormData: map[string]string{
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"grant_type":            "client_credentials",
			"client_id":             "$client_id",
			"scope":                 "$scope",
		},
	}
	tc := TestCase{Input: i, Context: ctx}
	res, err := i.CreateRequest(&tc, &ctx)
	require.NoError(t, err, "create request should succeed")
	assert.NotNil(t, res)
	jwtbearer, exists := ctx.Get("jwtbearer")
	assert.True(t, exists)
	assert.True(t, len(jwtbearer.(string)) > 20)
}

func TestJWTSignRS256(t *testing.T) {
	cert, err := authentication.NewCertificate(selfsignedDummypub, selfsignedDummykey)
	require.NoError(t, err)
	require.NotNil(t, cert)

	alg := jwt.GetSigningMethod("RS256")
	if alg == nil {
		t.Logf("Couldn't find signing method: %v\n", alg)
	}
	claims := jwt.MapClaims{}
	claims["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(alg, claims) // create new token
	token.Header["kid"] = "mykeyid"
	prikey := cert.PrivateKey()
	tokenString, err := token.SignedString(prikey) // sign the token - get as encoded string
	if err != nil {
		t.Log("error signing jwt: ", err)
	}
	assert.True(t, len(tokenString) > 30)

}

func TestBodyLiteral(t *testing.T) {
	ctx := Context{
		"replacebody":            "this is my body",
		"authorisation_endpoint": "https://example.com/authorisation",
	}

	i := Input{Method: "POST", Endpoint: "https://google.com", RequestBody: "This is my literal body"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	assert.Equal(t, "This is my literal body", req.Body)
}

func TestBodyReplacement(t *testing.T) {
	ctx := Context{
		"replacebody":            "this is my body",
		"authorisation_endpoint": "https://example.com/authorisation",
	}

	i := Input{Method: "POST", Endpoint: "https://google.com", RequestBody: "$replacebody"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	assert.Equal(t, "this is my body", req.Body)
}

func TestBodyTwoReplacements(t *testing.T) {
	ctx := Context{
		"replacebody":            "this is my body",
		"replace2":               "and this is my heart",
		"authorisation_endpoint": "https://example.com/authorisation",
	}

	i := Input{Method: "POST", Endpoint: "https://google.com", RequestBody: "$replacebody $replace2"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	assert.Equal(t, "this is my body and this is my heart", req.Body)
}

func TestPaymentBodyReplace(t *testing.T) {
	ctx := Context{
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
	}

	i := Input{Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	assert.Equal(t, "{\"Data\": {\"ConsentId\": \"sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d\",\"Initiation\":{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}} },\"Risk\":{}}", req.Body)
}

func TestPaymentBodyReplaceTestCase100300(t *testing.T) {
	ctx := Context{
		"x-fapi-financial-id":    "myfapiid",
		"thisSchemeName":         "myscheme",
		"thisIdentification":     "myid",
		"authorisation_endpoint": "https://example.com/authorisation",
	}
	_ = ctx
	var tc TestCase
	err := json.Unmarshal([]byte(paymentTestCaseData100300), &tc)
	assert.Nil(t, err)
	t.Logf("%#v\n", tc)
	ctx.PutContext(&tc.Context)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	ctx.DumpContext("Testcase")
	req, err := tc.Prepare(&ctx)

	assert.Nil(t, err)
	_ = req
	t.Logf("%#v\n", tc)
}

func TestJWSDetachedSignature(t *testing.T) {
	EnableJWS()
	ctx := Context{
		"signingPrivate":            selfsignedDummykey,
		"signingPublic":             selfsignedDummypub,
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
		"api-version":               "v3.0",
		"nonOBDirectory":            false,
		"requestObjectSigningAlg":   "PS256",
	}

	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
}

func TestJWSSignaturNotPOST(t *testing.T) {
	ctx := Context{
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
	}

	i := Input{JwsSig: true, Method: "GET", Endpoint: "https://google.com", RequestBody: ""}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.EqualError(t, err, "createRequest: cannot apply jws signature to method that isn't POST")
	assert.Nil(t, req)
}

func TestJWSSignatureEmptyBody(t *testing.T) {
	EnableJWS()
	ctx := Context{
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
	}

	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: ""}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.EqualError(t, err, "createRequest: cannot create x-jws-signature, as request body is empty")
	assert.Nil(t, req)
}

func TestJWSDetachedSignatureGET(t *testing.T) {
	ctx := Context{
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
	}

	i := Input{JwsSig: true, Method: "GET", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestJWSDetachedSignature312andBefore(t *testing.T) {
	EnableJWS()
	ctx := Context{
		"signingPrivate":            selfsignedDummykey,
		"signingPublic":             selfsignedDummypub,
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
		"api-version":               "v3.0",
		"nonOBDirectory":            false,
		"requestObjectSigningAlg":   "PS256",
	}

	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	fmt.Println(sig)
}

func TestJWSSign(t *testing.T) {
	EnableJWS()
	ctx := Context{
		"signingPrivate":            selfsignedDummykey,
		"signingPublic":             selfsignedDummypub,
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
		"api-version":               "v3.0",
		"nonOBDirectory":            false,
		"requestObjectSigningAlg":   "PS256",
	}

	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	fmt.Println(sig)

	hmacSampleSecret := []byte("hello")
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	fmt.Println(tokenString)

	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err = jwt.Parse(tokenString, keyFunc)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		//		time.Time(claims["nbf"]
		value := claims["nbf"]

		t2 := value.(float64)
		t1 := int64(t2)

		// i, err := strconv.ParseInt(tm1, 10, 64)
		// if err != nil {
		// 	panic(err)
		// }
		tm := time.Unix(t1, 0)
		fmt.Println(tm)

		fmt.Println(claims["foo"], tm)
	} else {
		fmt.Println("An Error has occurred: " + err.Error())
	}
}

type reportClaims struct {
	jwt.StandardClaims
	ReportDigest    string `json:"reportDigest,omitempty"`
	DiscoveryDigest string `json:"discoveryDigest,omitempty"`
	ManifestDigest  string `json:"manifestDigest,omitempty"`
}

func TestGenealSigs(t *testing.T) {
	pb, rest := pem.Decode([]byte(signing_private))
	require.NotNil(t, pb, "pem decode private key, block nil")
	require.Len(t, rest, 0, "rest should be zero length")
	privateKey, err := x509.ParsePKCS8PrivateKey(pb.Bytes)
	require.NoError(t, err, "parse private key")

	meta := map[string]string{
		"header-foo": "header-bar",
	}

	knownClaims := reportClaims{
		ReportDigest:    "report-hash-sum",
		DiscoveryDigest: "discovery-hash-sum",
		ManifestDigest:  "manifest-hash-sum",
		StandardClaims: jwt.StandardClaims{
			Issuer:    "https://openbanking.org.uk/fcs/reporting",
			Subject:   "openbanking.org.uk",
			Id:        "unique-jwt-id",
			ExpiresAt: 4077614132,
			NotBefore: 90000,
		},
	}

	signed, err := sign(knownClaims, meta, privateKey.(*rsa.PrivateKey))
	require.NoError(t, err, "Sign error")

	fmt.Println(signed)

	calcClaims := reportClaims{}
	tk, err := jwt.ParseWithClaims(signed, &calcClaims, keyFunc)
	require.NoError(t, err, "jwt.Parse")
	fmt.Printf("claims: %v", calcClaims)

	signedString, err := tk.SigningString()
	require.NoError(t, err, "tk.SigningString()")
	expectedString := fmt.Sprintf("%s.%s", signedString, tk.Signature)

	require.Equal(t, signed, expectedString)
	require.Equal(t, knownClaims.Issuer, calcClaims.Issuer, "claim Issuer")
	require.Equal(t, knownClaims.Subject, calcClaims.Subject, "claim Subject")
	require.Equal(t, knownClaims.Id, calcClaims.Id, "claim Id")
	require.Equal(t, knownClaims.NotBefore, calcClaims.NotBefore, "claim NotBefore")
	require.Equal(t, knownClaims.ExpiresAt, calcClaims.ExpiresAt, "claim ExpiresAt")
	require.Equal(t, knownClaims.ReportDigest, calcClaims.ReportDigest, "claim ReportDigest")
	require.Equal(t, knownClaims.DiscoveryDigest, calcClaims.DiscoveryDigest, "claim DiscoveryDigest")
	require.Equal(t, knownClaims.ManifestDigest, calcClaims.ManifestDigest, "claim ManifestDigest")

}

func keyFunc(_ *jwt.Token) (interface{}, error) {
	pb, _ := pem.Decode([]byte(signing_public))
	publicKey, err := x509.ParsePKIXPublicKey(pb.Bytes)
	return publicKey, err
}

func sign(claims jwt.Claims, meta map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	for k, v := range meta {
		t.Header[k] = v
	}

	signed, err := t.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func TestDirEncodedJwt(t *testing.T) {
	EnableJWS()
	ctx := Context{
		"signingPrivate":            signing_private,
		"signingPublic":             signing_public_cert,
		"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
		"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
		"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
		"authorisation_endpoint":    "https://example.com/authorisation",
		"api-version":               "v3.0",
		"nonOBDirectory":            false,
		"requestObjectSigningAlg":   "PS256",
	}

	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)

	keyFuncSelfSigned(nil)

	token, err := jwt.Parse(sig, keyFuncSelfSigned)
	if err != nil {
		fmt.Printf("Error verifying signature: %s", err)
		t.Fail()
	} else {
		fmt.Println("Verify signature ok")
	}

	_ = token

}

func keyFuncSelfSigned(_ *jwt.Token) (interface{}, error) {
	pb, _ := pem.Decode([]byte(signing_public))
	publicKey, err := x509.ParsePKIXPublicKey(pb.Bytes)
	return publicKey, err
}

func keyFuncSelfSigned1(_ *jwt.Token) (interface{}, error) {
	logrus.Warnln("Before spewing - keyfuncentry")
	fmt.Println("Parse Token")
	var cert tls.Certificate
	//var certPEMBlock []byte
	var skippedBlockTypes []string
	var certDERBlock *pem.Block
	certDERBlock, _ = pem.Decode([]byte(signing_public_cert))
	if certDERBlock == nil {
		return nil, errors.New("Empty der block")
	}
	if certDERBlock.Type == "CERTIFICATE" {
		cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
	} else {
		skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
	}

	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			logrus.Errorln("tls: failed to find any PEM data in certificate input")
		}
		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			logrus.Errorln(errors.New("tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
		}
		logrus.Errorln(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
	}

	//	spew.Dump(cert)

	//pb, _ := pem.Decode([]byte(cert.Certificate[0]))
	publicKey, err := x509.ParsePKIXPublicKey(cert.Certificate[0])
	return publicKey, err

}

func TestCertDNRetrieval(t *testing.T) {
	cert, err := loadSigningCert(t)
	if err != nil {
		t.Log("Certs not found, skip test")
		return
	}
	_ = cert

}

func TestReadCerts(t *testing.T) {
	certPEMBlock, err := ioutil.ReadFile("../../../certs/sig-xRWcKt4rSGqsIhqJ3xKC6DOjblY.pem")
	if err != nil {
		t.Log(err.Error())
	}
	var cert tls.Certificate
	var skippedBlockTypes []string
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		} else {
			skippedBlockTypes = append(skippedBlockTypes, certDERBlock.Type)
		}
	}
	if len(cert.Certificate) == 0 {
		if len(skippedBlockTypes) == 0 {
			t.Log("tls: failed to find any PEM data in certificate input")
		}
		if len(skippedBlockTypes) == 1 && strings.HasSuffix(skippedBlockTypes[0], "PRIVATE KEY") {
			t.Log(errors.New("tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched"))
		}
		t.Log(fmt.Errorf("tls: failed to find \"CERTIFICATE\" PEM block in certificate input after skipping PEM blocks of the following types: %v", skippedBlockTypes))
	}
	spew.Dump(cert)
}

func TestLoadX509Cert(t *testing.T) {
	cf, err := ioutil.ReadFile("../../../certs/sig-xRWcKt4rSGqsIhqJ3xKC6DOjblY.pem")
	if err != nil {
		t.Log("cfload:", err.Error())
		return
	}
	cpb, _ := pem.Decode(cf)
	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		t.Logf("cannont parse cert %s", err.Error())
		return
	}
	subject := crt.Subject
	c := subject.Country[0]
	o := subject.Organization[0]
	ou := subject.OrganizationalUnit[0]
	cn := subject.CommonName
	dn := fmt.Sprintf("C=%s, O=%s, OU=%s, CN=%s", c, o, ou, cn)
	t.Logf("DN is\n %s \n", dn)
}

func loadSigningCert(t *testing.T) (tls.Certificate, error) {
	certSigning, err := ioutil.ReadFile("../../../certstore/testcertSigning.pem")
	if err != nil {
		t.Log("cannot read signing certificate")
		return tls.Certificate{}, err
	}
	keySigning, err := ioutil.ReadFile("../../../certstore/testprivateKeySigning.key")
	if err != nil {
		t.Log("cannot read signing key")
		return tls.Certificate{}, err
	}

	cert, err := tls.X509KeyPair(certSigning, keySigning)
	if err != nil {
		t.Log("tls.X509KeyPair failed")
		return tls.Certificate{}, err
	}

	return cert, nil
}

func TestInputTest(t *testing.T) {
	t.Log("Running...")
	x := "123.1......45.789"
	firstPart := x[:strings.IndexByte(x, '.')]
	idx := strings.LastIndex(x, ".")
	lastPart := x[idx:]
	t.Logf("%s.%s\n", firstPart, lastPart)
}

const paymentTestCaseData100300 = `
{
    "@id": "OB-301-DOP-100300",
    "name": "Domestic Payment consents succeeds with minimal data set with additional schema checks.",
    "purpose": "Check that the resource succeeds posting a domestic payment consents with a minimal data set and checks additional schema.",
    "input": {
        "method": "POST",
        "endpoint": "/domestic-payment-consents",
        "headers": {
            "Content-Type": "application/json; charset=utf-8",
            "x-fapi-financial-id": "$x-fapi-financial-id",
            "x-fapi-interaction-id": "b4405450-febe-11e8-80a5-0fcebb1574e1",
            "x-fcs-testcase-id": "OB-301-DOP-100300"
        },
        "bodyData": "{\n    \"Data\": {\n        \"Initiation\": {\n            \"CreditorAccount\": {\n                \"Identification\": \"$thisIdentification\",\n                \"Name\": \"CF Tool\",\n                \"SchemeName\": \"$thisSchemeName\"\n            },\n            \"EndToEndIdentification\": \"$thisInstructionIdentification\",\n            \"InstructedAmount\": {\n                \"Amount\": \"1.00\",\n                \"Currency\": \"$thisCurrency\"\n            },\n            \"InstructionIdentification\": \"$thisInstructionIdentification\"\n        }\n    },\n    \"Risk\": {}\n}"
    },
    "context": {
        "baseurl": "http://mybaseurl",
        "requestConsent": "true",
        "thisCurrency": "GBP",
        "thisInstructionIdentification": "OB-301-DOP-100300",
        "tokenScope": "payments",
        "x-fapi-financial-id": "$x-fapi-financial-id"
    },
    "expect": {
        "status-code": 201,
        "schema-validation": true,
        "matches": [
            {
                "header-present": "x-fapi-interaction-id"
            },
            {
                "json": "Data.Status",
                "value": "AwaitingAuthorisation"
            },
            {
                "json": "Data.ConsentId"
            }
        ],
        "contextPut": {
            "matches": [
                {
                    "name": "OB-301-DOP-100300-ConsentId",
                    "json": "Data.ConsentId"
                }
            ]
        }
    }
}`

const selfsignedDummykey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA8Gl2x9KsmqwdmZd+BdZYtDWHNRXtPd/kwiR6luU+4w76T+9m
lmePXqALi7aSyvYQDLeffR8+2dSGcdwvkf6bDWZNeMRXl7Z1jsk+xFN91mSYNk1n
R6N1EsDTK2KXlZZyaTmpu/5p8SxwDO34uE5AaeESeM3RVqqOgRcXskmp/atwUMC+
qLfuXPoNSjKguWdcuZwJxbjJbqbvF5/zXISEoKly5iGK+11eDRcX2Rp8yRpOhO84
LtSpC21QTkMSK8VA4O3e1tOW+DXaJb3QtzwocTb+wXTw74agwvcAQP9yDilgR1t6
fdqeGrW4Y25bDulXlsD+S6PzhVo+EVPq1T3RJQIDAQABAoIBADZfQ9Pxm8PnhVJF
ZuUfEzS+nnOtH9jMmEooQel6s3xa2NXXSRZfGZfHDpVsl0p72CloJhQASxCs9jMu
HzwfnyWqq37SuRTA2VmPvjhcwasJWTt+ygrztvikz52SUMIuInYV6oNwCLnY2Qaz
k3rrh7nqg2j684tsS4p6lItoCaArA5xcQwxn6librK/NzHzLaXN0zLufn4WYuPMc
3NTuZWY6EYqbyHbuiwrsZGin8JKw+bqfG6OEtt5GVJbzmacjQrVTEnEcJNO8Pe3H
bC/ZczFBb9Vsznfp0EICKf/OZVen7zSZ58+l3zg/+h0A8Z1D4jbWkXWDS3dYiAQU
g2C9x8kCgYEA/sllVEZXyCduyUvP7nVYPasBHKbIuS8G/cIKfwy+Wd1ZhKg4JgIy
5nhERYfOJeDwSoQUYxJSZoCgnByc4jx3kSX4oyTdKdT+yj1Sma38GONRm3Erxany
aZvw40cj5vCn8iGl6hsSpqWWiHWizEyO+XvctfMaFq5vOQxgjTF4Yw8CgYEA8Y6L
VlZxByVO7kQwZXdJ1zEtu1YzZyw95kiHmnxHdOqhstDV9mwtcDLD12CYR2LVweT/
ndTTU4U3q1z/Zuo1t1HYvTHU4ss/4GBCskO0JIKFPDr2KdfUiDGn6eMWNmoQ35zU
Z2zfi3BMtX2dbobX+dBDyh+ZJf21zKsyujp/eIsCgYEA6poU9IGE6KbuiumEv6RL
KRVhg7lLD8Dupg/azFu2llaLy+t9L/pMVgydiIxg1F4HxAVUJFlFiF6eBMEP7/0P
d5ZIGCikgJVAOoY2nY0nmN8PUJrnXC19KaNOLmhd9ZLYgcpb1HEzPkEwl9wBmC5S
ZASaGOuMtR/PB++Oo9PObx8CgYEAjwcdH+kdEeMYYmKD2YCRe1bGQlefJib/G9y0
VlfiI6tORUf8eOXC3d1hMqUiZZpzAVTrufOrkZeex9vP6oshdUOENzpLWGKKlvvI
Yi9OehPCelBbM5l1YZMtXoK0w1F4Xj9JUVgY4UKEWS5gynITbfrQON0O3HzmaaKw
7a33jlMCgYBXA5OvGOav/uMQ45SiYf8KhyNx2gQqQlvSeF9z8yDqmHr2epoAVoEu
nuz0bzQ7ARNqBkWLQ4bqzwy0aKXlcvbIMBaVXyfQTiwzwWAZsAWr6WHPlvWDP6mP
1vAUje5xEMtsIwj6UnkJ3OpPVVeJ56aKQIxg6QU2ROrWYDccx4gg0g==
-----END RSA PRIVATE KEY-----`

const selfsignedDummypub = `-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAOze8GNkMIMMMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xOTAxMjExMzUyMjFaFw0yOTAxMTgxMzUy
MjFaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAPBpdsfSrJqsHZmXfgXWWLQ1hzUV7T3f5MIkepblPuMO
+k/vZpZnj16gC4u2ksr2EAy3n30fPtnUhnHcL5H+mw1mTXjEV5e2dY7JPsRTfdZk
mDZNZ0ejdRLA0ytil5WWcmk5qbv+afEscAzt+LhOQGnhEnjN0VaqjoEXF7JJqf2r
cFDAvqi37lz6DUoyoLlnXLmcCcW4yW6m7xef81yEhKCpcuYhivtdXg0XF9kafMka
ToTvOC7UqQttUE5DEivFQODt3tbTlvg12iW90Lc8KHE2/sF08O+GoML3AED/cg4p
YEdben3anhq1uGNuWw7pV5bA/kuj84VaPhFT6tU90SUCAwEAAaNQME4wHQYDVR0O
BBYEFG3WDJMv5wa4QvWwxcJpNU/RTBp/MB8GA1UdIwQYMBaAFG3WDJMv5wa4QvWw
xcJpNU/RTBp/MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBALDYGYv6
0KoSAbQamSOZT6h2LBJj/AbGV+W9ffUDW6OuJ1C7sa7sDki2HQgz7vfS0BtrKY/q
tszfZqWKlx8AFbLhusMcv3gc6Dv4Onod90EaIKuRU+sElo1BEak5asY4oHru5GIK
QxGi8GkcwKSwnxSrkKQz8xXcL+P3daOmaAUQDo6JPqxYE4DNsQ3HRtkCj9kTUk8+
ppJAzXoBrutQz7e2daEXHUNc+1+KcD+se5cmvK2cJg6vk1vpgY1kjXdLQr1CySxJ
XgfLm2jJfzMF/L5RX5Vdnon6x4ufi7e/3fOThjlhLRXMOkhlb0E+wSYP0NvLA12E
rjs761ndZ9Qrb0s=
-----END CERTIFICATE-----`

const selfsignedPubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8Gl2x9KsmqwdmZd+BdZY
tDWHNRXtPd/kwiR6luU+4w76T+9mlmePXqALi7aSyvYQDLeffR8+2dSGcdwvkf6b
DWZNeMRXl7Z1jsk+xFN91mSYNk1nR6N1EsDTK2KXlZZyaTmpu/5p8SxwDO34uE5A
aeESeM3RVqqOgRcXskmp/atwUMC+qLfuXPoNSjKguWdcuZwJxbjJbqbvF5/zXISE
oKly5iGK+11eDRcX2Rp8yRpOhO84LtSpC21QTkMSK8VA4O3e1tOW+DXaJb3Qtzwo
cTb+wXTw74agwvcAQP9yDilgR1t6fdqeGrW4Y25bDulXlsD+S6PzhVo+EVPq1T3R
JQIDAQAB
-----END PUBLIC KEY-----`

const samplePrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsK8mIapI9HetPmwfptmHr2+oYW5YGKhzq4xEl6zAISChNVk9
DSMMLALfnlOaAK02yPqSiVSOYpyPjlEK/mRwETMSsRQvO/i+pO5aI9NeSfTo0HAc
NZ4nEwzsdrgrL3vrEagBxA2UgJM397CYhij/kJNK7Gec/jvlJAZxvr/k0SPV1iik
mYPCk0jHdbuMCEYMJ769pDvNwZr/JtV/uwG4H6BJUFtfKj3T0kjLQC4Y/uGTjKfm
l39iEuEop8Gi05WniEvbcpcI64XZpKasnuj7nYpjBW57oZzRW8viok9/+UEW0nE6
whKnyfhQ9TfJEicyoKFeLEOVbq9UdrhC2vwcaQIDAQABAoIBAQCQ/vwFDrEWZwx2
uNb033n5oGGHq72CZuOeOdukubFmvldt55Exsbxwdd88GJG+0meuYexV5V2AUcmB
2sJx6M0LYGWLiuwEhGs4AR9aXUD44pMZU5fi7KpWePmpqBRQwJo2ADGKyjY/mhGJ
JJTXLNgmtqn6/kEZZt/yQ5OfHe3TLv21nwhjfk9BV7MTEdA1Rr2BXyHD7h9MjgxK
9Q7vHtmhwwV67fQZFXwwN7kgrSSNFVNXDz2S+8IOcDondGTQFTau+S6x/0ziY6rg
ENjkzyoV9JO2WqnQyKb+rtpaFnDNdkuJHw6oVdeZOZwIL6CQ/RhmWDQPeDoc/NDK
aD6mzwaBAoGBAMDgCctycsilaxqZlHY5l6Qt1V0NH57+RNpzFzHidS8tEjbwHii7
FFQusTTcCFnFPpiQzKTTrjbPqEH5dcC99ip0BA9CGIzVWaWHFSVUPI0comhZhLVX
Y5+uX309f3LrSqNQmuRkwAOCeNlH5l9r0ncrDX9DkXp3x6uNtmITvOoRAoGBAOqC
jo/z7X9XJRhWwgW39AGIpvxrJgK9lIWsxfQma+NglyqvUA/Dzw9Ou870oJCVBZ39
CfZkLiJAUAk3F6lOdeEByzRy4A7NL94O8B3lQ6huOayksgfr8A4ScCWMyrTIgtCh
zATedKa69QGm/TAg6KJA8eP3K0snAPRt38cvKnTZAoGAJxkDQ0eG9x95L6I0Uybn
k3NrDfrMDynSAUpVSFp0kMSdLZ/NLUqHG21/pIx58OCoCLtJkJwMc7XykLUl5pVb
Yk20SPeIDHxvOLvCUJfb0mscjPSgjzYQztzFJJkjzcLelW6Qh33Y4p0/LCSEEZHE
zz1d9g9XXTEMu7z1XLpNkFECgYAB0kPDMHTOwWGDX+Ef5D7b6DDL0xU3fjtyElZz
P/0khfKGnVf012N7TfQ9dj7tAItLn9R8+mg1UeSNPcVMRlS6C6aFYMMGumc9xUXu
JYKyAzElex3628VAhroiQIaugsQpVKhd/VBQnzEZ8y8SOZ8062Y1jAzlB4eFXnkX
dfFReQKBgB8uf54s2HOI+Yx1YuQ5bYzOUp/J4PjyzPECjylsKe5J4pBntV1TvA9m
GRw/287a1mi9hDUlbOZSIVNJHAzxArCnJJnrW6C29NDFqWGIAgUS066KRZIyEuB3
SEtoekWeoBbByz3ehuKpuBK5St7Mz1MHyqb7YHQlTFj7oRQ9uHNK
-----END RSA PRIVATE KEY-----`

const samplePublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAsK8mIapI9HetPmwfptmHr2+oYW5YGKhzq4xEl6zAISChNVk9DSMM
LALfnlOaAK02yPqSiVSOYpyPjlEK/mRwETMSsRQvO/i+pO5aI9NeSfTo0HAcNZ4n
EwzsdrgrL3vrEagBxA2UgJM397CYhij/kJNK7Gec/jvlJAZxvr/k0SPV1iikmYPC
k0jHdbuMCEYMJ769pDvNwZr/JtV/uwG4H6BJUFtfKj3T0kjLQC4Y/uGTjKfml39i
EuEop8Gi05WniEvbcpcI64XZpKasnuj7nYpjBW57oZzRW8viok9/+UEW0nE6whKn
yfhQ9TfJEicyoKFeLEOVbq9UdrhC2vwcaQIDAQAB
-----END RSA PUBLIC KEY-----`

const signing_private = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC61SZkZGSX7eAy
3mcqGxPkZf+yMOYIAv2iQmdHW69ssqHiRHxrDZmAx20QP9GjjJ8tbyQk2QORJ9ej
sAOVhdwMpq70QacjLc2AK9kgAO3LHdD1lnuN9onoWkZZFDctJgMYxUnujZi25PtQ
vuJl+Pi1UHnsosBf8HYou7SvP1HxoAYacWVuCL+PZJnuJfyevFyRJMyvrdTA3Qyu
I3+Hgj/nXrOtm6MEhmuWuXOExbopUgtzNbCN3MWxRu/CCS63YjBXmtwIY8zc3BkO
kdJbjOf1u4/z23CZyyAVEBjIBEeOwJcmxDqD2Ay4igpLmMUlAdKOQ90ihliufF1w
JqnhCQ1DAgMBAAECggEAKf1HoJ5zgTXEAoq7ctodEWLfIaQdvsU1TadQ4Ne5SFup
SFoOAF1RF4E6gMFnEzPCfoqQ+/sN8yyaKT6gv5UTDIDVpy2uK5jaq6ivJqMuzkyI
LvnAEPrMqbzIPLLvZ6U4YvPMFuIZ5Vj3JoGQDkzzUISisk0toSJA3Ay7oftAJmZn
Sosy3tLp8/MNkUfsCavtHVzFRrFAJ1N4i+UNpUE7hQ9YE1fPL11cmKYf3I6FQHPp
U8MCLPrGhOztnK9bV4rI9Rd9xFIQCDv2mpFwSTeh8lSTM8pUdxOAT0XbhCHMjQGu
xz3Wlrv66J0ebH8opE3ztnsq+h1WciOLjtHrRglHWQKBgQDcCp44dnpumLKox0Ew
kQSdE5zP4WdGZq4aJwE17w/o4upldaV9qXx3QmAs6eCwvvJOdKQXr/uMsZzIcMxt
VyQGMmBytAPdmGj0+qZ4E7DaM9Gfyluy+xf89pZVkxdQQS/jhUAoL+DKmUshWh1k
+9ngYUzGaWxooKQcD7XhWTVkpwKBgQDZXT9Pn6pj0vXW3dWQ6gkuIXe+LNWwSWZy
MzE1RaoixBiV2Y06Gul87I5h3Q5n94UGDRNlqdfGuda1Uv5J65lRTmUrnTPnO3s+
Pm6G9HxjPH2czEtnLcK9AvEVpLrl1rq/0Otv+4XoYBM/v8NYuy2cntP9PPIXabyO
TBjM+xX6BQKBgHGCyLw34lDLVN7cazSymr6tL2fNz4jxzz6OgIFiIcLxzBkq54Q7
uomLJDIHNHH5DuaKJVxS3GFn/okoJ00AdwT7V+XUF2ppBTvbUaUAA2uM78aOjV93
SJimXEco6g3skte8FaylhkD9c1RxOFiv02V8zC5OlC4lMIOJVzo42uJhAoGBAJ/Q
QHFRily0yc298n0GpdNGBh1MJ5ziirEiVGa/nrTLCux6NKzpBoyz/IeVmTb1tNdb
G8zekGhrUKKmr5I359Tw18+2WGgFwrpj+q286guoeQ6k4jetXIXNuOXZ5RSByXKo
r8H4416T7PMtEfqWPJXv7Rs/CRwPwPO6nW1wmprlAoGAN0pBAuw4/Ywj94oZpDBH
JYM2qdSstLt/KcYeYYOQw/N8NIYK0pKgWzn5QnPn24xrBFddO3GYt56KacYTv2nP
h3Tytkib/u85ZlNif8TB259RDjPJ3zmAfjQAWvhzTPmyXjBOttPMLUsLlBPNEkRL
hFVMxtUuf8fxNOFewsjOAMg=
-----END PRIVATE KEY-----`

const signing_public_cert = `-----BEGIN CERTIFICATE-----
MIIFLTCCBBWgAwIBAgIEWcWi3DANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJH
QjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFBy
ZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwOTE4MTM0NzExWhcNMjAxMDE4
MTQxNzExWjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZ
BgNVBAsTEjAwMTU4MDAwMDEwNDFSYkFBSTEfMB0GA1UEAxMWZkp1VVU2ZE50MHp4
bkRlNTllRzBZTjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALrVJmRk
ZJft4DLeZyobE+Rl/7Iw5ggC/aJCZ0dbr2yyoeJEfGsNmYDHbRA/0aOMny1vJCTZ
A5En16OwA5WF3AymrvRBpyMtzYAr2SAA7csd0PWWe432iehaRlkUNy0mAxjFSe6N
mLbk+1C+4mX4+LVQeeyiwF/wdii7tK8/UfGgBhpxZW4Iv49kme4l/J68XJEkzK+t
1MDdDK4jf4eCP+des62bowSGa5a5c4TFuilSC3M1sI3cxbFG78IJLrdiMFea3Ahj
zNzcGQ6R0luM5/W7j/PbcJnLIBUQGMgER47AlybEOoPYDLiKCkuYxSUB0o5D3SKG
WK58XXAmqeEJDUMCAwEAAaOCAfkwggH1MA4GA1UdDwEB/wQEAwIGwDAVBgNVHSUE
DjAMBgorBgEEAYI3CgMMMIHgBgNVHSAEgdgwgdUwgdIGCysGAQQBqHWBBgFkMIHC
MCoGCCsGAQUFBwIBFh5odHRwOi8vb2IudHJ1c3Rpcy5jb20vcG9saWNpZXMwgZMG
CCsGAQUFBwICMIGGDIGDVXNlIG9mIHRoaXMgQ2VydGlmaWNhdGUgY29uc3RpdHV0
ZXMgYWNjZXB0YW5jZSBvZiB0aGUgT3BlbkJhbmtpbmcgUm9vdCBDQSBDZXJ0aWZp
Y2F0aW9uIFBvbGljaWVzIGFuZCBDZXJ0aWZpY2F0ZSBQcmFjdGljZSBTdGF0ZW1l
bnQwbQYIKwYBBQUHAQEEYTBfMCYGCCsGAQUFBzABhhpodHRwOi8vb2IudHJ1c3Rp
cy5jb20vb2NzcDA1BggrBgEFBQcwAoYpaHR0cDovL29iLnRydXN0aXMuY29tL29i
X3BwX2lzc3VpbmdjYS5jcnQwOgYDVR0fBDMwMTAvoC2gK4YpaHR0cDovL29iLnRy
dXN0aXMuY29tL29iX3BwX2lzc3VpbmdjYS5jcmwwHwYDVR0jBBgwFoAUUHORxiFy
03f0/gASBoFceXluP1AwHQYDVR0OBBYEFHricuOycGqmX586P9epTpowEnx9MA0G
CSqGSIb3DQEBCwUAA4IBAQBjptS/LoJpSs35r5O6UkU7KsB2Aeemb18VDUvCfbbX
XI8nONI+74YzXV2V/vHOJBZL5lUD7BhDGr03M76YIvRalbecGVVfcKUJbTadjRSL
LAvsjVcnuHLj67gr3gTGcqkTk6QHKldlP38m1ZSjdGTqzlEDrcE1f5tXm3lP+RVV
9lKhDsrhn7PqNROMkP+XnD2oxovTIPNCpvyPad+lCFuK/DT7ZrnRgw+bG3cwP82+
4FKu6rKdQIlt25v3a4MnAs5rlqS8yONhPGrc0W3LcL8UhuTnb/Zb8Z774JO18/9I
W5KE7L11WczOWn5kMtn4juswluDQ02i8163BdDrsE9NC
-----END CERTIFICATE-----`

const signing_public = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAutUmZGRkl+3gMt5nKhsT
5GX/sjDmCAL9okJnR1uvbLKh4kR8aw2ZgMdtED/Ro4yfLW8kJNkDkSfXo7ADlYXc
DKau9EGnIy3NgCvZIADtyx3Q9ZZ7jfaJ6FpGWRQ3LSYDGMVJ7o2YtuT7UL7iZfj4
tVB57KLAX/B2KLu0rz9R8aAGGnFlbgi/j2SZ7iX8nrxckSTMr63UwN0MriN/h4I/
516zrZujBIZrlrlzhMW6KVILczWwjdzFsUbvwgkut2IwV5rcCGPM3NwZDpHSW4zn
9buP89twmcsgFRAYyARHjsCXJsQ6g9gMuIoKS5jFJQHSjkPdIoZYrnxdcCap4QkN
QwIDAQAB
-----END PUBLIC KEY-----
`
