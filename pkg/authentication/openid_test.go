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
	mockResponse, err := ioutil.ReadFile("../server/testdata/openid-configuration-mock.json")
	require.NoError(err)

	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK, string(mockResponse), nil)
	defer mockedServer.Close()

	config, err := OpenIdConfig(mockedServerURL)
	require.NoError(err)
	require.NotNil(config)
	authMethods := []string{
		"tls_client_auth",
		"client_secret_jwt",
		"client_secret_basic",
		"client_secret_post",
		"private_key_jwt",
	}
	responseTypesSupported := []string{
		"code",
		"code id_token",
	}
	expected := OpenIDConfiguration{
		TokenEndpoint:                     "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
		AuthorizationEndpoint:             "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>",
		Issuer:                            "https://modelobankauth2018.o3bank.co.uk:4101",
		TokenEndpointAuthMethodsSupported: authMethods,
		ResponseTypesSupported:            responseTypesSupported,
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
