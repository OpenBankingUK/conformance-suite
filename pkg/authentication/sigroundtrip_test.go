package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/stretchr/testify/assert"
)

var signingORG = `0015800001041RbAAI`
var signingSSA = `fJuUU6dNt0zxnDe59eG0YN`
var sigingJWKS = `https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks`
var signingKID = `kublR092tresXAcL2XFGF3OOsns`
var RTdetachedJWTOzone = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`
var RTRawHeaderB64True = `{
    "alg": "PS256",
    "crit": [
        "http://openbanking.org.uk/iat",
        "http://openbanking.org.uk/iss",
        "http://openbanking.org.uk/tan"
    ],
    "cty": "application/json",
    "http://openbanking.org.uk/iat": 1588587682.456,
    "http://openbanking.org.uk/iss": "0015800001041RbAAI",
    "http://openbanking.org.uk/tan": "openbanking.org.uk",
    "kid": "kublR092tresXAcL2XFGF3OOsns",
    "typ": "JOSE"
}
`
var RTRawHeaderB64FALSE = `{
	"alg": "PS256",
	"b64": false,
    "crit": [
		"b64",
        "http://openbanking.org.uk/iat",
        "http://openbanking.org.uk/iss",
        "http://openbanking.org.uk/tan"
    ],
    "cty": "application/json",
    "http://openbanking.org.uk/iat": 1588587682.456,
    "http://openbanking.org.uk/iss": "0015800001041RbAAI",
    "http://openbanking.org.uk/tan": "openbanking.org.uk",
    "kid": "kublR092tresXAcL2XFGF3OOsns",
    "typ": "JOSE"
}
`
var RTrawBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`

func TestSignB64True(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(RTRawHeaderB64True))
	_ = header

	body := base64.RawURLEncoding.EncodeToString([]byte(RTrawBody))
	//fmt.Println(header + "." + body + ".")
	toBeSigned := header + "." + body

	// jwk, err := getJwkFromJwks("kublR092tresXAcL2XFGF3OOsns", "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks")
	// assert.Nil(t, err)

	// certs, err := ParseCertificateChain(jwk.X5c)
	// assert.Nil(t, err)
	// cert := certs[0]

	//privateKey := loadPrivateKey()
	//privateKey, _ := pem.Decode([]byte(pembytes))
	privateKey, err := ParseRSAPrivateKeyFromPEM([]byte(pembytes))
	assert.Nil(t, err)

	fmt.Println("payload:" + toBeSigned)
	signature, err := jws.SignLiteral([]byte(RTrawBody), jwa.PS256, privateKey, []byte(RTRawHeaderB64True))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(signature))
}

func TestOzonePublicKey2EncodePayloadRT(t *testing.T) {
	kid, err := getKidFromToken(RTdetachedJWTOzone)
	assert.Nil(t, err)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RbAAI/fJuUU6dNt0zxnDe59eG0YN.jwks")
	assert.Nil(t, err)

	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]

	signedJWT, err := insertBodyIntoJWT(RTdetachedJWTOzone, RTrawBody, true)
	assert.Nil(t, err)

	verified, err := jws.Verify([]byte(signedJWT), jwa.PS256, cert.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}
	log.Printf("signed message verified! -> %s", verified)
}

func loadPrivateKey() *rsa.PrivateKey {

	priv, err := ioutil.ReadFile("/home/julianc/OpenBankingUK/notes/REFAPP-1030-Review-Waiver007/privatekey.pem")
	if err != nil {
		fmt.Println(err)
	}
	block, _ := pem.Decode([]byte(priv))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	fmt.Println("----------")
	return key
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

var pembytes = ``
