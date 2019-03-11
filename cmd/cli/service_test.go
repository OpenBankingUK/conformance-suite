package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestServiceVersion(t *testing.T) {
	response := `{"version": "KAOS in CONTROL"}`
	server, url := test.HTTPServer(http.StatusOK, response, nil)
	defer server.Close()
	conn := &Connection{Client: &http.Client{}}
	service := newService(url, url, conn)

	version, err := service.Version()

	assert.NoError(t, err)
	assert.Equal(t, VersionResponse{Version: "KAOS in CONTROL"}, version)
}

func TestSetDiscoveryModel(t *testing.T) {
	server, url := test.HTTPServer(http.StatusCreated, "", nil)
	defer server.Close()
	conn := &Connection{Client: &http.Client{}}
	service := newService(url, url, conn)

	err := service.SetDiscoveryModel("README.md")

	assert.NoError(t, err)
}

func TestSetConfig(t *testing.T) {
	server, url := test.HTTPServer(http.StatusCreated, "", nil)
	defer server.Close()
	conn := &Connection{Client: &http.Client{}}
	service := newService(url, url, conn)

	err := service.SetConfig("README.md")

	assert.NoError(t, err)
}

func TestTestCases(t *testing.T) {
	server, url := test.HTTPServer(http.StatusOK, "", nil)
	defer server.Close()
	conn := &Connection{Client: &http.Client{}}
	service := newService(url, url, conn)

	err := service.TestCases()

	assert.NoError(t, err)
}
