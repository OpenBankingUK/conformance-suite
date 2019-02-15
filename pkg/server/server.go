package server

import (
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// Server - wraps *echo.Echo.
type Server struct {
	*echo.Echo // Wrap (using composition) *echo.Echo, allows us to pretend Server is echo.Echo.
	logger     *logrus.Entry
	version    version.Checker
}

// NewServer returns new echo.Echo server.
func NewServer(logger *logrus.Entry, checker model.ConditionalityChecker, version version.Checker) *Server {
	server := &Server{
		Echo:    echo.New(),
		logger:  logger,
		version: version,
	}
	server.Validator = newEchoValidatorAdapter()

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
		Skipper: skipper,
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
	journey := NewJourney(logger, testGenerator, validatorEngine)

	// anything prefixed with api
	api := server.Group("/api")

	configHandlers := &configHandlers{logger: logger, journey: journey}
	// endpoint to post global configuration
	api.POST("/config/global", configHandlers.configGlobalPostHandler)

	// endpoints for discovery model
	discoveryHandlers := newDiscoveryHandlers(journey)
	api.POST("/discovery-model", discoveryHandlers.setDiscoveryModelHandler)

	// endpoints for test cases
	testCaseHandlers := newTestCaseHandlers(journey, NewWebSocketUpgrader(), logger)
	api.GET("/test-cases", testCaseHandlers.testCasesHandler)
	api.GET("/test-cases/ws", testCaseHandlers.listenCodeWebSocket)

	// endpoints for test runner
	runHandlers := newRunHandlers(journey, NewWebSocketUpgrader(), logger)
	api.POST("/run", runHandlers.runStartPostHandler)
	api.GET("/run/ws", runHandlers.listenResultWebSocket)
	api.DELETE("/run", runHandlers.stopRunHandler)

	// endpoints for utility function such as version/update checking.
	utilityEndpoints := newUtilityEndpoints(version)
	api.GET("/version", utilityEndpoints.versionCheck)

	// endpoints for validating and storing the token retrieved in `/conformancesuite/callback`
	// `pkg/server/assets/main.js` calls into this endpoint.
	redirectHandlers := &redirectHandlers{logger.WithField("module", "redirectHandlers")}
	api.POST("/redirect/fragment/ok", redirectHandlers.postFragmentOKHandler)
	api.POST("/redirect/query/ok", redirectHandlers.postQueryOKHandler)
	api.POST("/redirect/error", redirectHandlers.postErrorHandler)
}

// skipper - ensures that all requests not prefixed with any string in `pathsToSkip` is skipped.
// E.g., ensure that `/api/validation-runs` or `/swagger/docs` is not handled by the static middleware.
func skipper(c echo.Context) bool {
	pathsToSkip := []string{
		"/api",
		"/swagger",
	}

	path := c.Path()
	for _, pathToSkip := range pathsToSkip {
		if strings.HasPrefix(path, pathToSkip) {
			return true
		}
	}

	return false
}

// NewWebSocketUpgrader creates a new websocket.Ugprader.
func NewWebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}
