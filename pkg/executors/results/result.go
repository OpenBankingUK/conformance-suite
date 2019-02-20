package results

// TestCase result for a run
type TestCase struct {
	Id      string  `json:"id"`
	Pass    bool    `json:"pass"`
	Metrics Metrics `json:"metrics"`
	Fail    string  `json:"fail,omitempty"`
}

// NewTestCaseFail returns a failed test
func NewTestCaseFail(id string, metrics Metrics, err error) TestCase {
	return NewTestCaseResult(id, false, metrics, err)
}

// NewTestCaseResult return a new TestCase instance
func NewTestCaseResult(id string, pass bool, metrics Metrics, err error) TestCase {
	var failReason string
	if err != nil {
		failReason = err.Error()
	}
	return TestCase{
		Id:      id,
		Pass:    pass,
		Metrics: metrics,
		Fail:    failReason,
	}
}
