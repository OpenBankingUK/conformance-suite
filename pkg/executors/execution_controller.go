package executors

import (
	"fmt"
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// RunDefinition captures all the information required to run the test cases
type RunDefinition struct {
	DiscoModel    *discovery.Model
	TestCaseRun   generation.TestCasesRun
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

type TestCaseRunner struct {
	executor         *Executor
	definition       RunDefinition
	daemonController DaemonController
	logger           *logrus.Entry
	runningLock      *sync.Mutex
	running          bool
}

func NewTestCaseRunner(definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logrus.New().WithField("module", "TestCaseRunner"),
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

	go r.runTestCasesAsync()

	return nil
}

func (r *TestCaseRunner) runTestCasesAsync() {
	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx()
	ctxLogger := r.logger.WithField("id", uuid.New())

	ctxLogger.Debug("running pre-ExecutionCustomComponents")
	err := r.preExecuteCustomComponents(ruleCtx, ctxLogger)
	if err != nil {
		errResult := results.NewTestCaseFail("PreExecuteComponet", results.Metrics{}, err)
		r.daemonController.Results() <- errResult
		return
	}
	ctxLogger.Debug(fmt.Sprintf("context: %#v", ruleCtx), ctxLogger)

	ctxLogger.Info("running async test")
	for _, spec := range r.definition.TestCaseRun.TestCases {
		r.executeSpecTests(spec, ruleCtx, ctxLogger)
	}
	r.setNotRunning()
}

func (r *TestCaseRunner) preExecuteCustomComponents(ctx *model.Context, ctxLogger *logrus.Entry) error {
	ctxLogger.Debug("preExecution Custom Components")

	executionUnits, err := r.getExecutionUnits(ctxLogger)
	if err != nil {
		ctxLogger.Debug("failed to get Execution Units: " + err.Error())
		return err
	}
	if len(executionUnits) == 0 {
		ctxLogger.Debug("no ExecutionUnits")
		return nil
	}
	//TODO: move registry creation and component initialisation up a level or two
	//TODO: consider component reentrancy, use in multiple threads
	reg := model.NewRegistry()
	comp, err := model.LoadComponent("tokenProviderComponent.json")
	if err != nil {
		ctxLogger.Debug("Load Model Failed: " + err.Error())
		return err
	}
	reg.Add(comp.Name, &comp)
	for _, unitName := range executionUnits {
		// Execute components
		ctxLogger.Debug("process execution unit :" + unitName)
		comp, exists := reg.Get(unitName)
		if !exists {
			msg := "component execution error: component (%s) does not exist: " + unitName
			ctxLogger.Debug(msg)
			return errors.New(msg)
		}

		cmp, ok := comp.(*model.Component)
		if !ok {
			msg := "component execution error: component (%s) cannot cast registry value:" + unitName
			ctxLogger.Debug(msg)
			return errors.New(msg)
		}

		err := cmp.ValidateParameters(ctx) // correct parameters for component exist in context
		if err != nil {
			msg := fmt.Sprintf("component execution error: component (%s) cannot ValidateParameters: %s", unitName, err.Error())
			ctxLogger.Debug(msg)
			return errors.New(msg)
		}
		cmp.ProcessReplacementFields(*ctx)
		r.executeComponentTests(cmp, ctx, ctxLogger)
	}
	return nil
}

func (r *TestCaseRunner) getExecutionUnits(ctxLogger *logrus.Entry) ([]string, error) {
	var units []string
	for _, customTest := range r.definition.DiscoModel.DiscoveryModel.CustomTests {
		for _, executionUnit := range customTest.Execution {
			ctxLogger.Debug("Execution unit: " + executionUnit)
			units = append(units, executionUnit)
		}
	}
	return units, nil
}

func (r *TestCaseRunner) executeComponentTests(comp *model.Component, ruleCtx *model.Context, ctxLogger *logrus.Entry) {
	ctxLogger = ctxLogger.WithField("component", comp.Name)
	ctxLogger.Debugln("execute component Tests...")
	for _, testcase := range comp.Tests {
		if r.daemonController.ShouldStop() {
			ctxLogger.Debugln("stop component test run received, aborting runner")
			return
		}
		ctxLogger.Debugln("executing testcase: " + testcase.Name)
		testResult := r.executeTest(testcase, ruleCtx, ctxLogger)
		r.daemonController.Results() <- testResult
	}
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
