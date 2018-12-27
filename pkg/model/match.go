package model

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// MatchType enumeration
type MatchType int

// MatchType enumeration - this will be required when we extend to more than just BodyJSONValue
const (
	UnknownMatchType MatchType = iota
	HeaderValue
	HeaderRegex
	HeaderRegexContext
	HeaderPresent
	BodyRegex
	BodyJSONPresent
	BodyJSONCount
	BodyJSONValue
	BodyJSONRegex
	BodyLength
	Authorisation
)

// Match defines various types of response payload pattern and field checking.
// Match forms the basis for response validation outside of basic swagger/openapi schema validation
// Match is also used as the basis for field extraction and replacement which enable parameter passing
// between tests via the context.
// Match encapsulates a conditional statement that must 'match' in order to succeed.
// Matches can -
// - match a specified header field value for exact match
// - match a specified header field value using a regular expression
// - check that a specified header field exists in the response
// - check that a response body matches a regular expression
// - check that a response body has a particular json field present using json matching
// - check that a response body has a specific number of specified json array fields
// - check that a response body has a specific value of a specified json field
// - check that a response body has a specific json field and that the specific json field matches a regular expression
// - check that a response body is a specified length
// - allow for replacment of endpoint text ... e.g. {AccountId}
// - Authoriation: allow for manipulation of Bearer tokens in http headers
// - Result: allow for capturing of match values for further processing - like putting into a context
type Match struct {
	MatchType       MatchType `json:"match_type,omitempty"`        // Type of Match we're doing
	Description     string    `json:"description,omitempty"`       // Description of the purpose of the match
	ContextName     string    `json:"name,omitempty"`              // Context variable name
	Header          string    `json:"header,omitempty"`            // Header value to examine
	HeaderPresent   string    `json:"header-present,omitempty"`    // Header existence check
	Regex           string    `json:"regex,omitempty"`             // Regular expression to be used
	JSON            string    `json:"json,omitempty"`              // Json expression to be used
	Value           string    `json:"value,omitempty"`             // Value to match against (string)
	Numeric         int64     `json:"numeric,omitempty"`           //Value to match against - numeric
	Count           int64     `json:"count,omitempty"`             // Cont for JSON array match purposes
	BodyLength      *int64    `json:"body-length,omitempty"`       // Body payload length for matching
	ReplaceEndpoint string    `json:"replaceInEndpoint,omitempty"` // allows substituion of resourceIds
	Authorisation   string    `json:"authorisation,omitempty"`     // allows capturing of bearer tokens
	Result          string    `json:"result,omitempty"`            // capturing match values
}

// ContextAccessor - Manages access to matches for Put and Get value operations on a context
type ContextAccessor struct {
	Context *Context `json:"-"`
	Matches []Match  `json:"matches,omitempty"`
}

// PutValues is used by the 'contextPut' directive and essentially collects a set of matches whose purpose is
// To select values to put in a context. All the matches in this section must have a name (of the target context
// variable), a description (so if things go wrong we can accurately report) and an operation which results in
// the a selection which is copied into the context variable
// Note: the initial interation of this will just implement the JSON pattern/field matcher
func (c *ContextAccessor) PutValues(tc *TestCase, ctx *Context) error {
	for _, m := range c.Matches {
		success := m.PutValue(tc, ctx)
		if !success {
			return fmt.Errorf("error ContextAccessor PutValues - failed Match [%s]", m.String())
		}
	}
	return nil
}

// GetValues - checks for match elements in the contextGet section
// For each valid element there need to be a ContextName which is the name of the
// variable in the context we're trying to retrieve
// Once we have retrieved a context variable, we need to know what to do with it.
// Current we can ReplaceEndpoint - so basically do a string replace on our testcase endpoint
// which for example allows us to replace {AccountId} with a real account id
func (c *ContextAccessor) GetValues(tc *TestCase, ctx *Context) error {
	for _, match := range c.Matches {
		if len(match.ContextName) > 0 {
			value := ctx.Get(match.ContextName)
			if value != nil {
				contextValue := value.(string)
				if len(contextValue) > 0 {
					if len(match.ReplaceEndpoint) > 0 {
						result := strings.Replace(tc.Input.Endpoint, match.ReplaceEndpoint, contextValue, 1)
						tc.Input.Endpoint = result
					}
				}
			}
		}
	}
	return nil
}

