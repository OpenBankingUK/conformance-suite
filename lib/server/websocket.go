package server

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// WebSocketHandler for handling WebSocket connections.
type WebSocketHandler struct {
	upgrader *websocket.Upgrader
}

// NewWebSocketUpgrader creates a new websocket.Ugprader.
func NewWebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}
}

// MessageOut is the structure of the message that gets sent to the
// WebSocket client.
type MessageOut struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
}

// Handle a single WebSocket connection.
func (h *WebSocketHandler) Handle(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("WebSocketHandler:Handle -> Upgrade")
		return err
	}
	defer ws.Close()
	ws.EnableWriteCompression(true)

	logrus.WithFields(logrus.Fields{
		"RemoteAddr": ws.RemoteAddr(),
	}).Info("WebSocketHandler:Handle -> Upgrade")

	for {
		msgOut := &MessageOut{
			Number: rand.Int(),
			Type:   "update",
		}
		logrus.WithFields(logrus.Fields{
			"msgOut": msgOut,
		}).Info("WebSocketHandler:Handle -> WriteJSON")

		// Write
		if err := ws.WriteJSON(msgOut); err != nil {
			logrus.WithFields(logrus.Fields{
				"err":                    err,
				"IsCloseError":           websocket.IsCloseError(err),
				"IsUnexpectedCloseError": websocket.IsUnexpectedCloseError(err),
			}).Error("WebSocketHandler:Handle -> WriteJSON")
			c.Logger().Error(err)

			return err
		}

		// sleep for a bit before sending the next message
		time.Sleep(10 * time.Second)
	}
}
