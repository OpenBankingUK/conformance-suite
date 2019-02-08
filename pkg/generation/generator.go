package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// SpecificationTestCases - test cases generated for a specification
type SpecificationTestCases struct {
	Specification discovery.ModelAPISpecification `json:"apiSpecification"`
	TestCases     []model.TestCase                `json:"testCases"`
}

// Generator - generates test cases from discovery model
type Generator interface {
	GenerateSpecificationTestCases(discovery discovery.ModelDiscovery) []SpecificationTestCases
}

// NewGenerator - returns implementation of Generator interface
func NewGenerator() Generator {
	return generator{}
}

// generator - implements Generator interface
type generator struct {
}

// GenerateSpecificationTestCases - generates test cases
func (g generator) GenerateSpecificationTestCases(discovery discovery.ModelDiscovery) []SpecificationTestCases {
	results := []SpecificationTestCases{}
	globalReplacements := make(map[string]string)

	for _, customTest := range discovery.CustomTests { // assume ordering is prerun i.e. customtest run before other tests
		results = append(results, GetCustomTestCases(&customTest))
		for k, v := range customTest.Replacements {
			globalReplacements[k] = v
		}
	}

	nameGenerator := names.NewSententialPrefixedName("#t")
	for _, item := range discovery.DiscoveryItems {
		results = append(results, generateSpecificationTestCases(item, nameGenerator, globalReplacements))
	}
	return results
}

func generateSpecificationTestCases(item discovery.ModelDiscoveryItem, nameGenerator names.Generator, gobalReplacements map[string]string) SpecificationTestCases {
	return SpecificationTestCases{Specification: item.APISpecification, TestCases: GetImplementedTestCases(&item, nameGenerator, gobalReplacements)}
}
