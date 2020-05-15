package signature

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/test"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"
)

const OzoneTestSignature = `yJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`
const OzoneResponseBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`
const domesticPayBody = `{"Data":{"ConsentId":"sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d","Initiation":{"InstructionIdentification":"SIDP01","EndToEndIdentification":"FRESCO.21302.GFX.20","InstructedAmount":{"Amount":"15.00","Currency":"GBP"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Identification":"20000319470104","Name":"Messers Simplex & Co"}}},"Risk":{}}`

var ctx model.Context = model.Context{
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

func verify(signingMethod jwt.SigningMethod, token string) bool {
	segments := strings.Split(token, ".")
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], test.LoadRSAPublicKeyFromDisk("test/sample_key.pub"))
	return err == nil
}

func TestJWSDetachedSignature313andBefore(t *testing.T) {
	model.EnableJWS()
	cert, _ := SigningCertFromContext(ctx)
	pubKey := cert.PublicKey()

	i := model.Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := model.TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	fmt.Println(sig)
	fmt.Println("Now validate this")
	validateSignature(sig, authentication.SigningMethodPS256, pubKey)
}

func TestJWSDetachedSignature315andAfter(t *testing.T) {
	model.EnableJWS()
	cert, _ := SigningCertFromContext(ctx)
	pubKey := cert.PublicKey()

	i := model.Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := model.TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	fmt.Println(sig)
	fmt.Println("Now validate this")
	validateSignature(sig, authentication.SigningMethodPS256, pubKey)
}

func validateSignature(token string, signingMethod jwt.SigningMethod, pubKey *rsa.PublicKey) {
	segments := strings.Split(token, ".")
	segments[1] = domesticPayBody
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], pubKey)
	if err != nil {
		logrus.Errorln("failed to validate signature" + err.Error())
	} else {
		logrus.Infoln("Succeeded to validate signature")
	}
}

func getSignatureParameters() (SignatureParameters, error) {
	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		return SignatureParameters{}, nil
	}
	params := SignatureParameters{
		cert:        cert,
		issuer:      "issuer",
		kid:         "mykeylookup",
		trustAnchor: "mytrustanchor",
		apiVersion:  "v3.1.3",
		alg:         authentication.SigningMethodPS256,
	}
	return params, nil
}

func PubPrivFromPEMs(publicKeyPem, privateKeyPem string) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
	if err != nil {
		return nil, nil, errors.Wrap(err, "error with public key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		return nil, nil, errors.Wrap(err, "error with private key")
	}
	return publicKey, privateKey, nil
}

func TestSimpleSig(t *testing.T) {
	body := SignatureBodyClaim{Body: `{"Data":"mybody"}`}
	myjson, err := body.MarshalJSON()
	if err != nil {
		logrus.Error(err)
		t.Fail()
	}
	logrus.Printf("Json output: %s ", myjson)

	jsonvalue, err := json.Marshal(body)
	if err != nil {
		logrus.Error(err)
		t.Fail()
	}

	logrus.Printf("JsonValue output: %s ", jsonvalue)

}

const abody = `{
	"Data": {
	   "DebtorAccount": {
		  "Identification": "70000170000002",
		  "Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
		  "SchemeName": "UK.OBIE.SortCodeAccountNumber"
	   },
	   "ExpirationDateTime": "2021-01-01T00:00:00+01:00",
	   "ConsentId": "fcc-913556c0-0f03-43fe-8792-262beb7d2510",
	   "CreationDateTime": "2020-04-17T14:19:59.694Z",
	   "Status": "AwaitingAuthorisation",
	   "StatusUpdateDateTime": "2020-04-17T14:19:59.694Z"
	},
	"Links": {
	   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-913556c0-0f03-43fe-8792-262beb7d2510"
	},
	"Meta": {}
 }`

func minifiyJSONBody(body string) (string, error) {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minjson.Minify)
	minifiedBody, err := m.String("application/json", body)
	if err != nil {
		logrus.Error(err, `minifyJSONBody failed`)
		return "", err
	}
	return minifiedBody, nil
}

// Sign both Ways
func TestSignAndValidateB64True(t *testing.T) { // New Way
	parms, _ := getSignatureParameters()

	minifiedBody, err := minifiyJSONBody(abody)
	if err != nil {
		logrus.Error(err)
		t.Fail()
	}
	_ = minifiedBody

	tok := jwt.Token{
		Header: map[string]interface{}{
			"typ":                           "JOSE",
			"kid":                           parms.kid,
			"cty":                           "application/json",
			"http://openbanking.org.uk/iat": time.Now().Unix(),
			"http://openbanking.org.uk/iss": parms.issuer,      //ASPSP ORGID or TTP ORGID/SSAID
			"http://openbanking.org.uk/tan": parms.trustAnchor, //Trust anchor
			"alg":                           parms.alg.Alg(),
			"crit": []string{
				"http://openbanking.org.uk/iat",
				"http://openbanking.org.uk/iss",
				"http://openbanking.org.uk/tan",
			},
		},
		Method: authentication.SigningMethodPS256,
	}

	cert, _ := SigningCertFromContext(ctx)
	privKey := cert.PrivateKey()

	signature, err := tok.SignedString(privKey)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Printf("Sig: %s", signature)

}

func TestSignAndValidateB64False(t *testing.T) { // Old Way
	parms, _ := getSignatureParameters()
	token := authentication.GetSignatureToken313Minus(parms.kid, parms.issuer, parms.trustAnchor, parms.alg)
	var b64Encoded = false // explicitly set to false, just to be clear
	cert, _ := SigningCertFromContext(ctx)
	privKey := cert.PrivateKey()

	minifiedBody, err := minifiyJSONBody(abody)
	if err != nil {
		logrus.Error(err)
		t.Fail()
	}

	tokenString, err := authentication.SignedString(&token, privKey, minifiedBody, b64Encoded) // sign the token
	if err != nil {
		t.Fail()
	}
	detachedJWS := authentication.SplitJWSWithBody(tokenString) // remove the body from the signature string to form the detached signature
	logrus.Infoln(detachedJWS)

}

func validateB64True(token string) {

}

func validateB64False(token string) {

}

// Get Server Cert ... from JWKS !!! Drived from .wellKnownEndpoint
func TestGetServerPubCertFromJwks(t *testing.T) {

}

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
-----END CERTIFICATE-----
`
