package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMockHTTPServer_Minimal_Usage(t *testing.T) {
	server, client := MockHTTPServer(http.StatusOK, "body", nil, nil)

	response, err := client.Get("http://localhost")
	require.NoError(t, err)

	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "body", string(body))

	server.Close()
}

func TestMockHTTPServer_AddHeaders(t *testing.T) {
	headers := map[string]string{"key": "value"}
	server, client := MockHTTPServer(http.StatusOK, "body", headers, nil)

	response, err := client.Get("http://localhost")
	require.NoError(t, err)

	assert.Equal(t, "value", response.Header.Get("key"))

	server.Close()
}

func TestMockHTTPServer_ReturnsError(t *testing.T) {
	server, client := MockHTTPServer(http.StatusOK, "body", nil, errors.New("some error"))

	_, err := client.Get("http://localhost")
	assert.Error(t, err, "some error")

	server.Close()
}
