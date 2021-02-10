package authentication

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

var hsbcBody = `{"Data":{"ConsentId":"50122752-4061-4a03-8383-0b909a91d86b","Status":"AwaitingAuthorisation","StatusUpdateDateTime":"2020-10-05T18:24:59+00:00","CreationDateTime":"2020-10-05T18:24:59+00:00","Permission":"Create","Initiation":{"Frequency":"IntrvlWkDay:01:01","Reference":"JKHSJKHSKHK76799","FirstPaymentDateTime":"2020-10-08T00:00:00+01:00","FirstPaymentAmount":{"Amount":"0.01","Currency":"GBP"},"CreditorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"40179070001015","Name":"Business account"}}},"Links":{"Self":"https://api.ob.business.hsbc.co.uk/obie/open-banking/v3.1/pisp/domestic-standing-order-consents/50122752-4061-4a03-8383-0b909a91d86b"},"Meta":{"TotalPages":1},"Risk":{}}`
var hsbc_jwks_uri = "https://ob.business.hsbc.co.uk/jwks/public.jwks"

func TestHSBCTrustAnchor(t *testing.T) {

	encodedResult := make([]byte, base64.RawURLEncoding.EncodedLen(len(hsbcBody)))
	base64.RawURLEncoding.Encode(encodedResult, []byte(hsbcBody))

	head, _, sig := payloadSplit(hsbcResponseSignature)
	signingString := head + "." + string(encodedResult[:]) + "."

	valid := verifyHSBCSig(t, SigningMethodPS256, signingString, sig, true)
	fmt.Printf("Signature is valid: %t ", valid)
	assert.True(t, valid, "Signature verification with public key failed")
}

func verifyHSBCSig(t *testing.T, signingMethod jwt.SigningMethod, signingString, signature string, b64 bool) bool {
	cert, err := getCertForKid("external_2", hsbc_jwks_uri)
	assert.Nil(t, err)
	verified, err := JWSVerify(signingString+signature, jwa.PS256, cert.PublicKey, b64)
	if err != nil {
		logrus.Errorf("failed to verify message: %v", err)
		return false
	}
	logrus.Tracef("signed message verified! -> %s", verified)
	assert.Nil(t, err)
	return err == nil
}

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