func (m *Match) String() string {
	var b strings.Builder
	if m.MatchType == UnknownMatchType {
		m.MatchType = m.GetType()
	}

	b.WriteString(`MatchType: ` + matchTypeString[m.MatchType])
	if m.Description != `` {
		b.WriteString(` Description: ` + m.Description)
	}
	if m.ContextName != `` {
		b.WriteString(` ContextName: ` + m.ContextName)
	}
	if m.Header != `` {
		b.WriteString(` Header: ` + m.Header)
	}
	if m.HeaderPresent != `` {
		b.WriteString(` HeaderPresent: ` + m.HeaderPresent)
	}
	if m.Regex != `` {
		b.WriteString(` Regex: ` + m.Regex)
	}

	m.string2(&b)

	return b.String()
}

func (m *Match) string2(b *strings.Builder) *strings.Builder {
	if m.JSON != `` {
		b.WriteString(` JSON: ` + m.JSON)
	}
	if m.Value != `` {
		b.WriteString(` Value: ` + m.Value)
	}
	if m.Numeric > 0 {
		b.WriteString(` Numeric: ` + strconv.FormatInt(m.Numeric, 10))
	}
	if m.BodyLength != nil {
		b.WriteString(` BodyLength: ` + strconv.FormatInt(*m.BodyLength, 10))
	}
	if m.ReplaceEndpoint != `` {
		b.WriteString(` ReplaceEndpoint: ` + m.ReplaceEndpoint)
	}
	if m.Authorisation != `` {
		b.WriteString(` Authorisation: ` + m.Authorisation)
	}
	if m.Result != `` {
		b.WriteString(` Result: ` + m.Result)
	}
	return b
}

// Check a match function - figures out which match type we have and
// calls the appropraite match checking function
func (m *Match) Check(tc *TestCase) (bool, error) {
	matchType := m.GetType()
	return matchFuncs[matchType](m, tc)
}

// GetValue the value from the json match along with a context variable to put it into
func (m *Match) GetValue(inputBuffer string) (interface{}, string) {
	result := gjson.Get(inputBuffer, m.JSON)
	return result.String(), m.ContextName
}

// PutValue puts the value from the json match along with a context variable to put it into
func (m *Match) PutValue(tc *TestCase, ctx *Context) bool {
	success := false
	var err error
	switch m.GetType() {
	case BodyJSONPresent:
		success, err = m.setContextFromBodyPresent(tc, ctx)
		if err != nil {
			// future reporting hook
			return false
		}
	case BodyJSONValue:
		success, err = m.setContextFromCheckBodyJSONValue(tc, ctx)
		if err != nil {
			return false
		}
	case Authorisation:
		if strings.EqualFold("bearer", m.Authorisation) {
			success, err := checkAuthorisation(m, tc)
			if err != nil {
				return false
			}
			ctx.Put(m.ContextName, m.Result)
			return success
		}
	case HeaderRegexContext:
		success, err := checkHeaderRegexContext(m, tc)
		fmt.Println("HeaderRegexContext entered")
		if err != nil {
			return false
		}
		if success {
			ctx.Put(m.ContextName, m.Result)
			return true
		}
	}

	return success
}

func (m *Match) setContextFromBodyPresent(tc *TestCase, ctx *Context) (bool, error) {
	success, err := checkBodyJSONPresent(m, tc)
	if err != nil {
		return false, err
	}

	if success {
		if len(m.ContextName) > 0 {
			ctx.Put(m.ContextName, m.Result)
			return true, nil
		}
	}
	return false, nil
}

func (m *Match) setContextFromCheckBodyJSONValue(tc *TestCase, ctx *Context) (bool, error) {
	success, err := checkBodyJSONValue(m, tc)
	if err != nil {
		return false, err
	}
	if success {
		if len(m.ContextName) > 0 {
			ctx.Put(m.ContextName, m.Result)
			return true, nil
		}
	}
	return false, nil
}

