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
	BitBucketAPIRepository = serverURL
	defer mockedServer.Close()

	// Use an old version to test.
	version := "v0.0.0"
	// Asset that you get update boolean.
	_, flag, error := UpdateWarningVersion(version)
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
	BitBucketAPIRepository = serverURL
	defer mockedServer.Close()

	version := "v1000.0.0"
	message := ""
	flag := true
	message, flag, _ = UpdateWarningVersion(version)
	assert.Equal(t, false, flag)
	assert.Equal(t, "Conformance Suite is running the latest version "+GetHumanVersion(), message)
	version = Version
	_, flag, _ = UpdateWarningVersion(version)
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
	BitBucketAPIRepository = serverURL
	defer mockedServer.Close()

	message := ""
	// Check we get the appropriate error message.
	message, _, _ = UpdateWarningVersion(Version)
	assert.Equal(t, message, "Version check is unavailable at this time.")

}

// TestHTTPErrorUpdateWarningVersion asserts the correct error message
// is returned if BitBucket cannot return tags.
func TestHTTPErrorUpdateWarningVersion(t *testing.T) {
	// Update BitBucketAPIRepository to produce a no such host.
	BitBucketAPIRepository = "https://.com"
	message, flag, error := UpdateWarningVersion(Version)
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

	message, flag, error := UpdateWarningVersion(version)
	// Assert that update fag is false.
	assert.Equal(t, flag, false)
	// Assert the default UI/Human error message is returned.
	assert.Equal(t, message, "Version check is unavailable at this time.")
	// Asset that an error() is actually returned.
	assert.NotEqual(t, error, nil)
	// Assert error message is correct.
	assert.Errorf(t, error, "no version found")

}
