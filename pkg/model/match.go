package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

type ContextAccessor struct {
	Context *Context
	Matches []Match `json:"matches,omitempty"`
}

// // ContextPut is a directive used to PUT variables to the context from a test case definition
// type ContextPut struct {
// 	Context Context
// 	Matches []Match `json:"matches,omitempty"`
// }

// // ContextGet is a directive used to GET variables from a context and insert them into a test case
// // request payload
// type ContextGet struct {
// 	Context Context
// 	Matches []Match `json:"matches,omitempty"`
// }

// Match defines various types of request, and response payload pattern and field checking
// Match forms the basis for response validation outside of basic swagger/openapi schema validation
// Match is also used as the basic for field extraction and replacement which enable parameter passing
// between test cases via the context.Match// Match encapsulates a conditional statement that must 'match' in order to succeed.
// Matches should -
// - match using a specific JSON field and a value
// - match using a Regex expression
// - match a specific header field to a value
// - match using a Regex expression on a header field
type Match struct {
	Description  string `json:"description,omitempty"`   // Description of the purpose of the match
	ContextName  string `json:"name,omitempty"`          // Context variable name
	Header       string `json:"header,omitempty"`        // Header value to examine
	HeaderExists string `json:"header_exists,omitempty"` // Header existence check
	Regex        string `json:"regex,omitempty"`         // Regular expression to be used
	JSON         string `json:"json,omitempty"`          // Json expression to be used
	Value        string `json:"value,omitempty"`         // Value to match against (string)
	Numeric      int64  `json:"numeric,omitempty"`       //Value to match against - numeric
	Count        int64  `json:"count,omitempty"`         // Cont for JSON array match purposes
	Length       int64  `json:"length,omitempty"`        // Body payload length for matching
}

// Matcher captures the behaviour required to perform a match against a test case response using a number of different
// criteria - see /docs/matches.md for more information about the matching model
type Matcher interface {
	Execute(*http.Response) (interface{}, error)
}

type ContextPutter interface {
	PutValue(*http.Response) (interface{}, error)
}

type ContextGetter interface {
	GetValue(*http.Response) (interface{}, error)
}

// PutValues is used by the 'contextPut' directive and essentially collects a set of matches whose purpose is
// To select values to put in a context. All the matches in this section must have a name (of the target context
// variable), a description (so if things go wrong we can accurately report) and an operation which results in
// the a selection which is copied into the context variable
// Note: the initial interation of this will just implement the JSON pattern/field matcher
func (c *ContextAccessor) PutValues(resp *http.Response) error {
	for _, m := range c.Matches {
		// Figure out match type
		// Execute to move selected/data sample into

		// Must have name - target context variable
		// Must have description
		if len(m.ContextName) == 0 || len(m.Description) == 0 {
			return errors.New("ContextName or Description is empty")
		}

		if len(m.JSON) == 0 {
			return errors.New("JSON target is empty")
		}

		m.Execute(resp) // maybe not this ... but ...
		// maybe m.PutValues - returns variable name, value, error code --- generically
	}
	return nil
}

// Execute a match function
func (m *Match) Execute(resp *http.Response) (interface{}, error) {
	responseBody, _ := ioutil.ReadAll(resp.Body)
	result := gjson.Get(string(responseBody), m.JSON)
	fmt.Printf("Result: %s\n", result.String())
	return nil, nil

}

// GetValues -
func (c *ContextAccessor) GetValues() error {
	for _, match := range c.Matches {
		// Figure out match type
		// Execute to move selected/data sample into

		// Must have name
		// Must have description
		//
		// can have json - to start
		_ = match

	}
	return nil

}

func (c *ContextAccessor) PutValues() (string, interface{}, error) {

}

// Maybe put errors in the context for collection by the rule and reporting back.
// An slice of errors
// Also need a testsequence name - which means a test case can appear in many test
// sequences under the same rule
// Maybe have a testcase error object
