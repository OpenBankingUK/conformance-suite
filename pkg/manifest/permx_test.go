package manifest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/schema"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"

	"github.com/OpenBankingUK/conformance-suite/pkg/model"
)

const manifestPath = "file://manifests/ob_3.1_payment_fca.json"

func TestPermx(t *testing.T) {
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}

	specType, err := GetSpecType(apiSpec.SchemaVersion)

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1")
	context := model.Context{"apiversions": values}

	scripts, _, err := LoadGenerationResources(specType, manifestPath, &context)

	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &context,
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)

	assert.Nil(t, err)

	testcasePermissions, err := getTestCasePermissions(tests)
	assert.Nil(t, err)

	_, err = getRequiredTokens(testcasePermissions)
	assert.Nil(t, err)
}

func TestGetScriptConsentTokens(t *testing.T) {
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1")
	context := model.Context{"apiversions": values}

	specType, err := GetSpecType(apiSpec.SchemaVersion)
	scripts, _, err := LoadGenerationResources(specType, manifestPath, &context)
	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &context,
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.Nil(t, err)

	testcasePermissions, err := getTestCasePermissions(tests)
	assert.Nil(t, err)

	requiredTokens, err := getRequiredTokens(testcasePermissions)
	assert.Nil(t, err)

	populateTokens(t, requiredTokens)
}

func populateTokens(t *testing.T, gatherer []RequiredTokens) error {
	t.Helper()

	t.Logf("%d entries found\n", len(gatherer))
	for k, tokenGatherer := range gatherer {
		if len(tokenGatherer.Perms) == 0 {
			continue
		}
		token, err := getToken(tokenGatherer.Perms)
		if err != nil {
			return err
		}
		tokenGatherer.Token = token
		gatherer[k] = tokenGatherer

	}
	return nil
}

func getToken(perms []string) (string, error) {
	// for headless - get the tokens for the permissions here
	return "abigfattoken", nil
}

func TestApi311To312(t *testing.T) {
	apiversions := []string{"payments_v3.1.2", "accounts_v3.1.0", "cbpii_v3.1.0"}

	accounts := "accounts"
	payments := "payments"
	cpbii := "cbpii"

	for _, v := range apiversions {
		api := strings.Split(v, "_v")
		if len(api) > 1 {
			if strings.Compare(accounts, api[0]) == 0 {
				fmt.Println("Accounts " + api[1])
			} else if strings.Compare(payments, api[0]) == 0 {
				fmt.Println("Payments " + api[1])
			} else if strings.Compare(cpbii, api[0]) == 0 {
				fmt.Println("CBPII " + api[1])
			}
		}
	}

	fmt.Println("\n-----Done")
}

func TestCompareApiVersions(t *testing.T) {
	apiversions := []string{"payments_v3.1.2", "accounts_v3.1.2", "accounts_v3.1.0", "cbpii_v3.1.0"}
	str1 := apiversions[0]
	str2 := apiversions[1]
	str3 := apiversions[2]
	api1 := strings.Split(str1, "_v")
	api2 := strings.Split(str2, "_v")
	api3 := strings.Split(str3, "_v")

	fmt.Printf("%s, %s, %s, %s\n", str1, str2, api1, api2)

	s1, _ := semver.Make(api1[1])
	s2, _ := semver.Make(api2[1])
	s3, _ := semver.Make(api3[1])

	fmt.Printf("compare %s,%s = %d\n", api1[0], api2[0], s1.Compare(s2))
	fmt.Printf("compare %s,%s = %d\n", api1[0], api3[0], s1.Compare(s3))

}

func TestUpdateTokensFromConsentLinksOnlyUnresolvedConsentIDReference(t *testing.T) {
	rts := []RequiredTokens{{Name: "vrpsToken0001", ConsentParam: "OB-400-VRP-100100-ConsentId"}}

	fresh := model.MakeTestCase()
	fresh.ID = "OB-400-VRP-100600"
	fresh.Context = model.Context{"consentId": "$OB-400-VRP-100100-ConsentId"}

	linked, err := updateTokensFromConsent(rts, []model.TestCase{fresh})
	assert.NoError(t, err)
	assert.Equal(t, []string{"OB-400-VRP-100600"}, linked[0].IDs)

	// a consentId resolved at generation time from a previous run's context cannot be linked,
	// which is why journey state must be reset before test cases are regenerated
	stale := model.MakeTestCase()
	stale.ID = "OB-400-VRP-100600"
	stale.Context = model.Context{"consentId": "VRP_SWP_stale"}

	unlinked, err := updateTokensFromConsent([]RequiredTokens{{Name: "vrpsToken0001", ConsentParam: "OB-400-VRP-100100-ConsentId"}}, []model.TestCase{stale})
	assert.NoError(t, err)
	assert.Empty(t, unlinked[0].IDs)
}

func TestResetConsentJobs(t *testing.T) {
	GetConsentJobs().Add(model.TestCase{ID: "OB-400-VRP-100100"})
	_, exists := GetConsentJobs().Get("OB-400-VRP-100100")
	assert.True(t, exists)

	ResetConsentJobs()

	_, exists = GetConsentJobs().Get("OB-400-VRP-100100")
	assert.False(t, exists)
}
