package authentication

import (
	"fmt"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestOpenIdConfigWhenHttpResponseError(t *testing.T) {
	require := test.NewRequire(t)

	mockedBody := "<h1>503 Service Temporarily Unavailable</h1>"
	mockedBadServer, mockedBadServerURL := test.HTTPServer(http.StatusServiceUnavailable, mockedBody, nil)
	defer mockedBadServer.Close()

	_, err := NewOpenIdConfigGetter().Get(mockedBadServerURL)
	expected := fmt.Sprintf("failed to GET OpenIDConfiguration config: url=%+v, StatusCode=503, body=<h1>503 Service Temporarily Unavailable</h1>", mockedBadServerURL)
	require.EqualError(err, expected)
}

func TestOpenIdConfigWhenJsonParseFails(t *testing.T) {
	require := test.NewRequire(t)
	mockedBody := "<bad>json</bad>"
	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK, mockedBody, nil)
	defer mockedServer.Close()

	_, err := NewOpenIdConfigGetter().Get(mockedServerURL)
	expected := fmt.Sprintf("Invalid OpenIDConfiguration: url=%+v: invalid character '<' looking for beginning of value", mockedServerURL)
	require.EqualError(err, expected)
}
