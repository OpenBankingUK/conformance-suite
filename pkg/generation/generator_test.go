package generation

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
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

func TestPermissionsSetsEmpty(t *testing.T) {
	assert := test.NewAssert(t)

	generator := &generator{}
	specTestCases := []SpecificationTestCases{}

	results := generator.consentRequirements(specTestCases)

	assert.Len(results, 0)
}

func TestPermissionsShouldPassAllTestsToResolver(t *testing.T) {
	assert := test.NewAssert(t)

	g := generator{
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

	specTokens := g.consentRequirements(specTestCases)

	assert.Len(specTokens, 1)
	assert.Len(specTokens[0].NamedPermissions, 2)
	matchName := regexp.MustCompile(`to\w{4}`)
	assert.Regexp(matchName, specTokens[0].NamedPermissions[0].Name)
	assert.Regexp(matchName, specTokens[0].NamedPermissions[1].Name)
}

func TestShouldIgnoreDiscoveryItem(t *testing.T) {
	require := test.NewAssert(t)

	supportedSchemaVersions := []string{
		// Account and Transaction API Specification
		"https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
		// "Confirmation of Funds API Specification"
		"https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/confirmation-funds-swagger.json",
	}
	for _, supportedSchemaVersion := range supportedSchemaVersions {
		require.False(shouldIgnoreDiscoveryItem(discovery.ModelAPISpecification{
			SchemaVersion: supportedSchemaVersion,
		}))
	}

	unsupportedSchemaVersions := []string{
		// Payment Initiation API
		"https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json",
	}
	for _, unsupportedSchemaVersion := range unsupportedSchemaVersions {
		require.True(shouldIgnoreDiscoveryItem(discovery.ModelAPISpecification{
			SchemaVersion: unsupportedSchemaVersion,
		}))
	}
}
