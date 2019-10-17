package manifest

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

const manifestPath = "file://manifests/ob_3.1_payment_fca.json"

func TestPermx(t *testing.T) {
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}

	specType, err := GetSpecType(apiSpec.SchemaVersion)
	scripts, _, err := LoadGenerationResources(specType, manifestPath, nil)
	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &model.Context{},
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
	dumpJSON(requiredTokens)
}

func TestGetScriptConsentTokens(t *testing.T) {
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}
	specType, err := GetSpecType(apiSpec.SchemaVersion)
	scripts, _, err := LoadGenerationResources(specType, manifestPath, nil)
	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &model.Context{},
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)

	assert.Nil(t, err)

	assert.Nil(t, err)

	testcasePermissions, err := getTestCasePermissions(tests)
	assert.Nil(t, err)

	requiredTokens, err := getRequiredTokens(testcasePermissions)
	assert.Nil(t, err)

	populateTokens(t, requiredTokens)
	dumpJSON(requiredTokens)
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
	t.Error()
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

	t.Error()
}
