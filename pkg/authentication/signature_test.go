package authentication

import (
	"crypto/rsa"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var detachedJWT = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`

var rawBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`

func TestOzonePublicKey2EncodePayload(t *testing.T) {
	kid, err := getKidFromToken(detachedJWT)
	assert.Nil(t, err)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RHAAY/0015800001041RHAAY.jwks")
	assert.Nil(t, err)

	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]

	signedJWT, err := insertBodyIntoJWT(detachedJWT, rawBody, true)
	assert.Nil(t, err)

	verified, err := jws.Verify([]byte(signedJWT), jwa.PS256, cert.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}
	log.Printf("signed message verified! -> %s", verified)
}

/* build a test that
uses ozone rawbody
creates appropriate header
signs jwt string
verifies jwt string using ob jwks to lookup keystore - for our cert

*/

func TestOzoneHeaderValidation314(t *testing.T) {
	err := ValidateSignatureHeader(detachedJWT, true)
	assert.Nil(t, err)
}

func Test314HeaderFailsValidationWhenB64False(t *testing.T) {
	err := ValidateSignatureHeader(detachedJWT, false)
	assert.NotNil(t, err)
}

func Test314HeaderValidation(t *testing.T) {
	//invalid b64 in crit claim
	sig := signatureHeader{Alg: "PS256", Kid: "mykid", IssuedAt: decimal.NewFromInt(4000), Issuer: "ORG-ID",
		TrustAnchor: "openbanking.org.uk",
		Critical:    []string{"b64", "http://openbanking.org.uk/tan", "http://openbanking.org.uk/iat", "http://openbanking.org.uk/iss"}}
	err := sig.validateSignatureHeader(true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "crit claim must contain 3 elements")

	// valid crit claim
	sig.Critical = []string{"http://openbanking.org.uk/tan", "http://openbanking.org.uk/iat", "http://openbanking.org.uk/iss"}
	err = sig.validateSignatureHeader(true)
	assert.Nil(t, err)

	// b64=true should not be present
	b64 := true
	sig.B64 = &b64
	err = sig.validateSignatureHeader(true)
	assert.Nil(t, err)
}

var v313Sig = `eyJhbGciOiJQUzI1NiIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0IiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImN0eSI6ImFwcGxpY2F0aW9uL2pzb24iLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lhdCI6MTU5MDQ3NjM4MywiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUmJBQUkvZkp1VVU2ZE50MHp4bkRlNTllRzBZTiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIjoib3BlbmJhbmtpbmcub3JnLnVrIiwia2lkIjoia3VibFIwOTJ0cmVzWEFjTDJYRkdGM09Pc25zIiwidHlwIjoiSk9TRSJ9..OmLnEAFOPPS4VqUhAO-gXmTbjt3NsNKJKdXopk6_MOwoDGwuit9xzO5WQ80oO4LPOCB83NDZ-VR_AOXj-NX0j-fQ7iq-inwEM2Tps5OzcyJ9vXOhr5lb-2mSOxAOS_1_mW_Rs4hXQqa6wcXk0-x2pF1ubGuyBhGmqcx7ak0jcVHuTNb2VgJ-EmhT7EsCV4CEk0EUVKjhObiiq1f3SCm1U9t9UgeWyZduz0o9igNA6ZAKxV9jtZBmTQsDO8ZSAhSBY5l7P1Q-CKvi_1BFulD5oO7jPeByDvvxHoR9-pctqDrvlXZhWvviBRGyMGt1kZpNkK4lwPRGMqNTLeXdw4wTLQ`
var v313body = `{
	"Data": {
	   "Initiation": {
		  "CreditorAccount": {
			 "Identification": "70000170000002",
			 "Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
			 "SchemeName": "UK.OBIE.SortCodeAccountNumber"
		  },
		  "EndToEndIdentification": "e2e-domestic-pay",
		  "InstructedAmount": {
			 "Amount": "1.00",
			 "Currency": "GBP"
		  },
		  "InstructionIdentification": "64f4b7cc20524536ac2cd8f265220bbb"
	   }
	},
	"Risk": {}
 }
 `

func Test313HeaderValidation(t *testing.T) {
	err := ValidateSignatureHeader(v313Sig, false)
	assert.Nil(t, err)
	kid, err := getKidFromToken(v313Sig)
	assert.Nil(t, err)
	minbody, err := minifiyJSONBody(v313body)
	assert.Nil(t, err)
	_ = kid

	token, err := insertBodyIntoJWT(v313Sig, minbody, false)
	assert.Nil(t, err)

	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks")
	assert.Nil(t, err)

	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]

	fmt.Println(token)

	//	segments := strings.Split(token, ".")

	validatedOK, err := validateSignatureTest(v313Sig, domesticPayBody, SigningMethodPS256, cert.PublicKey.(*rsa.PublicKey))
	_ = validatedOK

	// verify2(fixedSigningMethodPS256, segments[0], minbody, segments[2], cert.PublicKey.(*rsa.PublicKey))

	// verified, err := jws.Verify([]byte(token), jwa.PS256, cert.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}
	log.Printf("signed message verified! -> %t", validatedOK)

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

func verify2(signingMethod jwt.SigningMethod, header, body, sig string, key *rsa.PublicKey) error {

	signedString := header + "." + body
	err := signingMethod.Verify(signedString, sig, key)
	fmt.Println(err)
	return err
}

func TestParsing(t *testing.T) {
	minbody, err := minifiyJSONBody(v313body)
	assert.Nil(t, err)

	kid, err := getKidFromToken(detachedJWT)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks")
	assert.Nil(t, err)

	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]

	//jwsMsg, err := jws.ParseString(minbody)
	jws.Verify([]byte(minbody), jwa.PS256, cert.PublicKey)

	kid, err = getKidFromToken(detachedJWT)
	assert.Nil(t, err)
	jwk, err = getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RHAAY/0015800001041RHAAY.jwks")
	assert.Nil(t, err)

	certs, err = ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert = certs[0]

	signedJWT, err := insertBodyIntoJWT(detachedJWT, rawBody, true)
	assert.Nil(t, err)

	verified, err := jws.Verify([]byte(signedJWT), jwa.PS256, cert.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}
	log.Printf("signed message verified! -> %s", verified)

}
