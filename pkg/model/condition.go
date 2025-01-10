package model

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
)

// ConditionEnum models endpoint conditionality based on:
// Account and Transaction API Specification - v3.1 - Section 4 Endpoints
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1#AccountandTransactionAPISpecification-v3.1-Endpoints
// Also see "Categorisation of Implementation Requirements" section of the following document
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937656404/Read+Write+Data+API+Specification+-+v3.1#Read/WriteDataAPISpecification-v3.1-CategorisationofImplementationRequirements
type ConditionEnum int

const (
	// Mandatory - required
	Mandatory ConditionEnum = iota
	// Conditional on a regulatory requirement
	Conditional
	// Optional at the implementors discretion
	Optional
	// UndefinedCondition -
	UndefinedCondition
)

// Conditionality - capture the conditionality of a method/endpoint
type Conditionality struct {
	Condition ConditionEnum `json:"condition,omitempty"`
	Method    string        `json:"method,omitempty"`
	Endpoint  string        `json:"endpoint,omitempty"`
}

// helper struct to load entries with string conditionalities
type conditionLoader struct {
	StringCondition string `json:"condition,omitempty"`
	Method          string `json:"method,omitempty"`
	Endpoint        string `json:"endpoint,omitempty"`
}

// ConditionalityChecker - interface to provide loose coupling
// between endpoint conditionality checks and invoking code
type ConditionalityChecker interface {
	IsPresent(method, endpoint string, specification string) (bool, error)
	IsOptional(method, endpoint string, specification string) (bool, error)
	IsMandatory(method, endpoint string, specification string) (bool, error)
	IsConditional(method, endpoint string, specification string) (bool, error)
	MissingMandatory(endpoints []Input, specification string) ([]Input, error)
}

// conditionalityChecker - implements ConditionalityChecker - for checking endpoint conditionality
type conditionalityChecker struct {
	data map[string][]Conditionality
}

// IsPresent - returns true if the method/endpoint mix exists for given specification
func (c *conditionalityChecker) IsPresent(method, endpoint string, specification string) (bool, error) {
	_, err := c.findCondition(method, endpoint, specification)
	return err == nil, nil
}

// IsOptional - returns true if the method/endpoint mix is optional
func (c *conditionalityChecker) IsOptional(method, endpoint string, specification string) (bool, error) {
	condition, err := c.findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	return condition.Condition == Optional, nil
}

// IsMandatory - returns true if the method/endpoint mix is mandatory in given specification
func (c *conditionalityChecker) IsMandatory(method, endpoint string, specification string) (bool, error) {
	condition, err := c.findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	return condition.Condition == Mandatory, nil
}

// IsConditional - returns true if the method/endpoint mix is conditional
func (c *conditionalityChecker) IsConditional(method, endpoint string, specification string) (bool, error) {
	condition, err := c.findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	return condition.Condition == Conditional, nil
}

// MissingMandatory - returns array of mandatory endpoint Inputs that are missing from given endpoints parameter
func (c *conditionalityChecker) MissingMandatory(endpoints []Input, specification string) ([]Input, error) {
	missingMandatoryEndpoints := []Input{}

	for _, condition := range GetEndpointConditionality(specification) {
		if condition.Condition == Mandatory {
			mandatory := condition
			isPresent := false
			for _, endpoint := range endpoints {
				isPresent = endpoint.Method == condition.Method && endpoint.Endpoint == condition.Endpoint
				if isPresent {
					break
				}
			}
			if !isPresent {
				missing := Input{Endpoint: mandatory.Endpoint, Method: mandatory.Method}
				missingMandatoryEndpoints = append(missingMandatoryEndpoints, missing)
			}
		}
	}

	return missingMandatoryEndpoints, nil
}

// NewConditionalityChecker - returns implementation of ConditionalityChecker interface
// for checking endpoint conditionality
func NewConditionalityChecker() ConditionalityChecker {
	checker := &conditionalityChecker{
		data: make(map[string][]Conditionality),
	}
	
	if err := checker.loadConditions(); err != nil {
		logrus.Fatal(err)
	}
	
	return checker
}

// GetEndpointConditionality - get a clone of conditions array for given specification identifier
func GetEndpointConditionality(specification string) []Conditionality {
	checker := NewConditionalityChecker()
	data := checker.(*conditionalityChecker).data
	clone := make([]Conditionality, len(data[specification]))
	copy(clone, data[specification])
	return clone
}

// GetConditionality - returns an indicator of the method/endpoint conditionality
// model.Mandatory - endpoint is Mandatory
// model.Conditional - endpoint is conditional
// model.Optional - endpoint is optional
// model.UndefineCondition - we don't recognise the endpoint
func GetConditionality(method, endpoint, specification string) (ConditionEnum, error) {
	checker := NewConditionalityChecker()
	condition, err := checker.(*conditionalityChecker).findCondition(method, endpoint, specification)
	if err != nil {
		return UndefinedCondition, err
	}
	return condition.Condition, nil
}

// findCondition - find a condition given the method and endpoint
// the condition can then be queried for optionality
func (c *conditionalityChecker) findCondition(method, endpoint string, specification string) (Conditionality, error) {
	for _, cond := range c.data[specification] {
		if cond.Method == method && cond.Endpoint == endpoint {
			return cond, nil
		}
	}
	return Conditionality{}, errors.New("method: " + method + " endpoint:" + endpoint + " not found in conditionality array")
}

// loadConditions - get Mandatory/Conditional/Optional data from json file
func (c *conditionalityChecker) loadConditions() error {
	var loader map[string][]conditionLoader
	if err := json.Unmarshal(conditionalityStaticData(), &loader); err != nil {
		return err
	}

	for specification, items := range loader {
		var list []Conditionality
		for _, item := range items {
			condition := Conditionality{}
			condition.Endpoint = item.Endpoint
			condition.Method = item.Method
			switch item.StringCondition {
			case "mandatory":
				condition.Condition = Mandatory
			case "conditional":
				condition.Condition = Conditional
			case "optional":
				condition.Condition = Optional
			default:
				logrus.StandardLogger().WithFields(logrus.Fields{
					"Condition": item.StringCondition,
					"Method":    item.Method,
					"Endpoint":  item.Endpoint,
				}).Warn("Load Conditions - unknown condition")
			}
			list = append(list, condition)
		}
		c.data[specification] = list
	}

	return nil
}