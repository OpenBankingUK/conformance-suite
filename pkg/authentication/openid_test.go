package authentication

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestOpenIdConfigWhenGetSuccessful(t *testing.T) {
	require := test.NewRequire(t)

	// https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration
	mockResponse, err := ioutil.ReadFile("../server/testdata/openid-configuration_ozone.json")
	require.NoError(err)

	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK, string(mockResponse), nil)
	defer mockedServer.Close()

	config, err := OpenIdConfig(mockedServerURL)
	require.NoError(err)
	require.NotNil(config)

	expected := OpenIDConfiguration{
		TokenEndpoint: "https://ob19-auth1.o3bank.co.uk:4201/token",
		TokenEndpointAuthMethodsSupported: []string{
			"client_secret_basic",
			"client_secret_jwt",
			"private_key_jwt",
			"tls_client_auth",
		},
		RequestObjectSigningAlgValuesSupported: []string{
			"none",
			"HS256",
			"RS256",
			"PS256",
		},
		AuthorizationEndpoint: "https://ob19-auth1-ui.o3bank.co.uk/auth",
		Issuer:                "https://ob19-auth1-ui.o3bank.co.uk",
		ResponseTypesSupported: []string{
			"code",
			"code id_token",
		},
	}

	require.Equal(expected, config)
}

func TestOpenIdConfigWhenHttpResponseError(t *testing.T) {
	require := test.NewRequire(t)

	mockedBody := "<h1>503 Service Temporarily Unavailable</h1>"
	mockedBadServer, mockedBadServerURL := test.HTTPServer(http.StatusServiceUnavailable, mockedBody, nil)
	defer mockedBadServer.Close()

	_, err := OpenIdConfig(mockedBadServerURL)
	expected := fmt.Sprintf("failed to GET OpenIDConfiguration config: url=%+v, StatusCode=503, body=<h1>503 Service Temporarily Unavailable</h1>", mockedBadServerURL)
	require.EqualError(err, expected)
}

func TestOpenIdConfigWhenJsonParseFails(t *testing.T) {
	require := test.NewRequire(t)
	mockedBody := "<bad>json</bad>"
	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK, mockedBody, nil)
	defer mockedServer.Close()

	_, err := OpenIdConfig(mockedServerURL)
	expected := fmt.Sprintf("Invalid OpenIDConfiguration: url=%+v: invalid character '<' looking for beginning of value", mockedServerURL)
	require.EqualError(err, expected)
}
