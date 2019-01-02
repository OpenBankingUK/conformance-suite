package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	jwt "github.com/dgrijalva/jwt-go"
	resty "gopkg.in/resty.v1"
)

// Input defines the content of the http request object used to execute the test case
// Input is built up typically from the openapi/swagger definition of the method/endpoint for a particualar
// specification. Additional properties/fields/headers can be added or change in order to setup the http
// request object of the specific test case. Once setup correctly,the testcase gives the http request object
// to the parent Rule which determine how to execute the requestion object. On execution an http response object
// is received and passed back to the testcase for validation using the Expects object.
type Input struct {
	Method      string            `json:"method,omitempty"`     // http Method that this test case uses
	Endpoint    string            `json:"endpoint,omitempty"`   // resource endpoint where the http object needs to be sent to get a response
	ContextGet  ContextAccessor   `json:"contextGet,omitempty"` // Allows retrieval of context variables an input parameters
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
		return nil, errors.New(i.AppErr(fmt.Sprintf("error CreateRequest - Testcase is nil")))
	}

	if ctx == nil {
		return nil, errors.New(i.AppErr(fmt.Sprintf("error CreateRequest - Context is nil")))
	}

	if i.Endpoint == "" || i.Method == "" { // we don't have a value input object
		return nil, errors.New(i.AppErr(fmt.Sprintf("error empty Endpoint(%s) or Method(%s)", i.Endpoint, i.Method)))
	}

	if err = i.ContextGet.GetValues(tc, ctx); err != nil { // look for endpoint replacment strings
		return nil, err
	}

	req := resty.R() // create basic request that will be sent to endpoint

	if err = i.setHeaders(req, ctx); err != nil {
		return nil, err
	}

	if err = i.setFormData(req, ctx); err != nil {
		return nil, err
	}

	if len(i.RequestBody) > 0 { // set any input raw request body ("bodyData")
		req.SetBody(i.RequestBody)
	}

	if err = i.setClaims(tc, ctx); err != nil {
		return nil, err
	}

	req.Method = tc.Input.Method
	req.URL = tc.Input.Endpoint

	return req, nil
}

func (i *Input) setClaims(tc *TestCase, ctx *Context) error {
	var err error
	for k, v := range i.Claims {
		i.Claims[k], err = i.expandContextVariable(v, ctx)
		if err != nil {
			return err
		}

		i.AppMsg(fmt.Sprintf("Claims [%s:%s]", k, i.Claims[k]))
	}

	if len(i.Claims) > 0 { // create JWT from claims - put in context?
		if i.Generation["strategy"] == "consenturl" {
			i.AppMsg("==> executing consenturl strategy")
			token, err := i.createAlgNoneJWT()
			if err != nil {
				i.AppMsg(err.Error())
			}
			_ = token
			i.AppMsg(fmt.Sprintf("jwt consent Token: %s", token))
			consent := i.Claims["aud"] + "/auth?" + "client_id=" + i.Claims["iss"] + "&response_type=" + i.Claims["responseType"] + "&scope=" + url.QueryEscape(i.Claims["scope"]) + "&request=" + token

			tc.Input.Endpoint = consent
			i.AppMsg("consent url: " + tc.Input.Endpoint)
		}
	}

	return nil
}

func (i *Input) setFormData(req *resty.Request, ctx *Context) error {
	if len(i.FormData) > 0 {
		i.AppMsg(fmt.Sprintf("AddFormData %v", i.FormData))
		for k, v := range i.FormData {
			v, err := i.expandContextVariable(v, ctx)
			if err != nil {
				i.AppErr("setFormdata - error setting contextVariable")
				return err
			}
			i.FormData[k] = v
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
		v, err := i.expandContextVariable(v, ctx)
		if err != nil {
			return errors.New(i.AppErr(fmt.Sprintf("setHeaders :%s", err.Error())))
		}
		req.SetHeader(k, v)
	}
	return nil
}

func (i *Input) expandContextVariable(v string, ctx *Context) (string, error) {
	if !strings.Contains(v, "$") {
		return v, nil
	}
	contextValue := strings.TrimLeft(v, "$")
	result := ctx.Get(contextValue)
	if result == nil {
		return v, errors.New(i.AppErr(fmt.Sprintf("Context value [%s] missing in context", contextValue)))
	}
	res, ok := result.(string)
	if !ok {
		return v, errors.New(i.AppErr(fmt.Sprintf("Context value [%s] - cannot convert result %v to string", contextValue, result)))
	}
	return res, nil
}

// AppMsg - application level trace
func (i *Input) AppMsg(msg string) string {
	tracer.AppMsg("Input", fmt.Sprintf("%s", msg), i.String())
	return msg
}

// AppErr - application level trace error msg
func (i *Input) AppErr(msg string) string {
	tracer.AppErr("Input", fmt.Sprintf("%s", msg), i.String())
	return msg
}

// String - object represetation
func (i *Input) String() string {
	bites, err := json.Marshal(i)
	if err != nil {
		return i.AppErr(fmt.Sprintf("error converting Input %s %s %s", i.Method, i.Endpoint, err.Error()))
	}
	return string(bites)
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
		i.AppErr(fmt.Sprintf("error signing jwt: %s", err.Error()))
		return "", err
	}
	tokenString = tokenString + "."
	return tokenString, nil
}

// take a JWT, generate a PSU consenturl
func (i *Input) generateConsentURI(jwt string) string {
	consent := i.Claims["aud"] + "/auth?" + "client_id=" + i.Claims["iss"] + "&response_type=" + i.Claims["response_type"] + "&scope=" + i.Claims["scope"] + "&request=" + jwt
	return consent
}
