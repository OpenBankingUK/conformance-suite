package results

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"gopkg.in/resty.v1"
	"time"
)

type Metrics struct {
	TestCase     *model.TestCase `json:"-"`
	ResponseTime time.Duration   `json:"response_time"` // Http Response Time
	ResponseSize int             `json:"response_size"` // Size of the HTTP Response body
}

var NoMetrics = Metrics{}

func NewMetricsFromRestyResponse(testCase *model.TestCase, response *resty.Response) Metrics {
	return NewMetrics(testCase, response.Time(), len(response.Body()))
}

func NewMetrics(testCase *model.TestCase, responseTime time.Duration, responseSize int) Metrics {
	return Metrics{
		TestCase:     testCase,
		ResponseTime: responseTime,
		ResponseSize: responseSize,
	}
}
