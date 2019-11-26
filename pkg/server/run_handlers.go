package server

import (
	"fmt"
	"net/http"

	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

const (
	// How often to send the ping event
	pingFrequency = time.Hour * 24
	// Deadline for write to timeout
	writeTimeout = time.Hour * 24
)

type runHandlers struct {
	journey  Journey
	upgrader *websocket.Upgrader
	logger   *logrus.Entry
}

func newRunHandlers(journey Journey, upgrader *websocket.Upgrader, logger *logrus.Entry) runHandlers {
	return runHandlers{
		journey:  journey,
		upgrader: upgrader,
		logger:   logger,
	}
}

// runStartPostHandler creates a new test run
func (h runHandlers) runStartPostHandler(c echo.Context) error {
	err := h.journey.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.NoContent(http.StatusCreated)
}

// listenResultWebSocket - /api/run/ws
// creates a socket connection to listen for test run results.
//
// Silence this linter error:
// `pkg/server/run_handlers.go:50:1: cyclomatic complexity 15 of func `(runHandlers).listenResultWebSocket` is high (> 13) (gocyclo)`
// nolint:gocyclo
func (h runHandlers) listenResultWebSocket(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	logger := h.logger.WithField("handler", "listenResultWebSocket").WithField("websocket", fmt.Sprintf("%p", ws))
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

	pingTicker := time.NewTicker(pingFrequency)
	daemon := h.journey.Results()
	events := h.journey.Events()
	for {
		if h.shouldStop(daemon, ws, logger) {
			break
		}

		select {
		case <-pingTicker.C:
			if !h.doSendPingMessage(ws, logger) {
				break
			}
		case isCompleted, ok := <-daemon.IsCompleted():
			if err := h.processTestCasesCompleted(ws, logger, isCompleted, ok); err != nil {
				break
			}
		case testCaseResult, ok := <-daemon.Results():
			if err := h.processTestCaseResult(ws, logger, testCaseResult, ok); err != nil {
				break
			}
		case event, ok := <-events.TokensChannel():
			if err := h.processAcquiredAccessTokenEvent(ws, logger, event, ok); err != nil {
				break
			}
		case event, ok := <-events.AllTokensChannel():
			if err := h.processAcquiredAllAccessTokensEvent(ws, logger, event, ok); err != nil {
				break
			}
		}
	}

	return nil
}

// shouldStop - True if caller should stop, false otherwise.
func (h runHandlers) shouldStop(daemon executors.DaemonController, ws *websocket.Conn, logger *logrus.Entry) bool {
	if daemon.ShouldStop() {
		daemon.Stopped()
		logger.Info("sending stop event")
		if err := ws.WriteJSON(newStoppedEvent()); err != nil {
			logger.WithError(err).Error("writing StoppedEvent")
			return true
		}
		return true
	}
	return false
}

// doSendPingMessage - If false, caller should terminate WebSocket connection.
func (h runHandlers) doSendPingMessage(ws *websocket.Conn, logger *logrus.Entry) bool {
	//logger.Debug("pinging websocket client")

	// We cannot return error here, if we do echo will try to write the error to conn
	// and we closed the ws with a defer func. So we return a bool, if it falsey the caller
	// should terminate the WebSocket connection.

	writeTimeout := time.Now().Add(writeTimeout)
	if err := ws.SetWriteDeadline(writeTimeout); err != nil {
		logger.WithError(err).Error("SetWriteDeadline failed")
		return false
	}
	if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
		logger.WithError(err).Error("WriteMessage failed")
		return false
	}
	return true
}

// stopHandler sends signal to stop running test
func (h runHandlers) stopRunHandler(c echo.Context) error {
	h.journey.StopTestRun()
	return nil
}

func (h runHandlers) processTestCasesCompleted(ws *websocket.Conn, logger *logrus.Entry, isCompleted bool, ok bool) error {
	if !ok {
		err := errors.New("error reading from daemon.IsCompleted channel")
		logger.Error(err)
		return err
	}

	wsEvent := newTestCasesCompletedWebSocketEvent(isCompleted)

	logger.WithFields(logrus.Fields{
		"wsEvent.Type": wsEvent.Type,
		"isCompleted":  isCompleted,
	}).Info("sending event")
	if err := ws.WriteJSON(wsEvent); err != nil {
		logger.WithError(err).Error("[processTestCasesCompleted] writing json to websocket")
		return err
	}

	return nil
}

