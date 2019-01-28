package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
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
	logger := h.logger.WithField("handler", "listen_results")

	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Error("list result websocket")
		return err
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			logger.WithError(err).Error("closing websocket")
		}
	}()

	daemon := h.journey.Results()
	for {
		if daemon.ShouldStop() {
			if err := ws.WriteJSON(newStoppedEvent()); err != nil {
				logger.WithError(err).Error("writing json to websocket")
			}
		}

		select {
		case result, ok := <-daemon.Results():
			if !ok {
				logger.Error("error reading from result channel")
				break
			}
			if err := ws.WriteJSON(newResultEvent(result)); err != nil {
				logger.WithError(err).Error("writing json to websocket")
				break
			}

		case err, ok := <-daemon.Errors():
			if !ok {
				logrus.Error("error reading from errors channel")
				break
			}
			if err != nil {
				if err := ws.WriteJSON(newErrorEvent(err)); err != nil {
					logger.WithError(err).Error("writing json to websocket")
					break
				}
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

type ErrorEvent struct {
	Error string `json:"error"`
}

func newErrorEvent(err error) ErrorEvent {
	return ErrorEvent{err.Error()}
}
