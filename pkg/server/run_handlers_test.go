package server

import (
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestServerRunStartPost - tests /api/run
func TestServerRunStartPost(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &versionmock.Version{})
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

	bodyExpected := `{ "error": "error discovery model not set" }`
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}
