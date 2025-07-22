package results

import "log"

// TestCase result for a run
type TestCase struct {
	Id         string   `json:"id"`
	Pass       bool     `json:"pass"`
	Metrics    Metrics  `json:"metrics"`
	Fail       []string `json:"fail,omitempty"`
	Detail     string   `json:"detail"`
	RefURI     string   `json:"refURI"`
	Endpoint   string   `json:"endpoint"`
	API        string   `json:"-"`
	APIVersion string   `json:"-"`
	HttpStatus string   `json:"httpStatusCode"`
	Warnings   []string `json:"warnings,omitempty"`
}

// NewTestCaseFail returns a failed test
func NewTestCaseFail(id string, metrics Metrics, errs []error, endpoint, api, apiVersion, detail, refURI, httpStatus string, isWarningTest bool) TestCase {
	// For a warning test, we want to pass the test but include the errors as warnings

	if isWarningTest {
		warnings := []string{}
		// Convert errors to warnings and mark as passed
		for _, err := range errs {
			warnings = append(warnings, err.Error())
		}

		log.Printf("Warning test case %s encountered errors: %v", id, warnings)
		return NewWarningTestCaseResult(id, true, metrics, nil, endpoint, api, apiVersion, detail, refURI, httpStatus, warnings)
	}
	return NewTestCaseResult(id, false, metrics, errs, endpoint, api, apiVersion, detail, refURI, httpStatus)
}

// NewTestCaseResult return a new TestCase instance
func NewTestCaseResult(id string, pass bool, metrics Metrics, errs []error, endpoint, apiName, apiVersion, detail, refURI, httpStatus string) TestCase {
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
		RefURI:     refURI,
		HttpStatus: httpStatus,
	}
}

// NewWarningTestCaseResult return a new TestCase instance
func NewWarningTestCaseResult(id string, pass bool, metrics Metrics, errs []error, endpoint, apiName, apiVersion, detail, refURI, httpStatus string, warnings []string) TestCase {
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
		RefURI:     refURI,
		HttpStatus: httpStatus,
		Warnings:   warnings,
	}
}

type ResultKey struct {
	APIName    string
	APIVersion string
}
