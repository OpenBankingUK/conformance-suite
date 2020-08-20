package authentication

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"
)

// Workaround for default PS256 signing parameter issue
// https://github.com/dgrijalva/jwt-go/issues/285
var SigningMethodPS256 = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	},
}

var b64Status bool      // for report export
var eidas_kid string    // key id when using eidas certificates
var eidas_issuer string // issuer for jwt signing for eidas certificates

func GetSigningAlg(alg string) (jwt.SigningMethod, error) {
	switch strings.ToUpper(alg) {
	case "PS256":
		return SigningMethodPS256, nil
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
	switch apiVersion {
	case "v3.1":
		issuer, err = cert.SignatureIssuer(true)
		if err != nil {
			logrus.Warn("cannot Issuer for Signature: ", err.Error())
			return "", errors.New("authentication.GetJWSIssuerString: cannot Issuer for Signature: " + err.Error())
		}
	case "v3.0":
		issuer, _, _, err = cert.DN()
		if err != nil {
			logrus.Warn("cannot get certificate DN: ", err.Error())
			return "", errors.New("authentication.GetJWSIssuerString: cert.DN() failed" + err.Error())
		}
	default:
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

// CreateSignature Get the complete, signed token for jws usage
// Takes the token object, private key, payload body and b64encoding indicator
// Create the signing string which includes the token header and payload body
// Then signs this string using the key provided - the signing algorithm is part of the jwt.Token object
func CreateSignature(t *jwt.Token, key interface{}, body string, b64encoded bool) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = SigningString(t, body, b64encoded); err != nil {
		return "", errors.Wrap(err, "authentication.CreateSignature: SigningString(t, body) failed")
	}

	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", errors.Wrap(err, "authentication.CreateSignature: t.Method.Sign(sstr, key failed")
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

// JWT SigningString
// takes the token, body string and b64 indicator
// if b64encoded=true - base64urlEncodes the payload string as part of the string to be signed
// if b64encoded=false - includes the payload unencoded (unmodified) in the string to be signed
func SigningString(t *jwt.Token, body string, b64encoded bool) (string, error) {
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
			if b64encoded { // b64=true so encode segment - Sign with payload B64 encoding true - default for v3.1.4 and above of apis
				parts[i] = jwt.EncodeSegment(jsonValue)
			} else { // b64=false so include unencoded - Sign with payload B64 encoding false - v3.1.3 and previous versions of apis
				parts[i] = string(jsonValue)
			}
		}
	}
	return strings.Join(parts, "."), nil
}

func NewJWSSignature(requestBody string, ctx ContextInterface, alg jwt.SigningMethod) (string, error) {

	minifiedBody, err := minifiyJSONBody(requestBody)
	if err != nil {
		return "", errors.Wrap(err, `NewJWSSignature: minifyBody failed`)
	}

	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "NewJWSSignature: unable to sign certificate from context")
	}

	var kid string
	if eidas_kid != "" {
		kid = eidas_kid
		logrus.Debugf("using Kid for eidas cert : %s", kid)
	} else {
		kid, err = getKidFromCertificate(cert)
		if err != nil {
			return "", errors.Wrap(err, "NewJWSSignature: getKidFromCertificate failed")
		}
	}

	var issuer string
	if eidas_issuer != "" {
		issuer = eidas_issuer
		logrus.Debugf("using issuer for eidas cert : %s", issuer)
	} else {
		issuer, err = GetJWSIssuerString(ctx, cert)
		if err != nil {
			return "", errors.Wrap(err, "NewJWSSignature: unable to retrieve issuer from context")
		}
	}

	trustAnchor := "openbanking.org.uk"
	useNonOBDirectory, exists := ctx.Get("nonOBDirectory")
	if !exists {
		return "", errors.New("NewJWSSignature: unable to retrieve nonOBDirectory from context")
	}
	useNonOBDirectoryAsBool, ok := useNonOBDirectory.(bool)
	if !ok {
		return "", errors.New("NewJWSSignature: unable to cast nonOBDirectory to bool")
	}
	if useNonOBDirectoryAsBool {
		kid, err = ctx.GetString("signingKid")
		if err != nil {
			return "", errors.Wrap(err, "NewJWSSignature: unable to retrieve singingKid from context")
		}
		issuer, err = ctx.GetString("issuer")
		if err != nil {
			return "", errors.Wrap(err, "NewJWSSignature: unable to retrieve issue from context")
		}
		trustAnchor, err = ctx.GetString("signatureTrustAnchor")
		if err != nil {
			return "", errors.Wrap(err, "NewJWSSignature: unable to retrieve signatureTrustAnchor from context")
		}
	}

	logrus.WithFields(logrus.Fields{
		"kid":    kid,
		"issuer": issuer,
		"alg":    alg.Alg(),
		"claims": minifiedBody,
	}).Trace("jws signature creation")

	b64encoding, err := GetB64Encoding(ctx)
	if err != nil {
		return "", errors.New("NewJWSSignature: cannot GetB64Encoding " + err.Error())
	}

	return buildSignature(b64encoding, kid, issuer, trustAnchor, minifiedBody, alg, cert.PrivateKey())
}

