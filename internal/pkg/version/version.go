	// Package version contains version information for Functional Conformance Suite.
package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

//  Version returns the semantic version (see http://semver.org).
const (
	// Version must conform to the format expected, major, minor and patch.
	major = "0"
	minor = "2"
	patch = "0"
	// Version is the full string version of Conformance Suite.
	FullVersion = major + "." + minor + "." + patch
	// VersionPrerelease is pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "alpha", "beta", "rc1", etc.
	Prerelease = "alpha"
)

// Version helper with capability to get release versions from source control repository
type Version struct {
	// bitBucketAPIRepository full URL of the TAG API 2.0 for the Conformance Suite.
	bitBucketAPIRepository string
}

// New returns a new instance of Version. bitBucketAPIRepository
func New(bitBucketAPIRepository string) Version {
		return Version{
		bitBucketAPIRepository: bitBucketAPIRepository,
	}
}

// Tag structure used map response of tags.
type Tag struct {
	Name          string `json:"name"`
	Date          string `json:"date"`
	CommitMessage string `json:"message"`
}

// TagsAPIResponse structure to map response.
type TagsAPIResponse struct {
	TagList []Tag `json:"values"`
}

func getTags(body []byte) (*TagsAPIResponse, error) {
	var s = new(TagsAPIResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func (v Version) GetHumanVersion() string {
	version := "v" + FullVersion
	release := Prerelease

	if release != "" {
		if !strings.HasSuffix(version, "-"+release) {
			// if we tagged a prerelease version then the release is in the version already.
			version += fmt.Sprintf("-%s", release)
		}
	}
	return version
}

// Versionformatter takes a string version number and returns just the numeric parts.
// This function is used when trying to compare two string versions that 'could'
// have non numerical properties.
func (v Version) Versionformatter(version string) (string, error) {
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			return "", fmt.Errorf("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	// Regex to remove all (non numeric OR period).
	reg, err := regexp.Compile("[^0-9.]")
	// Raise any errors running the expression.
	if err != nil {
		return "", errors.Wrap(err, "could not format version number.")

	}
	processedString := reg.ReplaceAllString(string(vo), "")

	return processedString, nil
}

// UpdateWarningVersion takes a version number and checks it against the
// latest tag version on Bitbucket, if a newer version is found it
// returns a message and bool value that can be used to inform a user
// a newer version is available for download.
func (v Version) UpdateWarningVersion(version string) (string, bool, error) {
	// A default message that can be presented to an end user.
	errorMessageUI := "Version check is unavailable at this time."

	// Some basic validation, check we have a version,
	if len(version) == 0 {
		return errorMessageUI, false, fmt.Errorf("no version found")
	}
	// Try to get the latest tag using the BitBucket API.
	resp, err := http.Get(v.bitBucketAPIRepository)
	if err != nil {
		// If network error then return message, flag to NOT update and actual error.
		return errorMessageUI, false, errors.Wrap(err, "Error: HTTP on GET to BitBucket API")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "cannot read body API error.")
		}

		s, _ := getTags([]byte(body))

		if len(s.TagList) == 0 {
			return errorMessageUI, false, fmt.Errorf("No Tags found")
		}

		latestTag := s.TagList[0].Name

		// Format version string to compare.
		versionLocal, err := v.Versionformatter(version)
		versionRemote, err := v.Versionformatter(latestTag)

		if versionLocal < versionRemote {
			errorMessageUI = fmt.Sprintf("Version v%s of the Conformance Suite is out-of-date, please update to v%s", versionLocal, versionRemote)
			return errorMessageUI, true, nil
		}
		// If local and remote version match or is higher then return false update flag.
		if versionLocal >= versionRemote {
			errorMessageUI = fmt.Sprintf("Conformance Suite is running the latest version %s", v.GetHumanVersion())
			return errorMessageUI, false, nil
		}

	} else {
		// handle anything else other than 200 OK.
		return errorMessageUI, false, nil
	}

	return errorMessageUI, false, nil

}
