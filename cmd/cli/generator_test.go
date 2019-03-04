package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGeneratorInput = `{}`
)

func TestGenerator(t *testing.T) {
	journeyConfig := server.JourneyConfig{}
	journey := &server.MockJourney{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("SetConfig", journeyConfig).Return(nil)
	testCaseRun := generation.TestCasesRun{TestCases: []generation.SpecificationTestCases{}, SpecConsentRequirements: []model.SpecConsentRequirements{}}
	journey.On("TestCases").Return(testCaseRun, nil)
	g := newGenerator(journey)
	input := testGeneratorInput
	output := &bytes.Buffer{}

	err := g.Generate(journeyConfig, strings.NewReader(input), output)

	require.NoError(t, err)
	assert.JSONEq(t, `{"specCases": [], "specTokens": []}`, output.String())
	journey.AssertExpectations(t)
}

func TestGeneratorHandlesWrongInput(t *testing.T) {
	journey := &server.MockJourney{}
	g := newGenerator(journey)
	input := `hannah montana`
	output := &bytes.Buffer{}

	err := g.Generate(server.JourneyConfig{}, strings.NewReader(input), output)

	assert.EqualError(t, err, "error parsing discovery model json: invalid character 'h' looking for beginning of value")
}

func TestGeneratorHandlesSetDiscoveryModelErr(t *testing.T) {
	journeyConfig := server.JourneyConfig{}
	journey := &server.MockJourney{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, errors.New("booboo"))
	journey.On("SetConfig", journeyConfig).Return(nil)
	g := newGenerator(journey)
	input := testGeneratorInput
	output := &bytes.Buffer{}

	err := g.Generate(server.JourneyConfig{}, strings.NewReader(input), output)

	assert.EqualError(t, err, "error setting discovery model: booboo")
}

func TestGeneratorHandlesFailuresFromSetDiscovery(t *testing.T) {
	journeyConfig := server.JourneyConfig{}
	journey := &server.MockJourney{}
	failures := discovery.ValidationFailures{{Key: "key", Error: "something wrong with this world"}}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(failures, nil)
	journey.On("SetConfig", journeyConfig).Return(nil)
	g := newGenerator(journey)
	input := testGeneratorInput
	output := &bytes.Buffer{}

	err := g.Generate(server.JourneyConfig{}, strings.NewReader(input), output)

	assert.EqualError(t, err, "error validating discovery model\nkey: something wrong with this world\n")
}

func TestGeneratorHandlesErrFromTestCases(t *testing.T) {
	journeyConfig := server.JourneyConfig{}
	journey := &server.MockJourney{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("SetConfig", journeyConfig).Return(nil)
	testCaseRun := generation.TestCasesRun{TestCases: []generation.SpecificationTestCases{}, SpecConsentRequirements: []model.SpecConsentRequirements{}}
	journey.On("TestCases").Return(testCaseRun, errors.New("more booboo"))
	g := newGenerator(journey)
	input := testGeneratorInput
	output := &bytes.Buffer{}

	err := g.Generate(server.JourneyConfig{}, strings.NewReader(input), output)
	assert.EqualError(t, err, "error generating test cases: more booboo")
}

func TestGeneratorHandlesErrWriteToOutput(t *testing.T) {
	journeyConfig := server.JourneyConfig{}
	journey := &server.MockJourney{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("SetConfig", journeyConfig).Return(nil)
	testCaseRun := generation.TestCasesRun{TestCases: []generation.SpecificationTestCases{}, SpecConsentRequirements: []model.SpecConsentRequirements{}}
	journey.On("TestCases").Return(testCaseRun, nil)
	g := newGenerator(journey)
	input := strings.NewReader(`{}`)
	output := &brokenBuffer{}

	err := g.Generate(server.JourneyConfig{}, input, output)

	assert.EqualError(t, err, "error writing results to output: booboo")
}

type brokenBuffer bytes.Buffer

func (bb *brokenBuffer) Write(p []byte) (n int, err error) {
	return 0, errors.New("booboo")
}
