package results

import (
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

func TestNewTestCaseResult123(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")
	result := NewTestCaseResult("123", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")

	assert.Equal("123", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

func TestNewTestCaseResult321(t *testing.T) {
	assert := test.NewAssert(t)

	err := errors.New("some error")

	result := NewTestCaseResult("321", true, NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")
	assert.Equal("321", result.Id)
	assert.True(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

func TestNewTestCaseFailResult(t *testing.T) {
	assert := test.NewAssert(t)
	err := errors.New("some error")

	result := NewTestCaseFail("id", NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")

	assert.Equal("id", result.Id)
	assert.False(result.Pass)
	assert.Equal(NoMetrics(), result.Metrics)
	assert.Equal(err.Error(), result.Fail[0])
}

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
