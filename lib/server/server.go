package server

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// Skipper ensures that all requests not prefixed with `/api` get sent
// to the `middleware.Static` or `middleware.StaticWithConfig`.
// E.g., ensure that `/api/validation-runs` does not get handled by the
// the static middleware.
//
// Anything not prefix by `/api` will get get handled by
// `middleware.Static` or `middleware.StaticWithConfig`
func Skipper(c echo.Context) bool {
	return strings.HasPrefix(c.Path(), "/api")
}

// ValidationRunsResponse is the response to the `/validation-runs` endpoint.
type ValidationRunsResponse struct {
	ID string `json:"id"`
}

// ValidationRunsIDResponse is the response to the `/validation-runs/:id` endpoint.
type ValidationRunsIDResponse struct {
	Status string `json:"status"`
}

// NewServer returns new echo.Echo server.
func NewServer() *echo.Echo {
	e := echo.New()

	// Write output to `/dev/null` if running under test mode.
	if flag.Lookup("test.v") == nil {
		// Normal run
		e.Use(middleware.Logger())
	} else {
		// Running under test
		// discard output when running in test mode
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: ioutil.Discard,
		}))
	}
	e.Use(middleware.Recover())
	// serve Vue.js site
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		// level between 1-9
		// where 1 indicates the fastest compression (less compression), and
		// 9 indicates the slowest compression method (best compression)
		Level: 5,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: Skipper,
		Root:    "web/dist",
		Index:   "index.html",
		HTML5:   true,
		Browse:  false,
	}))

	// anything prefixed with api
	api := e.Group("/api")

	wsHandler := &WebSocketHandler{
		upgrader: NewWebSocketUpgrader(),
	}
	// serve WebSocket
	api.GET("/ws", wsHandler.Handle)
	// health check endpoint
	api.GET("/health", healthHandler)
	api.POST("/validation-runs", validationRunsHandler)
	api.GET("/validation-runs/:id", validationRunsIDHandler)

	routes, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err == nil {
		logrus.Debugf("routes=%s", routes)
	}

	return e
}

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func validationRunsHandler(c echo.Context) error {
	id, err := uuid.NewUUID()

	if err != nil {
		return c.String(http.StatusNotAcceptable, err.Error())
	}

	// TODO: Looks a bit bad I know. I will clear it up later.
	return c.JSON(http.StatusAccepted, ValidationRunsResponse{
		ID: id.String(),
	})
}

func validationRunsIDHandler(c echo.Context) error {
	status := c.Param("id")

	// TODO: Looks a bit bad I know. I will clear it up later.
	return c.JSON(http.StatusOK, ValidationRunsIDResponse{
		Status: status,
	})
}
