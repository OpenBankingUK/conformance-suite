// Package version contains version information for Functional Conformance Suite
package version

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Version
const (
	// Version must conform to the format expected, major, minor and patch.
	major = "0"
	minor = "1"
	patch = "0"
	// Version is the full string version of conformance suite.
	Version = major + "." + minor + "." + patch
)

// VersionPrerelease is pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "alpha", "beta", "rc1", etc.
const VersionPrerelease = "pre-alpha"

// BitBucketAPIRepository full URL of the TAG API 2.0 for the Conformance Suite.
var BitBucketAPIRepository = "https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-suite/refs/tags"

// Tag structure.
type Tag struct {
	Name string `json:"name"`
}

// TagsAPIResponse structure.
type TagsAPIResponse struct {
	TagList []Tag `json:"values"`
}

func getTags(body []byte) (*TagsAPIResponse, error) {
	var s = new(TagsAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func GetHumanVersion() string {
	version := "v" + Version
	release := VersionPrerelease

	if release != "" {
		if !strings.HasSuffix(version, "-"+release) {
			// if we tagged a prerelease version then the release is in the version already.
			version += fmt.Sprintf("-%s", release)
		}
	}
	return version
}

// Versionformatter takes a version number and returns just the numeric parts.
func Versionformatter(version string) string {
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
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	// Regex to remove all (non numeric OR period).
	reg, err := regexp.Compile("[^0-9.]")
	// Raise any errors running the expression.
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
	}
	processedString := reg.ReplaceAllString(string(vo), "")

	return processedString
}

// UpdateWarningVersion takes a version number and checks it against the
// latest tag version on Bitbucket, if a newer version is found it
// returns a message and bool value that can be used to inform a user
// a newer version is available for download.
func UpdateWarningVersion(version string) (string, bool) {
	var buf bytes.Buffer

	// Some basic validation, check we have a version.
	if len(version) != 0 {
		fmt.Fprintf(&buf, " (%s)", Version)
	}

	resp, err := http.Get(BitBucketAPIRepository)

	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		s, err := getTags([]byte(body))

		// Format version string to compare.
		versionLocal := Versionformatter(version)
		versionRemote := Versionformatter(s.TagList[0].Name)

		if versionLocal < versionRemote {
			message := fmt.Sprintf("Version v%s of the Conformance Suite is out-of-date, please update to v%s", versionLocal, versionRemote)
			return message, true
		}
		fmt.Println(s.TagList[0].Name)

	} else {
		// handle anything else other than 200 OK.
		return "Version check is univailable at this time.", false
	}

	return "Error", false

}
