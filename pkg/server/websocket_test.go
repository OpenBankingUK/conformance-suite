package server

import (
	"context"
	"strings"
	"testing"

	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"github.com/stretchr/testify/mock"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"net/http/httptest"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// TestWebSocketHandler_Handle - tests that it handles an incoming connection and
// that it can receive and send messages.
func TestServerWebSocketHandlerHandle(t *testing.T) {
	require := require.New(t)

	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	echoServer := NewServer(nullLogger(), model.NewConditionalityChecker(), v)
	webServer := httptest.NewServer(echoServer)
	defer func() {
		require.NoError(echoServer.Shutdown(context.TODO()))
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
