package model

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// Context is intended to handle two types of object and make them available to various parts of the suite including
// testcases. The first set are objects created as a result of the discovery phase, which capture discovery model
// information like endpoints and conditional implementation indicators. The other set of data is information passed
// between a sequence of test cases, for example AccountId - extracted from the output of one testcase (/Accounts) and fed in
// as part of the input of another testcase for example (/Accounts/{AccountId}/transactions}
type Context map[string]interface{}

// Get the key form the Context map - currently assumes value converts easily to a string!
func (c Context) Get(key string) (interface{}, bool) {
	value, exist := c[key]
	return value, exist
}

var ErrNotFound = errors.New("error key not found")

// GetString get the string value associated with key
func (c Context) GetString(key string) (string, error) {
	value, exist := c[key]
	if !exist {
		return "", ErrNotFound
	}

	valueStr, ok := value.(string)
	if !ok {
		return "", errors.New("error casting key to string")
	}

	return valueStr, nil
}

// Put a value indexed by 'key' into the context. The value can be any type
func (c Context) Put(key string, value interface{}) {
	c[key] = value
}

// Put a value indexed by 'key' into the context. The value can be any type
func (c Context) PutString(key string, value string) {
	var interfaceValue interface{}
	interfaceValue = value
	c[key] = interfaceValue
}

// PutStringSlice puts a slice of strings into context
func (c Context) PutStringSlice(key string, values []string) {
	var valuesCasted []interface{}
	for _, value := range values {
		valuesCasted = append(valuesCasted, value)
	}
	c.Put(key, valuesCasted)
}

// GetStringSlice gets a slice of string from context
func (c Context) GetStringSlice(key string) ([]string, error) {
	var result []string
	stringsSlice, ok := c[key].([]interface{})
	if !ok {
		return nil, errors.New("cast error can't get string slice from context")
	}

	for _, value := range stringsSlice {
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("element cast error can't get string slice from context")
		}
		result = append(result, valueString)
	}

	return result, nil
}

// DumpContext - send the contents of a context to a logger
func (c *Context) DumpContext() {
	for k, v := range *c {
		if k == "client_secret" || k == "basic_authentication" { // skip potentially sensitive fields - likely need to be more robust
			continue
		}
		logrus.Debugf("key %s:%v", k, v)
	}
}