// GetType - returns the type of a match
func (m *Match) GetType() MatchType {

	if m.MatchType != UnknownMatchType { // only figure out match type if its the default
		return m.MatchType
	}

	if fieldsPresent(m.Authorisation) {
		m.MatchType = Authorisation
		return Authorisation
	}

	if fieldsPresent(m.Header) { // simply for cyclomatic tests!
		return m.getHeaderType()
	}

	if fieldsPresent(m.HeaderPresent) {
		m.MatchType = HeaderPresent
		return HeaderPresent
	}
	if fieldsPresent(m.JSON, m.Regex) {
		m.MatchType = BodyJSONRegex
		return BodyJSONRegex
	}

	if fieldsPresent(m.Regex) {
		m.MatchType = BodyRegex
		return BodyRegex
	}

	if fieldsPresent(m.JSON, m.Value) {
		m.MatchType = BodyJSONValue
		return BodyJSONValue
	}

	if fieldsPresent(m.JSON) {
		if m.Count > 0 {
			m.MatchType = BodyJSONCount
			return BodyJSONCount
		}
	}

	if fieldsPresent(m.JSON) {
		m.MatchType = BodyJSONPresent
		return BodyJSONPresent
	}

	if m.BodyLength != nil {
		m.MatchType = BodyLength
		return BodyLength
	}

	return UnknownMatchType
}

func (m *Match) getHeaderType() MatchType {
	if fieldsPresent(m.Header, m.Value) { // note: below ordering matters
		m.MatchType = HeaderValue
		return HeaderValue
	}

	if fieldsPresent(m.Header, m.Regex, m.ContextName) {
		m.MatchType = HeaderRegexContext
		return HeaderRegexContext
	}

	if fieldsPresent(m.Header, m.Regex) {
		m.MatchType = HeaderRegex
		return HeaderRegex
	}

	return UnknownMatchType
}

func fieldsPresent(str ...string) bool {
	result := true
	for _, v := range str {
		if len(v) == 0 {
			result = false
		}
	}
	return result
}

var matchFuncs = map[MatchType]func(*Match, *TestCase) (bool, error){
	UnknownMatchType:   defaultMatch,
	HeaderValue:        checkHeaderValue,
	HeaderRegexContext: checkHeaderRegexContext,
	HeaderRegex:        checkHeaderRegex,
	HeaderPresent:      checkHeaderPresent,
	BodyRegex:          checkBodyRegex,
	BodyJSONPresent:    checkBodyJSONPresent,
	BodyJSONCount:      checkBodyJSONCount,
	BodyJSONValue:      checkBodyJSONValue,
	BodyJSONRegex:      checkBodyJSONRegex,
	BodyLength:         checkBodyLength,
	Authorisation:      checkAuthorisation,
}

var matchTypeString = map[MatchType]string{
	UnknownMatchType:   "unknown",
	HeaderValue:        "HeaderValue",
	HeaderRegex:        "HeaderRegex",
	HeaderPresent:      "HeaderPresent",
	HeaderRegexContext: "HeaderRegexContext",
	BodyRegex:          "BodyRegex",
	BodyJSONPresent:    "BodyJSONPresent",
	BodyJSONCount:      "BodyJSONCount",
	BodyJSONValue:      "BodyJSONValue",
	BodyJSONRegex:      "BodyJSONRegex",
	BodyLength:         "BodyLength",
	Authorisation:      "Authorisation",
}

func defaultMatch(m *Match, tc *TestCase) (bool, error) {
	return false, errors.New("Unknown match type fails by default")
}

func checkHeaderValue(m *Match, tc *TestCase) (bool, error) {
	var success bool
	var actualHeader string
	for head := range tc.Header {
		success = strings.EqualFold(head, m.Header)
		if success {
			actualHeader = head
			break
		}
	}

	headerValue := tc.Header.Get(actualHeader)
	success = m.Value == headerValue
	if !success {
		return false, fmt.Errorf("Header Value Match Failed - expected (%s) got (%s)", m.Value, headerValue)
	}
	return success, nil
}

// Allows capturing of a regex subfield expression in a header
// For example with the following location header
// "Location:https://x.y.z/auth?code=12345&redir=https://redir"
// using the following match:
//{
//	"name": "xchange_code",
//	"description": "Get the xchange code from the location redirect",
//	"header": "Location",
//	"regex": "code=(.*)&?.*"
//  }
//
// Will extract the value of code "12345" and make it available in the match m.Result field
//
func checkHeaderRegexContext(m *Match, tc *TestCase) (bool, error) {
	var success bool
	var actualHeader string
	for head := range tc.Header {
		success = strings.EqualFold(head, m.Header)
		if success {
			actualHeader = head
			break
		}
	}
	headerValue := tc.Header.Get(actualHeader)
	regex := regexp.MustCompile(m.Regex)
	result := regex.FindStringSubmatch(headerValue)
	if len(result) < 2 {
		return false, fmt.Errorf("Header Regex Context Match Failed - regex (%s) failed to find anything on Header (%s) value (%s)", m.Regex, m.Header, headerValue)
	}
	m.Result = result[1]
	return success, nil
}

