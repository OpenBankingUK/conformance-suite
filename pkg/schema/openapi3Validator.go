package schema

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/pkg/errors"
)

// OpenAPI3Validator - type
type OpenAPI3Validator struct {
	router routers.Router
	doc    *openapi3.T
}

// RequestWrapper -
type RequestWrapper struct {
	Method      string
	URL         string
	ContentType string
	Body        string
}

// ResponseWrapper -
type ResponseWrapper struct {
	Status      int
	ContentType string
	Body        string
}

// internal validation parameters
type validateParams struct {
	httpReq    *http.Request
	route      *routers.Route
	pathParams map[string]string
	statusCode int
	header     http.Header
	body       []byte
}

var headerCT = http.CanonicalHeaderKey("Content-Type")

// NewOpenAPI3Validator - Create a router for OPenAPI3 based specifications
// preferring yaml for the spec file
func NewOpenAPI3Validator(specName, version string) (Validator, error) {
	return buildValidator(specName, version)
}

// NewRawOpenAPI3Validator -
func NewRawOpenAPI3Validator(specName, version string) (OpenAPI3Validator, error) {
	return buildValidator(specName, version)
}

func buildValidator(specName, version string) (OpenAPI3Validator, error) {
	router, doc, err := getRouterForSpec(specName, version)
	return OpenAPI3Validator{router: router, doc: doc}, err
}

// findOperation - Find operation in schema
func (v OpenAPI3Validator) findOperation(checkmethod, checkpath string) *openapi3.Operation {
	spec := v.doc
	for path, props := range spec.Paths {
		if checkMatchingMethod(checkpath, path) {
			for method, op := range getOas3Operations(props) {
				if method == checkmethod {
					return op
				}
			}
		}
	}
	return nil
}

func checkMatchingMethod(checkmethod, opMethod string) bool {
	// trimming "/" tails and heads is required, because one of strings might have one extra slash more
	checkmethod = strings.Trim(checkmethod, "/")
	opMethod = strings.Trim(opMethod, "/")

	checkmethodSlice := strings.Split(checkmethod, "/")
	opMethodSlice := strings.Split(opMethod, "/")
	if len(checkmethodSlice) == len(opMethodSlice) {
		for i := 0; i < len(opMethodSlice); i++ {
			if !strings.HasPrefix(checkmethodSlice[i], "$") && checkmethodSlice[i] != opMethodSlice[i] {
				return false
			}
		}
		return true
	}
	return false
}

// IsRequestProperty - Find param in schema and determines if it's part of request body
func (v OpenAPI3Validator) IsRequestProperty(checkmethod, checkpath, propertyPath string) (bool, string, error) {
	spec := v.doc
	for path, props := range spec.Paths {
		for method, op := range getOas3Operations(props) {
			if path == checkpath && method == checkmethod && op.RequestBody != nil {
				for _, param := range op.RequestBody.Value.Content {
					schema := param.Schema.Value
					found, objType := findPropertyInOas3Schema(schema, propertyPath, "")
					if found {
						return true, objType, nil
					}
				}
			}
		}
	}

	return false, "", nil
}

func getRouterForSpec(specName, version string) (routers.Router, *openapi3.T, error) {

	filenamePattern := getSpecFilePathPattern(specName)
	if filenamePattern == "" {
		return nil, nil, errors.New("cannot get router for spec: " + specName)
	}

	filename := fmt.Sprintf(filenamePattern, version)

	doc, err := loadSpecFromFile(filename)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot Load OpenApi Spec from file %s, %s", filename, err)
	}

	err = doc.Validate(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot Load OpenApi Spec from file %s, %s", filename, err)
	}

	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot Load OpenApi Router for %s file %s", specName, filename)
	}

	return router, doc, nil
}

