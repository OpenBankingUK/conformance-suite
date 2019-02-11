package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

// SpecificationTestCases - test cases generated for a specification
type SpecificationTestCases struct {
	Specification discovery.ModelAPISpecification `json:"apiSpecification"`
	TestCases     []model.TestCase                `json:"testCases"`
}

// Generator - generates test cases from discovery model
type Generator interface {
	GenerateSpecificationTestCases(discovery discovery.ModelDiscovery) TestCasesRun
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
func (g generator) GenerateSpecificationTestCases(discovery discovery.ModelDiscovery) TestCasesRun {
	specTestCases := []SpecificationTestCases{}
	globalReplacements := make(map[string]string)

	for _, customTest := range discovery.CustomTests { // assume ordering is prerun i.e. customtest run before other tests
		specTestCases = append(specTestCases, GetCustomTestCases(&customTest))
		for k, v := range customTest.Replacements {
			globalReplacements[k] = v
		}
	}

	nameGenerator := names.NewSequentialPrefixedName("#t")
	for _, item := range discovery.DiscoveryItems {
		specTestCases = append(specTestCases, generateSpecificationTestCases(item, nameGenerator, globalReplacements))
	}

	// calculate permission set required and update the header token in the test case request
	consentRequirements := g.consentRequirements(specTestCases)
	// @Julian glue with `updateSpecsBearer(consentRequirements, specTestCases)`

	return TestCasesRun{specTestCases, consentRequirements}

}

// consentRequirements calls resolver to get list of permission sets required to run all test cases
func (g generator) consentRequirements(specTestCases []SpecificationTestCases) []model.SpecConsentRequirements {
	nameGenerator := names.NewSequentialPrefixedName("to")
	var specConsentRequirements []model.SpecConsentRequirements
	for _, spec := range specTestCases {
		var groups []permissions.Group
		for _, tc := range spec.TestCases {
			groups = append(groups, model.NewPermissionGroup(tc))
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

func generateSpecificationTestCases(item discovery.ModelDiscoveryItem, nameGenerator names.Generator, globalReplacements map[string]string) SpecificationTestCases {
	return SpecificationTestCases{Specification: item.APISpecification, TestCases: GetImplementedTestCases(&item, nameGenerator, globalReplacements)}
}
