// Package version contains version information for Functional Conformance Suite.
package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/OpenBankingUK/conformance-suite/pkg/client"
	hashiVer "github.com/hashicorp/go-version"
	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

// Checker returns the semantic version (see http://semver.org).
const (
	// Checker must conform to the format expected, major, minor and patch.
	// @NEW-SPEC-RELEASE - make sure new version is accounted for
	// @NEW-RELEASE - make sure new version is accounted for
	//v1.7.6 - this comment allows searching
	major = "1"
	minor = "8"
	patch = "0"

	//FullVersion -  Checker is the full string version of Conformance Suite.
	FullVersion = major + "." + minor + "." + patch

	// Prerelease - is pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "alpha", "beta", "rc1", etc.
	Prerelease          = ""
	GitHubAPIRepository = "https://api.github.com/repos/OpenBankingUK/conformance-suite/tags"
)

// Checker defines functionality to reason about the current version of the software and if updates are available
type Checker interface {
	GetHumanVersion() string
	VersionFormatter(version string) (string, error)
	UpdateWarningVersion(version string) (string, bool, error)
}

// GitHub helper with capability to get release versions from source control repository
type GitHub struct {
	// gitHubAPIRepository full URL of the TAG API 2.0 for the Conformance Suite.
	gitHubAPIRepository string
}

// NewGitHub returns a new instance of Checker.
func NewGitHub(bitBucketAPIRepository string) GitHub {
	return GitHub{
		gitHubAPIRepository: bitBucketAPIRepository,
	}
}

// Tag structure used map response of tags.
type Tag struct {
	Name string `json:"name"`
}

// TagsAPIResponse structure to map response.
type TagsAPIResponse struct {
	TagList []Tag
}

func (t Tag) LessThan(subject string) bool {
	tv, err := hashiVer.NewVersion(t.Name)
	if err != nil {
		return false
	}
	sv, err := hashiVer.NewVersion(subject)
	if err != nil {
		return false
	}

	return tv.LessThan(sv)
}

type tagList []Tag

func (t tagList) Len() int {
	return len(t)
}

func (t tagList) Less(i, j int) bool {
	return t[i].LessThan(t[j].Name)
}

func (t tagList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func getTags(body []byte) (*TagsAPIResponse, error) {
	var tags = new([]Tag)
	err := json.Unmarshal(body, &tags)
	s := &TagsAPIResponse{
		*tags,
	}
	return s, err
}

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func (v GitHub) GetHumanVersion() string {
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

// VersionFormatter takes a string version number and returns just the numeric parts.
// This function is used when trying to compare two string versions that 'could'
// have non numerical properties.
func (v GitHub) VersionFormatter(version string) (string, error) {
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
	reg := regexp.MustCompile("[^0-9.]")
	processedString := reg.ReplaceAllString(string(vo), "")

	return processedString, nil
}

// UpdateWarningVersion takes a version number and checks it against the
// latest tag version on GitHub, if a newer version is found it
// returns a message and bool value that can be used to inform a user
// a newer version is available for download.
func (v GitHub) UpdateWarningVersion(version string) (string, bool, error) {
	// A default message that can be presented to an end user.
	errorMessageUI := "Version check is unavailable at this time."

	// Some basic validation, check we have a version,
	if len(version) == 0 {
		return errorMessageUI, false, fmt.Errorf("no version found")
	}

	// Try to get the latest tag using the GitHub API.
	var resp *http.Response
	var err error

	if viper.GetBool("proxy_version_check") {
		resp, err = client.NewHTTPClient(client.DefaultTimeout).Get(v.gitHubAPIRepository)
	} else {
		resp, err = client.NewHTTPClientWithoutProxy(client.DefaultTimeout).Get(v.gitHubAPIRepository)
	}

	if err != nil {
		// If network error then return message, flag to NOT update and actual error.
		return errorMessageUI, false, errors.Wrap(err, "HTTP on GET to GitHub API")
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "cannot read body API error.")
		}

		err = resp.Body.Close()
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "error on update warning version")
		}

		s, err := getTags(body)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "error on update warning version")
		}

		if len(s.TagList) == 0 {
			return errorMessageUI, false, fmt.Errorf("no Tags found")
		}

		// Convert the list of tags to tagList and sort
		tags := convertSortTags(s)

		// Get latest tag
		latestTag := tags[len(tags)-1].Name

		// Format version string to compare.
		versionLocal, err := v.VersionFormatter(version)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "error on update warning version")
		}
		versionRemote, err := v.VersionFormatter(latestTag)
		if err != nil {
			return errorMessageUI, false, errors.Wrap(err, "error on update warning version")
		}

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

func convertSortTags(tar *TagsAPIResponse) tagList {
	tags := tagList{}
	for _, v := range tar.TagList {
		tags = append(tags, v)
	}
	sort.Sort(tags)
	return tags
}