func loadSpecFromFile(filename string) (*openapi3.T, error) {
	prodDir := "pkg/schema/" + filename
	testDir := "../../pkg/schema/" + filename
	specPaths := []string{filename, prodDir, testDir}

	var doc *openapi3.T
	var err error
	loader := openapi3.NewLoader()

	for _, specPath := range specPaths {
		doc, err = loader.LoadFromFile(specPath)
		if err == nil {
			break
		}
	}
	return doc, err
}

func getSpecFilePathPattern(specName string) string {
	var filename string

	switch specName {
	case "Account and Transaction API Specification":
		filename = "spec/%s/account-info-openapi.json"

	case "Payment Initiation API":
		filename = "spec/%s/payment-initiation-openapi.json"

	case "Confirmation of Funds API Specification":
		filename = "spec/%s/confirmation-funds-openapi.json"

	case "OBIE VRP Profile":
		filename = "spec/%s/variable-recurring-payments-openapi.json"

	default:
		filename = ""
	}

	return filename
}

// Validate - validates the response
func (v OpenAPI3Validator) Validate(r HTTPResponse) ([]Failure, error) {
	failures := []Failure{}

	serverPath := v.doc.Servers[0].URL
	var path string
	serverIndex := strings.Index(r.Path, serverPath)
	if serverIndex != -1 {
		path = r.Path[serverIndex:]
	} else {
		path = serverPath + r.Path
	}

	httpReq, err := createHTTPReq(r.Method, path)
	if err != nil {
		return nil, err
	}

	route, pathParams, err := v.router.FindRoute(httpReq)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("OpenApi3Validator: error converting body %s", err)
	}

	// check body
	params := validateParams{
		httpReq:    httpReq,
		route:      route,
		pathParams: pathParams,
		statusCode: r.StatusCode,
		header:     r.Header,
		body:       body,
	}

	// accumulate failures
	err = v.validateResponse(params)
	if err != nil {
		return nil, fmt.Errorf("Validate error response:  %s", err.Error())
	}

	return failures, nil
}

// HTTPRequest represents a request object from a HTTP Call
type HTTPRequest struct {
	Method  string
	Path    string
	Header  http.Header
	StrBody string
}

// ValidateRequest - validates HTTPRequest based on OpenAPI3Validator
func (v OpenAPI3Validator) ValidateRequest(req HTTPRequest) error {
	errBase := fmt.Sprintf("%s %s -", req.Method, req.Path)

	reqParams, err := v.createRequestValidationParms(req)
	if err != nil {
		return fmt.Errorf("%s createRequestValidationParms: %s", errBase, err.Error())
	}

	requestValidationInput := createRequestValidationInput(reqParams)
	err = v.testRequest(requestValidationInput, req)
	if err != nil {
		return fmt.Errorf("%s testRequest: %s", errBase, err.Error())
	}
	return err
}

func createRequestValidationInput(reqParams validateParams) openapi3filter.RequestValidationInput {
	return openapi3filter.RequestValidationInput{
		Request:    reqParams.httpReq,
		PathParams: reqParams.pathParams,
		Route:      reqParams.route,
	}
}

func (v OpenAPI3Validator) createRequestValidationParms(req HTTPRequest) (validateParams, error) {
	params := validateParams{}
	serverPath := v.doc.Servers[0].URL
	var path string
	serverIndex := strings.Index(req.Path, serverPath)
	if serverIndex != -1 {
		path = req.Path[serverIndex:]
	} else {
		path = serverPath + req.Path
	}

	bodyReader := strings.NewReader(req.StrBody)

	httpReq, err := http.NewRequest(req.Method, path, bodyReader)
	httpReq.Header = req.Header
	if err != nil {
		return params, err
	}

	route, pathParams, err := v.router.FindRoute(httpReq)
	if err != nil {
		return params, err
	}

	params = validateParams{
		httpReq:    httpReq,
		route:      route,
		pathParams: pathParams,
		header:     req.Header,
	}
	return params, nil
}

