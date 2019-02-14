package discovery

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// CustomTest used to read and make sense of the custom test json
type CustomTest struct {
	ID           string            `json:"@id,omitempty"`                   // JSONLD ID Reference
	Name         string            `json:"name,omitempty"`                  // Name
	Description  string            `json:"description,omitempty"`           // Purpose of the testcase in simple words
	Replacements map[string]string `json:"replacementParameters,omitempty"` // replacement parameters
	Sequence     []model.TestCase  `json:"testSequence,omitempty"`          // TestCase to be run as part of this custom test
}

// SpecificationTestCases - test cases generated for a specification
type SpecificationTestCases struct {
	Specification ModelAPISpecification `json:"apiSpecification"`
	TestCases     []model.TestCase      `json:"testCases"`
}

func generateSpecificationTestCases(model *ModelDiscovery) []SpecificationTestCases {
	results := []SpecificationTestCases{}
	for _, customTest := range model.CustomTests { // assume ordering is prerun ...  ...
		results = append(results, getCustomTestCases(&customTest))
	}
	return results
}

func getCustomTestCases(discoReader *CustomTest) SpecificationTestCases {
	spec := ModelAPISpecification{Name: discoReader.Name}
	specTestCases := SpecificationTestCases{Specification: spec}
	testcases := []model.TestCase{}
	for _, testcase := range discoReader.Sequence {
		testcases = append(testcases, testcase)
	}
	specTestCases.TestCases = testcases
	return specTestCases
}
