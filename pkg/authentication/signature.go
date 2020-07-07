package authentication

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws/verify"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type JWKS struct {
	Keys []JWK
}

type JWK struct {
	Alg string   `json:"alg,omitempty"`
	Kty string   `json:"kty,omitempty"`
	X5c []string `json:"x5c,omitempty"`
	N   string   `json:"n,omitempty"`
	E   string   `json:"e,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5t string   `json:"x5t,omitempty"`
	X5u string   `json:"x5u,omitempty"`
	Use string   `json:"use,omitempty"`
}

// ValidateSignature
// take the signature JWT
// extract the kid
// used the kid to lookup the public key in the JWKS
//
func ValidateSignature(jwtToken, body, jwksUri string, b64 bool) (bool, error) {

	err := ValidateSignatureHeader(jwtToken, b64)
	if err != nil {
		return false, err
	}

	kid, err := getKidFromToken(jwtToken)
	if err != nil {
		return false, err
	}
	jwk, err := getJwkFromJwks(kid, jwksUri)
	if err != nil {
		return false, err
	}

	certs, err := ParseCertificateChain(jwk.X5c)
	if err != nil {
		return false, err
	}

	cert := certs[0]

	signature, err := insertBodyIntoJWT(jwtToken, body, b64) // b64claim
	if err != nil {
		logrus.Errorf("failed to insert body into signature message: %v", err)
		return false, err
	}
	logrus.Trace("Signature with payload: " + signature)

	verified, err := MyJwsVerify(signature, jwa.PS256, cert.PublicKey, b64)
	if err != nil {
		logrus.Errorf("failed to verify message: %v", err)
		return false, err
	}

	logrus.Tracef("signed message verified! -> %s", verified)

	return true, nil
}

// insertBodyB64False
// jwt contains "header..signature"
// insert body into jwt resulting in "header.body.signature"
// b64 parameter controls body encoding
// b64=true  = R/W Api 3.1.4 and after
// b64=false = R/W Api 3.1.3 and prior
func insertBodyIntoJWT(token, body string, b64 bool) (string, error) {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return "", errors.New("Signature Token does not have 3 segments: " + token)
	}
	if b64 {
		segments[1] = base64.RawURLEncoding.EncodeToString([]byte(body))
	} else {
		segments[1] = body
	}
	return strings.Join(segments, "."), nil
}

func getKidFromToken(token string) (string, error) {
	var tokenHeader map[string]interface{}
	segments := strings.Split(token, ".")

	decodedPayload, _ := base64.RawURLEncoding.DecodeString(segments[0])

	json.Unmarshal(decodedPayload, &tokenHeader)

	kid, ok := tokenHeader["kid"].(string)
	if !ok {
		return "", fmt.Errorf("GetKidFromToken: error getting kid string from header")
	}
	if len(kid) == 0 {
		return "", fmt.Errorf("GetKidFromToken: error kid header is zero length")
	}

	return kid, nil
}

var jwkCache = make(map[string]JWK)

// Get JWK from JWKS_URI - cache responses
func getJwkFromJwks(kid, url string) (JWK, error) {
	if jwk, ok := jwkCache[kid]; !ok {
		logrus.Traceln("Retrieving JWKS url: " + url)
		jwks, err := GetJwks(url)
		if err != nil {
			return JWK{}, fmt.Errorf("GetJwkFromJwks: errors: %v", err)
		}
		for _, k := range jwks.Keys {
			if k.Kid == kid {
				jwkCache[kid] = k
				return k, nil
			}
		}
	} else {
		logrus.Traceln("Using cached jwk")
		return jwk, nil
	}
	return JWK{}, nil
}

func ParseCertificateChain(chain []string) ([]*x509.Certificate, error) {
	certchain := make([]*x509.Certificate, len(chain))
	for i, cert := range chain {
		raw, err := base64.StdEncoding.DecodeString(cert)
		if err != nil {
			return nil, errors.New("ParseCertificateChain: decode cert: " + err.Error())
		}
		certchain[i], err = x509.ParseCertificate(raw)
		if err != nil {
			return nil, errors.New("ParseCertificateChain: parse certificate: " + err.Error())
		}
	}
	return certchain, nil
}

func GetJwks(url string) (JWKS, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return JWKS{}, fmt.Errorf("GetJwkss error retrieving url: %s, %v", url, err)
	}
	defer resp.Body.Close()
	jwksbytes := json.NewDecoder(resp.Body)
	var jwks JWKS
	if err := jwksbytes.Decode(&jwks); err != nil {
		return JWKS{}, fmt.Errorf("GetJwks: decoding error %s : %v", url, err)
	}

	return jwks, nil
}

// buildSignature - takes all the token parameters and assembles a detached header signed token string which is returned
// Handles api versions v3.1.4 and above, v3.1.3 and prior, plus v3.0 which has a slightly different JWT header
func buildSignature(b64 bool, kid, issuer, trustAnchor, body string, alg jwt.SigningMethod, privKey *rsa.PrivateKey) (string, error) {
	var token jwt.Token

	if b64 {
		token = GetSignatureToken314Plus(kid, issuer, trustAnchor, alg)
	} else {
		token = GetSignatureToken313Minus(kid, issuer, trustAnchor, alg)
	}

	tokenString, err := CreateSignature(&token, privKey, body, b64) // sign the token
	if err != nil {
		return "", errors.New("buildSignature: CreateSignature failed " + err.Error())
	}

	logrus.Tracef("Full Request JWT: %s", tokenString)

	detachedJWS := SplitJWSWithBody(tokenString) // remove the body from the signature string to form the detached signature

	return detachedJWS, nil
}

// Get Token with correct headers for v3.1.4 and above of the R/W Apis
func GetSignatureToken314Plus(kid, issuer, trustAnchor string, alg jwt.SigningMethod) jwt.Token {
	token := jwt.Token{
		Header: map[string]interface{}{
			"typ":                           "JOSE",
			"kid":                           kid,
			"cty":                           "application/json",
			"http://openbanking.org.uk/iat": time.Now().Unix(),
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
	return token
}

// Get Token with correct headers for v3.1.3 and previous versions of the R/W Apis
func GetSignatureToken313Minus(kid, issuer, trustAnchor string, alg jwt.SigningMethod) jwt.Token {
	token := jwt.Token{
		Header: map[string]interface{}{
			"typ":                           "JOSE",
			"kid":                           kid,
			"b64":                           false,
			"cty":                           "application/json",
			"http://openbanking.org.uk/iat": time.Now().Unix(),
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
	return token
}

// Read/Write Data API Specification - v3.0 Specification: https://openbanking.atlassian.net/wiki/spaces/DZ/pages/641992418/Read+Write+Data+API+Specification+-+v3.0.
// According to the spec this field `http://openbanking.org.uk/tan` should not be sent in the `x-jws-signature` header.
func GetSignatureToken30(kid, issuer, trustAnchor string, alg jwt.SigningMethod) jwt.Token {
	token := jwt.Token{
		Header: map[string]interface{}{
			"typ":                           "JOSE",
			"kid":                           kid,
			"b64":                           false,
			"cty":                           "application/json",
			"http://openbanking.org.uk/iat": time.Now().Unix(),
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
	return token
}

type signatureHeader struct {
	Type        string          `json:"typ,omitempty"`
	Kid         string          `json:"kid,omitemty"`
	Alg         string          `json:"alg,omitempty"`
	Ctype       string          `json:"cty,omitempty"`
	Issuer      string          `json:"http://openbanking.org.uk/iss,omitempty"`
	IssuedAt    decimal.Decimal `json:"http://openbanking.org.uk/iat,omitempty"`
	TrustAnchor string          `json:"http://openbanking.org.uk/tan,omitempty"`
	B64         *bool           `json:"b64,omitempty"`
	Critical    []string        `json:"crit,omitempty"`
}

func ValidateSignatureHeader(token string, b64 bool) error {
	var tokenHeader signatureHeader

	segments := strings.Split(token, ".")
	decodedPayload, _ := base64.RawURLEncoding.DecodeString(segments[0])

	logrus.Trace(string(decodedPayload))

	err := json.Unmarshal(decodedPayload, &tokenHeader)
	if err != nil {
		return fmt.Errorf("ValidateSignatureHeader: cannot convert header into JSON: " + err.Error())
	}

	dumpJSON(tokenHeader)

	err = tokenHeader.validateSignatureHeader(b64) // validate header depent on b64 setting for api true=3.1.4, false=3.1.3

	return err
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	logrus.Traceln(string(model))
}

// validate a signatureHeader structure
// according to b64=false v3.1.3 and older or
// b64=true v3.1.4 and newer
func (s signatureHeader) validateSignatureHeader(b64 bool) error {
	dumpJSON(s)

	if s.Type != "" { // Optional must be "JOSE" if present
		if s.Type != "JOSE" {
			return errors.New("Validate Signature - typ MUST equal`JOSE` if present")
		}
	}

	if s.Alg != "PS256" { // Mandatory must be "PS256"
		return errors.New("Validate Signature - alg claim MUST equal `PS256`")
	}

	if s.Kid == "" { // Mandatory - must be present
		return errors.New("Validate Signature - kid claim MUST be present")
	}

	if s.Ctype != "" { // Optional - if present must be json or application/json
		if s.Ctype != "json" && s.Ctype != "application/json" {
			return errors.New("Validate Signature - kid claim MUST be present")
		}
	}

	if b64 { // version 3.1.4 and newer
		if s.B64 != nil {
			return errors.New("Validate Signature - b64 claim MUST not be present for v3.1.4 and above APIs")
		}
		if len(s.Critical) != 3 {
			return errors.New("Validate Signature - `crit claim must contain 3 elements")
		}
		if !containsAllElements(s.Critical, "http://openbanking.org.uk/iss", "http://openbanking.org.uk/iat", "http://openbanking.org.uk/tan") {
			return errors.New("Validate Signature - `crit claim does not contain 3 values matching the specification")
		}

	} else { // version 3.1.3 and older
		if s.B64 == nil {
			return errors.New("Validate Signature - b64 claim MUST be present for v3.1.3 and older APIs")
		}
		if *s.B64 == true {
			return errors.New("Validate Signature - b64 claim MUST not equal`true` for v3.1.3 and older APIs")
		}
		if len(s.Critical) != 4 {
			return errors.New("Validate Signature - `crit claim must contain 4 elements for v3.1.3 and older APIs")
		}
		if !containsAllElements(s.Critical, "http://openbanking.org.uk/iss", "http://openbanking.org.uk/iat", "http://openbanking.org.uk/tan", "b64") {
			return errors.New("Validate Signature - `crit claim does not contain 4 values matching the specification")
		}
	}

	if s.IssuedAt == decimal.Zero {
		return errors.New("Validate Signature - `http://openbanking.org.uk/iat` claim must be a JSON number representing time")
	}
	if s.TrustAnchor != "openbanking.org.uk" {
		return errors.New("Validate Signature - `http://openbanking.org.uk/tan` claim MUST equal `openbanking.org.uk`")
	}

	if len(s.Issuer) == 0 {
		return errors.New("Validate Signature - `http://openbanking.org.uk/iss` claim MUST the ORG-ID for an ASPSP")
	}

	return nil
}

func containsAllElements(source []string, elements ...string) bool {
	for _, item := range elements {
		match := false
		for _, sourceElement := range source {
			if sourceElement == item {
				match = true
				break
			}
		}
		if match == false {
			return false
		}
	}
	return true
}

// Verify checks if the given JWS message is verifiable using `alg` and `key`.
// If the verification is successful, `err` is nil, and the content of the
// payload that was signed is returned.
func MyJwsVerify(buf string, alg jwa.SignatureAlgorithm, key interface{}, b64 bool) (ret []byte, err error) {
	verifier, err := verify.New(alg)
	if err != nil {
		return nil, errors.New("failed to create verifier")
	}
	protected, payload, signature := payloadSplit(buf)
	verifyBuf := []byte(protected + "." + payload)
	decodedSignature := make([]byte, base64.RawURLEncoding.DecodedLen(len(signature)))
	if _, err := base64.RawURLEncoding.Decode(decodedSignature, []byte(signature)); err != nil {
		return nil, errors.New(`failed to decode signature`)
	}
	if err := verifier.Verify(verifyBuf, decodedSignature, key); err != nil {
		return nil, errors.New(`failed to verify message`)
	}

	decodedPayload := make([]byte, base64.RawURLEncoding.DecodedLen(len(payload)))
	if b64 {
		if _, err := base64.RawURLEncoding.Decode(decodedPayload, []byte(payload)); err != nil {
			return nil, errors.New(`message verified, failed to decode payload`)
		}
	} else {
		decodedPayload = []byte(payload)
	}
	return decodedPayload, nil
}

// splits out a 3 part JWT into head, body, signature splitting by '.'
// Note the body may contain multiple '.' characters if its not base64 encoded (b64=false)
func payloadSplit(msg string) (head, body, sig string) {
	firstIdx := strings.IndexByte(msg, '.')
	firstPart := msg[:firstIdx]
	idx := strings.LastIndex(msg, ".")
	lastPart := msg[idx+1:]
	middle := msg[firstIdx+1 : idx]
	return firstPart, middle, lastPart
}