func (h runHandlers) processTestCaseResult(ws *websocket.Conn, logger *logrus.Entry, result results.TestCase, ok bool) error {
	if !ok {
		err := errors.New("error reading from daemon.Results channel")
		logger.Error(err)
		return err
	}

	wsEvent := newTestCaseResultWebSocketEvent(result)
	logger.WithFields(logrus.Fields{
		"wsEvent.Type": wsEvent.Type,
		"result.Id":    result.Id,
	}).Info("sending event")
	if err := ws.WriteJSON(wsEvent); err != nil {
		logger.WithError(err).Error("[processTestCaseResult] writing json to websocket")
		return err
	}

	return nil
}

func (h runHandlers) processAcquiredAccessTokenEvent(ws *websocket.Conn, logger *logrus.Entry, event events.AcquiredAccessToken, ok bool) error {
	if !ok {
		err := errors.New("error reading from events.Tokens channel")
		logger.Error(err)
		return err
	}

	wsEvent := newAcquiredAccessTokenWebSocketEvent(event)
	logger.WithFields(logrus.Fields{
		"wsEvent.Type":    wsEvent.Type,
		"event.TokenName": event.TokenName,
	}).Info("sending event")
	if err := ws.WriteJSON(wsEvent); err != nil {
		logger.WithError(err).Error("[processAcquiredAccessTokenEvent] writing json to websocket")
		return err
	}

	return nil
}

func (h runHandlers) processAcquiredAllAccessTokensEvent(ws *websocket.Conn, logger *logrus.Entry, event events.AcquiredAllAccessTokens, ok bool) error {
	if !ok {
		err := errors.New("error reading from events.AllTokens channel")
		logger.Error(err)
		return err
	}

	wsEvent := newAcquiredAllAccessTokensWebSocketEvent(event)
	logger.WithFields(logrus.Fields{
		"wsEvent.Type":     wsEvent.Type,
		"event.TokenNames": event.TokenNames,
	}).Info("sending event")
	if err := ws.WriteJSON(wsEvent); err != nil {
		logger.WithError(err).Error("[processAcquiredAllAccessTokensEvent] writing json to websocket")
		return err
	}

	return nil
}

type StoppedEvent struct {
	Stopped bool `json:"stopped"`
}

func newStoppedEvent() StoppedEvent {
	return StoppedEvent{
		Stopped: true,
	}
}

type TestCaseResultWebSocketEvent struct {
	Type string           `json:"type"`
	Test results.TestCase `json:"test"`
}

func newTestCaseResultWebSocketEvent(testCaseResult results.TestCase) TestCaseResultWebSocketEvent {
	return TestCaseResultWebSocketEvent{
		Type: "ResultType_TestCaseResult",
		Test: testCaseResult,
	}
}

type TestCasesCompletedWebSocketEvent struct {
	Type  string `json:"type"`
	Value bool   `json:"value"`
}

func newTestCasesCompletedWebSocketEvent(isCompleted bool) TestCasesCompletedWebSocketEvent {
	return TestCasesCompletedWebSocketEvent{
		Type:  "ResultType_TestCasesCompleted",
		Value: isCompleted,
	}
}

type AcquiredAccessTokenWebSocketEvent struct {
	Type  string                     `json:"type"`
	Value events.AcquiredAccessToken `json:"value"`
}

func newAcquiredAccessTokenWebSocketEvent(event events.AcquiredAccessToken) AcquiredAccessTokenWebSocketEvent {
	return AcquiredAccessTokenWebSocketEvent{
		Type:  "ResultType_AcquiredAccessToken",
		Value: event,
	}
}

type AcquiredAllAccessTokensWebSocketEvent struct {
	Type  string                         `json:"type"`
	Value events.AcquiredAllAccessTokens `json:"value"`
}

func newAcquiredAllAccessTokensWebSocketEvent(event events.AcquiredAllAccessTokens) AcquiredAllAccessTokensWebSocketEvent {
	return AcquiredAllAccessTokensWebSocketEvent{
		Type:  "ResultType_AcquiredAllAccessTokens",
		Value: event,
	}
}
