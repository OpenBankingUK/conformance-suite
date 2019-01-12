package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// /api/discovery-model/validate - POST - When valid model returns request payload
// Test Case Disabled Pending bug investigation whereby the following code sent to the backend
//  "regex": "code=(.*)&?.*"
// is changed on return to
//  "regex": "code=(.*)\u0026?.*"
// See REFAPPS-543 - for bug investigation
func disableTestServerDiscoveryModelPOSTValidateReturnsRequestPayloadWhenValid(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	assert.NoError(err)
	assert.NotNil(discoveryModel)

	code, body, headers := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(string(discoveryModel)), server)

	// we should get back the config
	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.JSONEq(string(discoveryModel), body.String())
	assert.Equal(http.StatusCreated, code)
}

// /api/discovery-model/validate - POST - When invalid JSON returns error message
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenInvalidJSON(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{})
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
// func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenIncomplete(t *testing.T) {
// 	assert := assert.New(t)

// 	server := NewServer(nullLogger(), conditionalityCheckerMock{})
// 	defer func() {
// 		require.NoError(t, server.Shutdown(context.TODO()))
// 	}()

// 	discoveryModel := `{}`
// 	expected := `{ "error":
// 					[
// 						{"key": "DiscoveryModel.Name", "error": "Field 'DiscoveryModel.Name' is required"},
// 						{"key": "DiscoveryModel.Description", "error": "Field 'DiscoveryModel.Description' is required"},
// 						{"key": "DiscoveryModel.DiscoveryVersion", "error": "Field 'DiscoveryModel.DiscoveryVersion' is required"},
// 						{"key": "DiscoveryModel.DiscoveryItems", "error": "Field 'DiscoveryModel.DiscoveryItems' is required"}
//                     ]
// 				}`

// 	code, body, _ := request(http.MethodPost, "/api/discovery-model",
// 		strings.NewReader(discoveryModel), server)

// 	assert.NotNil(body)
// 	assert.JSONEq(expected, body.String())
// 	assert.Equal(http.StatusBadRequest, code)
// }
