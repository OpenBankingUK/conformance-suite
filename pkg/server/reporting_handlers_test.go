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

func TestServerGetReportBadRequestIfNoDiscoveryModelSet(t *testing.T) {
	assert := assert.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/report", nil, server)

	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Equal(http.StatusOK, code)
	assert.JSONEq(`{"id":"00000000-0000-0000-0000-000000000000","specifications":null}`, body.String())
}

func TestServerGetReport(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	assert.NoError(err)
	assert.NotNil(discoveryModel)

	code, body, headers := request(http.MethodPost, "/api/discovery-model",
		strings.NewReader(string(discoveryModel)), server)
	require.Equal(http.StatusCreated, code)

	code, body, headers = request(http.MethodGet, "/api/test-cases",
		nil, server)
	assert.NoError(err)

	code, body, headers = request(http.MethodGet, "/api/report", nil, server)

	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	//	assert.Equal(http.StatusOK, code)
}
