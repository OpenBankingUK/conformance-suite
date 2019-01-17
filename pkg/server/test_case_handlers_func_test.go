package server

import (
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"context"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTestCases(t *testing.T) {
	assert := assert.New(t)

	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

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
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Contains(body.String(), `{"apiSpecification":{"name":"CustomTest-GetOzoneToken","url":""`)
	assert.Contains(body.String(), `testCases":[{"@id":"#co0001","name":"Post Account Consent"`)
}
