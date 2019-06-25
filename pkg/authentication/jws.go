package authentication

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"
)

type ContextInterface interface {
	// GetString get the string value associated with key
	GetString(key string) (string, error)
	// Get the key form the Context map - currently assumes value converts easily to a string!
	Get(key string) (interface{}, bool)
}

func GetSigningAlg(alg string) (jwt.SigningMethod, error) {
	switch strings.ToUpper(alg) {
	case "PS256":
		// Workaround
		// https://github.com/dgrijalva/jwt-go/issues/285
		return &jwt.SigningMethodRSAPSS{
			SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
			Options: &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthEqualsHash,
				Hash:       crypto.SHA256,
			},
		}, nil
	case "RS256":
		return jwt.SigningMethodRS256, nil
	case "NONE":
		fallthrough
	default:
		return nil, fmt.Errorf("authentication.GetSigningAlg: unable to find signing algorithm %q", alg)
	}
}

func SigningCertFromContext(ctx ContextInterface) (Certificate, error) {
	privKey, err := ctx.GetString("signingPrivate")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPrivate` in context")
	}
	pubKey, err := ctx.GetString("signingPublic")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPublic` in context")
	}
	cert, err := NewCertificate(pubKey, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "authentication.SigningCertFromContext: couldn't create `certificate` from pub/priv keys")
	}
	return cert, nil
}

func GetJWSIssuerString(ctx ContextInterface, cert Certificate) (string, error) {
	apiVersion, err := ctx.GetString("api-version")
	if err != nil {
		return "", errors.New("authentication.GetJWSIssuerString: cannot find api-version: " + err.Error())
	}
	var issuer string
	if apiVersion == "v3.1" {
		issuer, err = cert.SignatureIssuer(true)
		if err != nil {
			logrus.Warn("cannot Issuer for Signature: ", err.Error())
			return "", errors.New("authentication.GetJWSIssuerString: cannot Issuer for Signature: " + err.Error())
		}
	} else if apiVersion == "v3.0" {
		issuer, err = cert.DN()
		if err != nil {
			logrus.Warn("cannot get certificate DN: ", err.Error())
		}
	} else {
		return "", errors.New("authentication.GetJWSIssuerString: cannot get issuer for jws signature but api-version doesn't match 3.0.0 or 3.1.0")
	}

	return issuer, nil
}

func SplitJWSWithBody(token string) string {
	firstPart := token[:strings.IndexByte(token, '.')]
	idx := strings.LastIndex(token, ".")
	lastPart := token[idx:]
	return firstPart + "." + lastPart
}

// SignedString Get the complete, signed token for jws usage
func SignedString(t *jwt.Token, key interface{}, body string) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = SigningString(t, body); err != nil {
		return "", errors.Wrap(err, "authentication.SignedString: SigningString(t, body) failed")
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", errors.Wrap(err, "authentication.SignedString: t.Method.Sign(sstr, key failed")
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

// SigningString -
func SigningString(t *jwt.Token, body string) (string, error) {
	var err error
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return "", errors.Wrap(err, "authentication.SigningString: json.Marshal(t.Header) failed")
			}
		} else {
			jsonValue = []byte(body)
		}
		if i == 0 {
			parts[i] = jwt.EncodeSegment(jsonValue)
		} else {
			parts[i] = string(jsonValue)
		}
	}
	return strings.Join(parts, "."), nil
}

func NewJWSSignature(requestBody string, ctx ContextInterface, alg jwt.SigningMethod) (string, error) {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minjson.Minify)
	minifiedBody, err := m.String("application/json", requestBody)
	if err != nil {
		return "", err
	}
	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to sign certificate from context")
	}
	modulus := cert.PublicKey().N.Bytes()
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, _ := CalcKid(modulusBase64)

	issuer, err := GetJWSIssuerString(ctx, cert)
	if err != nil {
		return "", errors.Wrap(err, "authentication.NewJWSSignature: unable to retrieve issuer from context")
	}
	trustAnchor := "openbanking.org.uk"
	useNonOBDirectory, exists := ctx.Get("nonOBDirectory")
	if !exists {
		return "", errors.New("authentication.NewJWSSignature: unable to retrieve nonOBDirectory from context")
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

	tok := jwt.Token{
		Header: map[string]interface{}{
			"typ":                           "JOSE",
			"kid":                           kid,
			"b64":                           false,
			"cty":                           "application/json",
			"http://openbanking.org.uk/iat": time.Now().Unix(),
			"http://openbanking.org.uk/iss": issuer,      //ASPSP ORGID or TTP ORGID/SSAID
			"http://openbanking.org.uk/tan": trustAnchor, //Trust anchor
			"alg":                           alg.Alg(),
			"crit":                          []string{"b64", "http://openbanking.org.uk/iat", "http://openbanking.org.uk/iss", "http://openbanking.org.uk/tan"},
		},
		Method: alg,
	}

	tokenString, err := SignedString(&tok, cert.PrivateKey(), minifiedBody) // sign the token - get as encoded string

	logrus.Tracef("jws:  %v", tokenString)
	detachedJWS := SplitJWSWithBody(tokenString)
	logrus.Tracef("detached jws: %v", detachedJWS)

	return detachedJWS, nil
}
