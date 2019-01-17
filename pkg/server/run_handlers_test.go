package server

import (
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestServerRunStartPost - tests /api/run/start
func TestServerRunStartPost(t *testing.T) {
	require := require.New(t)

	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// make the request
	//
	// `?pretty` makes the JSON more readable in the event of a failure
	// see the example: https://echo.labstack.com/guide/response#json-pretty
	code, body, headers := request(
		http.MethodPost,
		"/api/run/start?pretty",
		nil,
		server)

	// do assertions
	require.Equal(http.StatusBadRequest, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])

	require.NotNil(body)

	bodyExpected := `{ "error": "error running test cases, test cases not set" }`
	bodyActual := body.String()
	// do not use `require.Equal`.
	require.JSONEq(bodyExpected, bodyActual)
}
