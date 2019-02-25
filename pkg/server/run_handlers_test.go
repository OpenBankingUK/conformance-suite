package server

import (
	"context"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
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
	require.Equal(http.StatusBadRequest, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])

	require.NotNil(body)

	bodyExpected := `{ "error": "error test cases not generated" }`
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}
