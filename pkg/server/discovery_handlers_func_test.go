package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
)

// /api/discovery-model/validate - POST - When invalid JSON returns error message
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenInvalidJSON(t *testing.T) {
	assert := test.NewAssert(t)

	server := NewServer(testJourney(), nullLogger(), &versionmock.Version{})
	defer func() {
		assert.NoError(server.Shutdown(context.TODO()))
	}()

	discoveryModel := `{ "bad-json" }`
	expected := `{"error": "code=400, message=Syntax error: offset=14, error=invalid character '}' after object key"}`

	code, body, headers := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(string(expected), body.String())
	assert.Equal(http.StatusBadRequest, code)
	assert.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}

// /api/discovery-model/validate - POST - When incomplete model returns validation failures messages
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenIncomplete(t *testing.T) {
	assert := test.NewAssert(t)

	server := NewServer(testJourney(), nullLogger(), &versionmock.Version{})
	defer func() {
		assert.NoError(server.Shutdown(context.TODO()))
	}()

	discoveryModel := `{}`
	expected := `{ "error":
					[
						{"key": "DiscoveryModel.Name", "error": "Field 'DiscoveryModel.Name' is required"},
						{"key": "DiscoveryModel.Description", "error": "Field 'DiscoveryModel.Description' is required"},
						{"key": "DiscoveryModel.DiscoveryVersion", "error": "Field 'DiscoveryModel.DiscoveryVersion' is required"},
						{"key": "DiscoveryModel.TokenAcquisition", "error": "Field 'DiscoveryModel.TokenAcquisition' is required"},
						{"key": "DiscoveryModel.DiscoveryItems", "error": "Field 'DiscoveryModel.DiscoveryItems' is required"}
                    ]
				}`

	code, body, headers := request(http.MethodPost, "/api/discovery-model", strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(expected, body.String())
	assert.Equal(http.StatusBadRequest, code)
	assert.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}

// TestServerDiscoveryModelPOSTResolvesValuesUsingOpenidConfigurationURIs - tests that a HTTP GET is called for each
// `discoveryItems` using the url `openidConfigurationUri`.
func TestServerDiscoveryModelPOSTResolvesValuesUsingOpenidConfigurationURIs(t *testing.T) {
	require := test.NewRequire(t)

	mockResponse, err := ioutil.ReadFile("./testdata/openid-configuration-mock.json")
	require.NoError(err)

	mockedServer, mockedServerURL := test.HTTPServer(http.StatusOK, string(mockResponse), nil)
	defer mockedServer.Close()

	expected := `
		{
      "token_endpoints": {
        "schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json":
          "https://modelobank2018.o3bank.co.uk:4201/<token_mock>"
      },
      "most_secure_token_endpoint_auth_method": {
        "schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json":
          "tls_client_auth"
      },
      "authorization_endpoints": {
        "schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json":
          "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>"
			},
      "issuers": {
				"schema_version=https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json":
          "https://modelobankauth2018.o3bank.co.uk:4101"
			}
		}`

	// modify `ob-v3.0-ozone.json` to make it point to mockedServerURL
	discoveryJSON, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-ozone.json")
	require.NoError(err)
	require.NotNil(discoveryJSON)

	discoveryModel := &discovery.Model{}
	require.NoError(json.Unmarshal(discoveryJSON, &discoveryModel))

	// make `openidConfigurationUri` point to `mockedServerURL`
	require.NotEmpty(discoveryModel.DiscoveryModel.DiscoveryItems)
	for index := range discoveryModel.DiscoveryModel.DiscoveryItems {
		discoveryItem := &discoveryModel.DiscoveryModel.DiscoveryItems[index]
		discoveryItem.OpenidConfigurationURI = mockedServerURL
	}

	// make new discoveryModel POST body
	postBody, err := json.Marshal(discoveryModel)
	require.NoError(err)
	require.NotNil(postBody)

	server := NewServer(testJourney(), nullLogger(), &versionmock.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodPost, "/api/discovery-model", bytes.NewReader(postBody), server)

	require.NotNil(body)
	require.JSONEq(expected, body.String())
	require.Equal(http.StatusCreated, code)
	require.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}

// TestServerDiscoveryModelPOSTReturnsErrorsWhenItCannotResolveOpenidConfigurationURIs - tests that errors are
// returned when HTTP GET to `openidConfigurationUri` fails.
func TestServerDiscoveryModelPOSTReturnsErrorsWhenItCannotResolveOpenidConfigurationURIs(t *testing.T) {
	require := test.NewRequire(t)

	mockedBadServer, mockedBadServerURL := test.HTTPServer(http.StatusInternalServerError, ``, nil)
	defer mockedBadServer.Close()

	expected := fmt.Sprintf(`
{
    "error": [
        {
			"key": "DiscoveryModel.DiscoveryItems[0].OpenidConfigurationURI",
			"error": "failed to GET OpenID config: %s : HTTP response status: 500"
        }
    ]
}
	`, mockedBadServerURL)

	// modify `ob-v3.0-ozone.json` to make it point to mockedServerURL
	discoveryJSON, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-ozone.json")
	require.NoError(err)
	require.NotNil(discoveryJSON)

	discoveryModel := &discovery.Model{}
	require.NoError(json.Unmarshal(discoveryJSON, &discoveryModel))

	// make `openidConfigurationUri` point to `mockedBadServerURL`
	require.NotEmpty(discoveryModel.DiscoveryModel.DiscoveryItems)
	for index := range discoveryModel.DiscoveryModel.DiscoveryItems {
		discoveryItem := &discoveryModel.DiscoveryModel.DiscoveryItems[index]
		discoveryItem.OpenidConfigurationURI = mockedBadServerURL
	}

	// make new discoveryModel POST body
	postBody, err := json.Marshal(discoveryModel)
	require.NoError(err)
	require.NotNil(postBody)

	server := NewServer(testJourney(), nullLogger(), &versionmock.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodPost, "/api/discovery-model", bytes.NewReader(postBody), server)

	require.NotNil(body)
	require.JSONEq(expected, body.String())
	require.Equal(http.StatusBadRequest, code)
	require.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
}
