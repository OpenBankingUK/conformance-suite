package test

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMockHTTPServer_Minimal_Usage(t *testing.T) {
	server, url := MockHTTPServer(http.StatusOK, "body", nil)

	response, err := http.Get(url)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, "body", string(body))
	server.Close()
}

func TestMockHTTPServer_Adds_Headers(t *testing.T) {
	headers := map[string]string{"key1": "value1", "key2": "value2"}
	server, url := MockHTTPServer(http.StatusOK, "body", headers)

	response, err := http.Get(url)
	require.NoError(t, err)

	assert.Equal(t, "value1", response.Header.Get("key1"))
	assert.Equal(t, "value2", response.Header.Get("key2"))
	server.Close()
}
