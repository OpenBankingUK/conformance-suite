package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"

	"crypto/rsa"

	"github.com/dgrijalva/jwt-go/test"
)

func TestJWSDetachedSignature312andBefore1(t *testing.T) {

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

	detachedJWS := authentication.SplitJWSWithBody(sig)

	token := &jwt.Token{Raw: detachedJWS}

	alg, err := authentication.GetSigningAlg("PS256")

	minified, err := MinifyBody(req.Body.(string))

	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fmt.Printf("First Part: %s\n", GetFirstPart(sig))

	// val, err := json.Unmarshal(headerBytes, &token.Header)
	// fmt.Printf("%v\n", val)

	tim := time.Now().Unix()
	token1, _ := MyNewJWSSignature(minified, ctx, alg, cert, tim)

	token2, _ := MyNewJWSSignature(minified, ctx, alg, cert, tim)
	assert.Equal(t, token1, token2)

	var headerBytes []byte
	json.Unmarshal(headerBytes, &token.Header)
	fmt.Printf("headers %v\n", headerBytes)

	assert.Equal(t, token, detachedJWS)

}

func SplitJWSWithBody(token string) string {
	firstPart := token[:strings.IndexByte(token, '.')]
	idx := strings.LastIndex(token, ".")
	lastPart := token[idx:]
	return firstPart + "." + lastPart
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

// //}

func MinifyBody(body string) (string, error) {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minjson.Minify)
	minifiedBody, err := m.String("application/json", body)

	return minifiedBody, err
}

func SplitBody(token string) string {
	firstPart := token[:strings.IndexByte(token, '.')]
	idx := strings.LastIndex(token, ".")
	lastPart := token[idx:]
	return firstPart + "." + lastPart
}

func GetFirstPart(token string) string {
	firstPart := token[:strings.IndexByte(token, '.')]
	//idx := strings.LastIndex(token, ".")
	//lastPart := token[idx:]
	return firstPart
}

func InsertBodyIn(token string, body interface{}) string {
	firstPart := token[:strings.IndexByte(token, '.')]
	idx := strings.LastIndex(token, ".")
	lastPart := token[idx:]
	return firstPart + "." + body.(string) + lastPart
}

type ContextInterface interface {
	// GetString get the string value associated with key
	GetString(key string) (string, error)
	// Get the key form the Context map - currently assumes value converts easily to a string!
	Get(key string) (interface{}, bool)
}

