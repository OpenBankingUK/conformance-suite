package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// ConditionEnum models endpoint conditionality based on:
// Account and Transaction API Specification - v3.0 - Section 4 Endpoints
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0#AccountandTransactionAPISpecification-v3.0-Endpoints
// Also see "Categorisation of Implementation Requirements" section of the following document
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/641992418/Read+Write+Data+API+Specification+-+v3.0#Read/WriteDataAPISpecification-v3.0-CategorisationofImplementationRequirements
type ConditionEnum int

const (
	// Mandatory - required
	Mandatory ConditionEnum = iota
	// Conditional on a regulartory requirement
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
	IsMandatory(method, endpoint string, specification string) (bool, error)
}

// conditionalityChecker - implements ConditionalityChecker - for checking endpoint conditionality
type conditionalityChecker struct {
}

// IsPresent - returns true if the method/endpoint mix exists for given specification
func (checker conditionalityChecker) IsPresent(method, endpoint string, specification string) (bool, error) {
	optional, err := isOptional(method, endpoint)
	if err != nil {
		return false, nil
	}
	mandatory, err := isMandatory(method, endpoint)
	if err != nil {
		return false, nil
	}
	conditional, err := isConditional(method, endpoint)
	if err != nil {
		return false, nil
	}
	return (optional || mandatory || conditional), nil
}

// IsMandatory - returns true if the method/endpoint mix is mandatory in given specification
func (checker conditionalityChecker) IsMandatory(method, endpoint string, specification string) (bool, error) {
	flag, err := isMandatory(method, endpoint)
	return flag, err
}

// NewConditionalityChecker - returns implementation of ConditionalityChecker interface
// for checking endpoint conditionality
func NewConditionalityChecker() ConditionalityChecker {
	checker := conditionalityChecker{}
	return checker
}

// EndpointConditionality - Store of endpoint conditionality
var endpointConditionality []Conditionality

func init() {
	err := loadConditions()
	if err != nil {
		logrus.Error(err)
		os.Exit(1) // Abort if we can't read the config correctly
	}
}

// GetEndpointConditionality - get a clone of the internal variable `endpointConditionality`.
func GetEndpointConditionality() []Conditionality {
	clone := make([]Conditionality, len(endpointConditionality))
	copy(clone, endpointConditionality)
	return clone
}

// isOptional - returns true if the method/endpoint mix is optional
func isOptional(method, endpoint string) (bool, error) {
	condition, err := findCondition(method, endpoint)
	if err != nil {
		return false, err
	}
	if condition.Condition == Optional {
		return true, nil
	}
	return false, nil
}

// isMandatory - returns true if the method/endpoint mix is mandatory
func isMandatory(method, endpoint string) (bool, error) {
	condition, err := findCondition(method, endpoint)
	if err != nil {
		return false, err
	}
	if condition.Condition == Mandatory {
		return true, nil
	}
	return false, nil
}

// isConditional - returns true if the method/endpoint mix is conditional
func isConditional(method, endpoint string) (bool, error) {
	condition, err := findCondition(method, endpoint)
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
func GetConditionality(method, endpoint string) (ConditionEnum, error) {
	condition, err := findCondition(method, endpoint)
	if err != nil {
		return UndefinedCondition, err
	}
	return condition.Condition, nil
}

// findCondition - find a condition given the method and endpoint
// the condition can then be queried for optionality
func findCondition(method, endpoint string) (Conditionality, error) {
	for _, cond := range endpointConditionality {
		if cond.Method == method && cond.Endpoint == endpoint {
			return cond, nil
		}
	}
	return Conditionality{}, errors.New("method: " + method + " endpoint:" + endpoint + " not found in conditionality array")
}

// loadConditions - get Mandatory/Conditional/Optional data from json file
func loadConditions() error {
	loader := []conditionLoader{}

	rawjson, err := ioutil.ReadFile("../../pkg/model/conditionality.json") // lives here for now until we figure out somewhere better
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawjson, &loader); err != nil {
		return err
	}

	for _, loaded := range loader { // map struct conditionality into enum conditionality
		condition := Conditionality{}
		condition.Endpoint = loaded.Endpoint
		condition.Method = loaded.Method
		switch loaded.StringCondition {
		case "mandatory":
			condition.Condition = Mandatory
		case "conditional":
			condition.Condition = Conditional
		case "optional":
			condition.Condition = Optional
		default:
			logrus.WithFields(logrus.Fields{
				"Condition": loaded.StringCondition,
				"Method":    loaded.Method,
				"Endpoint":  loaded.Endpoint,
			}).Warn("Load Conditions - unknown condition")
		}
		endpointConditionality = append(endpointConditionality, condition)
	}
	return nil
}
