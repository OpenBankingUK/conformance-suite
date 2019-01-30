package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type runHandlers struct {
	journey  Journey
	upgrader *websocket.Upgrader
	logger   *logrus.Entry
}

func newRunHandlers(journey Journey, upgrader *websocket.Upgrader, logger *logrus.Entry) *runHandlers {
	return &runHandlers{
		journey:  journey,
		upgrader: upgrader,
		logger:   logger,
	}
}

// runStartPostHandler creates a new test run
func (h *runHandlers) runStartPostHandler(c echo.Context) error {
	err := h.journey.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.NoContent(http.StatusCreated)
}

// listenResultWebSocket creates a socket connection to listen for test run results
func (h *runHandlers) listenResultWebSocket(c echo.Context) error {
	logger := h.logger.WithField("handler", "listenResultWebSocket")
	logger.Debug("client connected")

	var err error
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Error("list result websocket")
		return err
	}
	defer func() {
		logger.Debug("client disconnected")
		err := ws.Close()
		if err != nil {
			logger.WithError(err).Error("closing websocket")
		}
	}()

	pingFrequency := time.Second * 2
	pingTicker := time.NewTicker(pingFrequency)
	daemon := h.journey.Results()
	for {
		if daemon.ShouldStop() {
			daemon.Stopped()
			logger.Info("sending stop event")
			if err := ws.WriteJSON(newStoppedEvent()); err != nil {
				logger.WithError(err).Error("writing json to websocket")
			}
		}

		select {
		case <-pingTicker.C:
			logger.Debug("pinging websocket client")
			writeTimeout := time.Now().Add(time.Second)
			err := ws.SetWriteDeadline(writeTimeout)
			if err != nil {
				// we don't care about the error, just means the ws client has dropped the connection
				// and we want to close this gopher handler
				return err
			}
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				// same as above
				return err
			}
		case result, ok := <-daemon.Results():
			if !ok {
				logger.Error("error reading from result channel")
				break
			}
			logger.WithField("testId", result.Id).Info("sending result event")
			if err := ws.WriteJSON(newResultEvent(result)); err != nil {
				logger.WithError(err).Error("writing json to websocket")
				break
			}
		}
	}
}

// stopHandler sends signal to stop running test
func (h *runHandlers) stopRunHandler(c echo.Context) error {
	h.journey.StopTestRun()
	return nil
}

type StoppedEvent struct {
	Stopped bool `json:"stopped"`
}

func newStoppedEvent() StoppedEvent {
	return StoppedEvent{true}
}

type ResultEvent struct {
	Test results.TestCase `json:"test"`
}

func newResultEvent(testResult results.TestCase) ResultEvent {
	return ResultEvent{testResult}
}