// TEST
func MyNewJWSSignature(minifiedBody string, ctx ContextInterface, alg jwt.SigningMethod, cert authentication.Certificate, tim int64) (string, error) {
	modulus := cert.PublicKey().N.Bytes()
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, err := authentication.CalcKid(modulusBase64)
	if err != nil {
		return "", errors.Wrap(err, "authentication.NewJWSSignature: CalcKid(modulusBase64) failed")
	}

	issuer, err := authentication.GetJWSIssuerString(ctx, cert)
	if err != nil {
		return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to retrieve issuer from context")
	}
	trustAnchor := "openbanking.org.uk"
	useNonOBDirectory, exists := ctx.Get("nonOBDirectory")
	if !exists {
		return "", errors.New("authentication.NewJWSSiMyNewJWSSignaturegnature: unable to retrieve nonOBDirectory from context")
	}
	useNonOBDirectoryAsBool, ok := useNonOBDirectory.(bool)
	if !ok {
		return "", errors.New("authentication.NewJWSSignature: unable to cast nonOBDirectory to bool")
	}
	if useNonOBDirectoryAsBool {
		kid, err = ctx.GetString("signingKid")
		if err != nil {
			return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to retrieve singingKid from context")
		}
		issuer, err = ctx.GetString("issuer")
		if err != nil {
			return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to retrieve issue from context")
		}
		trustAnchor, err = ctx.GetString("signatureTrustAnchor")
		if err != nil {
			return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to retrieve signatureTrustAnchor from context")
		}
	}
	logrus.Tracef("jws issuer=%s", issuer)

	logrus.WithFields(logrus.Fields{
		"kid":    kid,
		"issuer": issuer,
		"alg":    alg.Alg(),
		"claims": minifiedBody,
	}).Trace("jws signature creation")

	apiVersion, err := ctx.GetString("api-version")
	if err != nil {
		return "", errors.New("authentication.NewJWSSignature: cannot find api-version: " + err.Error())
	}

	var tok jwt.Token
	switch apiVersion {
	case "v.3.1.3":
		// Contains `http://openbanking.org.uk/tan`.
		tok = jwt.Token{
			Header: map[string]interface{}{
				"typ":                           "JOSE",
				"kid":                           kid,
				"cty":                           "application/json",
				"http://openbanking.org.uk/iat": tim,
				"http://openbanking.org.uk/iss": issuer,      //ASPSP ORGID or TTP ORGID/SSAID
				"http://openbanking.org.uk/tan": trustAnchor, //Trust anchor
				"alg":                           alg.Alg(),
				"crit": []string{
					"http://openbanking.org.uk/iat",
					"http://openbanking.org.uk/iss",
					"http://openbanking.org.uk/tan",
				},
			},
			Method: alg,
		}
	case "v3.1":
		// Contains `http://openbanking.org.uk/tan`.
		tok = jwt.Token{
			Header: map[string]interface{}{
				"typ":                           "JOSE",
				"kid":                           kid,
				"b64":                           false,
				"cty":                           "application/json",
				"http://openbanking.org.uk/iat": tim,
				"http://openbanking.org.uk/iss": issuer,      //ASPSP ORGID or TTP ORGID/SSAID
				"http://openbanking.org.uk/tan": trustAnchor, //Trust anchor
				"alg":                           alg.Alg(),
				"crit": []string{
					"b64",
					"http://openbanking.org.uk/iat",
					"http://openbanking.org.uk/iss",
					"http://openbanking.org.uk/tan",
				},
			},
			Method: alg,
		}
	case "v3.0":
		// Does not contain `http://openbanking.org.uk/tan`.
		// Read/Write Data API Specification - v3.0 Specification: https://openbanking.atlassian.net/wiki/spaces/DZ/pages/641992418/Read+Write+Data+API+Specification+-+v3.0.
		// According to the spec this field `http://openbanking.org.uk/tan` should not be sent in the `x-jws-signature` header.
		tok = jwt.Token{
			Header: map[string]interface{}{
				"typ":                           "JOSE",
				"kid":                           kid,
				"b64":                           false,
				"cty":                           "application/json",
				"http://openbanking.org.uk/iat": tim,
				"http://openbanking.org.uk/iss": issuer, //ASPSP ORGID or TTP ORGID/SSAID
				"alg":                           alg.Alg(),
				"crit": []string{
					"b64",
					"http://openbanking.org.uk/iat",
					"http://openbanking.org.uk/iss",
				},
			},
			Method: alg,
		}
	default:
		return "", errors.New("authentication.GetJWSIssuerString: cannot get issuer for jws signature but api-version doesn't match 3.0.0 or 3.1.0")
	}

	//tokenString, err := tok.SignedString(cert.PrivateKey())
	tokenString, err := authentication.SignedString(&tok, cert.PrivateKey(), minifiedBody) // sign the token - get as encoded string
	fmt.Println("token string:" + tokenString)
	if err != nil {
		return "", errors.Wrap(err, "authentication.NewJWSSignature: SignedString(&tok, cert.PrivateKey(), minifiedBody) failed")
	}

	detachedJWS := authentication.SplitJWSWithBody(tokenString)
	return detachedJWS, nil
}

func SigningCertFromContext(ctx ContextInterface) (authentication.Certificate, error) {
	privKey, err := ctx.GetString("signingPrivate")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPrivate` in context")
	}
	pubKey, err := ctx.GetString("signingPublic")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPublic` in context")
	}
	cert, err := authentication.NewCertificate(pubKey, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "authentication.SigningCertFromContext: couldn't create `certificate` from pub/priv keys")
	}
	return cert, nil
}

func TestStuff(t *testing.T) {
	// Invalid signature on jwt.io: https://bit.ly/2FfYHLr
	before := makeToken(jwt.SigningMethodPS256)
	fmt.Printf("Before: %s\nAccepted before: %v\nAccepted after fix: %v\n",
		before, verify(jwt.SigningMethodPS256, before), verify(fixedSigningMethodPS256, before),
	)
	fmt.Println()
	// Valid signature on jwt.io: https://bit.ly/2FfYHLr (print some spaces to Encoded field to refresh signature status)
	after := makeToken(fixedSigningMethodPS256)
	fmt.Printf("After: %s\nAccepted before: %v\nAccepted after fix: %v\n",
		after, verify(jwt.SigningMethodPS256, after), verify(fixedSigningMethodPS256, after),
	)
}

func makeToken(method jwt.SigningMethod) string {
	token := jwt.NewWithClaims(method, jwt.StandardClaims{
		Issuer:   "example",
		IssuedAt: time.Now().Unix(),
	})
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")
	signed, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signed
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
