package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"

	"github.com/stretchr/testify/require"
)

func TestGetTestCases(t *testing.T) {
	assert := test.NewAssert(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &versionmock.Version{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-ozone.json")
	assert.NoError(err)
	assert.NotNil(discoveryModel)

	request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(string(discoveryModel)), server)

	code, body, headers := request(http.MethodGet, "/api/test-cases",
		nil, server)

	// we should get back the test cases
	assert.NotNil(body)
	assert.Equal(http.StatusOK, code)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Contains(body.String(), `{"apiSpecification":{"name":"CustomTest-GetOzoneToken","url":""`)
	assert.Contains(body.String(), `testCases":[{"@id":"#co0001","name":"Post Account Consent"`)
}
