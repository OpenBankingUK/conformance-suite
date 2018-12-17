package server

import (
	"strings"
	"testing"

	model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"net/http/httptest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketHandler_Handle - tests that it handles an incoming connection and
// that it can receive and send messages.
func TestWebSocketHandler_Handle(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(NewServer(NullLogger(), model.NewConditionalityChecker()))
	defer server.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	urlStr := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/ws"
	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	assert.NoError(err)
	defer ws.Close()

	// Send message to server, read response and check to see if it's what we expect.
	messagesToRead := 1
	for i := 0; i < messagesToRead; i++ {
		body := &MessageOut{}
		err := ws.ReadJSON(body)

		assert.NoError(err)
		assert.Equal("update", body.Type)
	}
}
