package authentication

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

var hsbcResponseSignature = "eyJodHRwOlwvXC9vcGVuYmFua2luZy5vcmcudWtcL2lhdCI6MTYwMTkyMjI5OSwiaHR0cDpcL1wvb3BlbmJhbmtpbmcub3JnLnVrXC90YW4iOiJzMy1ldS13ZXN0LTEuYW1hem9uYXdzLmNvbSIsImNyaXQiOlsiaHR0cDpcL1wvb3BlbmJhbmtpbmcub3JnLnVrXC9pYXQiLCJodHRwOlwvXC9vcGVuYmFua2luZy5vcmcudWtcL3RhbiIsImh0dHA6XC9cL29wZW5iYW5raW5nLm9yZy51a1wvaXNzIl0sImtpZCI6ImV4dGVybmFsXzIiLCJ0eXAiOiJKT1NFIiwiaHR0cDpcL1wvb3BlbmJhbmtpbmcub3JnLnVrXC9pc3MiOiJQb3N0YWxDb2RlPUIxIDFIUSwyLjUuNC45Nz1QU0RHQi1GQ0EtNzY1MTEyLENOPUhTQkMsU1RSRUVUPUJpcm1pbmdoYW0sTD1CaXJtaW5naGFtLE9VPTEgQ2VudGVuYXJ5IFNxdWFyZSxPPUhTQkMgVUssQz1VSyIsImFsZyI6IlBTMjU2In0..g3jvSLnCLo2x8E7LEsjLKjv6BVwctNBc3voHk6EhJ6v2gIuL5CYSIh4F0cJLGNEkz7jXNkXTilcSCeAYSaCkdumk6CosK-tdNj_AQXe0Ma1gQURJi5wfeNA_7uLAnSXW4nFzSe1wGjH4vUEf8nd72K5R-XGr3EOB41aYj37ON521c496IVQCDzsJ2aiS7KG4l-6-_IOIVto1utIaZfTJis2t1PDNHusFEOKq9tFCwVGz_cSEyhlBSl-blc6wik6Nket59UP3itUop1xNdaUecCA3-_CaqjWynvoA6ZH26h0tXtxczgk9BqKxweSn3VO7PEPRWD6_-GnBb6wSCev6VA"

/*
hsbcResponseSignature Header:
{
    "alg": "PS256",
    "crit": [
        "http://openbanking.org.uk/iat",
        "http://openbanking.org.uk/tan",
        "http://openbanking.org.uk/iss"
    ],
    "http://openbanking.org.uk/iat": 1601922299,
    "http://openbanking.org.uk/iss": "PostalCode=B1 1HQ,2.5.4.97=PSDGB-FCA-765112,CN=HSBC,STREET=Birmingham,L=Birmingham,OU=1 Centenary Square,O=HSBC UK,C=UK",
    "http://openbanking.org.uk/tan": "s3-eu-west-1.amazonaws.com",
    "kid": "external_2",
    "typ": "JOSE"
}
Claims:
{}
*/

// diabled - https://openbanking.atlassian.net/browse/REFAPP-1304
// var hsbcBody = `{"Data":{"ConsentId":"50122752-4061-4a03-8383-0b909a91d86b","Status":"AwaitingAuthorisation","StatusUpdateDateTime":"2020-10-05T18:24:59+00:00","CreationDateTime":"2020-10-05T18:24:59+00:00","Permission":"Create","Initiation":{"Frequency":"IntrvlWkDay:01:01","Reference":"JKHSJKHSKHK76799","FirstPaymentDateTime":"2020-10-08T00:00:00+01:00","FirstPaymentAmount":{"Amount":"0.01","Currency":"GBP"},"CreditorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"40179070001015","Name":"Business account"}}},"Links":{"Self":"https://api.ob.business.hsbc.co.uk/obie/open-banking/v3.1/pisp/domestic-standing-order-consents/50122752-4061-4a03-8383-0b909a91d86b"},"Meta":{"TotalPages":1},"Risk":{}}`

// var hsbc_jwks_uri = "https://ob.business.hsbc.co.uk/jwks/public.jwks"

// func TestHSBCTrustAnchor(t *testing.T) {

// 	encodedResult := make([]byte, base64.RawURLEncoding.EncodedLen(len(hsbcBody)))
// 	base64.RawURLEncoding.Encode(encodedResult, []byte(hsbcBody))

// 	head, _, sig := payloadSplit(hsbcResponseSignature)
// 	signingString := head + "." + string(encodedResult[:]) + "."

// 	valid := verifyHSBCSig(t, SigningMethodPS256, signingString, sig, true)
// 	fmt.Printf("Signature is valid: %t ", valid)
// 	assert.True(t, valid, "Signature verification with public key failed")
// }

// func verifyHSBCSig(t *testing.T, signingMethod jwt.SigningMethod, signingString, signature string, b64 bool) bool {
// 	cert, err := getCertForKid("external_2", hsbc_jwks_uri)
// 	assert.Nil(t, err)
// 	verified, err := JWSVerify(signingString+signature, jwa.PS256, cert.PublicKey, b64)
// 	if err != nil {
// 		logrus.Errorf("failed to verify message: %v", err)
// 		return false
// 	}
// 	logrus.Tracef("signed message verified! -> %s", verified)
// 	assert.Nil(t, err)
// 	return err == nil
// }

