package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	humanVersion = "0.1.2-RC1"
	formatted    = "0.1.2"
)

func makeVersionMock(warningMsg string, updateAvaiable bool) *versionmock.Version {
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, updateAvaiable, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)
	return v
}

// TestVersionCheckUpdateAvailable provided a version that is lower than the available
// version, a response denoting an available update should be returned.
func TestVersionCheckUpdateAvailable(t *testing.T) {
	assert := test.NewAssert(t)

	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	updateAvaiable := true
	v := makeVersionMock(warningMsg, updateAvaiable)
	server := NewServer(testJourney(), nullLogger(), v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	expected := fmt.Sprintf(`{"version":"%s", "message":"%s", "update":%t}`, humanVersion, warningMsg, updateAvaiable)
	assert.JSONEq(body.String(), expected)

	assert.Equal(http.StatusOK, code)
	assert.Equal(expectedJsonHeaders(), headers)
}

// TestVersionCheckNoUpdateAvailable provided a version that is greater than or equal to the available
// version, a response denoting an available update should be returned.
func TestVersionCheckNoUpdateAvailable(t *testing.T) {
	assert := test.NewAssert(t)

	warningMsg := "Conformance Suite is running the latest version 0.1.2-RC1"
	updateAvaiable := false
	v := makeVersionMock(warningMsg, updateAvaiable)
	server := NewServer(testJourney(), nullLogger(), v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	expected := fmt.Sprintf(`{"version":"%s", "message":"%s", "update":%t}`, humanVersion, warningMsg, updateAvaiable)
	assert.JSONEq(body.String(), expected)

	assert.Equal(http.StatusOK, code)
	assert.Equal(expectedJsonHeaders(), headers)
}

func TestVersionUpstreamUnavailableReturnsServerError(t *testing.T) {
	assert := test.NewAssert(t)

	formatted := "0.1.2"
	warningMsg := "Conformance Suite is running the latest version 0.1.2-RC1"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return("")
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, false, errors.New("service error"))
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(testJourney(), nullLogger(), v)
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()

	code, body, headers := request(http.MethodGet, "/api/version", nil, server)

	assert.NotNil(body)
	assert.Equal(http.StatusInternalServerError, code)
	assert.Equal(expectedJsonHeaders(), headers)
}
