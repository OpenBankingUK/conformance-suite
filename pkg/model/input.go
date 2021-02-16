package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tdewolff/minify/v2"
	minjson "github.com/tdewolff/minify/v2/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
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
	Method          string            `json:"method,omitempty"`          // http Method that this test case uses
	Endpoint        string            `json:"endpoint,omitempty"`        // resource endpoint where the http object needs to be sent to get a response
	Headers         map[string]string `json:"headers,omitempty"`         // Allows for provision of specific http headers
	RemoveHeaders   []string          `json:"removeheaders,omitempty"`   // Allows for removing specific http headers
	RemoveClaims    []string          `json:"removeClaims,omitempty"`    // Allows for removing specific signature claims
	FormData        map[string]string `json:"formData,omitempty"`        // Allow for provision of http form data
	QueryParameters map[string]string `json:"queryParameters,omitempty"` // Allow for provision of http URL query parameters
	RequestBody     string            `json:"bodyData,omitempty"`        // Optional request body raw data
	Generation      map[string]string `json:"generation,omitempty"`      // Allows for different ways of generating testcases
	Claims          map[string]string `json:"claims,omitempty"`          // collects claims for input strategies that require them
	JwsSig          bool              `json:"jws,omitempty"`             // controls inclusion of x-jws-signature header
	IdempotencyKey  bool              `json:"idempotency,omitempty"`     // specifices the inclusion of x-idempotency-key in the request
}

var disableJws = false // defaults to JWS disabled in line with waiver 007
var b64Status bool     // store b64 value for report

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
		req.SetBody(body)
	}

	for k, v := range i.QueryParameters {
		req.QueryParam.Add(k, v)
	}

	if i.JwsSig {
		// create jws detached signature - add to headers
		if i.Method == "POST" {
			err := i.createJWSDetachedSignature(ctx)
			if err != nil {
				logrus.Tracef("error creating detached signature: %s", err)
				return nil, err
			}
		} else {
			return nil, errors.New("cannot apply jws signature to method that isn't POST")
		}
	}

	if i.IdempotencyKey {
		i.SetHeader("x-idempotency-key", tc.ID+"-"+makeMiliSecondStringTimestamp()) // initial trivial x-idempotency-key implementation
	}

	if err = i.removeHeaders(); err != nil {
		return nil, err
	}

	if err = i.setHeaders(req, ctx); err != nil {
		return nil, err
	}

	req.Method = tc.Input.Method
	req.URL = tc.Input.Endpoint
	return req, nil
}

func (i *Input) setClaims(tc *TestCase, ctx *Context) error {
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
				tc.DoNotCallEndpoint = true
			}
		}
	}

	return nil
}

