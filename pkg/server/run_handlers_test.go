package server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

const (
	prefix = ""
	indent = "    "
)

// TestServerRunStartPost - tests /api/run
func TestServerRunStartPost(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &versionmock.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	code, body, headers := request(
		http.MethodPost,
		"/api/run",
		nil,
		server)

	// do assertions
	require.NotNil(body)

	expected := `{ "error": "error test cases not generated" }`
	actual := body.String()
	require.JSONEq(expected, actual)

	require.Equal(http.StatusBadRequest, code)
	require.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}

func TestServerRunHandlersnewTestCaseResultWebSocketEvent(t *testing.T) {
	require := test.NewRequire(t)

	testCaseResult := results.TestCase{
		Id:   "#t1025",
		Pass: true,
	}
	wsEvent := newTestCaseResultWebSocketEvent(testCaseResult)

	wsEventJson, err := json.MarshalIndent(wsEvent, prefix, indent)
	require.NoError(err)
	require.NotNil(wsEventJson)

	expected := `
{
    "test": {
        "id": "#t1025",
        "pass": true,
        "metrics": {
            "response_time": 0,
            "response_size": 0
        }
    }
}
	`
	actual := string(wsEventJson)

	require.JSONEq(expected, actual)
}