var hsbcTanTestList = []struct {
	tan      string
	expected bool
}{
	{"https://ob.hsbc.co.uk/jwks/public.jwks", true},
	{"https://ob.firstdirect.com/jwks/public.jwks", true},
	{"https://ob.mandsbank.com/jwks/public.jwks", true},
	{"https://ob.business.hsbc.co.uk/jwks/public.jwks", true},
	{"https://ob.hsbckinetic.co.uk/jwks/public.jwks", true},
	{"openbanking.org.uk", false},
	{"hsbc.co.uk", false},
	{"", false},
	{":", false},
}

func TestIsHSBCTrustAnchor(t *testing.T) {
	for _, tt := range hsbcTanTestList {
		actual := isHSBCTrustAnchor(tt.tan)
		if actual != tt.expected {
			t.Errorf("isHSBCTrustAnchor(%s): expected %t, actual %t", tt.tan, tt.expected, actual)
		}
	}

}

// TestGetJwksErrors
// performs unit tests for the getJwks function to verify the returned JWKS and errors.
func TestGetJwksErrors(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		wantNonEmptyJWKS bool
		wantErr          string
	}{
		{"valid", "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/0015800001041RbAAI.jwks", true, ""},
		{"not decodable", "https://www.google.com", false, "GetJwks: decoding error"},
		{"bad url", "www", false, "GetJwkss error retrieving url:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwks, err := getJwks(tt.url)

			if tt.wantNonEmptyJWKS == true && jwks.Keys == nil {
				t.Errorf("GetJwks(%s) returned empty JWKS", tt.url)
				t.Errorf("Error logs: %s", err.Error())
				return
			}

			if tt.wantNonEmptyJWKS == false && jwks.Keys != nil {
				t.Errorf("GetJwks(%s) returned non-empty JWKS", tt.url)
				t.Errorf("Error logs: %s", err.Error())
				return
			}

			if tt.wantErr == "" && err != nil {
				t.Errorf("GetJwks(%s) unexpected error: %s", tt.url, err)
				return
			}

			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("GetJwks(%s) error = %v, wantErr %s", tt.url, err, tt.wantErr)
				return
			}
		})
	}
}

// TestGetJwksContentTypeLogs
// performs unit tests for the getJwksContentType function to verify the content type returned from the JWKS url.
func TestGetJwksContentTypeLogs(t *testing.T) {
	tests := []struct {
		name                 string
		contentType          string
		expectedLog          string
		expectRecommendation bool
	}{
		{
			"application/jwk+json",
			"application/jwk+json",
			"Acceptable JWKS content type found: application/jwk+json",
			true,
		},
		{
			"application/jwk-set+json",
			"application/jwk-set+json",
			"",
			false,
		},
		{
			"application/json",
			"application/json",
			"Acceptable JWKS content type found: application/json",
			true,
		},
		{
			"text/plain",
			"text/plain",
			"Unexpected JWKS content type found: text/plain",
			true,
		},
		{
			"text/html",
			"text/html",
			"Unexpected JWKS content type found: text/html",
			true,
		},
	}

	var logBuffer bytes.Buffer
	prevOut := logrus.StandardLogger().Out
	prevLevel := logrus.GetLevel()
	logrus.SetOutput(&logBuffer)
	logrus.SetLevel(logrus.TraceLevel) // capture info and warn (and more)
	t.Cleanup(func() {
		logrus.SetOutput(prevOut)
		logrus.SetLevel(prevLevel)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				_, _ = w.Write([]byte(`{"Keys":[]}`))
			}))
			t.Cleanup(server.Close)

			logBuffer.Reset()

			_, _ = getJwks(server.URL)

			got := logBuffer.String()
			if tt.expectedLog != "" && !strings.Contains(got, tt.expectedLog) {
				t.Fatalf("expected %q in logs, got:\n%s", tt.expectedLog, got)
			}
			if tt.expectedLog == "" && strings.Contains(got, "JWKS content type found") {
				t.Fatalf("expected no content type classification log, got:\n%s", got)
			}
			recommendationLogFound := strings.Contains(got, "The recommended JWKS content type is")
			if tt.expectRecommendation != recommendationLogFound {
				t.Fatalf("expected recommendation log found to be %t, got logs:\n%s", tt.expectRecommendation, got)
			}
		})
	}
}

// TestGetJwksStatusCode
// performs unit tests for the getJwksStatusCode function to verify the status code returned from the JWKS url.
func TestGetJwksStatusCode(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		exptectedErr string
	}{
		{
			"200",
			"https://keystore.openbankingtest.org.uk/0015800001041RbAAI/0015800001041RbAAI.jwks",
			"",
		},
		{
			"404",
			"https://www.google.com/404",
			"GetJwks: bad status code (404) from JWKS url (https://www.google.com/404)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getJwks(tt.url)

			if tt.exptectedErr == "" && err != nil {
				t.Errorf("GetJwks(%s) unexpected error: %s", tt.url, err)
			}

			if tt.exptectedErr != "" && !strings.Contains(err.Error(), tt.exptectedErr) {
				t.Errorf("GetJwks(%s) error = %v, wantErr %s", tt.url, err, tt.exptectedErr)
			}
		})
	}
}

// TestInArray
// performs unit tests for the inArray function to verify string existence in an array with case insensitivity.
func TestInArray(t *testing.T) {
	arr := []string{"application/jwk+json", "application/json"}

	tests := []struct {
		name     string
		target   string
		expected bool
	}{
		{"exact match jwk+json", "application/jwk+json", true},
		{"mixed case jwk+json", "Application/JWK+JSON", true},
		{"upper case json", "APPLICATION/JSON", true},
		{"not present", "text/plain", false},
		{"empty target", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inArray(tt.target, arr)
			if got != tt.expected {
				t.Errorf("inArray(%q, %v) = %v, want %v", tt.target, arr, got, tt.expected)
			}
		})
	}

}
