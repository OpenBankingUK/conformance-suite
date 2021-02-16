package model

import (
	"crypto"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"

	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
)

/*
	High level signature tests, using Testcases, Input
	Specifically for Waiver007 testing - B64=true, B64=false
*/

const OzoneTestSignature = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`
const OzoneResponseBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`
const domesticPayBody = `{"Data":{"ConsentId":"sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d","Initiation":{"InstructionIdentification":"SIDP01","EndToEndIdentification":"FRESCO.21302.GFX.20","InstructedAmount":{"Amount":"15.00","Currency":"GBP"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Identification":"20000319470104","Name":"Messers Simplex & Co"}}},"Risk":{}}`

var ctx Context = Context{
	"signingPrivate":            selfsignedDummySigkey,
	"signingPublic":             selfsignedDummySigpub,
	"initiation":                "{\"InstructionIdentification\":\"SIDP01\",\"EndToEndIdentification\":\"FRESCO.21302.GFX.20\",\"InstructedAmount\":{\"Amount\":\"15.00\",\"Currency\":\"GBP\"},\"CreditorAccount\":{\"SchemeName\":\"SortCodeAccountNumber\",\"Identification\":\"20000319470104\",\"Name\":\"Messers Simplex & Co\"}}",
	"consent_id":                "sdp-1-b5bbdb18-eeb1-4c11-919d-9a237c8f1c7d",
	"domestic_payment_template": "{\"Data\": {\"ConsentId\": \"$consent_id\",\"Initiation\":$initiation },\"Risk\":{}}",
	"authorisation_endpoint":    "https://example.com/authorisation",
	"api-version":               "v3.0",
	"nonOBDirectoryTPP":         false,
	"requestObjectSigningAlg":   "PS256",
	"apiversions":               []interface{}{"payments_v3.1.3"},
}

// Create and validate b64=false signature
func TestSimpleb64falseSignature(t *testing.T) {
	ctx.PutStringSlice("apiversions", []string{"payments_v3.1.3"})
	ctx.PutString("tpp_signature_kid", "x")
	ctx.PutString("tpp_signature_issuer", "x/x")
	ctx.PutString("tpp_signature_tan", "openbanking.org.uk")
	cert, _ := authentication.SigningCertFromContext(ctx)
	pubKey := cert.PublicKey()
	_ = pubKey
	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	logrus.Infoln(sig)
	validatedOK, err := validateSignatureTest(sig, domesticPayBody, authentication.SigningMethodPS256, pubKey)
	assert.True(t, validatedOK)
}

// create and validate b64=true signature - apiversion 3.1.4
func TestSimpleb64trueSignature(t *testing.T) {
	ctx.PutStringSlice("apiversions", []string{"payments_v3.1.4"})
	ctx.PutString("tpp_signature_kid", "x")
	ctx.PutString("tpp_signature_issuer", "x/x")
	ctx.PutString("tpp_signature_tan", "openbanking.org.uk")
	cert, _ := authentication.SigningCertFromContext(ctx)
	pubKey := cert.PublicKey()
	_ = pubKey
	i := Input{JwsSig: true, Method: "POST", Endpoint: "https://google.com", RequestBody: "$domestic_payment_template"}
	tc := TestCase{Input: i}
	req, err := tc.Prepare(&ctx)
	assert.Nil(t, err)
	sig := req.Header.Get("x-jws-signature")
	assert.NotEmpty(t, sig)
	logrus.Infoln(sig)
	encodedBody := jwt.EncodeSegment([]byte(domesticPayBody))
	validatedOK, err := validateSignatureTest(sig, encodedBody, authentication.SigningMethodPS256, pubKey)
	assert.True(t, validatedOK)
}

// Test using ozone server certificate
func TestOzone314SignatureString(t *testing.T) {
	signingMethod := jwt.SigningMethodPS256.SigningMethodRSA

	segments := strings.Split(OzoneTestSignature, ".")
	segments[1] = jwt.EncodeSegment([]byte(OzoneResponseBody))
	logrus.Printf(strings.Join(segments, "."))

	pubkey, err := LoadRSAPublicKeyFromDisk("../../../certs/waiver007/DKePOLAOiXLwYhMfLS8aS6YU-d0.pem")
	if err != nil {
		fmt.Println("WARNING - keyfile not found")
		return
	}

	modulus := pubkey.N.Bytes()
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, err := authentication.CalcKid(modulusBase64)
	logrus.Infof("kid: " + kid)
	err = signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], pubkey)
	logrus.Error(err)
}

func LoadRSAPublicKeyFromDisk(location string) (*rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	return key, err
}

// GetPem -
func GetPem(url string) (*x509.Certificate, error) {
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := resty.R().Get(url)
	if err != nil {
		fmt.Printf("error accessing %s, %s\n", url, err)
		return &x509.Certificate{}, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("Cannot access certificate: http code: " + resp.Status())
	}

	block, _ := pem.Decode(resp.Body())
	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		fmt.Println("Cannot parse PEM - omg error: ", err)
		return nil, err
	}

	return cert, nil

}

func validateSignatureTest(token, body string, signingMethod jwt.SigningMethod, pubKey *rsa.PublicKey) (bool, error) {
	segments := strings.Split(token, ".")
	segments[1] = body
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], pubKey)
	if err != nil {
		logrus.Errorln("failed to validate signature" + err.Error())
		return false, err
	}
	logrus.Infoln("Signature validation succeeded")
	return true, nil
}

var fixedSigningMethodPS256 = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

func verify(signingMethod jwt.SigningMethod, token string) bool {
	segments := strings.Split(token, ".")
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], test.LoadRSAPublicKeyFromDisk("test/sample_key.pub"))
	return err == nil
}

func TestReversibleLocalSignature(t *testing.T) {
	signingMethod := jwt.SigningMethodPS256.SigningMethodRSA
	_ = signingMethod

}

func getPS256SigingAlg() jwt.SigningMethod {
	return &jwt.SigningMethodRSAPSS{
		SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
		Options: &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       crypto.SHA256,
		},
	}
}

const selfsignedDummySigkey = `-----BEGIN RSA PRIVATE KEY-----
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

const selfsignedDummySigpub = `-----BEGIN CERTIFICATE-----
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
