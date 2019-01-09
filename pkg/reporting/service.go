package reporting

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/google/uuid"
)

// Service represents a test case reporting service
type Service interface {
	Run([]generation.SpecificationTestCases) (Result, error)
}

type mockedService struct{}

// NewMockedService creates a mocked reporting service that returns always pass for all tests
func NewMockedService() Service {
	return mockedService{}
}

func (s mockedService) Run(testCases []generation.SpecificationTestCases) (Result, error) {
	return Result{
		Id:             uuid.New(),
		Specifications: mapSpecificationsTestToResults(testCases),
	}, nil
}

func mapSpecificationsTestToResults(testCases []generation.SpecificationTestCases) []Specification {
	specificationsResult := make([]Specification, len(testCases))
	for key, specification := range testCases {
		specificationsResult[key] = mapSpecificationTestToResult(specification)
	}
	return specificationsResult
}

func mapSpecificationTestToResult(cases generation.SpecificationTestCases) Specification {
	return Specification{
		Name:          cases.Specification.Name,
		Version:       cases.Specification.Version,
		SchemaVersion: cases.Specification.SchemaVersion,
		URL:           cases.Specification.URL,
		Pass:          true,
		Tests:         mapTestCasesToResults(cases),
	}
}

func mapTestCasesToResults(specification generation.SpecificationTestCases) []Test {
	testResults := make([]Test, len(specification.TestCases))
	for keyTest, test := range specification.TestCases {
		testResults[keyTest] = mapTestCaseToResult(test)
	}
	return testResults
}

func mapTestCaseToResult(test model.TestCase) Test {
	return Test{
		Name:     test.Name,
		Id:       test.ID,
		Endpoint: test.Input.Endpoint,
		Pass:     true,
	}
}
