package model

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/pkg/errors"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/resty.v1"
)

// Input defines the content of the http request object used to execute the test case
// Input is built up typically from the openapi/swagger definition of the method/endpoint for a particular
// specification. Additional properties/fields/headers can be added or change in order to setup the http
// request object of the specific test case. Once setup correctly,the testcase gives the http request object
// to the parent Rule which determine how to execute the requestion object. On execution an http response object
// is received and passed back to the testcase for validation using the Expects object.
type Input struct {
	Method         string            `json:"method,omitempty"`      // http Method that this test case uses
	Endpoint       string            `json:"endpoint,omitempty"`    // resource endpoint where the http object needs to be sent to get a response
	Headers        map[string]string `json:"headers,omitempty"`     // Allows for provision of specific http headers
	FormData       map[string]string `json:"formData,omitempty"`    // Allow for provision of http form data
	RequestBody    string            `json:"bodyData,omitempty"`    // Optional request body raw data
	Generation     map[string]string `json:"generation,omitempty"`  // Allows for different ways of generating testcases
	Claims         map[string]string `json:"claims,omitempty"`      // collects claims for input strategies that require them
	JwsSig         bool              `json:"jws,omitempty"`         // controls inclusion of x-jws-signature header
	IdempotencyKey bool              `json:"idempotency,omitempty"` // specifices the inclusion of x-idempotency-key in the request
}

var disableJws bool // Allows disabling of jws for ozone - effectively a development flag