func GetB64Encoding(ctx ContextInterface) (bool, error) {
	paymentApiVersion, err := getPaymentApiVersion(ctx)
	if err != nil {
		return false, errors.New("NewJWSSignature: cannot find payment apiversion: " + err.Error())
	}

	b64encoding, err := getB64Encoding(paymentApiVersion)
	if err != nil {
		return false, errors.New("NewJWSSignature: cannot getB64Encoding " + err.Error())
	}

	return b64encoding, nil
}

func getB64Encoding(paymentVersion string) (bool, error) {
	switch paymentVersion {
	case "v3.1.6":
		fallthrough
	case "v3.1.5":
		fallthrough
	case "v3.1.4":
		setB64Status(true) // record setting for report
		return true, nil
	case "v3.1.3":
		fallthrough
	case "v3.1.2":
		fallthrough
	case "v3.1.1":
		fallthrough
	case "v3.1.0":
		fallthrough
	case "v3.1":
		return false, nil
	case "v3.0":
		return false, errors.New("b64Encoding: Unsupported Payment api Version (" + paymentVersion + ")")
	}
	return false, errors.New("b64Encoding: unknown Payment apiVersion (" + paymentVersion + ")")
}

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

func getKidFromCertificate(cert Certificate) (string, error) {
	modulus := cert.PublicKey().N.Bytes()
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, err := CalcKid(modulusBase64)
	return kid, err
}

// Gets the payment api version from the context
// looks for the "apiversions" key
// requires payment version to be in the form similar to "payments_v3.1.0"
// apiversions is a string slice
func getPaymentApiVersion(ctx ContextInterface) (string, error) {
	apiVersions, err := ctx.GetStringSlice("apiversions")
	if err != nil {
		return "", errors.New("NewJWSSignature: cannot find apiversions: " + err.Error())
	}

	for _, str := range apiVersions {
		if strings.HasPrefix(str, "payments_") {
			paymentVersion := after(str, "payments_")
			if paymentVersion == "" {
				return "", errors.New("Cannot find payment api version: " + str)
			}
			return paymentVersion, nil
		}
	}
	return "", errors.New("Payment API version not found: " + strings.Join(apiVersions, ","))
}

// Get a string after the given string
func after(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

func setB64Status(status bool) {
	b64Status = status
}
func GetB64Status() bool {
	return b64Status
}

func SetEidasSigningParameters(issuer, kid string) {
	eidas_issuer = issuer
	eidas_kid = kid
	logrus.Debugf("Setting EIDAS Signing Parameters ssa: %s, kid: %s", issuer, kid)
}
