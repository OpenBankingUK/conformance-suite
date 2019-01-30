package results

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTestResult(t *testing.T) {
	err := errors.New("some error")
	result := NewTestCaseResult("123", true, NoMetrics, err)
	assert.Equal(t, "123", result.Id)
	assert.True(t, result.Pass)
	assert.Equal(t, NoMetrics, result.Metrics)
	assert.Equal(t, err.Error(), string(result.Fail))

	result = NewTestCaseResult("321", true, NoMetrics, err)
	assert.Equal(t, "321", result.Id)
	assert.True(t, result.Pass)
	assert.Equal(t, NoMetrics, result.Metrics)
	assert.Equal(t, err.Error(), string(result.Fail))
}

func TestNewTestFailResult(t *testing.T) {
	err := errors.New("some error")
	result := NewTestCaseFail("id", NoMetrics, err)
	assert.Equal(t, "id", result.Id)
	assert.False(t, result.Pass)
	assert.Equal(t, NoMetrics, result.Metrics)
	assert.Equal(t, err.Error(), string(result.Fail))
}
