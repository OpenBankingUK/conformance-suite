package schema

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/pkg/errors"
)

// regexp2Compiler is a RegexCompilerFunc that uses dlclark/regexp2 so that
// OpenAPI spec patterns using PCRE features (e.g. lookahead/lookbehind) are
// handled correctly instead of crashing Go's RE2-based stdlib regexp.
func regexp2Compiler(expr string) (openapi3.RegexMatcher, error) {
	re, err := regexp2.Compile(expr, 0)
	if err != nil {
		return nil, err
	}
	return &regexp2Matcher{re: re}, nil
}

type regexp2Matcher struct{ re *regexp2.Regexp }

func (m *regexp2Matcher) MatchString(s string) bool {
	ok, _ := m.re.MatchString(s)
	return ok
}

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

// IsRequestProperty - Find param in schema and determines if it's part of request body
func (v OpenAPI3Validator) IsRequestProperty(checkmethod, checkpath, propertyPath string) (bool, string, error) {
	spec := v.doc
	for path, props := range spec.Paths.Map() {
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

	err = doc.Validate(context.Background(), openapi3.SetRegexCompiler(regexp2Compiler))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot Load OpenApi Spec from file %s, %s", filename, err)
	}

	router, err := legacyrouter.NewRouter(doc, openapi3.SetRegexCompiler(regexp2Compiler))
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

	case "Variable Recurring Payments API Specification":
		filename = "spec/%s/variable-recurring-payments-openapi.json"

	case "Commercial Variable Recurring Payments API Specification":
		filename = "spec/%s/cvrp-openapi.json"

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
			RegexCompiler:         regexp2Compiler,
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

// typesToString converts an *openapi3.Types ([]string in v0.134.0+) to a
// single string, joining multiple types with "/" for display purposes.
func typesToString(t *openapi3.Types) string {
	if t == nil {
		return ""
	}
	return strings.Join(t.Slice(), "/")
}

// normalizePropertyType - Workaround to provide similar context to the one used in Swagger schema
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
			return true, normalizePropertyType(typesToString(j.Value.Type))
		}

		ret, propType := findPropertyInOas3Schema(j.Value, propertyPath, element)
		if ret {
			return true, propType
		}
	}

	return findItemInOas3Schema(sc, propertyPath, previousPath)
}

func findItemInOas3Schema(sc *openapi3.Schema, propertyPath, previousPath string) (bool, string) {
	if len(sc.Properties) == 0 {
		notFoundPath := strings.Replace(propertyPath, previousPath+".", "", 1)
		SplitedNotFoundPath := strings.Split(notFoundPath, ".")
		idx := SplitedNotFoundPath[0]
		if _, err := strconv.Atoi(idx); err == nil {
				if len(SplitedNotFoundPath) == 1 {
				return true, normalizePropertyType(typesToString(sc.Items.Value.Type))
			}
			element := previousPath + "." + idx
			ret, propType := findPropertyInOas3Schema(sc.Items.Value, propertyPath, element)
			if ret {
				return true, propType
			}
		}
	}

	return false, ""
}
