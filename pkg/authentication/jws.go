package authentication

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"
)

// SigningMethodPS256 is a workaround for default PS256 signing parameter issue
// https://github.com/dgrijalva/jwt-go/issues/285
var SigningMethodPS256 = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	},
}

var b64Status bool // for report export

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
		return nil, fmt.Errorf("authentication.SigningCertFromContext: couldn't create `certificate` from pub/priv keys: %w", err)
	}
	return cert, nil
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
		return "", fmt.Errorf("authentication.CreateSignature: SigningString(t, body) failed: %w", err)
	}

	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", fmt.Errorf("authentication.CreateSignature: t.Method.Sign(sstr, key failed: %w", err)
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

// SigningString takes the token, body string and b64 indicator
// if b64encoded=true - base64urlEncodes the payload string as part of the string to be signed
// if b64encoded=false - includes the payload unencoded (unmodified) in the string to be signed
func SigningString(t *jwt.Token, body string, b64encoded bool) (string, error) {
	headersJSON, err := json.Marshal(t.Header)
	if err != nil {
		return "", fmt.Errorf("authentication.SigningString: json.Marshal(t.Header) failed: %w", err)
	}
	headers := jwt.EncodeSegment(headersJSON)

	var payload string
	if b64encoded {
		payload = jwt.EncodeSegment([]byte(body))
	} else {
		payload = body
	}

	return strings.Join([]string{headers, payload}, "."), nil
}

// RemoveJWSHeader provides an option which modifies an existing JWT by
// deleting specified keys from its header.
func RemoveJWSHeader(removed []string) JWSHeaderOpt {
	substitutes := map[string]string{
		"iat": "http://openbanking.org.uk/iat",
		"iss": "http://openbanking.org.uk/iss",
		"tan": "http://openbanking.org.uk/tan",
	}

	return func(current map[string]interface{}) map[string]interface{} {
		result := map[string]interface{}{}
		for key, value := range current {
			result[key] = value
		}

		for _, key := range removed {
			substitute, ok := substitutes[key]
			if ok {
				key = substitute
			}

			delete(result, key)
		}
		return result
	}
}

// SetJWSHeader provides an option which modifies an existing JWT by
// setting specified keys on its header.
func SetJWSHeader(entries map[string]interface{}) JWSHeaderOpt {
	return func(current map[string]interface{}) map[string]interface{} {
		result := map[string]interface{}{}
		for key, value := range current {
			// not concurrently accessed, delete should be sufficient.
			result[key] = value
		}

		for key, value := range entries {
			result[key] = value
		}
		return result
	}
}

// JWSHeaderOpt is a function signature which is used for altering JWS header when passed to ModifyJWSHeaders
type JWSHeaderOpt func(map[string]interface{}) map[string]interface{}

// ModifyJWSHeaders allows the caller to mutate an existing JWS for testing purposes, re-signed with the new contents
func ModifyJWSHeaders(jws string, ctx ContextInterface, opts ...JWSHeaderOpt) (string, error) {
	if len(opts) == 0 {
		return jws, nil
	}

	// decode the headers
	segments := strings.Split(jws, ".")
	if len(segments) == 0 {
		return "", fmt.Errorf("failed to modify JWS: received invalid JWS as input")
	}

	// assuming b64 encoded body, meaning that this can't be used in tests prior 3.1.4 (!)
	b64Encoded := true
	body := []byte{}
	if len(segments) > 1 {
		body = []byte(segments[1])
	}

	headersB64Decoded, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		return "", fmt.Errorf("failed to modify JWS: %w", err)
	}
	header := map[string]interface{}{}
	err = json.Unmarshal(headersB64Decoded, &header)
	if err != nil {
		return "", fmt.Errorf("failed to modify JWS: %w", err)
	}

	// apply mutators (header opts)
	for _, opt := range opts {
		header = opt(header)
	}

	// The new signature is using the alg specified in the token;
	// if alg is not set then the signature can't be created
	alg, ok := header["alg"].(string)
	if !ok {
		return "", fmt.Errorf("failed to modify JWS: signing alg is undefined, 'alg' key is not set")
	}

	signingMethod, err := GetSigningAlg(alg)
	if err != nil {
		return "", fmt.Errorf("failed to modify JWS: %w", err)
	}

	token := &jwt.Token{
		Header: header,
		Method: signingMethod,
	}

	// sign anew and return the produce
	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to modify JWS: %w", err)
	}

	privKey := cert.PrivateKey()
	tokenString, err := CreateSignature(token, privKey, string(body), b64Encoded)
	if err != nil {
		return "", fmt.Errorf("failed to modify JWS: %w", err)
	}
	return SplitJWSWithBody(tokenString), nil
}

// NewJWSSignature creates a signature to be used with TPP API calls
func NewJWSSignature(requestBody string, ctx ContextInterface, alg jwt.SigningMethod) (string, error) {
	minifiedBody, err := minifiyJSONBody(requestBody)
	if err != nil {
		return "", fmt.Errorf("NewJWSSignature: minifyBody failed: %w", err)
	}

	cert, err := SigningCertFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("NewJWSSignature: unable to sign certificate from context: %w", err)
	}

	tppSignatureKID, err := ctx.GetString("tpp_signature_kid")
	if err != nil {
		return "", fmt.Errorf("missing configuration for key 'tpp_signature_kid': %w", err)
	}

	tppSignatureTAN, err := ctx.GetString("tpp_signature_tan")
	if err != nil {
		return "", fmt.Errorf("failed to populate TPP signature Trust Anchor: '%w'", err)
	}

	tppSignatureIssuer, err := ctx.GetString("tpp_signature_issuer")
	if err != nil {
		return "", fmt.Errorf("failed to populate TPP signature Issuer: '%w'", err)
	}

	version, _ := ctx.GetString("api-version")
	if version == "v3.0" {
		tppSignatureIssuer, err = legacyIssuerFromCert(cert)
	}

	logrus.WithFields(logrus.Fields{
		"kid":    tppSignatureKID,
		"issuer": tppSignatureIssuer,
		"alg":    alg.Alg(),
		"claims": minifiedBody,
		"tan":    tppSignatureTAN,
	}).Trace("jws signature creation")

	b64encoding, err := GetB64Encoding(ctx)
	if err != nil {
		return "", fmt.Errorf("NewJWSSignature: cannot GetB64Encoding: %w", err)
	}

	return buildSignature(b64encoding, tppSignatureKID, tppSignatureIssuer, tppSignatureTAN, minifiedBody, alg, cert.PrivateKey())
}

func legacyIssuerFromCert(cert Certificate) (string, error) {
	issuer, _, _, err := cert.DN()
	if err != nil {
		logrus.Warn("cannot get certificate DN: ", err.Error())
		return "", errors.New("authentication.GetJWSIssuerString: cert.DN() failed" + err.Error())
	}

	return issuer, nil
}

// GetB64Encoding returns - based on the API version - if the TPP signature should use base64 encoding for the payload
func GetB64Encoding(ctx ContextInterface) (bool, error) {
	paymentAPIVersion, err := getPaymentAPIVersion(ctx)
	if err != nil {
		return false, errors.New("NewJWSSignature: cannot find payment apiversion: " + err.Error())
	}

	b64encoding, err := getB64Encoding(paymentAPIVersion)
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
func getPaymentAPIVersion(ctx ContextInterface) (string, error) {
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
