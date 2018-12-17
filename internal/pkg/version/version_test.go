package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestTimeConsuming(t *testing.T) {

// 	fmt.Printf(GetHumanVersion())
// 	fmt.Printf("    Hello ============")

// 	assert.Equal(t, GetHumanVersion(), "")

// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// }

// TestOutOfDateUpdateWarningVersion asserts that given an outdated version of
// the suite an update boolean is returned.
func TestOutOfDateUpdateWarningVersion(t *testing.T) {
	// Use an old version to test.
	version := "v0.0.0"
	flag := false
	// Asset that you get update boolean.
	_, flag = UpdateWarningVersion(version)
	assert.Equal(t, true, flag)
}

// TestBadStatusUpdateWarningVersionFail asserts that an appropriate/correct
// error message is return if BitBucket 40x status code is given.
func TestBadStatusUpdateWarningVersionFail(t *testing.T) {
	message := ""
	// Modify the API URL to bad status 404.
	BitBucketAPIRepository = "https://api.bitbucket.org/2.0/repositories/openbankingteam/ooops/refs/tags"
	// Check we get the appropriate error message.
	message, _ = UpdateWarningVersion(Version)
	assert.Equal(t, message, "Version check is univailable at this time.")

}
