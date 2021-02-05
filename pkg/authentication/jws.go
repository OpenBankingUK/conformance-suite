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

// SigningString takes the token, body string and b64 indicator
// if b64encoded=true - base64urlEncodes the payload string as part of the string to be signed
// if b64encoded=false - includes the payload unencoded (unmodified) in the string to be signed
func SigningString(t *jwt.Token, body string, b64encoded bool) (string, error) {
	headersJSON, err := json.Marshal(t.Header)
	if err != nil {
		return "", errors.Wrap(err, "authentication.SigningString: json.Marshal(t.Header) failed")
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

// ModifyJWSHeaders allows the caller to mutate an existing JWS, re-signed with the new contents
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
	useNonOBDirectory, exists := ctx.Get("nonOBDirectoryTPP")
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
		"tan":    trustAnchor,
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

func SetEidasSigningParameters(issuer, kid string) error {
	eidas_issuer = issuer
	eidas_kid = kid
	logrus.Debugf("Setting EIDAS Signing Parameters iss: %s, kid: %s", issuer, kid)
	// Check relaxed to allow HSBC Trust Anchor issuers
	// if !checkSignatureIssuerTPP(eidas_issuer) {
	// 	return fmt.Errorf("Invalid EIDAS Issuer String (%s)", eidas_issuer)
	// }
	return nil
}
