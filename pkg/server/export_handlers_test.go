package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

func TestServerPostExport(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	exportRequest := ExportRequest{}
	exportRequestJSON, err := json.MarshalIndent(exportRequest, ``, `  `)
	require.NoError(err)
	require.NotNil(exportRequestJSON)

	// make the request
	code, body, headers := request(http.MethodPost, "/api/export?pretty", bytes.NewReader(exportRequestJSON), server)

	// do assertions
	results := []results.TestCase{}
	tokens := []events.AcquiredAccessToken{}
	exportResponse := ExportResponse{
		ExportRequest: exportRequest,
		HasPassed:     true,
		Results:       results,
		Tokens:        tokens,
	}
	exportResponseJSON, err := json.MarshalIndent(exportResponse, ``, `  `)
	require.NoError(err)
	require.NotNil(exportResponseJSON)

	require.NotNil(body)
	bodyExpected := string(exportResponseJSON)
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)

	require.Equal(http.StatusOK, code)
	require.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}
