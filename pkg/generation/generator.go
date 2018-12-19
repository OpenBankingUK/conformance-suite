package generation

import (
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
	GenerateSpecificationTestCases() []SpecificationTestCases
}

// NewGenerator - returns implementation of Generator interface
func NewGenerator(discovery discovery.ModelDiscovery) Generator {
	return generator{discovery: discovery}
}

// generator - implements Generator interface
type generator struct {
	discovery discovery.ModelDiscovery
}

// GenerateSpecificationTestCases - generates test cases
func (g generator) GenerateSpecificationTestCases() []SpecificationTestCases {
	results := []SpecificationTestCases{}
	// Assumes testNo is used as the base for all testcase IDs - to keep testcase IDs unique
	testNo := 1000

	for _, item := range g.discovery.DiscoveryItems {
		result := generateSpecificationTestCases(item, testNo)
		results = append(results, result)
		testNo += 1000
	}
	return results
}

func generateSpecificationTestCases(item discovery.ModelDiscoveryItem, testNo int) SpecificationTestCases {
	testCases := GetImplementedTestCases(&item, false, testNo)
	return SpecificationTestCases{Specification: item.APISpecification, TestCases: testCases}
}
