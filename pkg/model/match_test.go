package model

import (
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
	m := Match{Context: &ctx, JSON: "name.last", Description: "simple match test", ContextName: "LastName"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	assert.True(t, m.PutValue(resp))
	assert.Equal(t, ctx.Get(m.ContextName), ctx.Get(m.ContextName))
}

// Create a testcase that defines the basic matchers
// json matcher
func TestJSONMatcher(t *testing.T) {
	match := Match{JSON: "name.last", Description: "simple match test", ContextName: "NameInContext", Value: "Prichard"}
	resp := pkgutils.CreateHTTPResponse(200, "OK", simplejson)
	success, result := match.Check(resp)
	assert.Equal(t, "Prichard", result)
	assert.True(t, success)
}

func TestJSONPutValueInContext(t *testing.T) {
	match := Match{JSON: "name.last", Description: "simple match test", ContextName: "Lastname"}
	_ = match
}

/*
	body := string(bodyBytes)
	for _, match := range t.Expect.Matches {
		jsonMatch := match.JSON // check if there is a JSON match to be satisifed
		if len(jsonMatch) > 0 {
			matched := gjson.Get(body, jsonMatch)
			if matched.String() != match.Value { // check the value of the JSON match - is equal to the 'value' parameter of the match section within the testcase 'Expects' area
				return false, fmt.Errorf("(%s):%s: Json Match: expected %s got %s", t.ID, t.Name, match.Value, matched.String())
			}
		}
	}

*/

// Make a bunch of match fragments
// json field
// header match
// figure out a mock http response with the appropriate data in
//
// Two areas:-
//
// test data
// match definition
// creating match data can be done in two ways:-
//  programatically
//  declaratively
//
// Specifiying test data
//   Table driven? - maybe specify match json ?
//
//
//
//
