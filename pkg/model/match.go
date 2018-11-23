package model

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// ContextAccessor -
type ContextAccessor struct {
	Context *Context
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
	MatchType       MatchType `json:"match_type,omitempty"`        // Type of Match we're doing
	Description     string    `json:"description,omitempty"`       // Description of the purpose of the match
	ContextName     string    `json:"name,omitempty"`              // Context variable name
	Header          string    `json:"header,omitempty"`            // Header value to examine
	HeaderExists    string    `json:"header_exists,omitempty"`     // Header existence check
	Regex           string    `json:"regex,omitempty"`             // Regular expression to be used
	JSON            string    `json:"json,omitempty"`              // Json expression to be used
	Value           string    `json:"value,omitempty"`             // Value to match against (string)
	Numeric         int64     `json:"numeric,omitempty"`           //Value to match against - numeric
	Count           int64     `json:"count,omitempty"`             // Cont for JSON array match purposes
	Length          int64     `json:"length,omitempty"`            // Body payload length for matching
	ReplaceEndpoint string    `json:"replaceInEndpoint,omitempty"` // allows substituion of resourceIds
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
func (c *ContextAccessor) PutValues(tc *TestCase, ctx *Context) (string, error) {
	for _, m := range c.Matches {
		success := m.PutValue(tc.Body, ctx)
		if !success {
			return m.ContextName, errors.New("ContextPut variable check failed")
		}
	}
	return "", nil
}

// GetValues -
func (c *ContextAccessor) GetValues(tc *TestCase, ctx *Context) error {
	for _, match := range c.Matches {
		if len(match.ContextName) > 0 {
			value := ctx.Get(match.ContextName)
			if value != nil {
				contextValue := value.(string)
				if len(contextValue) > 0 {
					if len(match.ReplaceEndpoint) > 0 {
						result := strings.Replace(tc.Input.Endpoint, match.ReplaceEndpoint, contextValue, 1)
						fmt.Println("Old endpoint", tc.Input.Endpoint)
						fmt.Println("New endpoint", result)
						tc.Input.Endpoint = result
					}
				}
			}
		}
	}
	return nil

}

// Check a match function
func (m *Match) Check(inputBuffer string) (bool, string) {
	result := gjson.Get(inputBuffer, m.JSON)
	return result.String() == m.Value, result.String()
}

// GetValue the value from the json match along with a context variable to put it into
func (m *Match) GetValue(inputBuffer string) (interface{}, string) {
	result := gjson.Get(inputBuffer, m.JSON)
	return result.String(), m.ContextName
}

// PutValue puts the value from the json match along with a context variable to put it into
func (m *Match) PutValue(inputBuffer string, ctx *Context) bool {
	result := gjson.Get(inputBuffer, m.JSON)
	if len(m.ContextName) > 0 {
		ctx.Put(m.ContextName, result.String())
		return true
	}
	return false
}
