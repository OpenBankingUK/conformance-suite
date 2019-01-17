package server

import (
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// TestVersionCheckUpdateAvailable provided a version that is lower than the available
// version, a response denoting an available update should be returned.
func TestVersionCheckUpdateAvailable(t *testing.T) {
	assert := assert.New(t)

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

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Equal(http.StatusOK, code)
	expected := fmt.Sprintf(`{"version":"%s", "message":"%s", "update":%t}`, humanVersion, warningMsg, true)
	assert.JSONEq(body.String(), expected)
}

// TestVersionCheckNoUpdateAvailable provided a version that is greater than or equal to the available
// version, a response denoting an available update should be returned.
func TestVersionCheckNoUpdateAvailable(t *testing.T) {
	assert := assert.New(t)

	humanVersion := "0.1.2-RC1"
	warningMsg := "Conformance Suite is running the latest version 0.1.2-RC1"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, false, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Equal(http.StatusOK, code)
	expected := fmt.Sprintf(`{"version":"%s", "message":"%s", "update":%t}`, humanVersion, warningMsg, false)
	assert.JSONEq(body.String(), expected)
}

func TestVersionUpstreamUnavailableReturnsServerError(t *testing.T) {
	assert := assert.New(t)

	formatted := "0.1.2"
	warningMsg := "Conformance Suite is running the latest version 0.1.2-RC1"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return("")
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, false, errors.New("service error"))
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	assert.Equal(http.StatusInternalServerError, code)
}