func checkHeaderRegex(m *Match, tc *TestCase) (bool, error) {
	var success bool
	var actualHeader string
	for head := range tc.Header {
		success = strings.EqualFold(head, m.Header)
		if success {
			actualHeader = head
			break
		}
	}

	headerValue := tc.Header.Get(actualHeader)
	regex := regexp.MustCompile(m.Regex)
	success = regex.MatchString(headerValue)

	if !success {
		return false, fmt.Errorf("Header Regex Match Failed - regex (%s) failed on Header (%s) Value (%s)", m.Regex, m.Header, m.Value)
	}
	return success, nil
}

func checkHeaderPresent(m *Match, tc *TestCase) (bool, error) {
	var success bool
	for head := range tc.Header {
		success = strings.EqualFold(head, m.HeaderPresent)
		if success {
			return success, nil
		}
	}
	return false, fmt.Errorf("Header Present Match Failed - expected Header (%s) got nothing", m.HeaderPresent)
}

func checkBodyRegex(m *Match, tc *TestCase) (bool, error) {
	regex := regexp.MustCompile(m.Regex)
	success := regex.MatchString(tc.Body)
	if !success {
		return false, fmt.Errorf("Body Regex Match Failed - regex (%s) failed on Body", m.Regex)
	}
	return success, nil
}

func checkBodyJSONValue(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	success := result.String() == m.Value
	if !success {
		return false, fmt.Errorf("JSON Match Failed - expected (%s) got (%s)", m.Value, result)
	}
	return success, nil
}

func checkBodyJSONPresent(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	success := result.Exists()
	if !success {
		return false, fmt.Errorf("JSON Field Match Failed - no field present for pattern (%s)", m.JSON)
	}
	m.Result = result.String()
	return success, nil
}

func checkBodyJSONCount(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	if result.Int() != m.Count {
		return false, fmt.Errorf("JSON Count Field Match Failed - found (%d) not (%d) occurances of pattern (%s)",
			result.Int(), m.Count, m.JSON)
	}
	return true, nil
}

func checkBodyJSONRegex(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	if len(result.String()) == 0 {
		return false, fmt.Errorf("JSON Regex Match Failed - no field present for pattern (%s)", m.JSON)
	}
	regex := regexp.MustCompile(m.Regex)
	success := regex.MatchString(result.String())
	if !success {
		return false, fmt.Errorf("JSON Regex Match Failed - selected field (%s) does not match regex (%s)",
			result.String(), m.Regex)
	}
	return success, nil
}

func checkBodyLength(m *Match, tc *TestCase) (bool, error) {
	success := len(tc.Body) == int(*m.BodyLength)
	if !success {
		return false, fmt.Errorf("Check Body Length - body length (%d) does not match expected length (%d)",
			len(tc.Body), *m.BodyLength)
	}
	return success, nil
}

func checkAuthorisation(m *Match, tc *TestCase) (bool, error) {
	var success bool
	var actualHeader string
	for head := range tc.Header {
		success = strings.EqualFold(head, "Authorisation") // uk spelling
		if success {
			actualHeader = head
			break
		}
		success = strings.EqualFold(head, "Authorization") // us spelling
		if success {
			actualHeader = head
			break
		}
	}

	headerValue := tc.Header.Get(actualHeader)
	if len(headerValue) == 0 {
		return false, fmt.Errorf("Authorisation Bearer Match Failed - no header value found")
	}
	success = m.Value == headerValue
	idx := strings.Index(headerValue, "Bearer ")
	if idx == -1 {
		idx = strings.Index(headerValue, "bearer ")
	}
	if idx == -1 {
		return false, fmt.Errorf("Authorisation Bearer Match value Failed - no header bearer value found")
	}
	m.Result = headerValue[idx+7:] // copy the token after 7 chars "Bearer "...
	return true, nil
}

func checkUnimplemented(m *Match, tc *TestCase) (bool, error) {
	return false, errors.New("Unimplemented match type fails by default")
}