func (i *Input) GenerateRequestToken(ctx *Context) (string, error) {
	alg, err := ctx.GetString("requestObjectSigningAlg")
	if err != nil && err != ErrNotFound {
		return "", err
	}
	signingMethod, err := authentication.GetSigningAlg(alg)
	if err != nil {
		logrus.Warnln("Using Unsigned Jwt as Request Object")
		return i.generateUnsignedJWT(ctx) // not sure if this is still appropraite
	}
	return i.generateRequestJWT(ctx, signingMethod)
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

func (i *Input) removeHeaders() error {
	remainingHeaders := make(map[string]string, 0)
	if len(i.RemoveHeaders) > 0 {
		i.AppMsg(fmt.Sprintf("RemoveHeaders %v", i.RemoveHeaders))
	}

	var found bool
	for x, y := range i.Headers {
		for _, v := range i.RemoveHeaders {
			if strings.EqualFold(v, x) {
				found = true
				i.AppMsg("removing header: " + x)
				break
			}
		}
		if !found {
			remainingHeaders[x] = y
		}
		found = false
	}

	i.Headers = remainingHeaders
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

func (i *Input) createJWSDetachedSignature(ctx authentication.ContextInterface) error {
	if len(i.RequestBody) == 0 {
		return i.AppErr("cannot create x-jws-signature, as request body is empty")
	}

	if disableJws {
		i.AppMsg("x-jws-signature disabled")
		return nil
	}

	requestObjSigningAlg, err := ctx.GetString("requestObjectSigningAlg")
	if err != nil {
		return errors.Wrap(err, "input.createJWSDetachedSignature: unable to retrieve requestObjectSigningAlg")
	}

	alg, err := authentication.GetSigningAlg(requestObjSigningAlg)
	if err != nil {
		return errors.Wrapf(err, "input.createJWSDetachedSignature: unable to parse signing alg")
	}

	token, err := authentication.NewJWSSignature(i.RequestBody, ctx, alg)
	if err != nil {
		return i.AppErr(fmt.Sprintf("error generating jws signature %s", err.Error()))
	}

	if len(i.RemoveClaims) > 0 {
		token, err = authentication.ModifyJWSHeaders(token, ctx, authentication.RemoveJWSHeader(i.RemoveClaims))
		if err != nil {
			return err
		}
	}

	i.SetHeader("x-jws-signature", token)
	return nil
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

func (i *Input) generateRequestJWT(ctx *Context, alg jwt.SigningMethod) (string, error) {
	uuid := uuid.New()
	claims := jwt.MapClaims{}
	if iss, ok := i.Claims["iss"]; ok {
		claims["iss"] = iss
		claims["client_id"] = iss
	}

	if sub, ok := i.Claims["sub"]; ok {
		claims["sub"] = sub
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
	if state, ok := i.Claims["state"]; ok {
		claims["state"] = state
	}

	consentClaim := consentClaims{Essential: true, Value: i.Claims["consentId"]}
	myIdent := obIDToken{IntentID: consentClaim}
	acrValuesSupported, err := ctx.GetStringSlice("acrValuesSupported")
	if err == nil && len(acrValuesSupported) > 0 {
		myIdent = obIDToken{IntentID: consentClaim, Acr: &acr{Essential: true, Values: acrValuesSupported}}
	}
	var consentIDToken = consentIDTok{Token: myIdent}
	claims["claims"] = consentIDToken
	if responseType, ok := i.Claims["responseType"]; ok {
		claims["response_type"] = responseType
	}

	logrus.WithFields(logrus.Fields{
		"claims":   claims,
		"alg":      alg,
		"i.Claims": i.Claims,
	}).Debug("generateRequestJWT ...")

	token := jwt.NewWithClaims(alg, claims) // create new token

	cert, err := signingCertFromContext(ctx)
	if err != nil {
		return "", i.AppErr(errors.Wrap(err, "Create certificate from context").Error())
	}

	kid, err := ctx.GetString("tpp_signature_kid")
	if err != nil {
		return "", errors.Wrap(err, "generateRequestJWT failure: unable to get KID")
	}
	logrus.WithFields(logrus.Fields{
		"kid": kid,
	}).Debug("generateRequestJWT")
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(cert.PrivateKey()) // sign the token - get as encoded string
	if err != nil {
		return "", i.AppErr(fmt.Sprintf("error siging jwt: %s", err.Error()))
	}
	return tokenString, nil
}

// acr
// TPPs MAY provide a space-separated string that specifies the acr values that the Authorization Server is being requested to use for processing this Authentication Request, with the values appearing in order of preference.
// The values MUST be one or both of:
// urn:openbanking:psd2:sca: To indiciate that secure customer authentication must be carried out as mandated by the PSD2 RTS
// urn:openbanking:psd2:ca: To request that the customer is authenticated without using SCA (if permitted)
//
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/7046134/Open+Banking+Security+Profile+-+Implementer+s+Draft+v1.1.0
type acr struct {
	Essential bool     `json:"essential,omitempty"`
	Values    []string `json:"values,omitempty"`
}

type obIDToken struct {
	IntentID consentClaims `json:"openbanking_intent_id,omitempty"`
	Acr      *acr          `json:"acr,omitempty"`
}

type consentClaims struct {
	Essential bool   `json:"essential"`
	Value     string `json:"value"` // account-requestid
}

type consentIDTok struct {
	Token obIDToken `json:"id_token,omitempty"`
}

// setAdditonalClaims - read the comments in the function body, please.
func setAdditonalClaims(jwtClaims jwt.MapClaims, inputClaims map[string]string) {
	// https://openid.net/specs/openid-financial-api-part-2.html#authorization-server Pt 13
	// Ozone recently (the time I write this is 22/07/2019) went through FAPI certification so it is now more strict.
	//
	// We get an error when doing headless if the `exp` is not set. Here is the actual log:
	//
	// ---------------------- REQUEST LOG -----------------------
	// GET  /auth?client_id=081756dd-17f5-4543-a221-012e7ec8694e&redirect_uri=https%3A%2F%2F127.0.0.1%3A8443%2Fconformancesuite%2Fcallback&request=eyJhbGciOiJub25lIn0.eyJhdWQiOiJodHRwczovL29iMTktYXV0aDEtdWkubzNiYW5rLmNvLnVrIiwiY2xhaW1zIjp7ImlkX3Rva2VuIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJlc3NlbnRpYWwiOnRydWUsInZhbHVlIjoiYWFjLTdlNDE5ZTY3LWNiYTMtNGNlOS1iMzRkLTM3YjdmYjY4MDYzNyJ9fX0sImlzcyI6IjA4MTc1NmRkLTE3ZjUtNDU0My1hMjIxLTAxMmU3ZWM4Njk0ZSIsInJlZGlyZWN0X3VyaSI6Imh0dHBzOi8vMTI3LjAuMC4xOjg0NDMvY29uZm9ybWFuY2VzdWl0ZS9jYWxsYmFjayIsInNjb3BlIjoib3BlbmlkIGFjY291bnRzIn0.&response_type=code&scope=openid+accounts&state=  HTTP/1.1
	// HOST   : ob19-auth1-ui.o3bank.co.uk
	// HEADERS:
	// 				User-Agent: go-resty/1.10.3 (https://github.com/go-resty/resty)
	// BODY   :
	// ***** NO CONTENT *****
	// ----------------------------------------------------------
	// RESTY 2019/07/19 10:29:01
	// ---------------------- RESPONSE LOG -----------------------
	// STATUS 		: 400 Bad Request
	// RECEIVED AT	: 2019-07-19T10:29:01.086224+01:00
	// RESPONSE TIME	: 134.991846ms
	// HEADERS:
	// 				Connection: keep-alive
	// 			Content-Length: 194
	// 				Content-Type: application/json; charset=utf-8
	// 						Date: Fri, 19 Jul 2019 09:29:01 GMT
	// 						Etag: W/"c2-g632OX8LPgukAE/rB/n5eQtwF3w"
	// 				X-Powered-By: Express
	// BODY   :
	// {
	// 	"noRedirect": true,
	// 	"error": "invalid_request",
	// 	"error_description": "request_object_exp_undefined: The request object must have an exp claim",
	// 	"interactionId": "bd9fe4d9-9014-4c04-ae42-db053024f265"
	// }
	// ----------------------------------------------------------
	//
	// To fix this, if `claims` has the `exp` set to `true`, we set the `exp` claim to 30 minutes from now.
	// "claims": {
	// 	...
	// 	"exp": "true",
	// 	...
	// }
	//
	// We do the same thing for `nonce`
	logger := logrus.WithFields(logrus.Fields{
		"package":     "model",
		"function":    "authEndpoint",
		"inputClaims": inputClaims,
	})

	if exp, ok := inputClaims["exp"]; ok {
		if strings.ToLower(exp) == "true" {
			jwtClaims["exp"] = time.Now().Add(time.Minute * time.Duration(30)).Unix()
			logger.WithFields(logrus.Fields{
				`jwtClaims["exp"]`: jwtClaims["exp"],
			}).Debug(`setting "exp" claim because "exp" == "true"`)
		}
	}

	if nonce, ok := inputClaims["nonce"]; ok {
		if strings.ToLower(nonce) == "true" {
			uuid := uuid.New()
			jwtClaims["nonce"] = uuid
			logger.WithFields(logrus.Fields{
				`jwtClaims["nonce"]`: jwtClaims["nonce"],
			}).Debug(`setting "nonce" claim because "nonce" == "true"`)
		}
	}

	if state, ok := inputClaims["state"]; ok {
		if len(state) > 0 {
			jwtClaims["state"] = state
			logger.WithFields(logrus.Fields{
				`jwtClaims["state"]`: jwtClaims["state"],
			}).Debug(`setting "state" claim because len("state") > 0`)
		}
	}
}

func (i *Input) generateUnsignedJWT(ctx *Context) (string, error) {
	claims := jwt.MapClaims{}
	claims["iss"] = i.Claims["iss"]
	claims["scope"] = i.Claims["scope"]
	claims["aud"] = i.Claims["aud"]
	claims["redirect_uri"] = i.Claims["redirect_url"]
	setAdditonalClaims(claims, i.Claims)

	consentClaim := consentClaims{Essential: true, Value: i.Claims["consentId"]}
	myIdent := obIDToken{IntentID: consentClaim}
	acrValuesSupported, err := ctx.GetStringSlice("acrValuesSupported")
	if err == nil && len(acrValuesSupported) > 0 {
		myIdent = obIDToken{IntentID: consentClaim, Acr: &acr{Essential: true, Values: acrValuesSupported}}
	}
	var consentIDToken = consentIDTok{Token: myIdent}

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

func JWSStatus() string {
	if disableJws {
		return "disabled"
	}
	if authentication.GetB64Status() {
		return "true"
	}
	return "false"
}
