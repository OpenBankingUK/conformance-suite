package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/labstack/echo"
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
	testCases, err := d.journey.TestCases()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, testCases)
}

// listenCodeWebSocket upgrades to websocket and notifies when backend has collected all
// codes from PSU consent flow
func (d testCaseHandlers) listenCodeWebSocket(c echo.Context) error {
	ws, err := d.upgrader.Upgrade(c.Response(), c.Request(), nil)
	logger := d.logger.WithField("handler", "listenCodeWebSocket").WithField("websocket", fmt.Sprintf("%p", ws))

	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		logger.Debug("client disconnected")
		err := ws.Close()
		if err != nil {
			logger.WithError(err).Error("closing websocket")
		}
	}()

	logger.Debug("client connected")

	pingFrequency := time.Second * 2
	pingTicker := time.NewTicker(pingFrequency)

	for {
		if d.journey.AllTokenCollected() {
			logger.Info("sending all collected event")
			if err := ws.WriteJSON(newAllCollectedEvent()); err != nil {
				logger.WithError(err).Error("writing json to websocket")
			}
		}

		<-pingTicker.C
		logger.Debug("pinging websocket client")
		writeTimeout := time.Now().Add(time.Second)
		err := ws.SetWriteDeadline(writeTimeout)
		if err != nil {
			// we cannot return error here, if we do echo will try to write the error to conn
			// and we closed the ws with a defer func
			return nil
		}
		if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			// same as above
			return nil
		}
	}
}

type AllCollectedEvent struct {
	Stopped bool `json:"collected"`
}

func newAllCollectedEvent() AllCollectedEvent {
	return AllCollectedEvent{true}
}
