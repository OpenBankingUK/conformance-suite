package reporting

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMockedServiceReturnsPass(t *testing.T) {
	service := NewMockedService()

	results, err := service.Run([]generation.SpecificationTestCases{
		{
			TestCases: []model.TestCase{
				{},
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, results.Specifications[0].Pass)
	assert.True(t, results.Specifications[0].Tests[0].Pass)
}

func TestMockedServiceMapsAllTests(t *testing.T) {
	service := NewMockedService()

	results, err := service.Run([]generation.SpecificationTestCases{
		{
			TestCases: []model.TestCase{
				{}, {}, {},
			},
		},
		{
			TestCases: []model.TestCase{
				{},
			},
		},
	})

	require.NoError(t, err)
	assert.Len(t, results.Specifications, 2)
	assert.Len(t, results.Specifications[0].Tests, 3)
	assert.Len(t, results.Specifications[1].Tests, 1)
}

func TestMockedServiceMapsProperties(t *testing.T) {
	service := NewMockedService()

	results, err := service.Run([]generation.SpecificationTestCases{
		{
			Specification: discovery.ModelAPISpecification{
				Name:          "name",
				Version:       "version",
				SchemaVersion: "schemaversion",
				URL:           "url",
			},
			TestCases: []model.TestCase{
				{
					Name:  "testname",
					ID:    "testid",
					Input: model.Input{Endpoint: "testendpoint"},
				},
			},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "name", results.Specifications[0].Name)
	assert.Equal(t, "url", results.Specifications[0].URL)
	assert.Equal(t, "schemaversion", results.Specifications[0].SchemaVersion)
	assert.Equal(t, "version", results.Specifications[0].Version)
	assert.Equal(t, "testname", results.Specifications[0].Tests[0].Name)
	assert.Equal(t, "testid", results.Specifications[0].Tests[0].Id)
	assert.Equal(t, "testendpoint", results.Specifications[0].Tests[0].Endpoint)
}
