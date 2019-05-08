package results

// TestCase result for a run
type TestCase struct {
	Id         string   `json:"id"`
	Pass       bool     `json:"pass"`
	Metrics    Metrics  `json:"metrics"`
	Fail       []string `json:"fail,omitempty"`
	Detail     string   `json:"detail"`
	Endpoint   string   `json:"endpoint"`
	API        string   `json:"-"`
	APIVersion string   `json:"-"`
}

// NewTestCaseFail returns a failed test
func NewTestCaseFail(id string, metrics Metrics, errs []error, endpoint, api, apiVersion, detail string) TestCase {
	return NewTestCaseResult(id, false, metrics, errs, endpoint, api, apiVersion, detail)
}

// NewTestCaseResult return a new TestCase instance
func NewTestCaseResult(id string, pass bool, metrics Metrics, errs []error, endpoint, apiName, apiVersion, detail string) TestCase {
	reasons := []string{}
	for _, err := range errs {
		reasons = append(reasons, err.Error())
	}
	return TestCase{
		API:        apiName,
		APIVersion: apiVersion,
		Id:         id,
		Pass:       pass,
		Metrics:    metrics,
		Fail:       reasons,
		Endpoint:   endpoint,
		Detail:     detail,
	}
}

type ResultKey struct {
	APIName    string
	APIVersion string
}
