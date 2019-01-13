package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
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
	tracer.Silent = false
	executor := MakeExecutor()
	executor.SetCertificates(defn.SigningCert, defn.TransportCert)
	return runTests(defn, executor)
}

func runTests(defn *RunDefinition, executor *Executor) (reporting.Result, error) {
	rulectx := &model.Context{}

	reportTestResults := []reporting.Test{}
	reportSpecs := []reporting.Specification{reporting.Specification{Tests: reportTestResults}}
	reportResult := reporting.Result{Specifications: reportSpecs}

	fmt.Println("RUN TEST LOOP :- ")

	for _, spec := range defn.SpecTests {
		for _, testcase := range spec.TestCases {
			req, err := testcase.Prepare(rulectx)
			if err != nil {
				return reporting.Result{}, err
			}

			resp, err := executor.ExecuteTestCase(req, &testcase, rulectx)
			if err != nil {
				return reporting.Result{}, err
			}

			result, err := testcase.Validate(resp, rulectx)
			if err != nil {
				return reporting.Result{}, err
			}
			reportTestResults = append(reportTestResults, makeTestResult(&testcase, result))
		}
	}

	fmt.Printf("%#v\n", reportResult)
	return reportResult, nil
}

func makeTestResult(tc *model.TestCase, result bool) reporting.Test {
	return reporting.Test{Name: tc.Name, Endpoint: tc.Input.Method + " " + tc.Input.Endpoint, Pass: result}
}
