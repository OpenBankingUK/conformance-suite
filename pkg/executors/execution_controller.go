package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
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
	rulectx.Put("SigingCert", defn.SigningCert)

	reportTestResults := []reporting.Test{}
	reportSpecs := []reporting.Specification{reporting.Specification{Tests: reportTestResults}}
	reportResult := reporting.Result{Specifications: reportSpecs}

	for _, spec := range defn.SpecTests {
		logrus.Println("running " + spec.Specification.Name)
		for _, testcase := range spec.TestCases {
			req, err := testcase.Prepare(rulectx)
			if err != nil {
				reportTestResults = append(reportTestResults, makeTestResult(&testcase, false))
				logrus.Error(err)
				return reportResult, err
			}
			resp, err := executor.ExecuteTestCase(req, &testcase, rulectx)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"testcase":     testcase.Name,
					"method":       testcase.Input.Method,
					"endpoint":     testcase.Input.Endpoint,
					"err":          err.Error(),
					"statuscode":   testcase.Expect.StatusCode,
					"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
					"responsesize": testcase.ResponseSize,
				}).Info("FAIL")
				reportTestResults = append(reportTestResults, makeTestResult(&testcase, false))
				continue
			}

			result, err := testcase.Validate(resp, rulectx)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"testcase":     testcase.Name,
					"method":       testcase.Input.Method,
					"endpoint":     testcase.Input.Endpoint,
					"err":          err.Error(),
					"statuscode":   testcase.Expect.StatusCode,
					"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
					"responsesize": testcase.ResponseSize,
				}).Info("FAIL")

				reportTestResults = append(reportTestResults, makeTestResult(&testcase, false))
				continue
			}
			reportTestResults = append(reportTestResults, makeTestResult(&testcase, result))
			logrus.WithFields(logrus.Fields{
				"testcase":     testcase.Name,
				"method":       testcase.Input.Method,
				"endpoint":     testcase.Input.Endpoint,
				"statuscode":   testcase.Expect.StatusCode,
				"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
				"responsesize": testcase.ResponseSize,
			}).Info("PASS")
		}
	}

	logrus.Println("runTests OK")
	return reportResult, nil
}

func makeTestResult(tc *model.TestCase, result bool) reporting.Test {
	return reporting.Test{Name: tc.Name, Endpoint: tc.Input.Method + " " + tc.Input.Endpoint, Pass: result}
}
