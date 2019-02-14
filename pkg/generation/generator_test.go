package generation

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"regexp"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/require"
)

func testLoadDiscoveryModel(t *testing.T) *discovery.ModelDiscovery {
	t.Helper()
	template, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-generic.json")
	require.NoError(t, err)
	require.NotNil(t, template)
	json := string(template)
	model, err := discovery.UnmarshalDiscoveryJSON(json)
	require.NoError(t, err)
	return &model.DiscoveryModel
}

func TestGenerateSpecificationTestCases(t *testing.T) {
	discovery := *testLoadDiscoveryModel(t)
	generator := NewGenerator()
	config := GeneratorConfig{}
	testCasesRun := generator.GenerateSpecificationTestCases(config, discovery)
	cases := testCasesRun.TestCases

	t.Run("returns slice of SpecificationTestCases, one per discovery item", func(t *testing.T) {
		assert := test.NewAssert(t)

		require.NotNil(t, cases)
		assert.Equal(len(discovery.DiscoveryItems), len(cases))
	})

	t.Run("returns each SpecificationTestCases with a Specification matching discovery item", func(t *testing.T) {
		assert := test.NewAssert(t)

		require.Equal(t, len(discovery.DiscoveryItems), len(cases))

		for i, specificationCases := range cases {
			if specificationCases.Specification.Name == "CustomTest-GetOzoneToken" {
				continue
			}
			expectedSpec := discovery.DiscoveryItems[i].APISpecification
			assert.Equal(expectedSpec, specificationCases.Specification)
		}
	})
}

func TestPermissionsSetsEmpty(t *testing.T) {
	assert := test.NewAssert(t)

	generator := &generator{}
	specTestCases := []SpecificationTestCases{}

	results := generator.consentRequirements(specTestCases)

	assert.Len(results, 0)
}

func TestPermissionsShouldPassAllTestsToResolver(t *testing.T) {
	assert := test.NewAssert(t)

	generator := generator{
		resolver: func(groups []permissions.Group) permissions.CodeSetResultSet {
			return permissions.CodeSetResultSet{
				{
					TestIds: []permissions.TestId{"1"},
				},
				{
					TestIds: []permissions.TestId{"2"},
				},
			}
		},
	}
	specTestCases := []SpecificationTestCases{
		{
			TestCases: []model.TestCase{
				{ID: "1"},
			},
		},
	}

	specTokens := generator.consentRequirements(specTestCases)

	assert.Len(specTokens, 1)
	assert.Len(specTokens[0].NamedPermissions, 2)
	matchName := regexp.MustCompile(`to\w{4}`)
	assert.Regexp(matchName, specTokens[0].NamedPermissions[0].Name)
	assert.Regexp(matchName, specTokens[0].NamedPermissions[1].Name)
}

func TestConsentRequirementsWithConsentUrlMapsEmpty(t *testing.T) {
	var consentRequirements []model.SpecConsentRequirements
	config := GeneratorConfig{}

	result := withConsentUrl(config, consentRequirements)

	var expected []model.SpecConsentRequirements
	assert.Equal(t, expected, result)
}

func TestConsentRequirementsWithConsentUrlMapsAllProperties(t *testing.T) {
	consentRequirements := []model.SpecConsentRequirements{
		{
			Identifier: "spec name",
			NamedPermissions: []model.NamedPermission{
				{
					Name:       "set1",
					CodeSet:    permissions.CodeSetResult{CodeSet: []permissions.Code{"a"}},
					ConsentUrl: "",
				},
			},
		},
	}
	config := GeneratorConfig{AuthorizationEndpoint: "/auth"}

	result := withConsentUrl(config, consentRequirements)

	expected := []model.SpecConsentRequirements{
		{
			Identifier: "spec name",
			NamedPermissions: []model.NamedPermission{
				{
					Name:       "set1",
					CodeSet:    permissions.CodeSetResult{CodeSet: []permissions.Code{"a"}},
					ConsentUrl: "/auth?client_id=&request=eyJhbGciOiJub25lIn0.eyJhdWQiOiIiLCJjbGFpbXMiOnsiaWRfdG9rZW4iOnsib3BlbmJhbmtpbmdfaW50ZW50X2lkIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIifX19LCJpc3MiOiIiLCJyZWRpcmVjdF91cmkiOiIiLCJzY29wZSI6IiJ9.&response_type=&scope=&state=set1",
				},
			},
		},
	}
	assert.Equal(t, expected, result)
}
