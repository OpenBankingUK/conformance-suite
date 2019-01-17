package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/mocks"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	journey := &mocks.Journey{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("TestCases").Return([]generation.SpecificationTestCases{}, nil)
	g := newGenerator(journey)

	input := `{}`
	output := &bytes.Buffer{}

	err := g.Generate(strings.NewReader(input), output)
	require.NoError(t, err)
	assert.Equal(t, "[]", output.String())
	journey.AssertExpectations(t)
}

func TestGeneratorHandlesWrongInput(t *testing.T) {
	journey := &mocks.Journey{}
	g := newGenerator(journey)

	input := `hannah montana`
	output := &bytes.Buffer{}

	err := g.Generate(strings.NewReader(input), output)
	assert.EqualError(t, err, "error parsing discovery model json: invalid character 'h' looking for beginning of value")
}

func TestGeneratorHandlesSetDiscoveryModelErr(t *testing.T) {
	journey := &mocks.Journey{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, errors.New("booboo"))
	g := newGenerator(journey)

	input := `{}`
	output := &bytes.Buffer{}

	err := g.Generate(strings.NewReader(input), output)
	assert.EqualError(t, err, "error setting discovery model: booboo")
}

func TestGeneratorHandlesFailuresFromSetDiscovery(t *testing.T) {
	journey := &mocks.Journey{}
	failures := discovery.ValidationFailures{{Key: "key", Error: "something wrong with this world"}}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(failures, nil)
	g := newGenerator(journey)

	input := `{}`
	output := &bytes.Buffer{}

	err := g.Generate(strings.NewReader(input), output)
	assert.EqualError(t, err, "error validating discovery model\nkey: something wrong with this world\n")
}

func TestGeneratorHandlesErrFromTestCases(t *testing.T) {
	journey := &mocks.Journey{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("TestCases").Return([]generation.SpecificationTestCases{}, errors.New("more booboo"))
	g := newGenerator(journey)

	input := `{}`
	output := &bytes.Buffer{}

	err := g.Generate(strings.NewReader(input), output)
	assert.EqualError(t, err, "error generating test cases: more booboo")
}

func TestGeneratorHandlesErrWriteToOutput(t *testing.T) {
	journey := &mocks.Journey{}
	journey.On("SetDiscoveryModel", &discovery.Model{}).Return(discovery.NoValidationFailures, nil)
	journey.On("TestCases").Return([]generation.SpecificationTestCases{}, nil)
	g := newGenerator(journey)

	input := strings.NewReader(`{}`)
	output := &brokenBuffer{}

	err := g.Generate(input, output)
	assert.EqualError(t, err, "error writing results to output: booboo")
}

type brokenBuffer bytes.Buffer

func (bb *brokenBuffer) Write(p []byte) (n int, err error) {
	return 0, errors.New("booboo")
}
