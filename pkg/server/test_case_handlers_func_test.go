package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetTestCases(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(NullLogger(), conditionalityCheckerMock{})

	discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	assert.NoError(err)
	assert.NotNil(discoveryModel)

	request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(string(discoveryModel)), server)

	code, body, headers := request(http.MethodGet, "/api/test-cases",
		nil, server)

	// we should get back the test cases
	assert.NotNil(body)
	assert.Equal(http.StatusOK, code)
	assert.Equal(headers["Content-Type"][0], "application/json; charset=UTF-8")
	assert.Contains(body.String(), `[{"apiSpecification":{"name":"Account and Transaction API Specification"`)
	assert.Contains(body.String(), `"testCases":[{"@id":"#t1000","name":"Create Account Access Consents"`)

	server.Shutdown(context.TODO())
}
