package server_test

import (
	"strings"

	"net/http/httptest"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/lib/server"
)

var _ bool = Describe("Server Websocket", func() {
	var (
		the_server *server.Server
	)

	It("handles connection and receives and sends messages", func() {
		the_server = server.NewServer()
		srv := httptest.NewServer(the_server)
		defer srv.Close()

		// Convert http://127.0.0.1 to ws://127.0.0.
		urlStr := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws"
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
		assert.NoError(GinkgoT(), err)
		defer ws.Close()

		// Send message to server, read response and check to see if it's what we expect.
		messagesToRead := 1
		for i := 0; i < messagesToRead; i++ {
			var body server.MessageOut
			err := ws.ReadJSON(&body)

			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), "update", body.Type)
		}
	})
})
