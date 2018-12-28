package server

import (
	"strings"
	"testing"

	model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"net/http/httptest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// TestWebSocketHandler_Handle - tests that it handles an incoming connection and
// that it can receive and send messages.
func TestServerWebSocketHandlerHandle(t *testing.T) {
	require := require.New(t)

	echoServer := NewServer(nullLogger(), model.NewConditionalityChecker())
	webServer := httptest.NewServer(echoServer)
	defer func() {
		require.NoError(echoServer.Shutdown(nil))
	}()
	defer webServer.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	urlStr := "ws" + strings.TrimPrefix(webServer.URL, "http") + "/api/ws"
	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	require.NoError(err)
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	messagesToRead := 1
	for i := 0; i < messagesToRead; i++ {
		body := &MessageOut{}

		require.NoError(ws.ReadJSON(body))
		require.Equal("update", body.Type)
	}
}
