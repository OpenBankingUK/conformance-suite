package manifest

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Scripts -
type Scripts struct {
	Scripts []Script `json:"scripts,omitempty"`
}

// Script represents a high level test definition
type Script struct {
	Description       string            `json:"description,omitempty"`
	Detail            string            `json:"detail,omitempty"`
	ID                string            `json:"id,omitempty"`
	RefURI            string            `json:"refURI,omitempty"`
	Parameters        map[string]string `json:"parameters,omitempty"`
	Headers           map[string]string `json:"headers,omitempty"`
	Resource          string            `json:"resource,omitempty"`
	Asserts           []string          `json:"asserts,omitempty"`
	Method            string            `json:"method,omitempty"`
	URI               string            `json:"uri,omitempty"`
	URIImplementation string            `json:"uri_implementation,omitempty"`
	SchemaCheck       bool              `json:"schemaCheck,omitempty"`
}

// References - reference collection
type References struct {
	References map[string]Reference `json:"references,omitempty"`
}

// Reference is an item referred to by the test script list an assert of token requirement
type Reference struct {
	Expect      model.Expect `json:"expect,omitempty"`
	Permissions []string     `json:"permissions,omitempty"`
}

// AccountData stores account number to be used in the test scripts
type AccountData struct {
	Ais           map[string]string `json:"ais,omitempty"`
	AisConsentIds []string          `json:"ais.ConsentAccountId,omitempty"`
	Pis           PisData           `json:"pis,omitempty"`
}

// PisData contains information about PIS accounts required for the test scrips
type PisData struct {
	Currency        string            `json:"Currency,omitempty"`
	DebtorAccount   map[string]string `json:"DebtorAccount,omitempty"`
	MADebtorAccount map[string]string `json:"MADebtorAccount,omitempty"`
}

// LoadScripts loads the scripts from JSON encoded contents of filename
// and returns Scripts objects
func LoadScripts(filename string) (Scripts, error) {
	plan, err := ioutil.ReadFile(filename)
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
type DiscoveryPathsTestIDs map[string]map[string][]string

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
func MapDiscoveryEndpointsToManifestTestIDs(disco *discovery.Model, mf Scripts, mpParams map[string]string) DiscoveryPathsTestIDs {
	mapURLTests := make(DiscoveryPathsTestIDs)

	// Iterate the discoveryModel.discoveryItems.endpoints
	for _, v1 := range disco.DiscoveryModel.DiscoveryItems {
		for _, v2 := range v1.Endpoints {
			v2PathFixed := fixReferenceStructure(v2.Path)

			// For each discovery item, iterate all the `uri` fields and see if there is a match.
			for _, v3 := range mf.Scripts {
				v2Method := strings.ToUpper(v2.Method)
				v3Method := strings.ToUpper(v3.Method)

				// Need to prepend '/' as manifest URIs do not contain it.
				v3URIFixed := "/" + v3.URI
				if strings.ToLower(v2PathFixed) == strings.ToLower(v3URIFixed) &&
					v2Method == v3Method {

					v2PathFixed = replaceParams(v2PathFixed, mpParams)

					if _, ok := mapURLTests[v2PathFixed]; !ok {
						mapURLTests[v2PathFixed] = map[string][]string{}
					}

					if _, ok := mapURLTests[v2PathFixed]; !ok {
						mapURLTests[v2PathFixed][v3Method] = make([]string, 0)
					}
					mapURLTests[v2PathFixed][v3Method] = append(mapURLTests[v2PathFixed][v3Method], v3.ID)
				}
			}
		}
	}
	return mapURLTests
}

// FindUnmatchedManifestTests
// Find all the TestIDs from Manifest that have not been matched against an endpoint in the discovery model
func FindUnmatchedManifestTests(mf Scripts, mpParams DiscoveryPathsTestIDs) []string {
	var result []string
	for _, script := range mf.Scripts {

		uriFixed := strings.Replace(script.URI, "$accountId", "500000000000000000000004", -1)
		uriFixed = "/" + uriFixed

		if methods, ok := mpParams[uriFixed]; ok {
			for method, testIDs := range methods {
				if strings.ToLower(method) == strings.ToLower(script.Method) {
					if !isInArray(script.ID, testIDs) {
						result = append(result, script.ID)
					}
				}
			}
		} else {
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

func fixReferenceStructure(url string) string {
	// Need to manipulate URI placeholders in one file to match the other.
	// For example:
	// `/accounts/{AccountId}` in the discovery file
	// is not equal the
	// `accounts/$accountId` in manifest file.

	urlFixed := strings.Replace(url, "{", "$", -1)
	urlFixed = strings.Replace(urlFixed, "}", "", -1)

	return urlFixed
}

func replaceParams(url string, params map[string]string) string {
	var result string

	for key, val := range params {
		result = strings.Replace(strings.ToLower(url), strings.ToLower(key), val, -1)
	}
	return result
}
