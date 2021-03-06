package server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/executors/results"
	"github.com/OpenBankingUK/conformance-suite/pkg/test"
	versionmock "github.com/OpenBankingUK/conformance-suite/pkg/version/mocks"
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

	code, body, headers := request(http.MethodPost, "/api/run", nil, server)

	// do assertions
	require.NotNil(body)

	expected := `{ "error": "error test cases not generated" }`
	actual := body.String()
	require.JSONEq(expected, actual)

	require.Equal(http.StatusBadRequest, code)
	require.Equal(expectedJSONHeaders(), headers)
}

func TestServerRunHandlersnewTestCaseResultWebSocketEvent(t *testing.T) {
	require := test.NewRequire(t)

	testCaseResult := results.TestCase{
		Id:         "#t1025",
		Pass:       true,
		Detail:     "Example Test Case",
		RefURI:     "https://openbanking.org.uk/ref/uri",
		Endpoint:   "/foobar",
		HttpStatus: "200",
	}
	wsEvent := newTestCaseResultWebSocketEvent(testCaseResult)

	wsEventJSON, err := json.MarshalIndent(wsEvent, prefix, indent)
	require.NoError(err)
	require.NotNil(wsEventJSON)

	expected := `
{
	"type": "ResultType_TestCaseResult",
    "test": {
        "id": "#t1025",
        "pass": true,
        "metrics": {
            "response_time": 0,
            "response_size": 0
        },
		"detail": "Example Test Case",
		"refURI": "https://openbanking.org.uk/ref/uri",
		"endpoint": "/foobar",
		"httpStatusCode": "200"
    }
}
	`
	actual := string(wsEventJSON)

	require.JSONEq(expected, actual)
}

func TestServerRunHandlersnewTestCasesCompletedWebSocketEvent(t *testing.T) {
	require := test.NewRequire(t)

	wsEvent := newTestCasesCompletedWebSocketEvent(true)

	wsEventJSON, err := json.MarshalIndent(wsEvent, prefix, indent)
	require.NoError(err)
	require.NotNil(wsEventJSON)

	expected := `
{
    "type": "ResultType_TestCasesCompleted",
    "value": true
}
	`
	actual := string(wsEventJSON)

	require.JSONEq(expected, actual)
}
