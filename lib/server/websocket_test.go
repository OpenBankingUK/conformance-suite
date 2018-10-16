package server

import (
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Test WebSocket connection are being handled
// and that messages get sent and received.
func TestServer_WebSocket(t *testing.T) {
	t.Parallel()

	server := NewServer()
	assert.NotNil(t, server)

	srv := httptest.NewServer(server)
	defer srv.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	urlStr := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws"
	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	messagesToRead := 1
	for i := 0; i < messagesToRead; i++ {
		var body MessageOut
		err := ws.ReadJSON(&body)

		assert.NoError(t, err)
		assert.Equal(t, "update", body.Type)
	}
}
