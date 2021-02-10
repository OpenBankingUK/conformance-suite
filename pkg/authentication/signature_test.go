package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
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

var detachedJWT = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTk2MDMwMjczLjU4MiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiemtib0tGalFSd0JkOVVFblBDNXdsdjU3aWc0In0..ioBQDgKbY04pjS3LF4ezFuB-so4DwAobnLJPhn4uCLyUjN2JEWQTCkkKbilMTSKq3mYuFywu8Nc3eELpZfK50wxPOdKyGt5fBH89_F0OzAw-9xvGWLlAyubIhIMnTe05sSXEi-6pOti6SdoKP4KabxeBustwMzUoH0Fq0UTPel0pgJam9aSRqG8y-MfueeAHWE0icrAzsb1Wtprpinn62EmZyfYCgWWIIPgk323L4ETptBvn6PHBpybCIQHF8omxRw9mjcyLlq0mdI-JeVyXjikHXjRHLbx2ZtHpGuiwloCOhdyrslH3kpIGAAOr9sny1JLljPy-dGZ04H8WYVEzsw`
var expiredDetachedJWT = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`

var rawBody = `{"Data":{"Initiation":{"CreditorAccount":{"Identification":"70000170000002","Name":"Mr. Roberto Rastapopoulos & Ivan Sakharine & mits","SchemeName":"UK.OBIE.SortCodeAccountNumber"},"EndToEndIdentification":"e2e-domestic-pay","InstructedAmount":{"Amount":"1.00","Currency":"GBP"},"InstructionIdentification":"b470d6561b4147c29e519bb9dfca40af"},"ConsentId":"sdp-1-007c81e6-76be-4592-9fb5-72691751eb32","CreationDateTime":"2020-07-29T13:44:33.577Z","Status":"AwaitingAuthorisation","StatusUpdateDateTime":"2020-07-29T13:44:33.577Z"},"Risk":{},"Links":{"Self":"https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-007c81e6-76be-4592-9fb5-72691751eb32"},"Meta":{}}`
var expiredRawBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`

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
	pembytes, err := getPemBytes()
	if err != nil {
		fmt.Println("WARNING - missing OB Test private key file: " + err.Error())
		return
	}
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
	pembytes, err := getPemBytes()
	if err != nil {
		fmt.Println("WARNING - missing OB Test private key file: " + err.Error())
		return
	}
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
	certs, err := parseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]
	verified, err := JWSVerify(signingString+signature, jwa.PS256, cert.PublicKey, b64)
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

	if jwk.X5c == nil {
		fmt.Println("certificate chain not found in jwk.X5c claim")
		return
	}
	certs, err := parseCertificateChain(jwk.X5c)
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

func TestOzoneExpiredPublicKey2EncodePayload(t *testing.T) {
	kid, err := getKidFromToken(expiredDetachedJWT)
	assert.Nil(t, err)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RHAAY/0015800001041RHAAY.jwks")
	assert.Nil(t, err)

	if jwk.X5c == nil {
		fmt.Println("certificate chain not found in jwk.X5c claim")
		return
	}
	certs, err := parseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]
	signedJWT, err := insertBodyIntoJWT(expiredDetachedJWT, expiredRawBody, true)
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

var aspspIssuerTests = []struct {
	iss      string // input
	expected bool   // expected result
}{
	{"0014H00002LmnSNQAZ", true},
	{"", false},
	{"0014H00002LmnSNQAZz", false},
	{"0014H000 2LmnSNQAZ", false},
	{"0014H00002LmnSNQA=", false},
	{"0014H00002LmnSNQA,", false},
	{"0014H00002LmnSNQA,", false},
	{"0014H00002LmnSNQA14H00002Lmn", false},
	{"001=93939d9dkkdmde", false},
	{"000000", false},
	{"CN=AORlKI9oUF85BB2WUr7ZK2,OU=0015800000Af7GZAAI,O=OpenBanking,C=GB", false},
	{"0014H0", false},
}

var tppIssuerTests = []struct {
	iss      string // input
	expected bool   // expected result
}{
	{"0014H00002LmnSNQAZ/AORlKI7oUF85BB2WUr7ZK2", true},
	{"", false}, // empty
	{"0014H00002LmnSNQA/AORlKI7oUF85BB2WUr7ZK2", false},   // orgid too short
	{"0014H00002LmnSNQAZ/ARlKI7oUF85BB2WUr7ZK2", false},   // ssa too short
	{"0014H00002LmnSNQAZ/AORlKI7oUF85BB2WUr7ZK22", false}, // ssa too long
	{"0014H00002LmnSNQAZZ/AORlKI7oUF85BB2WUr7ZK2", false}, // orgid too long
	{"0014H00002LmnSNQAZz/", false},                       // ssa missing
	{"/0014H00002LmnSNQA=", false},                        // orgid missing
	{"0014H00002LmnSNQA,/AORlKI7oUF85BB2WUr7ZK2", false},  // orgid included ,
	{"0014H00002LmnSNQA=/AORlKI7oUF85BB2WUr7ZK2", false},  // orgid included =
	{"0014H00002LmnSNQAA/AOR,KI7oUF85BB2WUr7ZK2", false},  // orgid in
	{"0014H00002LmnSNQAA/AORlKI7:UF85BB2WU,7ZK2", false},
	{"0014H00002AmnSNQAA/AORlKI7oUF85BB2WU=7ZK2", false},
	{"0014H00002AmnSNQAA/AORlKI7oUF85BB2WU 7ZK2", false},
	{"0014H00002A nSNQAA/AORlKI7oUF85BB2WU 7ZK2", false},
	{"0014H00002LmnSNQA,", false},
	{"000000", false},
	{"CN=AORlKI7oUF85BB2WUrGZK2,OU=0015800000X67GgAAI,O=OpenBanking,C=GB", false},
	{"0014H0", false},
}

func TestInvalidASPSPIssuerString(t *testing.T) {
	for _, tt := range aspspIssuerTests {
		actual := checkSignatureIssuerASPSP(tt.iss)
		if actual != tt.expected {
			t.Errorf("checkSignatureIssuerASPSP(%s): expected %t, actual %t", tt.iss, tt.expected, actual)
		}
	}
}

func TestInvalidTPPIssuerString(t *testing.T) {
	for _, tt := range tppIssuerTests {
		actual := checkSignatureIssuerTPP(tt.iss)
		if actual != tt.expected {
			t.Errorf("checkSignatureIssuerTPP(%s): expected %t, actual %t", tt.iss, tt.expected, actual)
		}
	}
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

func getPemBytes() (string, error) {
	raw, err := ioutil.ReadFile("../../../certs/testprivatekey.pem")
	return string(raw), err
}
