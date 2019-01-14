package version

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOutOfDateUpdateWarningVersion asserts that given an outdated version of
// the suite an update boolean is returned.
func TestOutOfDateUpdateWarningVersion(t *testing.T) {
	mockResponse := `
{
	"values": [{
		"name": "0.0.2",
		"date": "2019-01-11T13:56:34+0000",
		"message": "mocked response"
	}]
}`
	mockedServer, serverURL := test.HTTPServer(http.StatusOK, mockResponse, nil)
	defer mockedServer.Close()

	// Version helper
	v := New(serverURL)

	// Use an old version to test.
	version := "v0.0.0"
	// Asset that you get update boolean.
	_, flag, error := v.UpdateWarningVersion(version)
	assert.Equal(t, true, flag)
	assert.Equal(t, error, nil)
}

// TestNoUpdateUpdateWarningVersion asserts no updated required boolean when
// local version maches or is higher.
func TestNoUpdateUpdateWarningVersion(t *testing.T) {
	mockResponse := `
{
	"values": [{
		"name": "0.0.2",
		"date": "2019-01-11T13:56:34+0000",
		"message": "mocked response"
	}]
}`
	mockedServer, serverURL := test.HTTPServer(http.StatusOK, mockResponse, nil)
	defer mockedServer.Close()

	// Version helper
	v := New(serverURL)

	version := "v1000.0.0"
	message := ""
	flag := true
	message, flag, _ = v.UpdateWarningVersion(version)
	assert.Equal(t, false, flag)
	assert.Equal(t, "Conformance Suite is running the latest version "+v.GetHumanVersion(), message)
	version = FullVersion
	_, flag, _ = v.UpdateWarningVersion(version)
	assert.Equal(t, false, flag)

}

// TestBadStatusUpdateWarningVersionFail asserts that an appropriate/correct
// error message is return if BitBucket 40x status code is given.
func TestBadStatusUpdateWarningVersionFail(t *testing.T) {
	mockResponse := `
{
	"values": [{
		"name": "0.0.2",
		"date": "2019-01-11T13:56:34+0000",
		"message": "mocked response"
	}]
}`
	mockedServer, serverURL := test.HTTPServer(http.StatusBadRequest, mockResponse, nil)
	defer mockedServer.Close()

	// Version helper
	v := New(serverURL)

	message := ""
	// Check we get the appropriate error message.
	message, _, _ = v.UpdateWarningVersion(FullVersion)
	assert.Equal(t, message, "Version check is unavailable at this time.")

}

// TestHTTPErrorUpdateWarningVersion asserts the correct error message
// is returned if BitBucket cannot return tags.
func TestHTTPErrorUpdateWarningVersion(t *testing.T) {
	// Version helper
	// Update BitBucketAPIRepository to produce a no such host.
	v := New("https://.com")
	message, flag, error := v.UpdateWarningVersion(FullVersion)
	// Assert that update fag is false.
	assert.Equal(t, flag, false)
	// Assert the default UI/Human error message is returned.
	assert.Equal(t, message, "Version check is unavailable at this time.")
	// Asset that an error() is actually returned.
	assert.NotEqual(t, error, nil)

}

// TestHaveVersionUpdateWarningVersion assert that if a version has
// no length the correct error message is returned.
func TestHaveVersionUpdateWarningVersion(t *testing.T) {
	version := ""

	// Version helper
	// Update BitBucketAPIRepository to produce a no such host.
	v := New("https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-suite/refs/tags")

	message, flag, error := v.UpdateWarningVersion(version)
	// Assert that update fag is false.
	assert.Equal(t, flag, false)
	// Assert the default UI/Human error message is returned.
	assert.Equal(t, message, "Version check is unavailable at this time.")
	// Asset that an error() is actually returned.
	assert.NotEqual(t, error, nil)
	// Assert error message is correct.
	assert.Errorf(t, error, "no version found")

}