// CreateRequest is the main Input work horse which examines the various Input parameters and generates an
// http.Request object which represents the request
func (i *Input) CreateRequest(tc *TestCase, ctx *Context) (*resty.Request, error) {
	var err error

	if tc == nil {
		return nil, i.AppErr(fmt.Sprintf("error CreateRequest - Testcase is nil"))
	}

	if ctx == nil {
		return nil, i.AppErr(fmt.Sprintf("error CreateRequest - Context is nil"))
	}

	if i.Endpoint == "" || i.Method == "" { // we don't have a value input object
		return nil, i.AppErr(fmt.Sprintf("error empty Endpoint(%s) or Method(%s)", i.Endpoint, i.Method))
	}

	req := resty.R() // create basic request that will be sent to endpoint

	tc.Input.Endpoint, err = replaceContextField(tc.Input.Endpoint, ctx)
	if err != nil {
		return nil, err
	}

	if err = i.setClaims(tc, ctx); err != nil {
		return nil, err
	}

	if err = i.setFormData(req, ctx); err != nil {
		return nil, err
	}

	if len(i.RequestBody) > 0 { // set any input raw request body ("bodyData")
		body, err := i.getBody(req, ctx)
		if err != nil {
			return nil, err
		}
		i.RequestBody = body
		logrus.Tracef("setting resty body to: %s", body)
		req.SetBody(body)
	}

	if i.JwsSig {
		// create jws detached signature - add to headers
		if i.Method == "POST" {
			err := i.createJWSDetachedSignature(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("cannot apply jws signature to method that isn't POST")
		}
	}

	if i.IdempotencyKey {
		i.SetHeader("x-idempotency-key", tc.ID+"-"+makeMiliSecondStringTimestamp()) // initial trivial x-idempotency-key implementation
	}

	if err = i.setHeaders(req, ctx); err != nil {
		return nil, err
	}

	req.Method = tc.Input.Method
	req.URL = tc.Input.Endpoint
	return req, nil
}

func (i *Input) setClaims(tc *TestCase, ctx *Context) error {
	logrus.WithFields(logrus.Fields{
		"i.Claims": i.Claims,
	}).Debug("Input.setClaims before ...")
	ctx.DumpContext()

	for k, v := range i.Claims {
		value, err := replaceContextField(v, ctx)
		if err != nil {
			return i.AppErr(fmt.Sprintf("setClaims Replace Context value %s :%s", v, err.Error()))
		}
		i.Claims[k] = value
		i.AppMsg(fmt.Sprintf("Claims [%s:%s]", k, i.Claims[k]))
	}

	if len(i.Claims) > 0 { // create JWT from claims - put in context?
		switch i.Generation["strategy"] {
		case "psuConsenturl":
			fallthrough
		case "consenturl":
			i.AppMsg("==> executing consenturl strategy")
			token, err := i.GenerateRequestToken(ctx)
			if err != nil {
				return i.AppErr(fmt.Sprintf("error creating request token %s", err.Error()))
			}
			i.AppMsg(fmt.Sprintf("jwt consent Token: %s", token))

			authEndpoint, _ := ctx.Get("authorisation_endpoint")
			consent := consentURL(authEndpoint.(string), i.Claims, token)

			tc.Input.Endpoint = consent           // Result - set jwt token in endpoint url
			ctx.PutString("consent_url", consent) // make consent available in context
			logrus.Tracef("===> consentURL: " + consent)
			i.AppMsg("consent url: " + tc.Input.Endpoint)
			if i.Generation["strategy"] == "psuConsenturl" {
				tc.Input.Endpoint = i.Claims["aud"] + "/PsuDummyURL"
			}
		case "jwt-bearer":
			i.AppMsg("==> executing jwt-bearer strategy")
			token, err := i.GenerateSignedJWT(ctx, jwt.SigningMethodRS256)
			if err != nil {
				return i.AppErr(fmt.Sprintf("error creating AlgRS256JWT %s", err.Error()))
			}
			i.AppMsg(fmt.Sprintf("jwt-bearer Token: %s", token))
			ctx.Put("jwtbearer", token) // Result - set jwt-bearer token in context
		}
	}

	logrus.WithFields(logrus.Fields{
		"i.Claims": i.Claims,
	}).Debug("Input.setClaims after ...")
	ctx.DumpContext()

	return nil
}

func (i *Input) GenerateRequestToken(ctx *Context) (string, error) {
	alg, err := ctx.GetString("requestObjectSigningAlg")
	if err != nil && err != ErrNotFound {
		return "", err
	}
	signingMethod, err := getSigningMethod(alg)
	if err != nil {
		return i.generateUnsignedJWT()
	}
	return i.GenerateSignedJWT(ctx, signingMethod)
}

func consentURL(authEndpoint string, claims map[string]string, token string) string {
	queryString := url.Values{}
	queryString.Set("client_id", claims["iss"])
	queryString.Set("response_type", claims["responseType"])
	queryString.Set("scope", claims["scope"])
	queryString.Set("request", token)
	queryString.Set("state", claims["state"])
	queryString.Set("redirect_uri", claims["redirect_url"])

	consentURL := fmt.Sprintf("%s?%s", authEndpoint, queryString.Encode())

	logrus.WithFields(logrus.Fields{
		"claims":       claims,
		"authEndpoint": authEndpoint,
		"token":        token,
		"consentURL":   consentURL,
	}).Trace("Generating consentURL")

	return consentURL
}

func (i *Input) setFormData(req *resty.Request, ctx *Context) error {
	if len(i.FormData) > 0 {
		i.AppMsg(fmt.Sprintf("AddFormData %v", i.FormData))
		for k, v := range i.FormData {
			value, err := replaceContextField(v, ctx)
			if err != nil {
				return i.AppErr(fmt.Sprintf("setFormdata Replace Context value %s :%s", v, err.Error()))
			}
			if len(value) == 0 {
				return i.AppErr(fmt.Sprintf("setFormdata Replace Context value %s - empty", v))
			}
			i.FormData[k] = value
		}
		req.SetFormData(i.FormData)
	}
	return nil
}

func (i *Input) setHeaders(req *resty.Request, ctx *Context) error {
	if len(i.Headers) > 0 {
		i.AppMsg(fmt.Sprintf("SetHeaders %v", i.Headers))
	}
	for k, v := range i.Headers { // set any input headers ("headers")
		value, err := replaceContextField(v, ctx)
		if err != nil {
			return i.AppErr(fmt.Sprintf("setHeaders Replaced Context value %s :%s", v, err.Error()))
		}
		if len(value) == 0 {
			return i.AppErr(fmt.Sprintf("setHeaders Replaced Context value %s:%s not found in context", k, v))
		}
		req.SetHeader(k, value)
	}
	return nil
}

func (i *Input) createJWSDetachedSignature(ctx *Context) error {

	if len(i.RequestBody) > 0 && !disableJws {
		requestObjSigningAlg, err := ctx.GetString("requestObjectSigningAlg")
		if err != nil {
			return errors.Wrap(err, "input.createJWSDetachedSignature failure: unable to retrieve requestObjectSigningAlg")
		}
		alg, err := getSigningMethod(requestObjSigningAlg)
		if err != nil {
			return errors.Wrapf(err, "input.createJWSDetachedSignature failure: unable to parse signing alg")
		}
		token, err := i.generateJWSSignature(ctx, alg)

		if err != nil {
			return i.AppErr(fmt.Sprintf("error generating jws signature %s", err.Error()))
		}
		i.SetHeader("x-jws-signature", token)

		return nil
	}

	if disableJws {
		i.AppMsg("x-jws-signature disabled")
		return nil
	}
	return i.AppErr("cannot create x-jws-signature, as request body is empty")

}

func getSigningMethod(alg string) (jwt.SigningMethod, error) {
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
		return nil, fmt.Errorf("unable to find signing algorithm %s", alg)
	}
}

