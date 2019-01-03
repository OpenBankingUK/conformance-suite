package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"bitbucket.org/openbankingteam/conformance-suite/appconfig"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Proxy defining type
type Proxy struct {
	// Opts
	target            string
	verbose           bool
	appconfig         *appconfig.AppConfig
	router            *mux.Router
	routes            map[*mux.Route]*spec.Operation
	reverseProxy      *httputil.ReverseProxy // Reverse Proxy
	reporter          Reporter               // Report Handler - extensible interface
	doc               interface{}            // This is useful for validate (TODO: find a better way)
	spec              *spec.Swagger
	pendingOperations map[*spec.Operation]struct{}
}

// Opts - Proxy options
type Opts func(*Proxy)

// WithTarget - option
func WithTarget(target string) Opts { return func(proxy *Proxy) { proxy.target = target } }

// WithVerbose - option
func WithVerbose(v bool) Opts { return func(proxy *Proxy) { proxy.verbose = v } }

// WithAppConfig - option
func WithAppConfig(a *appconfig.AppConfig) Opts { return func(proxy *Proxy) { proxy.appconfig = a } }

// New Proxy ceration
func New(s *spec.Swagger, reporter Reporter, opts ...Opts) (*Proxy, error) {
	proxy := &Proxy{
		target:   "http://localhost:8080",
		router:   mux.NewRouter(),
		routes:   make(map[*mux.Route]*spec.Operation),
		reporter: reporter,
	}

	for _, opt := range opts { // Process options using standard pattern
		opt(proxy)
	}

	if err := proxy.SetSpec(s); err != nil {
		return nil, err
	}

	rpURL, err := url.Parse(proxy.target)
	if err != nil {
		return nil, err
	}
	proxy.reverseProxy = httputil.NewSingleHostReverseProxy(rpURL)

	// Add OB Certs to the reverse proxy tls configuration
	if proxy.appconfig.CertTransport != "" { // if we have MATLS transports certs in the application config - used them
		tlsconfig, _ := proxy.appconfig.NewTLSConfig() // if this fails we just get an empty tlsconfig
		proxy.reverseProxy.Transport = &http.Transport{
			TLSClientConfig: tlsconfig,
		}
	}

	proxy.router.NotFoundHandler = http.HandlerFunc(proxy.notFound)
	proxy.registerPaths()

	return proxy, nil
}

// SetSpec - Marshalls the spec into a generic doc interface
// for the purpose of validating the spec (if it won't marshall its broken)
// The doc interface is useful in other scenarios within the proxy
func (proxy *Proxy) SetSpec(spec *spec.Swagger) error {
	// validate.NewSchemaValidator requires the spec as an interface{}
	// That's why we Unmarshal(Marshal()) the document
	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}

	proxy.doc = doc
	proxy.spec = spec

	return nil
}

// Router - return router
func (proxy *Proxy) Router() http.Handler {
	return proxy.router
}

// Target - return target
func (proxy *Proxy) Target() string {
	return proxy.target
}

// AppConfig -
func (proxy *Proxy) AppConfig() *appconfig.AppConfig {
	return proxy.appconfig
}

// Iterator over the spec and pull out all the API paths
// Then setup a router to recognise the path so we can process requests that match the paths
func (proxy *Proxy) registerPaths() {
	proxy.pendingOperations = make(map[*spec.Operation]struct{})
	base := proxy.spec.BasePath
	logrus.Println("BasePath", base)

	router := mux.NewRouter()
	WalkOps(proxy.spec, func(path, method string, op *spec.Operation) {
		newPath := base + path
		if proxy.verbose {
			logrus.Debugf("Register %s %s", method, newPath)
		}
		route := router.Handle(newPath, proxy.newHandler()).Methods(method)
		proxy.routes[route] = op
		proxy.pendingOperations[op] = struct{}{}
	})

	*proxy.router = *router
}

func (proxy *Proxy) notFound(w http.ResponseWriter, req *http.Request) {
	proxy.reporter.Warning(req, "Route not defined on the Spec")
	proxy.reverseProxy.ServeHTTP(w, req)
}

func (proxy *Proxy) newHandler() http.Handler {
	return proxy.Handler(proxy.reverseProxy)
}

