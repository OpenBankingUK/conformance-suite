package authentication

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestSortAuthMethodsMostSecureFirstAllMethods(t *testing.T) {
	sorted := sortAuthMethodsMostSecureFirst([]string{
		"client_secret_basic",
		"client_secret_post",
		"client_secret_jwt",
		"private_key_jwt",
		"tls_client_auth",
	}, test.NullLogger())

	expected := AUTH_METHODS_SORTED_MOST_SECURE_FIRST

	test.NewRequire(t).Equal(expected, sorted)
}

func TestSortAuthMethodsMostSecureFirstSubsetOfMethods(t *testing.T) {
	sorted := sortAuthMethodsMostSecureFirst([]string{
		"client_secret_basic",
		"tls_client_auth",
	}, test.NullLogger())

	expected := []string{
		tls_client_auth,
		client_secret_basic,
	}
	test.NewRequire(t).Equal(expected, sorted)
}

func TestSortAuthMethodsMostSecureFirstDuplicateMethods(t *testing.T) {
	sorted := sortAuthMethodsMostSecureFirst([]string{
		"client_secret_basic",
		"client_secret_basic",
	}, test.NullLogger())

	expected := []string{
		client_secret_basic,
		client_secret_basic,
	}
	test.NewRequire(t).Equal(expected, sorted)
}

func TestSortAuthMethodsMostSecureFirstNonMatching(t *testing.T) {
	sorted := sortAuthMethodsMostSecureFirst([]string{
		"client_secret_basic",
		"private_key_jwt_bad_match",
		"tls_client_auth",
	}, test.NullLogger())

	expected := []string{
		tls_client_auth,
		client_secret_basic,
		"",
	}
	test.NewRequire(t).Equal(expected, sorted)
}

func TestOpenIdConfigWhenGetSuccessful(t *testing.T) {
	require := test.NewRequire(t)
	mockResponse, err := ioutil.ReadFile("../server/testdata/openid-configuration-mock.json")
	require.NoError(err)

	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK,
		string(mockResponse), nil)
	defer mockedServer.Close()

	config, err := OpenIdConfig(mockedServerURL, test.NullLogger())
	require.NoError(err)
	require.NotNil(config)
	authMethodsMostSecureFirst := []string{
		tls_client_auth,
		private_key_jwt,
		client_secret_jwt,
		client_secret_post,
		client_secret_basic,
	}
	expected := OpenIDConfiguration{
		TokenEndpoint:                     "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
		AuthorizationEndpoint:             "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>",
		Issuer:                            "https://modelobankauth2018.o3bank.co.uk:4101",
		TokenEndpointAuthMethodsSupported: authMethodsMostSecureFirst,
	}

	require.Equal(expected, config)
}

func TestOpenIdConfigWhenHttpResponseError(t *testing.T) {
	require := test.NewRequire(t)

	mockedBadServer, mockedBadServerURL := test.HTTPServer(http.StatusServiceUnavailable,
		"<h1>503 Service Temporarily Unavailable</h1>", nil)
	defer mockedBadServer.Close()

	_, err := OpenIdConfig(mockedBadServerURL, test.NullLogger())
	require.EqualError(err, fmt.Sprintf("failed to GET OpenID config: %s : HTTP response status: 503", mockedBadServerURL))
}

func TestOpenIdConfigWhenJsonParseFails(t *testing.T) {
	require := test.NewRequire(t)
	mockedBody := "<bad>json</bad>"
	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK,
		mockedBody, nil)
	defer mockedServer.Close()

	_, err := OpenIdConfig(mockedServerURL, test.NullLogger())
	require.EqualError(err, fmt.Sprintf("Invalid OpenID config JSON returned: %s : invalid character '<' looking for beginning of value", mockedServerURL))
}