func (i *Input) getBody(req *resty.Request, ctx *Context) (string, error) {
	value := i.RequestBody
	for {
		val2, err := replaceContextField(value, ctx)
		if err != nil {
			return "", i.AppErr(fmt.Sprintf("setBody Replaced Context value %s :%s", val2, err.Error()))
		}
		if len(val2) == 0 {
			return "", i.AppErr(fmt.Sprintf("setBody Replaced Context value %s : %s not found in context", value, i.RequestBody))
		}
		if val2 == value {
			break
		}
		value = val2
	}

	body := value
	contentType := i.contentTypeHeader()
	if strings.Contains(contentType, "application/json") {
		m := minify.New()
		m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minjson.Minify)
		var err error
		body, err = m.String("application/json", value)
		if err != nil {
			return "", err
		}
	}
	i.RequestBody = body
	req.SetBody(body)

	return body, nil
}

func (i *Input) contentTypeHeader() string {
	for key, value := range i.Headers {
		if strings.ToLower(key) == "content-type" {
			return value
		}
	}
	return ""
}

// AppMsg - application level trace
func (i *Input) AppMsg(msg string) string {
	tracer.AppMsg("Input", msg, i.String())
	return msg
}

// AppErr - application level trace error msg
func (i *Input) AppErr(msg string) error {
	tracer.AppErr("Input", msg, i.String())
	return errors.New(msg)
}

// String - object represetation
func (i *Input) String() string {
	bites, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		// String() doesn't return error but still want to log as error to tracer ...
		return i.AppErr(fmt.Sprintf("error converting Input %s %s %s", i.Method, i.Endpoint, err.Error())).Error()
	}
	return string(bites)
}

// Clone an Input
func (i *Input) Clone() Input {
	in := Input{}
	in.Endpoint = i.Endpoint
	in.FormData = i.FormData
	in.Generation = i.Generation
	in.Headers = i.Headers
	in.Method = i.Method
	in.RequestBody = i.RequestBody
	in.Claims = i.Claims

	return in
}

func signingCertFromContext(ctx *Context) (authentication.Certificate, error) {
	privKey, err := ctx.GetString("signingPrivate")
	if err != nil {
		return nil, errors.New("input, couldn't find `SigningPrivate` in context")
	}
	pubKey, err := ctx.GetString("signingPublic")
	if err != nil {
		return nil, errors.New("input, couldn't find `SigningPublic` in context")
	}
	cert, err := authentication.NewCertificate(pubKey, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "input, couldn't create `certificate` from pub/priv keys")
	}
	return cert, nil
}

