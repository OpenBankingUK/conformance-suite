//go:generate mockery -name Generator -inpkg
package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
)

// Generator generates testcases from a discovery model using io.Reader/io.Writer pair
type Generator interface {
	Generate(config server.JourneyConfig, input io.Reader, output io.Writer) error
}

type generator struct {
	journey server.Journey
}

func newGenerator(journey server.Journey) generator {
	return generator{journey: journey}
}

// generate test cases from a discovery model, steps
//		1. read input
//		2. unmarshal json input to model struct
//		3. sets discovery model in Journey
//		4. generate test cases
//		5. marshal results to json
//		6. write to output marshaled result
func (g generator) Generate(config server.JourneyConfig, input io.Reader, output io.Writer) error {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		return errors.Wrap(err, "error reading discovery model file")
	}

	discoveryModel := &discovery.Model{}
	err = json.Unmarshal(content, discoveryModel)
	if err != nil {
		return errors.Wrap(err, "error parsing discovery model json")
	}

	err = g.journey.SetConfig(config)
	if err != nil {
		return errors.Wrap(err, "error generating testcases on setConfig")
	}

	failures, err := g.journey.SetDiscoveryModel(discoveryModel)
	if err != nil {
		return errors.Wrap(err, "error setting discovery model")
	}

	if !failures.Empty() {
		var errMsg strings.Builder
		errMsg.WriteString("error validating discovery model\n")
		for _, failure := range failures {
			errMsg.WriteString(fmt.Sprintf("%s: %s\n", failure.Key, failure.Error))
		}
		return errors.New(errMsg.String())
	}

	testCases, err := g.journey.TestCases()
	if err != nil {
		return errors.Wrap(err, "error generating test cases")
	}

	response, err := json.MarshalIndent(testCases, "", "    ")
	if err != nil {
		return errors.Wrap(err, "error making json from test cases")
	}

	_, err = fmt.Fprint(output, string(response))
	if err != nil {
		return errors.Wrap(err, "error writing results to output")
	}

	return nil
}
