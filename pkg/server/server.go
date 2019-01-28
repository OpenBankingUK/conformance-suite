package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// Server wraps *echo.Echo and stores the proxy once configured.
type Server struct {
	*echo.Echo // Wrap (using composition) *echo.Echo, allows us to pretend Server is echo.Echo.
	proxy   *http.Server
	logger  *logrus.Entry
	version version.Checker
}

// NewServer returns new echo.Echo server.
func NewServer(
	logger *logrus.Entry,
	checker model.ConditionalityChecker,
	version version.Checker,
) *Server {

	server := &Server{
		Echo:    echo.New(),
		proxy:   nil,
		logger:  logger,
		version: version,
	}
	server.HideBanner = true

	// Use custom logger config so that we can control where log lines like below get sent to - either /dev/null or stdout.
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

	registerRoutes(server, logger, checker, version)
	
	return server
}

func registerRoutes(server *Server, logger *logrus.Entry, checker model.ConditionalityChecker, version version.Checker) {
	// swagger ui endpoints
	for path, handler := range swaggerHandlers(logger) {
		server.GET(path, handler)
	}

	validatorEngine := discovery.NewFuncValidator(checker)
	testGenerator := generation.NewGenerator()
	journey := NewJourney(testGenerator, validatorEngine)

	server.Validator = newCustomValidator()

	// anything prefixed with api
	api := server.Group("/api")

	configHandlers := &configHandlers{server, journey}
	// endpoints to post a config and setup the proxy server
	api.POST("/config", configHandlers.configPostHandler)
	api.DELETE("/config", configHandlers.configDeleteHandler)
	// endpoint to post global configuration
	api.POST("/config/global", configHandlers.configGlobalPostHandler)

	// endpoints for discovery model
	discoveryHandlers := newDiscoveryHandlers(journey)
	api.POST("/discovery-model", discoveryHandlers.setDiscoveryModelHandler)

	// endpoints for test cases
	testCaseHandlers := newTestCaseHandlers(journey)
	api.GET("/test-cases", testCaseHandlers.testCasesHandler)

	// endpoints for test runner
	runHandlers := newRunHandlers(journey, NewWebSocketUpgrader())
	api.POST("/run", runHandlers.runStartPostHandler)
	api.GET("/run", runHandlers.listenResultWebSocket)
	api.DELETE("/run", runHandlers.stopRunHandler)

	// endpoints for utility function such as version/update checking.
	utilityEndpoints := newUtilityEndpoints(version)
	api.GET("/version", utilityEndpoints.versionCheck)
}

// Shutdown the server and the proxy if it is alive
func (s *Server) Shutdown(ctx context.Context) error {
	if s.proxy != nil {
		if err := s.proxy.Shutdown(ctx); err != nil {
			s.logger.Errorln("Server:Shutdown -> s.proxy.Shutdown err=", err)
			return err
		}
	}

	if s.Echo == nil {
		s.logger.Errorf("Server:Shutdown -> s.Echo=%p\n", s.Echo)
		return errors.New(`e.Echo == nil in Server:Shutdown`)
	}

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.logger.Errorln("Server:Shutdown -> s.Echo.Shutdown err=", err)
		return err
	}

	return nil
}

// Skipper ensures that all requests not prefixed with `/api` or `/swagger` get sent
// to the `middleware.Static` or `middleware.StaticWithConfig`.
// E.g., ensure that `/api/validation-runs` or `/swagger/docs` does not get
// handled by the the static middleware.
func (s *Server) skipper(c echo.Context) bool {
	skip := strings.HasPrefix(c.Path(), "/api") || strings.HasPrefix(c.Path(), "/swagger")
	return skip
}

// NewWebSocketUpgrader creates a new websocket.Ugprader.
func NewWebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}
