package authentication

import (
	"encoding/json"
	"testing"
	"net/http"
	"io/ioutil"
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestOpenIDUnmarshal(t *testing.T) {
	require := test.NewRequire(t)

	data := `
{
	"token_endpoint": "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
	"authorization_endpoint": "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>",
	"issuer": "https://modelobankauth2018.o3bank.co.uk:4101"
}
	`
	expected := OpenIDConfiguration{
		TokenEndpoint:         "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
		AuthorizationEndpoint: "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>",
		Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
	}
	actual := OpenIDConfiguration{}
	require.NoError(json.Unmarshal([]byte(data), &actual))
	require.Equal(expected, actual)
}

func TestOpenIdConfigWhenGetSuccessful(t *testing.T) {
	require := test.NewRequire(t)

	mockResponse, err := ioutil.ReadFile("../server/testdata/openid-configuration-mock.json")
	require.NoError(err)

	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK,
		string(mockResponse), nil)
	defer mockedServer.Close()

	config, err := OpenIdConfig(mockedServerURL)
	require.NoError(err)
	require.NotNil(config)
}

func TestOpenIdConfigWhenHttpResponseError(t *testing.T) {
	require := test.NewRequire(t)

	mockedBadServer, mockedBadServerURL := test.HTTPServer(http.StatusServiceUnavailable,
		"<h1>503 Service Temporarily Unavailable</h1>", nil)
	defer mockedBadServer.Close()

	_, err := OpenIdConfig(mockedBadServerURL)
	require.EqualError(err, fmt.Sprintf("Failed to GET OpenID config: %s : HTTP response status: 503", mockedBadServerURL))
}

func TestOpenIdConfigWhenJsonParseFails(t *testing.T) {
	require := test.NewRequire(t)
	mockedBody := "<bad>json</bad>"
	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK,
		mockedBody, nil)
	defer mockedServer.Close()

	_, err := OpenIdConfig(mockedServerURL)
	require.EqualError(err, fmt.Sprintf("Invalid OpenID config JSON returned: %s : invalid character '<' looking for beginning of value", mockedServerURL))
}
