package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"regexp"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
)

func testLoadDiscoveryModel(t *testing.T) *discovery.ModelDiscovery {
	t.Helper()
	template, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-generic.json")
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
	testCasesRun := generator.GenerateSpecificationTestCases(discovery)
	cases := testCasesRun.TestCases

	t.Run("returns slice of SpecificationTestCases, one per discovery item", func(t *testing.T) {
		require.NotNil(t, cases)
		assert.Equal(t, len(discovery.DiscoveryItems), len(cases))
	})

	t.Run("returns each SpecificationTestCases with a Specification matching discovery item", func(t *testing.T) {
		require.Equal(t, len(discovery.DiscoveryItems), len(cases))

		for i, specificationCases := range cases {
			if specificationCases.Specification.Name == "CustomTest-GetOzoneToken" {
				continue
			}
			expectedSpec := discovery.DiscoveryItems[i].APISpecification
			assert.Equal(t, expectedSpec, specificationCases.Specification)
		}
	})
}

func TestPermissionsSetsEmpty(t *testing.T) {
	generator := &generator{}
	specTestCases := []SpecificationTestCases{}

	results := generator.consentRequirements(specTestCases)

	assert.Len(t, results, 0)
}

func TestPermissionsShouldPassAllTestsToResolver(t *testing.T) {
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

	assert.Len(t, specTokens, 1)
	assert.Len(t, specTokens[0].NamedPermissions, 2)
	// token name expected in format to-{3 letters}-{sequence number}
	matchName := regexp.MustCompile(`to\w{4}`)
	assert.Regexp(t, matchName, specTokens[0].NamedPermissions[0].Name)
	assert.Regexp(t, matchName, specTokens[0].NamedPermissions[1].Name)
}
