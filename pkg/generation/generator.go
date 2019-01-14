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
	results := make([]SpecificationTestCases, len(discovery.DiscoveryItems))

	for _, customTest := range discovery.CustomTests { // assume ordering is prerun ...  ...
		results = append(results, GetCustomTestCases(&customTest))
	}

	// Assumes testNo is used as the base for all testcase IDs - to keep testcase IDs unique
	testNo := 1000
	for _, item := range discovery.DiscoveryItems {
		results = append(results, generateSpecificationTestCases(item, testNo))
		testNo += 1000
	}
	return results
}

func generateSpecificationTestCases(item discovery.ModelDiscoveryItem, testNo int) SpecificationTestCases {
	testCases := GetImplementedTestCases(&item, testNo)
	return SpecificationTestCases{Specification: item.APISpecification, TestCases: testCases}
}
