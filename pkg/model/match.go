package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
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
	CustomCheck
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
// - allow for replacement of endpoint text ... e.g. {AccountId}
// - Authorization: allow for manipulation of Bearer tokens in http headers
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
	Numeric         int64     `json:"numeric,omitempty"`           // Value to match against - numeric
	Count           int64     `json:"count,omitempty"`             // Cont for JSON array match purposes
	BodyLength      *int64    `json:"body-length,omitempty"`       // Body payload length for matching
	ReplaceEndpoint string    `json:"replaceInEndpoint,omitempty"` // allows substitution of resourceIds
	Authorisation   string    `json:"authorisation,omitempty"`     // allows capturing of bearer tokens
	Result          string    `json:"result,omitempty"`            // capturing match values
	Custom          string    `json:"custom,omitempty"`            // specifies custom matching routine
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
			return c.AppErr(fmt.Sprintf("error ContextAccessor PutValues - failed Match [%s]", m.String()))
		}
	}
	return nil
}

// AppMsg - application level trace
func (c *ContextAccessor) AppMsg(msg string) string {
	tracer.AppMsg("ContextAccessor", msg, "")
	return msg
}

// AppErr - application level trace error msg
func (c *ContextAccessor) AppErr(msg string) error {
	tracer.AppErr("ContextAccessor", msg, "")
	return errors.New(msg)
}

func (m *Match) String() string {
	if m.MatchType == UnknownMatchType {
		m.MatchType = m.GetType()
	}

	btes, err := json.Marshal(m)
	if err != nil {
		return m.AppMsg(fmt.Sprintf("error converting %s %s %s", matchTypeString[m.MatchType], m.Description, err.Error()))
	}
	return string(btes)
}

// Check a match function - figures out which match type we have and
// calls the appropriate match checking function
func (m *Match) Check(tc *TestCase) (bool, error) {
	matchType := m.GetType()
	return matchFuncs[matchType](m, tc)
}

// PutValue puts the value from the json match along with a context variable to put it into
func (m *Match) PutValue(tc *TestCase, ctx *Context) bool {
	switch m.GetType() {
	case BodyJSONPresent:
		success, err := m.setContextFromBodyPresent(tc, ctx)
		if err != nil {
			m.AppMsg(err.Error())
			return false
		}
		return success
	case BodyJSONValue:
		success, err := m.setContextFromCheckBodyJSONValue(tc, ctx)
		if err != nil {
			return false
		}
		return success
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
		if err != nil {
			return false
		}
		if success {
			ctx.Put(m.ContextName, m.Result)
			return true
		}
	case BodyRegex:
		return handleBodyRegex(tc, m, ctx)
	}
	return false
}

