package model

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	jwt "github.com/dgrijalva/jwt-go"
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
	ctx := Context{
		"consent_id":        "aac-fee2b8eb-ce1b-48f1-af7f-dc8f576d53dc",
		"xchange_code":      "10e9d80b-10d4-4abd-9fe0-15789cc512b5",
		"baseurl":           "https://matls-sso.openbankingtest.org.uk",
		"access_token":      "18d5a754-0b76-4a8f-9c68-dc5caaf812e2",
		"client_id":         "12312",
		"scope":             "AuthoritiesReadAccess ASPSPReadAccess TPPReadAll",
		"PrivateSigningKey": KeySigning,
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
	for k, v := range res.FormData {
		fmt.Printf("%27.27s:\t%s\n", k, v)
	}
	assert.True(t, len(res.FormData["client_assrtion"]) > 0, "client assertion should be populated with signed jwt")
}

func TestCertStuff(t *testing.T) {

	cert, err := authentication.NewCertificate(signpub, signkey)
	require.NoError(t, err)
	require.NotNil(t, cert)

	//privateKey := []byte(KeySigning)
	alg := jwt.GetSigningMethod("RS256")
	if alg == nil {
		fmt.Printf("Couldn't find signing method: %v\n", alg)
	}
	// key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signkey))
	// if err != nil {
	// 	fmt.Println("error : ", err)
	// }
	claims := jwt.MapClaims{}
	claims["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(alg, claims) // create new token
	token.Header["kid"] = "mykeyid"

	prikey := cert.PrivateKey()
	tokenString, err := token.SignedString(prikey) // sign the token - get as encoded string
	if err != nil {
		fmt.Println("error signing jwt: ", err)
	}
	fmt.Println(tokenString)

}

// KeySigning -
var KeySigning = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDAfrw8ClTBLkzD
lCEQlrpZgS5ERl5bihQUiiNy6JxRtozwF/sQqp4T4IoqMR4FD37nmHjtIK/SPLnN
Bwn/PzBPpHm7U1kKSjoAM/FnS6O+2FrrdkTDmrwyg04J29EXH2ckmIAAUsEtUMcN
GBBCga2Qa2+hq2Weom/imrmRQDe4gSljVZrL2AgeZnO/6jHdQSaslMsiRpIWOwvV
ZubdQ2vNocQdgElKu9afWRZya0vF+HufqpqvZ52kPyq9RHe/YE+OY72MmJF+bymi
r3ZYCyy9UW3wmM8soYjM+7uly7Z6hcDSy2Kpu+plHn/Qv9gd4sBXznA4iMeaGUcP
TJb7WkThAgMBAAECggEAR9UDWURhrFUiwDkevZoBoDTclw3LWE2GgMOrxs2Wx8df
gJjyT53br38zD9uUYD8QFEyJk7OG6OVQUHo3+NATrySpaIYJzBU236yCgRFw4V7L
TuKrdnLfl9n33SXyOLa3PqjJ21UGUWq7XN+F8cuCgUoWNjZHjZMAPYePh+x23ppu
Fb96S5ZP6nPPEbN1f8G2n4Ea1tfYpFw4Iq6gPea1OMGII3xl5ngkMeOaDNwqYhYb
ixBs2kI0eAIVVDcPx8X5CbAcWeCznGc2zW4aiAt3MLPjWqq9hlypoig+8GjRCGzT
Ut/lBJhMRgSe1YMuhBLeGmbr5cNWzUcENAvJ0CLCoQKBgQDx3Z/6V30HP06Kb/08
Nbjv1JIRDH3O+H5iF5I979HYeMzBrBcTw4eDdctfIUDqubpEOyVvhl11S/V2kclA
Pibf+0wTh9Tae1wGlrxlhSycrrzu7xG+oBazsj3hX70CAa9o5NWdvh40bgWOr7lA
81lO2Isgrrti9uVqxhxQwA4I4wKBgQDLvoMhN82xjz+BagbYJLwmwKZAPZ4PHC3A
cnnE2ZuU2CcmvA/CBIDWXXnT87pXBF26b94ajR8bQgcSph4jKujeCLz96FKwbTr2
8xxEbaVf8f3ilzHlk+OD+HirieQmGOQfJ8rPgjZykgLd7gFwpwwO9gpiP10MCCKC
pk9deIyaawKBgHBE9N6Kv+GeVEHUjBLnyQmifY7mYnuxQ1EbKeoQKTM3l6wKyseE
yqGOCzIESJLsVXcYkV78WuN4t98q+uUUNI1ho8WpFne4LVZtn9PsBnJQdije1jjL
LN6KzUiRXTXSPG8PUc0gE/s4WuIJ1Y89pmYABEzObvMYMhPnE/uzupALAoGAGc9C
cTzOc8W/t7ckstDEfOw+ozirAyMAsLZPsp4WVV6kZwW/wUYsw/sHadAgNNG6xdlR
+28RF7TfjH86ph3Tbf0RY+DASNUteQcG96wkHOlczg11Jq37TkZ1ktVe72yLyV6T
FIJcP1s7vb1etVST9Hk6i4OXV+TX6lEDEMYqmY0CgYEAvmfQ7D9ygrK62ysv5zCD
CFLMJDQVz/msm49UrBSK8wlLjCVPEsJe1/l0dESufES3NsbEXRZTqgQAZH15T0Mu
QrrRwdoLowUS4kEFKYH/emkjJKZQlCu0ccPJdiwZcuhePCyDo9uwSboSIhGoyvvb
12X8lr0o/RRVim7kHI9PI44=
-----END PRIVATE KEY-----
`

var signkey = "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCrFv8OQRO8f7xS\n19fvpZ5rHFYU8odwFQWZ/78aTa+F0yhfFce8yTydgHi2UsDvY3ehUopWXmXOst+B\n00/0m+QTqNSIi9+2UymOvOOJDyrfJS572Mf5LLx+WOjNeMFtxGPElYY2e94icsmm\ngTWtwsaMyGrW7I2sdjHYcH/+qQhJNblOZLFcJBCkL9ORz48Y4wqZ+HGyk2TC4fWJ\nirt403sWmC8XQt0WJFedcDqafZxWkJ0sJK3JoIboXftIkMj5z5Yryd1+Bhrxav2W\naqnzLHbXYekfT++l5rEgkI3abeJeBzZmKlV/juREV+Hwu4gatTgW4YL8IopI+eoc\nIQUw3RlvAgMBAAECggEBAJckcZ3+D5lunsfwtmqXPSQSnFlVCCET8Sbir8hk6LKo\nn/mgHBvDCzF41Sr8YEUa8gwqBtvV+Mppatod+3x0W0Ci3V7jcnZ3cTcP11K1e4I2\nLqJqF/8gbkSP9tnN29NEs35vOWnYc5yrG0lkzC786rpkMz47K803fUFf4TLv0Moa\noUfPQQljI0T2BfATPryElmN3lq6BTp17KcUVlsnCrAFLBtrBFtsZpfQPKuq+DMH8\n/z1brPTo1ZBCHTMd2+bo7NiEvuG89KO8zo8oLOqe8F5xPnQDsU2Uq2AaJUyZsKsr\neI76ttI7EV2eHlyHwev9oN7gHEcpEGGz5HLbAnrIi4ECgYEA38s2ISatEVMwz/4E\nRyEkLplkp6XWrioX3E/FMyvPFf1JygHn8rhcF6e5MBL4ApLr1Gdz4WHaQbdhVRLW\nccyavG0Kzv9vvzSNSNoPgRAaQ+QgHiThyfEyewID8fUn6CQ5gsADTqDeXHgmbw5v\nc+ROrPO4/N5C/xFdkXyXi1clcYkCgYEAw7YeJWmfaLHF6GsuA7uYbjy4vs8v1Q8b\nWvwDB8Fj8wlJQ9eFo9AAfU7eoQA/kBJXswnVXV1QSa3Sh8PqvHdLDHb3RKiPHSiS\nofhbCRvmUoPkqt74gYfS10LajIN4Jkev5oqC1N1oYy7IOsT7C2NKt92AxHvfgq4P\nADG9PVXvzTcCgYEAuQTiXYoCL36drneNxdiqdzQuOUQsNpVqYKQ6ntGrRbzAUpg8\n0TiGOrBZtFsaW9ZnzpUxArbJoOchOxp13GOR0hI8i2I3Wtbxr7dIdiV/8X0a6JEJ\nctFMMNI7vMA4G/5G5cglc84fyEc1Tz+Z+TBZszdUSwreTM5oky10hKiptjECgYEA\nvhLEqm8va32kCPr3AJcUDpQYlPAhs1ntpmq1ArY2vRYaurG5UAQ2RXzwyQq1sNWv\nqOl2+CslS7lui36iHpH5KEzuDxdpjtcVugq7V1hqU19XGQBd92cTRQ7ftLIGYZ8j\n3dJOCDBULmeD/VfLvR6ctX+BjNIFnCQx221zLfulXvcCgYBTMw8aVGKLPcwsxma9\nLC72wxhJJUyfhUaENKDDfzsVlQhDRAbqDkAwOKiNOFwI49rIIl1lk7hbfUki3HFh\nT8y309FOGeoxKsaNxM/kgChfWTimu74Dg9TY9bnfYSBgE2kQjBnHVwmB+z52cHPZ\nEw5MsQIXT0BRKZCkg9oGdAQ+JA==\n-----END PRIVATE KEY-----\n"
var signpub = "-----BEGIN CERTIFICATE-----\nMIIFLTCCBBWgAwIBAgIEWcVcgTANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJH\nQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFBy\nZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTgxMTI5MTI0NzU0WhcNMTkxMjI5\nMTMxNzU0WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZ\nBgNVBAsTEjAwMTU4MDAwMDEwNDFSZEFBSTEfMB0GA1UEAxMWN2w0WWR3VmdVRVR4\nVjJQakxZaDNzYzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKsW/w5B\nE7x/vFLX1++lnmscVhTyh3AVBZn/vxpNr4XTKF8Vx7zJPJ2AeLZSwO9jd6FSilZe\nZc6y34HTT/Sb5BOo1IiL37ZTKY6844kPKt8lLnvYx/ksvH5Y6M14wW3EY8SVhjZ7\n3iJyyaaBNa3CxozIatbsjax2Mdhwf/6pCEk1uU5ksVwkEKQv05HPjxjjCpn4cbKT\nZMLh9YmKu3jTexaYLxdC3RYkV51wOpp9nFaQnSwkrcmghuhd+0iQyPnPlivJ3X4G\nGvFq/ZZqqfMsdtdh6R9P76XmsSCQjdpt4l4HNmYqVX+O5ERX4fC7iBq1OBbhgvwi\nikj56hwhBTDdGW8CAwEAAaOCAfkwggH1MA4GA1UdDwEB/wQEAwIGwDAVBgNVHSUE\nDjAMBgorBgEEAYI3CgMMMIHgBgNVHSAEgdgwgdUwgdIGCysGAQQBqHWBBgFkMIHC\nMCoGCCsGAQUFBwIBFh5odHRwOi8vb2IudHJ1c3Rpcy5jb20vcG9saWNpZXMwgZMG\nCCsGAQUFBwICMIGGDIGDVXNlIG9mIHRoaXMgQ2VydGlmaWNhdGUgY29uc3RpdHV0\nZXMgYWNjZXB0YW5jZSBvZiB0aGUgT3BlbkJhbmtpbmcgUm9vdCBDQSBDZXJ0aWZp\nY2F0aW9uIFBvbGljaWVzIGFuZCBDZXJ0aWZpY2F0ZSBQcmFjdGljZSBTdGF0ZW1l\nbnQwbQYIKwYBBQUHAQEEYTBfMCYGCCsGAQUFBzABhhpodHRwOi8vb2IudHJ1c3Rp\ncy5jb20vb2NzcDA1BggrBgEFBQcwAoYpaHR0cDovL29iLnRydXN0aXMuY29tL29i\nX3BwX2lzc3VpbmdjYS5jcnQwOgYDVR0fBDMwMTAvoC2gK4YpaHR0cDovL29iLnRy\ndXN0aXMuY29tL29iX3BwX2lzc3VpbmdjYS5jcmwwHwYDVR0jBBgwFoAUUHORxiFy\n03f0/gASBoFceXluP1AwHQYDVR0OBBYEFBa6IGX7vYGp8TzyKXpeuWmK6S7YMA0G\nCSqGSIb3DQEBCwUAA4IBAQAsRlvcU19yW5Yh6SyeeUT/dbwFyh2fyclbbKToeawt\nkCxh4Icbl5hGocjaYuKyKOvTx2Y3FtkUYL+gnPeCMkNiZTN/08Bp+fqqk/a/6F65\nNP6AZIKzqkk5D2Dgl2gOsyS7SOCIMFTIkwm/aya4Bq3h9DFs8DZSyviujxux/LT8\nZTFlfnMhKWHLzgyANj86Nni8E/ParrJvrf0lieH1ObdMZfcm1Z/JwQSfYN9emyDr\nfjtMZdIbd/u9vN4nJs744MfJvqJP7nTXnGczMRBLIHNSqacQqmQ+RtkxV4l+YoPf\nBJD1GWVYohnEWIe05b1dQbV1s/x2k0kUdWjRFNHkwodc\n-----END CERTIFICATE-----\n"
