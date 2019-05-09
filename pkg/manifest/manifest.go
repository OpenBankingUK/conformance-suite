package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
)

// LoadScripts loads the scripts from JSON encoded contents of filename
// and returns Scripts objects
func LoadScripts(filename string) (Scripts, error) {
	if strings.HasPrefix(filename, "https://") {
		return Scripts{}, errors.New("https:// manifest loading not implemented")
	}
	path := strings.TrimPrefix(filename, "file://")
	plan, err := ioutil.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		plan, err = ioutil.ReadFile("../../" + path)
	}
	if err != nil {
		return Scripts{}, err
	}
	var m Scripts
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return Scripts{}, err
	}
	return m, nil
}

// DiscoveryPathsTestIDs -
type DiscoveryPathsTestIDs map[string][]string

// MapDiscoveryEndpointsToManifestTestIDs creates a mapping such that:
// - For each [endpoint + method] in the discovery file
// - Find all of the tests that exist in the manifest file, which contain the same [endpoint + method] combination
// - For each match, store that match in a map, which uses the endpoint as the map pair key and the map pair value
// is a list of each of the tests in the manifest relating to specified endpoint.
// - The value from the previous should be further broken down into another map, containing a list of each test id,
// where the keys in the second map are the http methods.
// Example output:
// 3 tests for "GET" method on the "/accounts" endpoint and 1 test for "HEAD" method.
//"/accounts": {
//	"GET": [
//		"OB-301-ACC-811741",
//		"OB-301-ACC-431102",
//		"OB-301-ACC-880736"
//	],
//	"HEAD": [
//		"HEAD-OB-301-ACC-431102"
//	]
//}
func MapDiscoveryEndpointsToManifestTestIDs(disco *discovery.Model, mf Scripts) DiscoveryPathsTestIDs {
	mapURLTests := make(DiscoveryPathsTestIDs)

	// Iterate the discoveryModel.discoveryItems.endpoints
	for _, discoItem := range disco.DiscoveryModel.DiscoveryItems {
		for _, discoEndpoint := range discoItem.Endpoints {
			// For each discovery item, iterate all the `uri` fields and see if there is a match.
			for _, mfScript := range mf.Scripts {
				if strings.EqualFold(discoEndpoint.Path, mfScript.URI) &&
					strings.EqualFold(discoEndpoint.Method, mfScript.Method) {
					key := fmt.Sprintf("%s %s", strings.ToUpper(mfScript.Method), discoEndpoint.Path)
					if _, ok := mapURLTests[key]; !ok {
						mapURLTests[key] = []string{}
					}
					mapURLTests[key] = append(mapURLTests[key], mfScript.ID)
				}
			}
		}
	}
	return mapURLTests
}

// FindUnmatchedManifestTests
// Find all the TestIDs from Manifest that have not been matched against an endpoint in the discovery model
func FindUnmatchedManifestTests(mf Scripts, mappedTests DiscoveryPathsTestIDs) []string {
	var result []string
	for _, script := range mf.Scripts {
		i := 0
		for _, v := range mappedTests {
			if isInArray(script.ID, v) {
				continue
			}
			i++
		}
		if i == len(mappedTests) {
			result = append(result, script.ID)
		}
	}

	return result
}

func isInArray(s string, arr []string) bool {
	for _, b := range arr {
		if b == s {
			return true
		}
	}
	return false
}
