package model

import (
	"encoding/json"
	"errors"
	"os"

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
}

// IsPresent - returns true if the method/endpoint mix exists for given specification
func (checker conditionalityChecker) IsPresent(method, endpoint string, specification string) (bool, error) {
	optional, err := isOptional(method, endpoint, specification)
	if err != nil {
		return false, nil
	}
	mandatory, err := isMandatory(method, endpoint, specification)
	if err != nil {
		return false, nil
	}
	conditional, err := isConditional(method, endpoint, specification)
	if err != nil {
		return false, nil
	}
	return optional || mandatory || conditional, nil
}

// IsOptional - returns true if the method/endpoint mix is optional
func (checker conditionalityChecker) IsOptional(method, endpoint string, specification string) (bool, error) {
	flag, err := isOptional(method, endpoint, specification)
	return flag, err
}

// IsMandatory - returns true if the method/endpoint mix is mandatory in given specification
func (checker conditionalityChecker) IsMandatory(method, endpoint string, specification string) (bool, error) {
	flag, err := isMandatory(method, endpoint, specification)
	return flag, err
}

func (checker conditionalityChecker) IsConditional(method, endpoint string, specification string) (bool, error) {
	flag, err := isConditional(method, endpoint, specification)
	return flag, err
}

// MissingMandatory - returns array of mandatory endpoint Inputs that are missing from given endpoints parameter
func (checker conditionalityChecker) MissingMandatory(endpoints []Input, specification string) ([]Input, error) {
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
	checker := conditionalityChecker{}
	return checker
}

// EndpointConditionality - Store of endpoint conditionality by specification key
var endpointConditionality map[string][]Conditionality

func init() {
	err := loadConditions()
	if err != nil {
		logrus.StandardLogger().Error(err)
		os.Exit(1) // Abort if we can't read the config correctly
	}
}

// GetEndpointConditionality - get a clone of `endpointConditionality` array for given specification identifier
func GetEndpointConditionality(specification string) []Conditionality {
	clone := make([]Conditionality, len(endpointConditionality[specification]))
	copy(clone, endpointConditionality[specification])
	return clone
}

// isOptional - returns true if the method/endpoint mix is optional
func isOptional(method, endpoint string, specification string) (bool, error) {
	condition, err := findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	if condition.Condition == Optional {
		return true, nil
	}
	return false, nil
}

// isMandatory - returns true if the method/endpoint mix is mandatory
func isMandatory(method, endpoint string, specification string) (bool, error) {
	condition, err := findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	if condition.Condition == Mandatory {
		return true, nil
	}
	return false, nil
}

// isConditional - returns true if the method/endpoint mix is conditional
func isConditional(method, endpoint string, specification string) (bool, error) {
	condition, err := findCondition(method, endpoint, specification)
	if err != nil {
		return false, err
	}
	if condition.Condition == Conditional {
		return true, nil
	}
	return false, nil
}

// GetConditionality - returns and indicator in the following for of the method/endpoint conditionality
// model.Mandatory - endpoint is Mandatory
// model.Conditional - endpoint is conditional
// model.Optional - endpoint is optional
// model.UndefineCondition - we don't recognise the endpoint
func GetConditionality(method, endpoint, specification string) (ConditionEnum, error) {
	condition, err := findCondition(method, endpoint, specification)
	if err != nil {
		return UndefinedCondition, err
	}
	return condition.Condition, nil
}

// findCondition - find a condition given the method and endpoint
// the condition can then be queried for optionality
func findCondition(method, endpoint string, specification string) (Conditionality, error) {
	for _, cond := range endpointConditionality[specification] {
		if cond.Method == method && cond.Endpoint == endpoint {
			return cond, nil
		}
	}
	return Conditionality{}, errors.New("method: " + method + " endpoint:" + endpoint + " not found in conditionality array")
}

// loadConditions - get Mandatory/Conditional/Optional data from json file
func loadConditions() error {
	var loader map[string][]conditionLoader
	if err := json.Unmarshal(conditionalityStaticData(), &loader); err != nil {
		return err
	}

	endpointConditionality = make(map[string][]Conditionality)

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
		endpointConditionality[specification] = list
	}

	return nil
}
