package server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// /api/discovery-model/validate - POST - When invalid JSON returns error message
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenInvalidJSON(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &versionmock.Version{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	discoveryModel := `{ "bad-json" }`
	expected := `{"error": "code=400, message=Syntax error: offset=14, error=invalid character '}' after object key"}`

	code, body, _ := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(string(expected), body.String())
	assert.Equal(http.StatusBadRequest, code)
}

// /api/discovery-model/validate - POST - When incomplete model returns validation failures messages
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenIncomplete(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &versionmock.Version{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
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

	code, body, _ := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(expected, body.String())
	assert.Equal(http.StatusBadRequest, code)
}
