//go:generate mockery -name Generator
package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

// SpecificationTestCases - test cases generated for a specification
type SpecificationTestCases struct {
	Specification discovery.ModelAPISpecification `json:"apiSpecification"`
	TestCases     []model.TestCase                `json:"testCases"`
}

type GeneratorConfig struct {
	ClientID              string
	Aud                   string
	ResponseType          string
	Scope                 string
	AuthorizationEndpoint string
	RedirectURL           string
}

// Generator - generates test cases from discovery model
type Generator interface {
	GenerateSpecificationTestCases(GeneratorConfig, discovery.ModelDiscovery) TestCasesRun
}

// NewGenerator - returns implementation of Generator interface
func NewGenerator() Generator {
	return generator{
		resolver: permissions.Resolver,
	}
}

// generator - implements Generator interface
type generator struct {
	resolver func(groups []permissions.Group) permissions.CodeSetResultSet
}

// GenerateSpecificationTestCases - generates test cases
func (g generator) GenerateSpecificationTestCases(config GeneratorConfig, discovery discovery.ModelDiscovery) TestCasesRun {
	specTestCases := []SpecificationTestCases{}
	customReplacements := make(map[string]string)
	originalEndpoints := make(map[string]string, 0)

	for _, customTest := range discovery.CustomTests { // assume ordering is prerun i.e. customtest run before other tests
		specTestCases = append(specTestCases, GetCustomTestCases(&customTest))
		for k, v := range customTest.Replacements {
			customReplacements[k] = v
		}
	}

	nameGenerator := names.NewSequentialPrefixedName("#t")
	for _, item := range discovery.DiscoveryItems {
		specTests, endpoints := generateSpecificationTestCases(item, nameGenerator, customReplacements)
		specTestCases = append(specTestCases, specTests)
		for k, v := range endpoints {
			originalEndpoints[k] = v
		}
	}

	tmpSpecTestCases := []SpecificationTestCases{}
	for _, specTest := range specTestCases {
		tmpSpecTestCases = append(tmpSpecTestCases, specTest)
		for x, y := range specTest.TestCases {
			y.Input.Endpoint = originalEndpoints[y.ID]
			specTest.TestCases[x] = y
		}
	}

	// calculate permission set required and update the header token in the test case request
	consentRequirements := g.consentRequirements(tmpSpecTestCases) // uses pre-modified swagger urls

	// generate PSU consent URL onto the perm set structure
	consentRequirements = withConsentUrl(config, consentRequirements)

	return TestCasesRun{specTestCases, consentRequirements}

}

// withConsentUrl copies the full requirement consent structure into a new one with the Consent url populated
func withConsentUrl(config GeneratorConfig, consentRequirements []model.SpecConsentRequirements) []model.SpecConsentRequirements {
	var withUrlSpecs []model.SpecConsentRequirements
	for _, spec := range consentRequirements {
		var namedPermsWithUrl model.NamedPermissions
		for _, namedPerm := range spec.NamedPermissions {
			claims := authentication.PSUConsentClaims{
				AuthorizationEndpoint: config.AuthorizationEndpoint,
				Aud:                   config.Aud,
				Iss:                   config.ClientID,
				ResponseType:          config.ResponseType,
				Scope:                 config.Scope,
				RedirectURI:           config.RedirectURL,
				ConsentId:             "",
				State:                 namedPerm.Name,
			}
			consentUrl, _ := authentication.PSUURLGenerate(claims)
			namedPermsWithUrl = append(
				namedPermsWithUrl,
				model.NamedPermission{
					Name:       namedPerm.Name,
					CodeSet:    namedPerm.CodeSet,
					ConsentUrl: consentUrl.String(),
				},
			)
		}
		specWithUrl := model.SpecConsentRequirements{
			Identifier:       spec.Identifier,
			NamedPermissions: namedPermsWithUrl,
		}
		withUrlSpecs = append(withUrlSpecs, specWithUrl)
	}
	return withUrlSpecs
}

// consentRequirements calls resolver to get list of permission sets required to run all test cases
func (g generator) consentRequirements(specTestCases []SpecificationTestCases) []model.SpecConsentRequirements {
	nameGenerator := names.NewSequentialPrefixedName("to")
	var specConsentRequirements []model.SpecConsentRequirements
	for _, spec := range specTestCases {
		var groups []permissions.Group
		for _, tc := range spec.TestCases {
			g := model.NewDefaultPermissionGroup(tc)
			groups = append(groups, g)
		}
		resultSet := g.resolver(groups)
		consentRequirements := model.NewSpecConsentRequirements(nameGenerator, resultSet, spec.Specification.Name)
		specConsentRequirements = append(specConsentRequirements, consentRequirements)
	}
	return specConsentRequirements
}

// TestCasesRun represents all specs and their test and a list of tokens
// required to run those tests
type TestCasesRun struct {
	TestCases               []SpecificationTestCases        `json:"specCases"`
	SpecConsentRequirements []model.SpecConsentRequirements `json:"specTokens"`
}

func generateSpecificationTestCases(item discovery.ModelDiscoveryItem, nameGenerator names.Generator, globalReplacements map[string]string) (SpecificationTestCases, map[string]string) {
	testcases, originalEndpoints := GetImplementedTestCases(&item, nameGenerator, globalReplacements)
	return SpecificationTestCases{Specification: item.APISpecification, TestCases: testcases}, originalEndpoints
}
