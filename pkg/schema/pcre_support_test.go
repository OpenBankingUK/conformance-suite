package schema

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegexp2CompilerPCREPattern verifies that a PCRE-only pattern (negative
// lookahead) compiles without error. This is the pattern used by the OB v4.0.0
// x-idempotency-key header and the root cause of the original bug.
func TestRegexp2CompilerPCREPattern(t *testing.T) {
	_, err := regexp2Compiler(`^(?!\s)(.*)(\\S)$`)
	require.NoError(t, err)
}

// TestRegexp2CompilerSetsMatchTimeout verifies that the compiled regexp has a
// finite MatchTimeout set, protecting against ReDoS via catastrophic backtracking.
// This is a security property, so white-box inspection of the internal field is
// acceptable.
func TestRegexp2CompilerSetsMatchTimeout(t *testing.T) {
	m, err := regexp2Compiler(`^(?!\s)(.*)(\\S)$`)
	require.NoError(t, err)
	matcher, ok := m.(*regexp2Matcher)
	require.True(t, ok, "expected *regexp2Matcher")
	assert.Equal(t, matchTimeout, matcher.re.MatchTimeout, "MatchTimeout must be finite to prevent ReDoS")
}

// TestRegexp2CompilerRE2Pattern verifies that a standard RE2-compatible pattern
// still compiles correctly under regexp2, confirming there is no regression for
// existing v3.x spec patterns.
func TestRegexp2CompilerRE2Pattern(t *testing.T) {
	_, err := regexp2Compiler(`^[A-Z]{2,2}$`)
	require.NoError(t, err)
}

// TestRegexp2CompilerMalformedPattern verifies that a genuinely malformed
// pattern (invalid even as PCRE) is rejected with an error, preserving the
// hard-fail behaviour for bad spec files.
func TestRegexp2CompilerMalformedPattern(t *testing.T) {
	_, err := regexp2Compiler(`[unclosed`)
	assert.Error(t, err)
}

// TestRegexp2MatcherPCREPatternMatches verifies that a value satisfying the
// exact OB v4.0.0 x-idempotency-key pattern is correctly matched.
// Pattern: ^(?!\s)(.*)(\S)$  — two capture groups, not the folded form .*\S.
func TestRegexp2MatcherPCREPatternMatches(t *testing.T) {
	m, err := regexp2Compiler(`^(?!\s)(.*)(\S)$`)
	require.NoError(t, err)
	// "abc" does not start with whitespace and does not end with whitespace — should match.
	assert.True(t, m.MatchString("abc"))
}

// TestRegexp2MatcherPCREPatternRejects verifies that a value violating the
// exact OB v4.0.0 x-idempotency-key pattern is correctly rejected.
// Pattern: ^(?!\s)(.*)(\S)$  — two capture groups, not the folded form .*\S.
func TestRegexp2MatcherPCREPatternRejects(t *testing.T) {
	m, err := regexp2Compiler(`^(?!\s)(.*)(\S)$`)
	require.NoError(t, err)
	// " abc" starts with whitespace — should not match.
	assert.False(t, m.MatchString(" abc"))
}

// --- v4.0.0 spec loading integration tests ---
// Each test directly asserts the fix for the bug reported in:
//   ERROR cannot Load OpenApi Spec from file spec/v4.0.0/account-info-openapi.json,
//   invalid components: parameter "x-idempotency-key" schema is invalid:
//   cannot compile pattern "^(?!\s)(.*)(\\S)$"

func TestV4AccountsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v4.0.0")
	require.NoError(t, err)
}

func TestV4PaymentInitiationSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Payment Initiation API", "v4.0.0")
	require.NoError(t, err)
}

func TestV4ConfirmationOfFundsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Confirmation of Funds API Specification", "v4.0.0")
	require.NoError(t, err)
}

func TestV4VariableRecurringPaymentsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Variable Recurring Payments API Specification", "v4.0.0")
	require.NoError(t, err)
}

func TestV4CommercialVRPSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Commercial Variable Recurring Payments API Specification", "v4.0.0")
	require.NoError(t, err)
}

// TestV4RuntimeValidationWithPCRESpec verifies that calling Validate against
// a v4.0.0 spec (which contains PCRE patterns in parameter definitions) does
// not produce errors during response validation.
//
// This test makes a narrower, honest claim: PCRE patterns present in the spec
// do not cause the runtime validation path to fail. The x-idempotency-key
// pattern is a request header parameter on POST endpoints; the FCS runtime
// validator validates responses only (ExcludeRequestBody: true), so the
// pattern is never applied to a value at runtime. The Options.RegexCompiler
// wiring in validateResponse is forward-looking — correct and necessary for
// any future PCRE patterns placed on response body fields.
func TestV4RuntimeValidationWithPCRESpec(t *testing.T) {
	validator, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v4.0.0")
	require.NoError(t, err)

	// Minimal valid GET /accounts response body for v4.0.0.
	const body = `{
		"Data": {"Account": []},
		"Links": {"Self": "https://example.com/open-banking/v4.0/aisp/accounts"},
		"Meta": {"TotalPages": 1}
	}`
	r := HTTPResponse{
		Method:     "GET",
		Path:       "/open-banking/v4.0/aisp/accounts",
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(body),
		Header: http.Header{
			"Content-Type":          []string{"application/json; charset=utf-8"},
			"X-Fapi-Interaction-Id": []string{"test-interaction-id"},
		},
	}
	_, err = validator.Validate(r)
	assert.NoError(t, err)
}

// --- v4.0.1 spec loading integration tests ---
// Each test verifies that the v4.0.1 spec file loads without error,
// preventing filename and validation regressions.

func TestV4_0_1_AccountsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v4.0.1")
	require.NoError(t, err)
}

func TestV4_0_1_PaymentInitiationSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Payment Initiation API", "v4.0.1")
	require.NoError(t, err)
}

func TestV4_0_1_ConfirmationOfFundsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Confirmation of Funds API Specification", "v4.0.1")
	require.NoError(t, err)
}

func TestV4_0_1_VariableRecurringPaymentsSpecLoads(t *testing.T) {
	_, err := NewRawOpenAPI3Validator("Variable Recurring Payments API Specification", "v4.0.1")
	require.NoError(t, err)
}