func handleBodyRegex(tc *TestCase, m *Match, ctx *Context) bool {
	success, err := checkBodyRegex(m, tc)
	if err != nil {
		return false
	}
	if success {
		if len(m.ContextName) > 0 {
			ctx.Put(m.ContextName, m.Result)
			return true
		}
	}
	return true
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

	if fieldsPresent(m.Custom) {
		m.MatchType = CustomCheck
		return CustomCheck
	}

	if fieldsPresent(m.Authorisation) {
		m.MatchType = Authorisation
		return Authorisation
	}

	if fieldsPresent(m.Header) {
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

// AppMsg - application level trace
func (m *Match) AppMsg(msg string) string {
	tracer.AppMsg("Match", msg, m.String())
	return msg
}

// AppErr - application level trace error msg
func (m *Match) AppErr(msg string) error {
	tracer.AppErr("Match", msg, m.String())
	return errors.New(msg)
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
	CustomCheck:        checkCustom,
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
	CustomCheck:        "Custom",
}

func defaultMatch(m *Match, _ *TestCase) (bool, error) {
	return false, m.AppErr("Unknown match type fails by default")
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
		return false, m.AppErr(fmt.Sprintf("Header Value Match Failed - expected (%s) got (%s)", m.Value, headerValue))
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
	regex, err := regexp.Compile(m.Regex)
	if err != nil {
		return false, err
	}
	result := regex.FindStringSubmatch(headerValue)
	if len(result) < 2 {
		return false, m.AppErr(fmt.Sprintf("Header Regex Context Match Failed - regex (%s) failed to find anything on Header (%s) value (%s)", m.Regex, m.Header, headerValue))
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
	regex, err := regexp.Compile(m.Regex)
	if err != nil {
		return false, err
	}

	success = regex.MatchString(headerValue)

	if !success {
		return false, m.AppErr(fmt.Sprintf("Header Regex Match Failed - regex (%s) failed on Header (%s) Value (%s)", m.Regex, m.Header, m.Value))
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
	return false, m.AppErr(fmt.Sprintf("Header Present Match Failed - expected Header (%s) got nothing", m.HeaderPresent))
}

func checkBodyRegex(m *Match, tc *TestCase) (bool, error) {
	regex, err := regexp.Compile(m.Regex)
	if err != nil {
		return false, err
	}
	success := regex.MatchString(tc.Body)
	if !success {
		return false, m.AppErr(fmt.Sprintf("Body Regex Match Failed - regex (%s) failed on Body", m.Regex))
	}
	if len(m.ContextName) > 0 {
		regexMatch := regex.FindStringSubmatch(tc.Body)
		if len(regexMatch) > 0 {
			m.Result = regexMatch[0]
		}
	}
	return success, nil
}

func checkBodyJSONValue(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	success := result.String() == m.Value
	if !success {
		return false, m.AppErr(fmt.Sprintf("JSON Match Failed - expected (%s) got (%s)", m.Value, result))
	}
	return success, nil
}

func checkBodyJSONPresent(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	success := result.Exists()
	if !success {
		return false, m.AppErr(fmt.Sprintf("JSON Field Match Failed - no field present for pattern (%s)", m.JSON))
	}
	m.Result = result.String()
	return success, nil
}

func checkBodyJSONCount(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	if result.Int() != m.Count {
		return false, m.AppErr(fmt.Sprintf("JSON Count Field Match Failed - found (%d) not (%d) occurrences of pattern (%s)", result.Int(), m.Count, m.JSON))
	}
	return true, nil
}

func checkBodyJSONRegex(m *Match, tc *TestCase) (bool, error) {
	result := gjson.Get(tc.Body, m.JSON)
	if len(result.String()) == 0 {
		return false, m.AppErr(fmt.Sprintf("JSON Regex Match Failed - no field present for pattern (%s)", m.JSON))
	}
	regex, err := regexp.Compile(m.Regex)
	if err != nil {
		return false, err
	}

	success := regex.MatchString(result.String())
	if !success {
		return false, m.AppErr(fmt.Sprintf("JSON Regex Match Failed - selected field (%s) does not match regex (%s)", result.String(), m.Regex))
	}
	return success, nil
}

func checkBodyLength(m *Match, tc *TestCase) (bool, error) {
	success := len(tc.Body) == int(*m.BodyLength)
	if !success {
		return false, m.AppErr(fmt.Sprintf("Check Body Length - body length (%d) does not match expected length (%d)", len(tc.Body), *m.BodyLength))
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
		return false, m.AppErr(fmt.Sprintf("Authorisation Bearer Match Failed - no header value found"))
	}

	idx := strings.Index(headerValue, "Bearer ")
	if idx == -1 {
		idx = strings.Index(headerValue, "bearer ")
	}
	if idx == -1 {
		return false, m.AppErr(fmt.Sprintf("Authorisation Bearer Match value Failed - no header bearer value found"))
	}
	m.Result = headerValue[idx+7:] // copy the token after 7 chars "Bearer "...
	return true, nil
}

func checkCustom(m *Match, tc *TestCase) (bool, error) {
	return true, nil //TODO: implement custom checks
}

// ProcessReplacementFields allows parameter replacement within match string fields
func (m *Match) ProcessReplacementFields(ctx *Context) {
	m.Header, _ = replaceContextField(m.Header, ctx)
	m.HeaderPresent, _ = replaceContextField(m.HeaderPresent, ctx)
	m.JSON, _ = replaceContextField(m.JSON, ctx)
	m.Value, _ = replaceContextField(m.Value, ctx)
	m.ContextName, _ = replaceContextField(m.ContextName, ctx)
}

// Clone duplicates a Match into a separate independent object
// TODO: consider cloning the contextPut
// Omit bodylength for now...
func (m *Match) Clone() Match {
	ma := Match{Authorisation: m.Authorisation,
		ContextName:     m.ContextName,
		Count:           m.Count,
		Custom:          m.Custom,
		Description:     m.Description,
		Header:          m.Header,
		HeaderPresent:   m.HeaderPresent,
		JSON:            m.JSON,
		MatchType:       m.MatchType,
		Numeric:         m.Numeric,
		Regex:           m.Regex,
		Result:          m.Result,
		ReplaceEndpoint: m.ReplaceEndpoint,
		Value:           m.Value,
	}
	return ma
}
