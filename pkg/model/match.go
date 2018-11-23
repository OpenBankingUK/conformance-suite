package model

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

// ContextAccessor -
type ContextAccessor struct {
	Matches []Match `json:"matches,omitempty"`
}

// MatchType enumeration
type MatchType int

// MatchType enumeration
const (
	StatusCode MatchType = iota
	HeaderValue
	HeaderRegex
	HeaderPresent
	BodyRegex
	BodyJSONPresent
	BodyJSONCount
	BodyJSONValue
	BodyJSONRegex
	BodyLength
)

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
	MatchType    MatchType // Type of Match we're doing
	Description  string    `json:"description,omitempty"`   // Description of the purpose of the match
	ContextName  string    `json:"name,omitempty"`          // Context variable name
	Header       string    `json:"header,omitempty"`        // Header value to examine
	HeaderExists string    `json:"header_exists,omitempty"` // Header existence check
	Regex        string    `json:"regex,omitempty"`         // Regular expression to be used
	JSON         string    `json:"json,omitempty"`          // Json expression to be used
	Value        string    `json:"value,omitempty"`         // Value to match against (string)
	Numeric      int64     `json:"numeric,omitempty"`       //Value to match against - numeric
	Count        int64     `json:"count,omitempty"`         // Cont for JSON array match purposes
	Length       int64     `json:"length,omitempty"`        // Body payload length for matching
	Context      *Context  // allows easily getting/setting context variables as required
}

// Matcher captures the behaviour required to perform a match against a test case response using a number of different
// criteria - see /docs/matches.md for more information about the matching model
type Matcher interface {
	Match(*http.Response) (interface{}, error)
}

// PutValues is used by the 'contextPut' directive and essentially collects a set of matches whose purpose is
// To select values to put in a context. All the matches in this section must have a name (of the target context
// variable), a description (so if things go wrong we can accurately report) and an operation which results in
// the a selection which is copied into the context variable
// Note: the initial interation of this will just implement the JSON pattern/field matcher
func (c *ContextAccessor) PutValues(resp *http.Response) (string, interface{}, error) {
	for _, m := range c.Matches {
		// Figure out match type
		// Execute to move selected/data sample into

		// Must have name - target context variable
		// Must have description
		if len(m.ContextName) == 0 || len(m.Description) == 0 {
			return "", nil, errors.New("ContextName or Description is empty")
		}

		if len(m.JSON) == 0 {
			return "", nil, errors.New("JSON target is empty")
		}

		m.Check(resp) // maybe not this ... but ...
		// maybe m.PutValues - returns variable name, value, error code --- generically
	}
	return "", nil, nil
}

// Check a match function
func (m *Match) Check(resp *http.Response) (bool, string) {
	responseBody, _ := ioutil.ReadAll(resp.Body)
	result := gjson.Get(string(responseBody), m.JSON)
	return result.String() == m.Value, result.String()
}

// GetValue the value from the json match along with a context variable to put it into
func (m *Match) GetValue(resp *http.Response) (interface{}, string) {
	responseBody, _ := ioutil.ReadAll(resp.Body)
	result := gjson.Get(string(responseBody), m.JSON)
	return result.String(), m.ContextName
}

// PutValue puts the value from the json match along with a context variable to put it into
func (m *Match) PutValue(resp *http.Response) bool {
	responseBody, _ := ioutil.ReadAll(resp.Body)
	result := gjson.Get(string(responseBody), m.JSON)
	if m.Context != nil {
		if len(m.ContextName) > 0 {
			m.Context.Put(m.ContextName, result.String())
			return true
		}
	}
	return false
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

// func (c *ContextAccessor) PutValues() (string, interface{}, error) {

// }

// Maybe put errors in the context for collection by the rule and reporting back.
// An slice of errors
// Also need a testsequence name - which means a test case can appear in many test
// sequences under the same rule
// Maybe have a testcase error object
