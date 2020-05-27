package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var detachedJWT = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`

var rawBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`

var tokenb64true = jwt.Token{
	Header: map[string]interface{}{
		"typ":                           "JOSE",
		"kid":                           "kublR092tresXAcL2XFGF3OOsns",
		"cty":                           "application/json",
		"http://openbanking.org.uk/iat": 1588587682.456,
		"http://openbanking.org.uk/iss": "0015800001041RbAAI", //ASPSP ORGID or TTP ORGID/SSAID
		"http://openbanking.org.uk/tan": "openbanking.org.uk", //Trust anchor
		"alg":                           "PS256",
		"crit": []string{
			"http://openbanking.org.uk/iat",
			"http://openbanking.org.uk/iss",
			"http://openbanking.org.uk/tan",
		},
	},
	Method: SigningMethodPS256,
}

var tokenb64false = jwt.Token{
	Header: map[string]interface{}{
		"typ":                           "JOSE",
		"b64":                           false,
		"kid":                           "kublR092tresXAcL2XFGF3OOsns",
		"cty":                           "application/json",
		"http://openbanking.org.uk/iat": 1588587682.456,
		"http://openbanking.org.uk/iss": "0015800001041RbAAI", //ASPSP ORGID or TTP ORGID/SSAID
		"http://openbanking.org.uk/tan": "openbanking.org.uk", //Trust anchor
		"alg":                           "PS256",
		"crit": []string{
			"b64",
			"http://openbanking.org.uk/iat",
			"http://openbanking.org.uk/iss",
			"http://openbanking.org.uk/tan",
		},
	},
	Method: SigningMethodPS256,
}

func TestSigningStringB64True(t *testing.T) {
	signingString, err := SigningString(&tokenb64true, rawBody, true)
	assert.Nil(t, err)
	fmt.Println("signing string: b64=true: ", signingString)
	key, err := ParseRSAPrivateKeyFromPEM([]byte(pembytes))
	assert.Nil(t, err)
	sig, err := tokenb64true.Method.Sign(signingString, key)
	assert.Nil(t, err)
	fullSig := strings.Join([]string{signingString, sig}, ".")
	fmt.Println("fullsig: " + fullSig)
	detachedJwt := SplitJWSWithBody(fullSig)
	fmt.Println("detached JWT: " + detachedJwt)
	idx := strings.LastIndex(detachedJwt, ".")
	signature := detachedJwt[idx:]
	valid := verifySig(t, SigningMethodPS256, signingString, signature, true) // This line works
	fmt.Printf("Signature is valid: %t ", valid)
}

func TestSigningStringB64False(t *testing.T) {
	signingString, err := SigningString(&tokenb64false, rawBody, false)
	assert.Nil(t, err)
	fmt.Println("signing string b64=false: ", signingString)
	key, err := ParseRSAPrivateKeyFromPEM([]byte(pembytes))
	assert.Nil(t, err)
	sig, err := tokenb64false.Method.Sign(signingString, key)
	assert.Nil(t, err)
	fullSig := strings.Join([]string{signingString, sig}, ".")
	fmt.Println("fullsig: " + fullSig)
	detachedJwt := SplitJWSWithBody(fullSig)
	fmt.Println("detached JWT: " + detachedJwt)
	idx := strings.LastIndex(detachedJwt, ".")
	signature := detachedJwt[idx:]
	fmt.Println("Signature: " + signature)
	valid := verifySig(t, SigningMethodPS256, signingString, signature, false) // !!!! This is the line !!!!!
	fmt.Printf("Signature is valid: %t ", valid)
}

func verifySig(t *testing.T, signingMethod jwt.SigningMethod, signingString, signature string, b64 bool) bool {
	kid, err := getKidFromToken(signingString)
	fmt.Println("kid: " + kid)
	assert.Nil(t, err)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks")
	assert.Nil(t, err)
	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]
	verified, err := MyJwsVerify(signingString+signature, jwa.PS256, cert.PublicKey, b64)
	if err != nil {
		logrus.Errorf("failed to verify message: %v", err)
		return false
	}
	logrus.Tracef("signed message verified! -> %s", verified)
	assert.Nil(t, err)
	return err == nil
}

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
func TestOzoneHeaderValidation314(t *testing.T) {
	err := ValidateSignatureHeader(detachedJWT, true)
	assert.Nil(t, err)
}

func Test314HeaderFailsValidationWhenB64False(t *testing.T) {
	err := ValidateSignatureHeader(detachedJWT, false)
	assert.NotNil(t, err)
}

func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("ErrKeyMustBePEMEncoded")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, errors.New("ErrNotRSAPrivateKey")
	}

	return pkey, nil
}

var pembytes = `-----BEGIN PRIVATE KEY-----
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
-----END PRIVATE KEY-----
`
