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
}

func newRunHandlers(journey Journey, upgrader *websocket.Upgrader) *runHandlers {
	return &runHandlers{
		journey:  journey,
		upgrader: upgrader,
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
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("WebSocketHandler:Handle -> Upgrade")
		return err
	}

	daemon := h.journey.Results()
	for {
		if daemon.ShouldStop() {
			if err := ws.WriteJSON(newStoppedEvent()); err != nil {
				logrus.WithError(err).Error("WebSocketHandler:Handle -> WriteJSON")
			}
			return ws.Close()
		}

		select {
		case result, ok := <-daemon.Results():
			if ok == false {
				logrus.Error("error reading from result channel")
				break
			}
			if err := ws.WriteJSON(newResultEvent(result)); err != nil {
				logrus.WithError(err).Error("WebSocketHandler:Handle -> WriteJSON")
				break
			}

		case err, ok := <-daemon.Errors():
			if ok == false {
				logrus.Error("error reading from errors channel")
				break
			}
			if err := ws.WriteJSON(newErrorEvent(err)); err != nil {
				logrus.WithError(err).Error("WebSocketHandler:Handle -> WriteJSON")
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

type ErrorEvent struct {
	Error error `json:"error"`
}

func newErrorEvent(err error) ErrorEvent {
	return ErrorEvent{err}
}
