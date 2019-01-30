package results

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	tc := &model.TestCase{}

	metrics := NewMetrics(tc, time.Second, 1)

	assert.Equal(t, tc, metrics.TestCase)
	assert.Equal(t, time.Second, metrics.ResponseTime)
	assert.Equal(t, 1, metrics.ResponseSize)
}
