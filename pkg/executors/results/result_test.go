package results

import (
	"encoding/json"

	"github.com/OpenBankingUK/conformance-suite/pkg/test"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

// TestNewTestCaseResult123 verifies the creation of a TestCaseResult and its properties using given parameters and assertions.
func TestNewTestCaseResult123(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")
	result := NewTestCaseResult("123", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")

	assert.Equal("123", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

// TestNewTestCaseResult321 tests the creation of a TestCaseResult with predefined parameters and validates the expected attributes.
func TestNewTestCaseResult321(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")

	result := NewTestCaseResult("321", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")
	assert.Equal("321", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

// TestNewTestCaseFailResult verifies that NewTestCaseFail creates a TestCase with expected failure details and properties.
func TestNewTestCaseFailResult(t *testing.T) {
	assert := test.NewAssert(t)
	err := errors.New("some error")

	result := NewTestCaseFail("id", NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200", false)

	assert.Equal("id", result.Id)
	assert.False(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

// TestTestCaseResultJsonMarshal tests the JSON marshalling of a TestCaseResult to ensure expected structure and content.
func TestTestCaseResultJsonMarshal(t *testing.T) {
	result := NewTestCaseResult("123", true, NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")

	expected := `
{
	"endpoint": "endpoint",
	"id": "123",
	"pass": true,
	"metrics": {
		"response_time": 0,
		"response_size": 0
	},
	"detail": "detailed description",
	"refURI": "https://openbanking.org.uk/ref/uri",
	"httpStatusCode":"200"
}
	`
	actual, err := json.Marshal(result)
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.JSONEq(t, expected, string(actual))
}

// TestNewTestCaseWarningFail verifies the behaviour of NewTestCaseFail function for a warning test case scenario.
func TestNewTestCaseWarningFail(t *testing.T) {
	assert := test.NewAssert(t)
	err := errors.New("some error")

	result := NewTestCaseFail("id", NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200", true)

	assert.Equal("id", result.Id)
	assert.True(result.Pass) // Should pass test even if fail.
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Warnings[0])
}
