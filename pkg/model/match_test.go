package model

import (
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const simplejson = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

// create a context
// create a json match object
// create an http response
// send response to match object
// have match object parse response and extract json field into context parameter
func TestContextPutFromMatch(t *testing.T) {
	ctx := Context{}
	m := Match{JSON: "name.last", Description: "simple match test", ContextName: "LastName"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, m.PutValue(string(bodyBytes), &ctx))
	assert.Equal(t, ctx.Get(m.ContextName), "Prichard")
}

func TestContextGetFromContext(t *testing.T) {
	ctx := Context{}
	m := Match{JSON: "name.first", Description: "simple match test", ContextName: "FirstName"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	value, variable := m.GetValue(string(bodyBytes))
	assert.Equal(t, "Janet", value)
	assert.Equal(t, ctx.Get(m.ContextName), ctx.Get(variable))
}

// Create a testcase that defines the basic matchers
// json matcher
func TestJSONMatcher(t *testing.T) {
	match := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	success, result := match.Check(string(bodyBytes))
	assert.Equal(t, "Prichard", result)
	assert.True(t, success)
}
