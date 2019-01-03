package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"bitbucket.org/openbankingteam/conformance-suite/appconfig"
	"bitbucket.org/openbankingteam/conformance-suite/proxy"

	"github.com/go-openapi/loads"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// GlobalConfiguration holds:
// * private signing key
// * public signing key
// * private transport key
// * public transport key
type GlobalConfiguration struct {
	SigningPrivate   string `json:"signing_private"`
	SigningPublic    string `json:"signing_public"`
	TransportPrivate string `json:"transport_private"`
	TransportPublic  string `json:"transport_public"`
}

// CustomValidator used to validate incoming payloads (for now).
// https://echo.labstack.com/guide/request#validate-data
type CustomValidator struct {
	validator *validator.Validate
}

// Validate incoming payloads (for now) that contain the struct tag `validate:"required"`.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Server wraps *echo.Echo and stores the proxy once configured.
type Server struct {
	*echo.Echo // Wrap (using composition) *echo.Echo, allows us to pretend Server is echo.Echo.
	proxy      *http.Server
	logger     *logrus.Entry
}

// NewServer returns new echo.Echo server.
func NewServer(
	logger *logrus.Entry,
	checker model.ConditionalityChecker,
) *Server {

	server := &Server{
		Echo:   echo.New(),
		proxy:  nil,
		logger: logger,
	}

	// Use custom logger config so that log lines like below don't appear in the output:
	// {"time":"2018-12-18T13:00:40.291032Z","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"POST","uri":"/api/config/global?pretty","status":400, "latency":627320,"latency_human":"627.32µs","bytes_in":0,"bytes_out":137}
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	server.Use(middleware.Recover())
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		// level between 1-9
		// where 1 indicates the fastest compression (less compression), and
		// 9 indicates the slowest compression method (best compression)
		Level: 5,
	}))
	// serve Vue.js site
	server.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: server.skipper,
		Root:    "web/dist",
		Index:   "index.html",
		HTML5:   true,
		Browse:  false,
	}))

	server.HideBanner = true

	validatorEngine := discovery.NewFuncValidator(checker)
	testGenerator := generation.NewGenerator()
	webJourney := web.NewWebJourney(testGenerator, validatorEngine)

	// https://echo.labstack.com/guide/request#validate-data
	validator := validator.New()
	server.Validator = &CustomValidator{validator}

	// anything prefixed with api
	api := server.Group("/api")

	wsHandler := &WebSocketHandler{
		upgrader: NewWebSocketUpgrader(),
	}
	// serve WebSocket
	api.GET("/ws", wsHandler.Handle)

	configHandlers := &configHandlers{server}
	// endpoints to post a config and setup the proxy server
	api.POST("/config", configHandlers.configPostHandler)
	api.DELETE("/config", configHandlers.configDeleteHandler)
	// endpoint to post global configuration
	api.POST("/config/global", configHandlers.configGlobalPostHandler)

	// endpoints for discovery model
	discoveryHandlers := newDiscoveryHandlers(webJourney)
	api.POST("/discovery-model", discoveryHandlers.setDiscoveryModelHandler)

	// endpoints for test cases
	testCaseHandlers := newTestCaseHandlers(webJourney)
	api.GET("/test-cases", testCaseHandlers.testCasesHandler)

	server.logRoutes()

	return server
}

func (s *Server) logRoutes() {
	for _, route := range s.Routes() {
		s.logger.Debugf("route -> path=%+v, method=%+v", route.Path, route.Method)
	}
}

// Shutdown the server and the proxy if it is alive
func (s *Server) Shutdown(ctx context.Context) error {
	if s.proxy != nil {
		if err := s.proxy.Shutdown(nil); err != nil {
			s.logger.Errorln("Server:Shutdown -> s.proxy.Shutdown err=", err)
			return err
		}
	}

	if s.Echo == nil {
		s.logger.Errorf("Server:Shutdown -> s.Echo=%p\n", s.Echo)
	}

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.logger.Errorln("Server:Shutdown -> s.Echo.Shutdown err=", err)
		return err
	}

	return nil
}

// Skipper ensures that all requests not prefixed with `/api` get sent
// to the `middleware.Static` or `middleware.StaticWithConfig`.
// E.g., ensure that `/api/validation-runs` does not get handled by the
// the static middleware.
//
// Anything not prefix by `/api` will get get handled by
// `middleware.Static` or `middleware.StaticWithConfig`
func (s *Server) skipper(c echo.Context) bool {
	return strings.HasPrefix(c.Path(), "/api")
}

// Run the proxy at the address specified by "bind"
// Requests get sent to the target server identified by proxy.Target()
// configure some channels to handle shutdown/interrupts
//
// Return channel so that caller can block waiting
// to see if we managed to start the server or not.
func serveProxy(proxy *proxy.Proxy, bind string) (*http.Server, chan error) {
	server := &http.Server{
		Addr:    bind,
		Handler: proxy.Router(),
	}

	// Run server s.ListenAndServe on a goroutine
	serveErr := make(chan error)
	go func() {
		err := server.ListenAndServe()
		serveErr <- err
	}()

	return server, serveErr
}

// createProxy - kick off proxy by:
// loading the spec,
// creating a new proxy configured with
//    - bind address
//    - swagger specification location
//    - target host (aspsp resource server)
//    - verbosity
// configure an default logreport
func createProxy(appConfig *appconfig.AppConfig) (*http.Server, error) {
	logrus.Info("Server:createProxy -> Proxy")

	appConfig.PrintAppConfig()
	doc, err := loads.Spec(appConfig.Spec)
	if err != nil {
		logrus.Errorln("Server:createProxy -> loads.Spec err=", err)
		return nil, err
	}

	proxy, err := proxy.New(
		doc.Spec(),
		&proxy.LogReporter{},
		proxy.WithTarget(appConfig.TargetHost),
		proxy.WithVerbose(appConfig.Verbose),
		proxy.WithAppConfig(appConfig),
	)
	if err != nil {
		logrus.Errorln("Server:createProxy -> proxy.New err=", err)
		return nil, err
	}

	// start serving the proxy - and don't return unless there is a problem/exit
	// also sleep for a bit until it starts...
	server, serveErr := serveProxy(proxy, appConfig.Bind)
	time.Sleep(200 * time.Millisecond)

	// block until serveErr has an error value or the specified timeout has elapsed.
	// we might need to bump up this timeout to something a bit larger.
	timeout := time.After(1 * time.Second)
	select {
	case err := <-serveErr: // Error from listen&serve - exit
		return nil, err
	case <-timeout:
	}

	logrus.WithFields(logrus.Fields{
		"bind":   appConfig.Bind,
		"target": proxy.Target(),
	}).Info("Server:createProxy -> Proxy is listening")

	// Report PendingOperations - part of shutdown tidyup
	logrus.Debugln("Pending Operations:")
	for i, op := range proxy.PendingOperations() {
		logrus.Debugf("%03d) id=%s", i+1, op.ID)
	}

	return server, nil
}
