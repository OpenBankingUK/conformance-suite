package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// TestCaseExecutor defines an interface capable of executing a testcase
type TestCaseExecutor interface {
	ExecuteTestCase(r *resty.Request, t *model.TestCase, ctx *model.Context) (*resty.Response, error)
	SetCertificates(certificateSigning, certificationTransport authentication.Certificate) error
}

// RunDefinition captures all the information required to run the test cases
type RunDefinition struct {
	DiscoModel    *discovery.Model
	SpecTests     []generation.SpecificationTestCases
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

// RunTestCases runs the testCases
func RunTestCases(defn *RunDefinition) (reporting.Result, error) {
	executor := MakeExecutor()
	executor.SetCertificates(defn.SigningCert, defn.TransportCert)

	rulectx := &model.Context{}
	reportSpecs := []reporting.Specification{}
	for _, spec := range defn.SpecTests {
		reportTestResults := []reporting.Test{}
		logrus.Println("running " + spec.Specification.Name)
		for _, testcase := range spec.TestCases {
			req, err := testcase.Prepare(rulectx)
			if err != nil {
				logrus.Error(err)
				reportTestResults = append(reportTestResults, makeTestResult(testcase, false))
				reportSpecs = append(reportSpecs, makeSpecResult(spec.Specification, reportTestResults))
				return makeReportResult(reportSpecs), err
			}
			resp, err := executor.ExecuteTestCase(req, &testcase, rulectx)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"testcase":   testcase.Name,
					"method":     testcase.Input.Method,
					"endpoint":   testcase.Input.Endpoint,
					"err":        err.Error(),
					"statuscode": testcase.Expect.StatusCode,
				}).Info("FAIL")
				reportTestResults = append(reportTestResults, makeTestResult(testcase, false))
				reportSpecs = append(reportSpecs, makeSpecResult(spec.Specification, reportTestResults))
				continue
			}

			result, err := testcase.Validate(resp, rulectx)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"testcase":   testcase.Name,
					"method":     testcase.Input.Method,
					"endpoint":   testcase.Input.Endpoint,
					"err":        err.Error(),
					"statuscode": testcase.Expect.StatusCode,
				}).Info("FAIL")

				reportTestResults = append(reportTestResults, makeTestResult(testcase, false))
				reportSpecs = append(reportSpecs, makeSpecResult(spec.Specification, reportTestResults))
				continue
			}
			reportTestResults = append(reportTestResults, makeTestResult(testcase, result))
			logrus.WithFields(logrus.Fields{
				"testcase":   testcase.Name,
				"method":     testcase.Input.Method,
				"endpoint":   testcase.Input.Endpoint,
				"statuscode": testcase.Expect.StatusCode,
			}).Info("PASS")
		}
		reportSpecs = append(reportSpecs, makeSpecResult(spec.Specification, reportTestResults))
	}
	logrus.Println("runTests OK")
	return makeReportResult(reportSpecs), nil
}

func makeTestResult(tc model.TestCase, result bool) reporting.Test {
	return reporting.Test{
		Id:       tc.ID,
		Name:     tc.Name,
		Endpoint: tc.Input.Method + " " + tc.Input.Endpoint,
		Pass:     result,
	}
}

func makeSpecResult(spec discovery.ModelAPISpecification, testResults []reporting.Test) reporting.Specification {
	return reporting.Specification{
		Name:          spec.Name,
		Version:       spec.Version,
		URL:           spec.URL,
		SchemaVersion: spec.SchemaVersion,
		Tests:         testResults,
	}
}

func makeReportResult(specsResults []reporting.Specification) reporting.Result {
	return reporting.Result{
		Id:             uuid.New(),
		Specifications: specsResults,
	}
}
