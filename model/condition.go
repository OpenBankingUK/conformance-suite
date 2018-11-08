package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// ConditionEnum models endpoint conditionality based on:
// Account and Transaction API Specification - v3.0 - Section 4 Endpoints
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0#AccountandTransactionAPISpecification-v3.0-Endpoints
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

// Conditionaity - capture the conditionaliy of a method/endpoint
type Conditionaity struct {
	Condition ConditionEnum `json:"condition,omitempty"`
	Method    string        `json:"method,omitempty"`
	Endpoint  string        `json:"endpoint,omitempty"`
}

// EndpointConditionality - Store of endpoint conditionality
var endpointConditionality []Conditionaity

func init() {
	err := loadConditions()
	if err != nil {
		logrus.Error(err)
	}
}

// IsOptional - returns true if the method/endpoint mix is optional
func IsOptional(method, endpoint string) (bool, error) {
	condition, err := findCondition(method, endpoint)
	if err != nil {
		return false, err
	}
	if condition.Condition == Optional {
		return true, nil
	}
	return false, nil
}

// IsMandatory - returns true if the method/endpoint mix is mandatory
func IsMandatory(method, endpoint string) (bool, error) {
	condition, err := findCondition(method, endpoint)
	if err != nil {
		return false, err
	}
	if condition.Condition == Mandatory {
		return true, nil
	}
	return false, nil
}

// IsConditional - returns true if the method/endpoint mix is conditional
func IsConditional(method, endpoint string) (bool, error) {
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
func findCondition(method, endpoint string) (Conditionaity, error) {
	for _, cond := range endpointConditionality {
		if cond.Method == method && cond.Endpoint == endpoint {
			return cond, nil
		}
	}
	return Conditionaity{}, errors.New("method: " + method + " endpoint:" + endpoint + " not found in conditionality array")
}

// loadConditions - get Mandatory/Conditional/Optional data from json file
func loadConditions() error {
	rawjson, _ := ioutil.ReadFile("../config/conditionality.json") // lives here for now until we figure out somewhere better
	err := json.Unmarshal(rawjson, &endpointConditionality)
	if err != nil {
		return err
	}
	return nil
}