func (i *Input) GenerateSignedJWT(ctx *Context, alg jwt.SigningMethod) (string, error) {
	uuid := uuid.New()
	claims := jwt.MapClaims{}
	if iss, ok := i.Claims["iss"]; ok {
		claims["iss"] = iss
	}
	if iss, ok := i.Claims["iss"]; ok {
		claims["sub"] = iss
	}
	if scope, ok := i.Claims["scope"]; ok {
		claims["scope"] = scope
	}
	if aud, ok := i.Claims["aud"]; ok {
		claims["aud"] = aud
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(30)).Unix()
	claims["jti"] = uuid
	claims["nonce"] = uuid

	if redirectURI, ok := i.Claims["redirect_url"]; ok {
		claims["redirect_uri"] = redirectURI
	}
	if iss, ok := i.Claims["iss"]; ok {
		claims["client_id"] = iss
	}
	if state, ok := i.Claims["state"]; ok {
		claims["state"] = state
	}

	consentClaim := consentClaims{Essential: true, Value: i.Claims["consentId"]}
	myIdent := obintentID{IntentID: consentClaim}
	var consentIDToken = consentIDTok{Token: myIdent}
	claims["claims"] = consentIDToken
	if responseType, ok := i.Claims["responseType"]; ok {
		claims["response_type"] = responseType
	}

	logrus.WithFields(logrus.Fields{
		"claims":   claims,
		"alg":      alg,
		"i.Claims": i.Claims,
	}).Debug("Input.generateSignedJWT ...")

	token := jwt.NewWithClaims(alg, claims) // create new token

	cert, err := signingCertFromContext(ctx)
	if err != nil {
		return "", i.AppErr(errors.Wrap(err, "Create certificate from context").Error())
	}
	kid, err := GetKID(ctx, cert.PublicKey().N.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to get KID")
	}
	logrus.WithFields(logrus.Fields{
		"kid": kid,
	}).Debug("GenerateSignedJWT")
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(cert.PrivateKey()) // sign the token - get as encoded string
	if err != nil {
		return "", i.AppErr(fmt.Sprintf("error siging jwt: %s", err.Error()))
	}
	return tokenString, nil
}

type payload []byte

func (p payload) Valid() error {
	return nil
}