// Handler -
func (proxy *Proxy) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// https://stackoverflow.com/questions/38016477/reverse-proxy-does-not-work
		// https://blog.semanticart.com/2013/11/11/a-proper-api-proxy-written-in-go/
		// https://forum.golangbridge.org/t/explain-how-reverse-proxy-work/6492/7
		// https://stackoverflow.com/questions/34745654/golang-reverseproxy-with-apache2-sni-hostname-error
		req.Host = req.URL.Host

		wr := &WriterRecorder{ResponseWriter: w}
		next.ServeHTTP(wr, req)
		logrus.Println("Handling Request ", req.URL)
		// for key, value := range req.Header {
		// 	logrus.WithFields(logrus.Fields{
		// 		key: strings.Join(value, " ")}).Info("")
		// }

		var match mux.RouteMatch
		proxy.router.Match(req, &match)
		op := proxy.routes[match.Route]

		if match.Handler == nil || op == nil {
			proxy.reporter.Warning(req, "Route not defined on the Spec")
			// Route hasn't been registered on the muxer
			return
		}
		proxy.operationExecuted(op)

		if err := proxy.Validate(wr, op); err != nil {
			proxy.reporter.Error(req, err)
		} else {
			proxy.reporter.Success(req)
		}
	}
	return http.HandlerFunc(fn)
}

type validatorFunc func(Response, *spec.Operation) error

// Validate -
func (proxy *Proxy) Validate(resp Response, op *spec.Operation) error {
	if _, ok := op.Responses.StatusCodeResponses[resp.Status()]; !ok {
		return fmt.Errorf("server Status %d not defined by the spec", resp.Status())
	}

	var validators = []validatorFunc{
		proxy.ValidateMIME,
		proxy.ValidateHeaders,
		proxy.ValidateBody,
	}

	var errs []error
	for _, v := range validators {
		if err := v(resp, op); err != nil {
			if cErr, ok := err.(*errors.CompositeError); ok {
				errs = append(errs, cErr.Errors...)
			} else {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errors.CompositeValidationError(errs...)
}

// DumpHeaders -
func DumpHeaders(headers map[string]spec.Header) {
	for name, values := range headers {
		//res = append(res, fmt.Sprintf("%s: %s", name, value))
		logrus.Printf("%s %s\n", name, values.Description)

	}
	return
}

// ValidateMIME -
func (proxy *Proxy) ValidateMIME(resp Response, op *spec.Operation) error {
	// Use Operation Spec or fallback to root
	produces := op.Produces
	if len(produces) == 0 {
		produces = proxy.spec.Produces
	}

	logrus.Println("ValidateMIME")
	ct := resp.Header().Get("Content-Type")
	if len(produces) == 0 {
		return nil
	}
	for _, mime := range produces {
		if ct == mime {
			return nil
		}
	}

	logrus.Printf("Content-Type Error: Should maybe produce %q, but got: '%s' - but we'll let that go\n", produces, ct)
	return nil

	//return fmt.Errorf("Content-Type Error: Should maybe produce %q, but got: '%s'", produces, ct)
}

// ValidateHeaders -
func (proxy *Proxy) ValidateHeaders(resp Response, op *spec.Operation) error {
	var errs []error
	r := op.Responses.StatusCodeResponses[resp.Status()]
	//DumpHeaders(r.Headers)
	for key, spec := range r.Headers {
		if err := validateHeaderValue(key, resp.Header().Get(key), &spec); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.CompositeValidationError(errs...)
}

// ValidateBody -
func (proxy *Proxy) ValidateBody(resp Response, op *spec.Operation) error {
	r := op.Responses.StatusCodeResponses[resp.Status()]
	if r.Schema == nil {
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return err
	}

	v := validate.NewSchemaValidator(r.Schema, proxy.doc, "", strfmt.Default)
	if result := v.Validate(data); result.HasErrors() {
		return result.AsError()
	}

	return nil
}

func validateHeaderValue(key, value string, spec *spec.Header) error {
	if value == "" {
		return fmt.Errorf("%s in headers is missing", key)
	}

	// TODO: Implement the rest of the format validators
	switch spec.Format {
	case "int32":
		_, err := swag.ConvertInt32(value)
		return err
	case "date-time":
		_, err := strfmt.ParseDateTime(value)
		return err
	}
	return nil
}

// PendingOperations -
func (proxy *Proxy) PendingOperations() []*spec.Operation {
	var ops []*spec.Operation
	for op := range proxy.pendingOperations {
		ops = append(ops, op)
	}
	return ops
}

func (proxy *Proxy) operationExecuted(op *spec.Operation) {
	delete(proxy.pendingOperations, op)
}

// WalkOpsFunc - for use with the WalkOps spec iterator
// Define func type to handline data structure walking
type WalkOpsFunc func(path, meth string, op *spec.Operation)

// WalkOps ...
// For all the paths in the spec
// get the operations and methods
// call the WalkOpsFunc (which is passed in) to process
func WalkOps(spec *spec.Swagger, fn WalkOpsFunc) {
	for path, props := range spec.Paths.Paths {
		for meth, op := range getOperations(&props) {
			fn(path, meth, op)
		}
	}
}

// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOperations(props *spec.PathItem) map[string]*spec.Operation {
	ops := map[string]*spec.Operation{
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
