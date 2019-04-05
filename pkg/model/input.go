package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

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
	Method      string            `json:"method,omitempty"`     // http Method that this test case uses
	Endpoint    string            `json:"endpoint,omitempty"`   // resource endpoint where the http object needs to be sent to get a response
	Headers     map[string]string `json:"headers,omitempty"`    // Allows for provision of specific http headers
	FormData    map[string]string `json:"formData,omitempty"`   // Allow for provision of http form data
	RequestBody string            `json:"bodyData,omitempty"`   // Optional request body raw data
	Generation  map[string]string `json:"generation,omitempty"` // Allows for different ways of generating testcases
	Claims      map[string]string `json:"claims,omitempty"`     // collects claims for input strategies that require them
}

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

	if err = i.setHeaders(req, ctx); err != nil {
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

	req.Method = tc.Input.Method
	req.URL = tc.Input.Endpoint
	return req, nil
}

func (i *Input) setClaims(tc *TestCase, ctx *Context) error {
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
			token, err := i.createAlgNoneJWT()
			if err != nil {
				return i.AppErr(fmt.Sprintf("error creating AlgNoneJWT %s", err.Error()))
			}
			i.AppMsg(fmt.Sprintf("jwt consent Token: %s", token))
			consent := i.Claims["aud"] + "/auth?" + "client_id=" + i.Claims["iss"] + "&response_type=" + i.Claims["responseType"] + "&scope=" + url.QueryEscape(i.Claims["scope"]) + "&request=" + token + "&state=" + i.Claims["state"]

			tc.Input.Endpoint = consent           // Result - set jwt token in endpoint url
			ctx.PutString("consent_url", consent) // make consent available in context
			logrus.Tracef("===> consentURL: " + consent)
			i.AppMsg("consent url: " + tc.Input.Endpoint)
			if i.Generation["strategy"] == "psuConsenturl" {
				tc.Input.Endpoint = i.Claims["aud"] + "/PsuDummyURL"
			}
		case "jwt-bearer":
			i.AppMsg("==> executing jwt-bearer strategy")
			token, err := i.createAlgRS256JWT(ctx)
			if err != nil {
				return i.AppErr(fmt.Sprintf("error creating AlgRS256JWT %s", err.Error()))
			}
			i.AppMsg(fmt.Sprintf("jwt-bearer Token: %s", token))
			ctx.Put("jwtbearer", token) // Result - set jwt-bearer token in context
		}
	}

	return nil
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

func (i *Input) getBody(_ *resty.Request, ctx *Context) (string, error) {
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
	return value, nil
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

func (i *Input) createAlgRS256JWT(ctx *Context) (string, error) {
	uuid := uuid.New()
	claims := jwt.MapClaims{}
	claims["iss"] = i.Claims["iss"]
	claims["sub"] = i.Claims["iss"]
	claims["scope"] = i.Claims["scope"]
	claims["aud"] = i.Claims["aud"]
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(60)).Unix()
	claims["jti"] = uuid

	alg := jwt.GetSigningMethod("RS256")
	if alg == nil {
		msg := fmt.Sprintf("couldn't find RS256 signing method: %v", alg)
		logrus.StandardLogger().Error(msg)
		return "", errors.New(msg)
	}
	token := jwt.NewWithClaims(alg, claims) // create new token
	token.Header["kid"] = i.Claims["kid"]

	pk, ok := ctx.Get("SigningCert")
	if !ok {
		return "", i.AppErr(fmt.Sprintf("input, couldn't find `SigningCert` in context"))
	}
	cert, ok := pk.(authentication.Certificate)
	if !ok {
		return "", i.AppErr(fmt.Sprintf("input, cannot convert `SigningCert` to certificate"))
	}
	tokenString, err := token.SignedString(cert.PrivateKey()) // sign the token - get as encoded string
	if err != nil {
		return "", i.AppErr(fmt.Sprintf("error siging jwt: %s", err.Error()))
	}
	logrus.StandardLogger().Debugf("\nCreated JWT:\n-------------\n%s\n", tokenString)
	return tokenString, nil
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

// Initial implementation of JWT creation with algorithm 'None'
// Used only to support the PSU consent URL generation for headless consent flow
func (i *Input) createAlgNoneJWT() (string, error) {
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
