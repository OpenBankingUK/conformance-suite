package model

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"

	"crypto/rsa"

	"github.com/dgrijalva/jwt-go/test"
)

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
