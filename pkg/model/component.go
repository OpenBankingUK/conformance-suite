package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Component - a reusable test case building block
type Component struct {
	ID               string            `json:"@id,omitempty"`              // JSONLD ID Reference
	Name             string            `json:"name,omitempty"`             // Name
	Description      string            `json:"description,omitempty"`      // Purpose of the testcase in simple words
	Documentation    string            `json:"documentation,omitempty"`    // What input parameters do, what output parameters are
	InputParameters  map[string]string `json:"inputParameters,omitempty"`  // input parameters
	OutputParameters map[string]string `json:"outputParameters,omitempty"` // output parameters
	Tests            []TestCase        `json:"testcases,omitempty"`        // TestCase to be run as part of this custom test
	Execution        []string          `json:"execution,omitempty"`
	Components       []string          `json:"components,omitempty"`
}

// MakeComponent -
func MakeComponent(name string) Component {
	return Component{Name: name}
}

const (
	productionComponentDirectory = "components/"
	relativeComponentDirectory   = "../../components/"
	testComponentDirectory       = "../model/component/testdata/"
)

// LoadComponent - Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func LoadComponent(filename string) (Component, error) {
	var c Component

	// handle varying location of component directory depending on test/prod
	fileContents, err := ioutil.ReadFile(productionComponentDirectory + filename)
	if err != nil {
		fileContents, err = ioutil.ReadFile(relativeComponentDirectory + filename)
		if err != nil {
			fileContents, err = ioutil.ReadFile(testComponentDirectory + filename)
			if err != nil {
				fileContents, err = ioutil.ReadFile(filename)
				if err != nil {
					return c, err
				}
			}
		}
	}

	err = json.Unmarshal(fileContents, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// ProcessReplacementFields - performance context/testcase parameter substitution on
// component test cases before they are run
func (c *Component) ProcessReplacementFields(ctx *Context) {
	for _, testcase := range c.Tests {
		testcase.ProcessReplacementFields(ctx, true)
	}
}

// ValidateParameters - check that the components required input and output parameters are present in the supplied context
func (c *Component) ValidateParameters(ctx *Context) error {
	err := c.checkInputParamsInContext(ctx)
	if err != nil {
		return err
	}
	return c.checkOutputParamsInContext(ctx)
}

func (c *Component) checkInputParamsInContext(ctx *Context) error {
	for k := range c.InputParameters {
		_, exists := ctx.Get(k)
		if !exists {
			return fmt.Errorf("input parameter (%s) not present in context for component (%s)", k, c.Name)
		}
	}
	return nil
}

func (c *Component) checkOutputParamsInContext(ctx *Context) error {
	for k := range c.OutputParameters {
		_, exists := ctx.Get(k)
		if !exists {
			return fmt.Errorf("output parameter (%s) not present in context for component (%s)", k, c.Name)
		}
	}
	return nil
}

// GetTests - returns the tests that need to be run for this component
func (c *Component) GetTests() []TestCase {
	for key := range c.Tests {
		c.Tests[key].Validator = schema.NewNullValidator()
	}
	return c.Tests
}
