package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	resty "gopkg.in/resty.v1"
)

// TestCaseExecutor defines an interface capable of executing a testcase
type TestCaseExecutor interface {
	ExecuteTestCase(r *resty.Request, t *model.TestCase, ctx *model.Context) (*resty.Response, error)
	SetCertificates(certificateSigning, certificationTransport authentication.Certificate)
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
	fmt.Println("Up and Running -- execution those test cases!!!")

	return reporting.Result{}, nil
}
