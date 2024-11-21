package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		d.logger.Println("user_id not found in context")
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.New("user_id not found in context")))
	}
	testCases, err := d.journey.TestCases(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, testCases)
}
