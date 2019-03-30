package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type testCaseHandlers struct {
	journey  Journey
	upgrader *websocket.Upgrader
	logger   *logrus.Entry
}

func newTestCaseHandlers(journey Journey, upgrader *websocket.Upgrader, logger *logrus.Entry) testCaseHandlers {
	return testCaseHandlers{
		journey:  journey,
		upgrader: upgrader,
		logger:   logger,
	}
}

func (d testCaseHandlers) testCasesHandler(c echo.Context) error {
	d.journey.NewDaemonController() // fix for not sending events to correct websocket after a websocket reconnect
	testCases, err := d.journey.TestCases()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, testCases)
}