func (v OpenAPI3Validator) testRequest(requestInput openapi3filter.RequestValidationInput, req HTTPRequest) error {
	schemaBody, err := v.findSchemaRequestBody(req)
	if err != nil {
		return err
	}

	err = openapi3filter.ValidateRequestBody(context.Background(), &requestInput, schemaBody)
	if err != nil {
		return fmt.Errorf("ValidateRequestBody: %s", err.Error())
	}

	return nil
}

func (v OpenAPI3Validator) findSchemaRequestBody(req HTTPRequest) (*openapi3.RequestBody, error) {
	op := v.findOperation(req.Method, req.Path)
	if op == nil {
		return nil, fmt.Errorf("findOperation: cannot find operation in schema")
	}

	// no body
	if op.RequestBody == nil {
		return openapi3.NewRequestBody(), nil
	}

	return op.RequestBody.Value, nil
}

func (v OpenAPI3Validator) GetEndpoints() map[string]map[string]bool {
	endpoints := make(map[string]map[string]bool)
	for path, endpoint := range v.doc.Paths {
		endpoints[path] = make(map[string]bool)
		if endpoint.Get != nil {
			endpoints[path]["get"] = true
		}
		if endpoint.Head != nil {
			endpoints[path]["head"] = true
		}
		if endpoint.Post != nil {
			endpoints[path]["post"] = true
		}
		if endpoint.Put != nil {
			endpoints[path]["put"] = true
		}
		if endpoint.Delete != nil {
			endpoints[path]["delete"] = true
		}
		if endpoint.Connect != nil {
			endpoints[path]["connect"] = true
		}
		if endpoint.Options != nil {
			endpoints[path]["options"] = true
		}
		if endpoint.Trace != nil {
			endpoints[path]["trace"] = true
		}
		if endpoint.Patch != nil {
			endpoints[path]["patch"] = true
		}
	}
	return endpoints
}

func (v OpenAPI3Validator) validateResponse(params validateParams) error {
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    params.httpReq,
		PathParams: params.pathParams,
		Route:      params.route,
	}

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 params.statusCode,
		Header:                 params.header,
		Options: &openapi3filter.Options{
			ExcludeRequestBody:    true,
			IncludeResponseStatus: true,
			MultiError:            false,
		},
	}

	if len(params.body) > 0 {
		responseValidationInput.SetBodyBytes(params.body)
	}

	return openapi3filter.ValidateResponse(context.Background(), responseValidationInput)
}

func (v OpenAPI3Validator) findTestRoute(req *http.Request) (*routers.Route, map[string]string, error) {
	route, pathParams, err := v.router.FindRoute(req)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %s - findTestRoute:  %s", req.Method, req.URL, err)
	}
	return route, pathParams, err
}

func createHTTPReq(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, strings.NewReader(""))
	req.Header = http.Header{"Content-type": []string{"application/json; charset=utf-8"}}
	return req, err
}

// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOas3Operations(props *openapi3.PathItem) map[string]*openapi3.Operation {
	ops := map[string]*openapi3.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}

	// Keep those != nil
	for key, op := range ops {
		if op == nil {
			delete(ops, key)
		}
	}
	return ops
}

//normalizePropertyType - Workaround to provide similar context to the one used in Swagger schema
func normalizePropertyType(propertyType string) string {
	return fmt.Sprintf("[%s]", propertyType)
}

func findPropertyInOas3Schema(sc *openapi3.Schema, propertyPath, previousPath string) (bool, string) {
	for k, j := range sc.Properties {
		var element string
		if len(previousPath) == 0 {
			element = k
		} else {
			element = previousPath + "." + k
		}

		if element == propertyPath {
			return true, fmt.Sprintf("%s", normalizePropertyType(j.Value.Type))
		}

		ret, propType := findPropertyInOas3Schema(j.Value, propertyPath, element)
		if ret {
			return true, propType
		}
	}

	return false, ""
}
