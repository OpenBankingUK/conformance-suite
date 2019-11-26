package server

import (
	"strings"

	"math"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/version"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// ListenHost defines the name/address by which the service can be accessed.
const ListenHost = "0.0.0.0"

// Server - wraps *echo.Echo.
type Server struct {
	*echo.Echo // Wrap (using composition) *echo.Echo, allows us to pretend Server is echo.Echo.
	logger     *logrus.Entry
	version    version.Checker
}

// NewServer returns new echo.Echo server.
func NewServer(journey Journey, logger *logrus.Entry, version version.Checker) *Server {
	server := &Server{
		Echo:    echo.New(),
		logger:  logger,
		version: version,
	}
	server.Validator = newEchoValidatorAdapter()
	server.HideBanner = true

	// Use custom logger config so that we can control where log lines like below get sent to - either /dev/null or stdout.
	// {"time":"2018-12-18T13:00:40.291032Z","id":"","remote_ip":"192.0.2.1","host":"example.com","method":"POST","uri":"/api/config/global?pretty","status":400, "latency":627320,"latency_human":"627.32µs","bytes_in":0,"bytes_out":137}
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	server.Use(middleware.Recover())
	// TODO(mbana): figure out if this will break the downloading the report.zip file from the `/api/export` route.
	// https://github.com/labstack/echo/issues/873. If it doesn't break it, enable Gzip compression.
	// Another approach is to skip `/api/export` route using a skipper, see https://github.com/labstack/echo/issues/964.
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		// level between 1-9
		// where 1 indicates the fastest compression (less compression), and
		// 9 indicates the slowest compression method (best compression)
		Level:   5,
		Skipper: skipperGzip,
	}))

	// serve Vue.js site
	server.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: skipperSwagger,
		Root:    "web/dist",
		Index:   "index.html",
		HTML5:   true,
		Browse:  false,
	}))

	registerRoutes(journey, server, logger, version)

	return server
}

func registerRoutes(journey Journey, server *Server, logger *logrus.Entry, version version.Checker) {
	// swagger ui endpoints
	for path, handler := range swaggerHandlers(logger) {
		server.GET(path, handler)
	}

	// anything prefixed with api
	api := server.Group("/api")

	api.GET("/ping", func(c echo.Context) error { return nil })

	importHandlers := newImportHandlers(journey, logger)
	api.POST("/import/review", importHandlers.postImportReview)
	api.POST("/import/rerun", importHandlers.postImportRerun)

	configHandlers := newConfigHandlers(journey, logger)
	// endpoint to post global configuration
	api.POST("/config/global", configHandlers.configGlobalPostHandler)
	api.GET("/config/conditional-property", configHandlers.configConditionalPropertyHandler)

	// endpoints for discovery model
	discoveryHandlers := newDiscoveryHandlers(journey, logger)
	api.POST("/discovery-model", discoveryHandlers.setDiscoveryModelHandler)

	// endpoints for test cases
	testCaseHandlers := newTestCaseHandlers(journey, NewWebSocketUpgrader(), logger)
	api.GET("/test-cases", testCaseHandlers.testCasesHandler)

	// endpoints for test runner
	runHandlers := newRunHandlers(journey, NewWebSocketUpgrader(), logger)
	api.POST("/run", runHandlers.runStartPostHandler)
	api.GET("/run/ws", runHandlers.listenResultWebSocket)
	api.DELETE("/run", runHandlers.stopRunHandler)

	// endpoints for validating and storing the token retrieved in `/conformancesuite/callback`
	// `pkg/server/assets/main.js` calls into this endpoint.
	redirectHandlers := newRedirectHandlers(journey, logger)
	api.POST("/redirect/fragment/ok", redirectHandlers.postFragmentOKHandler)
	api.POST("/redirect/query/ok", redirectHandlers.postQueryOKHandler)
	api.POST("/redirect/error", redirectHandlers.postErrorHandler)

	exportHandlers := newExportHandlers(journey, logger)
	api.POST("/export", exportHandlers.postExport)

	// endpoints for utility function such as version/update checking.
	utilityEndpoints := newUtilityEndpoints(version)
	api.GET("/version", utilityEndpoints.versionCheck)
}

// skipperSwagger - ensures that all requests not prefixed with any string in `pathsToSkip` is skipped.
// E.g., ensure that `/api/validation-runs` or `/swagger/docs` is not handled by the static middleware.
func skipperSwagger(c echo.Context) bool {
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

// skipperGzip - ensures that gzip compression is not turned on for the `/api/export` and `/api/import` paths.
// I.e., don't run the Gzip middleware for certain paths.
func skipperGzip(c echo.Context) bool {
	pathsToSkip := []string{
		"/api/export",
		"/api/import",
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
	maxMessageSize := math.MaxInt32
	return &websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
	}
}