func (i *Input) generateJWSSignature(ctx *Context, alg jwt.SigningMethod) (string, error) {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minjson.Minify)
	minifiedBody, err := m.String("application/json", i.RequestBody)
	if err != nil {
		return "", err
	}
	cert, err := signingCertFromContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to sign certificate from context")
	}
	modulus := cert.PublicKey().N.Bytes()
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, _ := authentication.CalcKid(modulusBase64)

	issuer, err := i.getJWSIssuerString(ctx, cert)
	if err != nil {
		return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to retrieve issuer from context")
	}
	trustAnchor := "openbanking.org.uk"
	useNonOBDirectory, exists := ctx.Get("nonOBDirectory")
	if !exists {
		return "", errors.New("model.Input.generateJWSSignature failure: unable to retrieve nonOBDirectory from context")
	}
	useNonOBDirectoryAsBool, ok := useNonOBDirectory.(bool)
	if !ok {
		return "", errors.New("model.Input.generateJWSSignature failure: unable to cast nonOBDirectory to bool")
	}
	if useNonOBDirectoryAsBool {
		kid, err = ctx.GetString("signingKid")
		if err != nil {
			return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to retrieve singingKid from context")
		}
		issuer, err = ctx.GetString("issuer")
		if err != nil {
			return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to retrieve issue from context")
		}
		trustAnchor, err = ctx.GetString("signatureTrustAnchor")
		if err != nil {
			return "", errors.Wrap(err, "model.Input.generateJWSSignature failure: unable to retrieve signatureTrustAnchor from context")
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
	detachedJWS := splitJwsWithBody(tokenString)
	logrus.Tracef("detached jws: %v", detachedJWS)

	return detachedJWS, nil
}

func (i *Input) getJWSIssuerString(ctx *Context, cert authentication.Certificate) (string, error) {

	apiVersion, err := ctx.GetString("api-version")
	if err != nil {
		return "", errors.New("generate JWS Signature - cannot find api-version: " + err.Error())
	}
	var issuer string
	if apiVersion == "v3.1" {
		issuer, err = cert.SignatureIssuer(true)
		if err != nil {
			logrus.Warn("cannot Issuer for Signature: ", err.Error())
			return "", errors.New("cannot Issuer for Signature: " + err.Error())
		}
	} else if apiVersion == "v3.0" {
		issuer, err = cert.DN()
		if err != nil {
			logrus.Warn("cannot get certificate DN: ", err.Error())
		}
	} else {
		return "", errors.New("cannot get issuer for jws signature but api-version doesn't match 3.0.0 or 3.1.0")
	}

	return issuer, nil
}

func splitJwsWithBody(token string) string {
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
		return "", err
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", err
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
				return "", err
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

type obintentID struct {
	IntentID consentClaims `json:"openbanking_intent_id,omitempty"`
}

type consentClaims struct {
	Essential bool   `json:"essential"`
	Value     string `json:"value"` // account-requestid
}

type consentIDTok struct {
	Token obintentID `json:"id_token,omitempty"`
}

func (i *Input) generateUnsignedJWT() (string, error) {
	claims := jwt.MapClaims{}
	claims["iss"] = i.Claims["iss"]
	claims["scope"] = i.Claims["scope"]
	claims["aud"] = i.Claims["aud"]
	claims["redirect_uri"] = i.Claims["redirect_url"]

	consentClaim := consentClaims{Essential: true, Value: i.Claims["consentId"]}
	myident := obintentID{IntentID: consentClaim}
	var consentIDToken = consentIDTok{Token: myident}

	claims["claims"] = consentIDToken

	alg := jwt.SigningMethodNone
	if alg == nil {
		return "", errors.New(i.AppMsg(fmt.Sprintf("no signing method: %v", alg)))
	}

	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": alg.Alg(),
		},
		Claims: claims,
		Method: alg,
	}

	tokenString, err := token.SigningString() // sign the token - get as encoded string
	if err != nil {
		i.AppMsg(fmt.Sprintf("error signing jwt: %s", err.Error()))
		return "", err
	}
	tokenString += "."
	return tokenString, nil
}

// SetHeader - on the testcase input object
func (i *Input) SetHeader(key, value string) {
	if i.Headers == nil {
		i.Headers = map[string]string{}
	}
	i.Headers[key] = value
}

// SetFormField - sets a field in the form
func (i *Input) SetFormField(key, value string) {
	if i.FormData == nil {
		i.FormData = map[string]string{}
	}
	i.FormData[key] = value
}

func makeMiliSecondStringTimestamp() string {
	a := time.Now().UnixNano() / int64(time.Millisecond)
	milliString := fmt.Sprintf("%d", a)
	return milliString
}

// DisableJWS - disable jws-signature for ozone
func DisableJWS() {
	disableJws = true
}

// GetKID determines the value of the JWS Key ID
func GetKID(ctx *Context, modulus []byte) (string, error) {
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, err := authentication.CalcKid(modulusBase64)
	if err != nil {
		return "", errors.Wrap(err, "could not calculate kid")
	}
	nonOBDirectory, exists := ctx.Get("nonOBDirectory")
	if !exists {
		return "", errors.New("unable get nonOBDirectory value from context")
	}
	nonOBDirectoryAsBool, ok := nonOBDirectory.(bool)
	if !ok {
		return "", errors.New("unable to cast nonOBDirectory value to bool")
	}
	if nonOBDirectoryAsBool {
		kid, err = ctx.GetString("signingKid")
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve signingKid from context")
		}
	}

	return kid, nil
}
