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

// GetBool get the bool value associated with key
func (c Context) GetBool(key string) (bool, error) {
	value, exist := c[key]
	if !exist {
		return false, ErrNotFound
	}

	valueBool, ok := value.(bool)
	if !ok {
		return false, errors.New("error casting key to bool")
	}

	return valueBool, nil
}

// PutContext - puts another context into this one
func (c Context) PutContext(ctx *Context) {
	for k, v := range *ctx {
		c.Put(k, v)
	}
}

// Put a value indexed by 'key' into the context. The value can be any type
func (c Context) Put(key string, value interface{}) {
	c[key] = value
}

// PutString Put a value indexed by 'key' into the context. The value can be any type
func (c Context) PutString(key string, value string) {
	var interfaceValue interface{} = value
	c[key] = interfaceValue
}

// PutMap of strings - into context
func (c Context) PutMap(mymap map[string]string) {
	for k, v := range mymap {
		c.PutString(k, v)
	}
}

// PutStringSlice puts a slice of strings into context
func (c Context) PutStringSlice(key string, values []string) {
	valuesCasted := []interface{}{}
	for _, value := range values {
		valuesCasted = append(valuesCasted, value)
	}
	c.Put(key, valuesCasted)
}

// GetStringSlice gets a slice of string from context
func (c Context) GetStringSlice(key string) ([]string, error) {
	result := []string{}
	_, ok := c[key]
	if !ok {
		return nil, ErrNotFound
	}

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

// Delete Key from Context
func (c *Context) Delete(delKey string) {
	delete(*c, delKey)
}

var dumpContexts bool

func EnableContextDumps() {
	dumpContexts = true
}

// DumpContext - send the contents of a context to a logger
func (c *Context) DumpContext(text ...string) {

	if !dumpContexts {
		return
	}

	if len(text) > 0 {
		logrus.StandardLogger().Trace("[Context] |=== " + text[0] + "===|")
	}

	if len(text) > 1 {
		for i := 1; i < len(text); i++ {
			key := text[i]
			value, _ := c.Get(key)
			logrus.StandardLogger().Tracef("[Context] %s : %v", key, value)
		}
	} else {
		for k, v := range *c {
			if k == "client_secret" || k == "basic_authentication" || k == "signingPublic" || k == "signingPrivate" { // skip potentially sensitive fields - likely need to be more robust
				continue
			}
			logrus.StandardLogger().Tracef("[Context] %s:%v", k, v)
		}
	}
}

// IsSet returns true if the key exists and is not set to zero value (nil or empty string)
func (c *Context) IsSet(key string) bool {
	val, exists := c.Get(key)
	if !exists {
		return false
	} else if val, ok := val.(string); ok && val == "" {
		return false
	}
	if val == nil {
		return false
	}
	return true
}

// GetStrings - given a list of strings, returns a map of the strings values from context
func (c *Context) GetStrings(text ...string) (map[string]string, error) {
	stringsMap := map[string]string{}
	for i := 0; i < len(text); i++ {
		key := text[i]
		value, err := c.GetString(key)
		if err != nil {
			return nil, errors.New("cannot get value of " + key + " from context")
		}
		if len(value) < 1 {
			return nil, errors.New("cannot get value of " + key + " from context is empty")
		}
		stringsMap[key] = value
	}
	return stringsMap, nil
}
