package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

// RunDefinition captures all the information required to run the test cases
type RunDefinition struct {
	DiscoModel    *discovery.Model
	SpecTests     []generation.SpecificationTestCases
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

type TestCaseRunner struct {
	executor         *Executor
	definition       RunDefinition
	daemonController DaemonController
	logger           *logrus.Entry
	resolver         func(groups []permissions.Group) permissions.CodeSetResultSet
	runningLock      *sync.Mutex
	running          bool
}

func NewTestCaseRunner(definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logrus.New().WithField("module", "TestCaseRunner"),
		resolver:         permissions.Resolver,
		runningLock:      &sync.Mutex{},
		running:          false,
	}
}

// RunTestCases runs the testCases
func (r *TestCaseRunner) RunTestCases() error {
	r.runningLock.Lock()
	defer r.runningLock.Unlock()
	if r.running {
		return errors.New("test cases runner already running")
	}
	r.running = true

	codeSets := r.permissionSets()
	r.logger.WithField("size", len(codeSets)).Info("code sets required")

	go r.runTestCasesAsync()

	return nil
}

// permissionSets calls resolver to get list of permission sets required to run all test cases
func (r *TestCaseRunner) permissionSets() permissions.CodeSetResultSet {
	var groups []permissions.Group
	for _, spec := range r.definition.SpecTests {
		for _, tc := range spec.TestCases {
			groups = append(groups, model.NewPermissionGroup(tc))
		}
	}
	return r.resolver(groups)
}

func (r *TestCaseRunner) runTestCasesAsync() {
	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx()
	ctxLogger := r.logger.WithField("id", uuid.New())
	ctxLogger.Info("running async test")
	for _, spec := range r.definition.SpecTests {
		r.executeSpecTests(spec, ruleCtx, ctxLogger)
	}
	r.setNotRunning()
}

func (r *TestCaseRunner) setNotRunning() {
	r.runningLock.Lock()
	r.running = false
	r.runningLock.Unlock()
}

func (r *TestCaseRunner) makeRuleCtx() *model.Context {
	ruleCtx := &model.Context{}
	ruleCtx.Put("SigningCert", r.definition.SigningCert)
	for _, customTest := range r.definition.DiscoModel.DiscoveryModel.CustomTests {
		for k, v := range customTest.Replacements {
			ruleCtx.Put(k, v)
		}
	}
	return ruleCtx
}

func (r *TestCaseRunner) executeSpecTests(spec generation.SpecificationTestCases, ruleCtx *model.Context, ctxLogger *logrus.Entry) {
	ctxLogger = ctxLogger.WithField("spec", spec.Specification.Name)
	for _, testcase := range spec.TestCases {
		if r.daemonController.ShouldStop() {
			ctxLogger.Info("stop test run received, aborting runner")
			return
		}
		testResult := r.executeTest(testcase, ruleCtx, ctxLogger)
		r.daemonController.Results() <- testResult
	}
}

func (r *TestCaseRunner) executeTest(tc model.TestCase, ruleCtx *model.Context, ctxLogger *logrus.Entry) results.TestCase {
	ctxLogger = logWithTestCase(r.logger, tc)

	req, err := tc.Prepare(ruleCtx)
	if err != nil {
		ctxLogger.WithError(err).Error("preparing executing test")
		return results.NewTestCaseFail(tc.ID, results.NoMetrics, err)
	}

	resp, metrics, err := r.executor.ExecuteTestCase(req, &tc, ruleCtx)
	ctxLogger = logWithMetrics(ctxLogger, metrics)
	if err != nil {
		ctxLogger.WithError(err).WithField("result", "FAIL").Info("test result")
		return results.NewTestCaseFail(tc.ID, metrics, err)
	}

	result, err := tc.Validate(resp, ruleCtx)
	if err != nil {
		ctxLogger.WithError(err).WithField("result", passText[result]).Info("test result")
		return results.NewTestCaseFail(tc.ID, metrics, err)
	}

	ctxLogger.WithError(err).WithField("result", passText[result]).Info("test result")
	return results.NewTestCaseResult(tc.ID, result, metrics, err)
}

var passText = map[bool]string{
	true:  "PASS",
	false: "FAIL",
}

func logWithTestCase(logger *logrus.Entry, tc model.TestCase) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"testcase":   tc.Name,
		"method":     tc.Input.Method,
		"endpoint":   tc.Input.Endpoint,
		"statuscode": tc.Expect.StatusCode,
	})
}

func logWithMetrics(logger *logrus.Entry, metrics results.Metrics) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"responsetime": fmt.Sprintf("%v", metrics.ResponseTime),
		"responsesize": metrics.ResponseSize,
	})
}
