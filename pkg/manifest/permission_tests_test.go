package manifest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
)

func TestManifestScriptIDsAreUnique(t *testing.T) {
	// Pre-existing duplicates, left as-is so existing test ids are not changed. Any other duplicate fails.
	knownDuplicates := map[string]int{
		"OB-301-CBPII-000009": 2,
		"OB-400-CBPII-000009": 2,
	}

	files, err := filepath.Glob("../../manifests/*.json")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	checked := 0
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		require.NoError(t, err)

		var manifest struct {
			Scripts []struct {
				ID string `json:"id"`
			} `json:"scripts"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil || len(manifest.Scripts) == 0 {
			continue // not a script manifest (e.g. assertions.json, data.json)
		}
		checked++

		seen := map[string]int{}
		for _, script := range manifest.Scripts {
			assert.NotEmpty(t, script.ID, "%s contains a script without an id", file)
			seen[script.ID]++
		}
		for id, count := range seen {
			if count > 1 && count != knownDuplicates[id] {
				t.Errorf("%s contains test id %s %d times", file, id, count)
			}
		}
	}
	require.NotZero(t, checked, "no script manifests found")
}

// Tests that check a valid access token lacking the required permission is rejected must be
// mapped to a consent token, otherwise they are sent without an Authorization header.
func TestInsufficientPermissionTestsSendAuthorizationHeader(t *testing.T) {
	newTests := map[string][]string{
		"BAL-101501": {"ReadBalances"},
		"BEN-102001": {"ReadBeneficiariesBasic", "ReadBeneficiariesDetail"},
		"BEN-102101": {"ReadBeneficiariesBasic", "ReadBeneficiariesDetail"},
		"TRA-105301": {"ReadTransactionsBasic", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
		"TRA-105401": {"ReadTransactionsBasic", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
		"PAR-103109": {"ReadParty"},
		"PAR-103110": {"ReadParty"},
		"PAR-103111": {"ReadPartyPSU"},
	}
	// Legacy tests intentionally left unchanged: they have no permissions so receive no token.
	legacyTests := []string{"BAL-101500", "BEN-102000", "BEN-102100", "PAR-103103", "TRA-105300", "TRA-105400"}

	cases := []struct {
		version      string
		manifestPath string
		discovery    string
		apiVersion   string
	}{
		{"400", "file://manifests/ob_4.0_accounts_transactions_fca.json", "../discovery/templates/ob-v4.0.1-ozone.json", "accounts_v4.0.1"},
		{"301", "file://manifests/ob_3.1_accounts_transactions_fca.json", "../discovery/templates/ob-v3.1-ozone.json", "accounts_v3.1.1"},
	}

	for _, c := range cases {
		t.Run(c.version, func(t *testing.T) {
			data, err := ioutil.ReadFile(c.discovery)
			require.NoError(t, err)
			disco := &discovery.Model{}
			require.NoError(t, json.Unmarshal(data, disco))
			item := disco.DiscoveryModel.DiscoveryItems[0]
			item.APISpecification.Version = strings.TrimPrefix(c.apiVersion, "accounts_v")

			ctx := model.Context{"apiversions": []interface{}{c.apiVersion}, "consentedAccountId": "account-id"}
			params := GenerationParameters{
				Spec:         item.APISpecification,
				Baseurl:      "http://mybaseurl",
				Ctx:          &ctx,
				Endpoints:    allEndpointsFor(item.Endpoints),
				ManifestPath: c.manifestPath,
				Validator:    schema.NewNullValidator(),
			}
			tcs, _, err := GenerateTestCases(&params)
			require.NoError(t, err)

			tcps, err := getTestCasePermissions(tcs)
			require.NoError(t, err)
			rts, err := getRequiredTokens(tcps)
			require.NoError(t, err)
			MapTokensToTestCases(rts, tcs)
			t.Logf("%d consent tokens required", len(rts))

			byID := map[string]model.TestCase{}
			for _, tc := range tcs {
				byID[tc.ID] = tc
			}

			newTokens := map[string]bool{}
			for suffix, excluded := range newTests {
				id := "OB-" + c.version + "-" + suffix
				tc, ok := byID[id]
				require.True(t, ok, "%s not generated", id)

				auth := tc.Input.Headers["Authorization"]
				require.True(t, strings.HasPrefix(auth, "Bearer $account"), "%s has no Authorization header", id)
				tokenName := strings.TrimPrefix(auth, "Bearer $")
				newTokens[tokenName] = true

				groups := 0
				for _, rt := range rts {
					if !containsString(rt.IDs, id) {
						continue
					}
					groups++
					assert.Equal(t, tokenName, rt.Name, "%s mapped to unexpected token", id)
					assert.Contains(t, rt.Perms, "ReadAccountsBasic", "%s token must include ReadAccountsBasic", id)
					for _, x := range excluded {
						assert.NotContains(t, rt.Perms, x, "%s token must not grant %s", id, x)
					}
				}
				assert.Equal(t, 1, groups, "%s should be in exactly one token group", id)
			}
			// keep the extra PSU authorisation cost down: transactions tests need their own consent,
			// everything else should share a consent
			assert.LessOrEqual(t, len(newTokens), 2, "new tests spread across too many consents: %v", newTokens)

			for _, suffix := range legacyTests {
				id := "OB-" + c.version + "-" + suffix
				tc, ok := byID[id]
				require.True(t, ok, "%s not generated", id)
				assert.Empty(t, tc.Input.Headers["Authorization"], "%s is expected to remain unchanged", id)
			}
		})
	}
}

// allEndpointsFor adds the bulk endpoints exercised by the tests so none are filtered out by discovery.
func allEndpointsFor(endpoints []discovery.ModelEndpoint) []discovery.ModelEndpoint {
	required := []string{
		"/balances", "/beneficiaries", "/accounts/{AccountId}/beneficiaries", "/transactions",
		"/accounts/{AccountId}/transactions", "/party", "/accounts/{AccountId}/party", "/accounts/{AccountId}/parties",
	}
	present := map[string]bool{}
	for _, e := range endpoints {
		present[e.Path] = true
	}
	for _, path := range required {
		if !present[path] {
			endpoints = append(endpoints, discovery.ModelEndpoint{Method: "GET", Path: path})
		}
	}
	return endpoints
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
