package results

import (
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
)

func TestNewMetricsFromRestyResponse(t *testing.T) {
	tc := &model.TestCase{}
	response := &resty.Response{Request: &resty.Request{Time: time.Now()}}

	metrics := NewMetricsFromRestyResponse(tc, response)

	assert.Equal(t, tc, metrics.TestCase)
	assert.True(t, metrics.ResponseTime < time.Second)
	assert.Equal(t, 0, metrics.ResponseSize)
}

func TestNewMetrics(t *testing.T) {
	tc := &model.TestCase{}

	metrics := NewMetrics(tc, time.Second, 1)

	assert.Equal(t, tc, metrics.TestCase)
	assert.Equal(t, time.Second, metrics.ResponseTime)
	assert.Equal(t, 1, metrics.ResponseSize)
}
