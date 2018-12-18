package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// /api/discovery-model/validate - POST - When valid model returns request payload
func TestServerDiscoveryModelPOSTValidateReturnsRequestPayloadWhenValid(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(NullLogger(), conditionalityCheckerMock{})

	discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	assert.NoError(err)
	assert.NotNil(discoveryModel)

	code, body, headers := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(string(discoveryModel)), server)

	// we should get back the config
	assert.NotNil(body)
	assert.Equal(headers["Content-Type"][0], "application/json; charset=UTF-8")
	assert.JSONEq(string(discoveryModel), body.String())
	assert.Equal(http.StatusCreated, code)

	server.Shutdown(context.TODO())
}

// /api/discovery-model/validate - POST - When invalid JSON returns error message
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenInvalidJSON(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(NullLogger(), conditionalityCheckerMock{})

	discoveryModel := `{ "bad-json" }`
	expected := `{"error": "code=400, message=Syntax error: offset=14, error=invalid character '}' after object key"}`

	code, body, _ := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(string(expected), body.String())
	assert.Equal(http.StatusBadRequest, code)

	server.Shutdown(context.TODO())
}

// /api/discovery-model/validate - POST - When incomplete model returns validation failures messages
func TestServerDiscoveryModelPOSTValidateReturnsErrorsWhenIncomplete(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(NullLogger(), conditionalityCheckerMock{})

	discoveryModel := `{}`
	expected := `{ "error": 
					[ 
						{"key": "DiscoveryModel.Name", "error": "Field 'DiscoveryModel.Name' is required"},
						{"key": "DiscoveryModel.Description", "error": "Field 'DiscoveryModel.Description' is required"},
						{"key": "DiscoveryModel.DiscoveryVersion", "error": "Field 'DiscoveryModel.DiscoveryVersion' is required"},
						{"key": "DiscoveryModel.DiscoveryItems", "error": "Field 'DiscoveryModel.DiscoveryItems' is required"} 
                    ]
				}`

	code, body, _ := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(discoveryModel), server)

	assert.NotNil(body)
	assert.JSONEq(expected, body.String())
	assert.Equal(http.StatusBadRequest, code)

	server.Shutdown(context.TODO())
}
