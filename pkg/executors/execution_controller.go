package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
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
	mux              *sync.Mutex
	running          bool
}

func NewTestCaseRunner(definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logrus.New().WithField("module", "TestCaseRunner"),
		mux:              &sync.Mutex{},
		running:          false,
	}
}

// RunTestCases runs the testCases
func (r *TestCaseRunner) RunTestCases() error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if r.running == true {
		return errors.New("test cases runner already running")
	}
	r.running = true

	go r.runTestCasesAsync()

	return nil
}

func (r *TestCaseRunner) runTestCasesAsync() {
	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx()
	ctxLogger := r.logger.WithField("id", uuid.New())
	ctxLogger.Info("running async test")
	for _, spec := range r.definition.SpecTests {
		err := r.executeSpecTests(spec, ruleCtx, ctxLogger)
		if err != nil {
			r.daemonController.Errors() <- err
		}
	}
	r.setNotRunning()
}

func (r *TestCaseRunner) setNotRunning() {
	r.mux.Lock()
	r.running = false
	r.mux.Unlock()
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

func (r *TestCaseRunner) executeSpecTests(spec generation.SpecificationTestCases, ruleCtx *model.Context, ctxLogger *logrus.Entry) error {
	ctxLogger = ctxLogger.WithField("spec", spec.Specification.Name)
	for _, testcase := range spec.TestCases {
		if r.daemonController.ShouldStop() {
			ctxLogger.Info("stop test run received, aborting runner")
			return nil
		}
		testResult, err := r.executeTest(testcase, ruleCtx, ctxLogger)
		r.daemonController.Results() <- testResult
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TestCaseRunner) executeTest(testcase model.TestCase, ruleCtx *model.Context, ctxLogger *logrus.Entry) (results.Test, error) {
	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"testcase": testcase.Name,
		"method":   testcase.Input.Method,
		"endpoint": testcase.Input.Endpoint,
	})

	req, err := testcase.Prepare(ruleCtx)
	if err != nil {
		logrus.Error(err)
		return results.NewTestFailResult(testcase.ID), err
	}

	resp, err := r.executor.ExecuteTestCase(req, &testcase, ruleCtx)
	if err != nil {
		ctxLogger.WithFields(logrus.Fields{
			"err":          err.Error(),
			"statuscode":   testcase.Expect.StatusCode,
			"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
			"responsesize": testcase.ResponseSize,
			"result":       "FAIL",
		}).Info("test result")
		return results.NewTestFailResult(testcase.ID), err
	}

	result, err := testcase.Validate(resp, ruleCtx)
	if err != nil {
		ctxLogger.WithFields(logrus.Fields{
			"err":          err.Error(),
			"statuscode":   testcase.Expect.StatusCode,
			"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
			"result":       "FAIL",
		}).Info("test result")
		return results.NewTestFailResult(testcase.ID), err
	}

	ctxLogger.WithFields(logrus.Fields{
		"statuscode":   testcase.Expect.StatusCode,
		"responsetime": fmt.Sprintf("%v", testcase.ResponseTime),
		"result":       "PASS",
	}).Info("test result")

	return results.NewTestResult(testcase.ID, result), nil
}
