package test

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockHTTPServer_Minimal_Usage(t *testing.T) {
	server, url := HTTPServer(http.StatusOK, "body", nil)

	response, err := client.NewHTTPClient(client.DefaultTimeout).Get(url)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, "body", string(body))
	server.Close()
}

func TestMockHTTPServer_Adds_Headers(t *testing.T) {
	headers := map[string]string{"key1": "value1", "key2": "value2"}
	server, url := HTTPServer(http.StatusOK, "body", headers)

	response, err := client.NewHTTPClient(client.DefaultTimeout).Get(url)
	require.NoError(t, err)

	assert.Equal(t, "value1", response.Header.Get("key1"))
	assert.Equal(t, "value2", response.Header.Get("key2"))

	server.Close()
}
