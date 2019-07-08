package results

import (
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"errors"
	"testing"
)

func TestNewTestCaseResult123(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")
	result := NewTestCaseResult("123", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri")

	assert.Equal("123", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

func TestNewTestCaseResult321(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")

	result := NewTestCaseResult("321", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri")
	assert.Equal("321", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

func TestNewTestCaseFailResult(t *testing.T) {
	assert := test.NewAssert(t)
	err := errors.New("some error")

	result := NewTestCaseFail("id", NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri")

	assert.Equal("id", result.Id)
	assert.False(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

func TestTestCaseResultJsonMarshal(t *testing.T) {
	require := test.NewRequire(t)

	result := NewTestCaseResult("123", true, NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri")

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
	"refURI": "https://openbanking.org.uk/ref/uri"
}
	`
	actual, err := json.Marshal(result)
	require.NoError(err)
	require.NotEmpty(actual)

	require.JSONEq(expected, string(actual))
}
