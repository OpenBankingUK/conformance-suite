package executors

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
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

// NewTestCaseRunner -
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

// NewConsentAcquisitionRunner -
func NewConsentAcquisitionRunner(definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logrus.New().WithField("module", "ConsentAcquisitionRunner"),
		runningLock:      &sync.Mutex{},
		running:          false,
	}
}

// RunTestCases runs the testCases
func (r *TestCaseRunner) RunTestCases(ctx *model.Context) error {
	r.runningLock.Lock()
	defer r.runningLock.Unlock()
	if r.running {
		return errors.New("test cases runner already running")
	}
	r.running = true

	go r.runTestCasesAsync(ctx)

	return nil
}

// RunConsentAcquisition -
func (r *TestCaseRunner) RunConsentAcquisition(tokenName string, permissionList string, ctx *model.Context, consentType string) error {
	r.runningLock.Lock()
	defer r.runningLock.Unlock()
	if r.running {
		return errors.New("consent acquisition test cases runner already running")
	}
	r.running = true

	go r.runConsentAcquisitionAsync(tokenName, permissionList, ctx, consentType)

	return nil
}

func (r *TestCaseRunner) runTestCasesAsync(ctx *model.Context) {

	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx(ctx)

	ctxLogger := r.logger.WithField("id", uuid.New())
	r.AppMsg("runTestCasesAsync()")
	r.AppMsg("determine Token Requirements")

	r.AppMsg("running async test")
	for _, spec := range r.definition.TestCaseRun.TestCases {
		r.executeSpecTests(spec, ruleCtx, ctxLogger)
	}
	r.setNotRunning()
}

func (r *TestCaseRunner) runConsentAcquisitionAsync(tokenName string, permissionList string, ctx *model.Context, consentType string) {

	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx(ctx)
	ruleCtx.PutString("consent_id", tokenName)
	ruleCtx.PutString("permission_list", permissionList)

	ctxLogger := r.logger.WithField("id", uuid.New())
	r.AppMsg("runTokenAcquisitionAsync)")
	var err error
	var comp model.Component
	if consentType == "psu" {
		comp, err = model.LoadComponent("PSUConsentProviderComponent.json")
		if err != nil {
			r.AppErr("Load PSU Component Failed: " + err.Error())
			r.setNotRunning()
			return
		}
	} else {
		comp, err = model.LoadComponent("headlessTokenProviderProviderComponent.json")
		if err != nil {
			r.AppErr("Load HeadlessConsent Component Failed: " + err.Error())
			r.setNotRunning()
			return
		}
	}

	err = comp.ValidateParameters(ruleCtx) // correct parameters for component exist in context
	if err != nil {
		msg := fmt.Sprintf("component execution error: component (%s) cannot ValidateParameters: %s", comp.Name, err.Error())
		r.AppErr(msg)
		r.setNotRunning()
		return
	}

	for k, v := range comp.GetTests() {
		v.ProcessReplacementFields(ruleCtx)
		comp.Tests[k] = v
		logrus.Debugln(v.String())
	}
	r.executeComponentTests(&comp, ruleCtx, ctxLogger)

	r.setNotRunning()
}

func (r *TestCaseRunner) executeComponentTests(comp *model.Component, ruleCtx *model.Context, ctxLogger *logrus.Entry) {
	ctxLogger = ctxLogger.WithField("component", comp.Name)
	ctxLogger.Debugln("execute component Tests...")
	logrus.Debug("component:" + comp.Name)
	for _, testcase := range comp.Tests {
		if r.daemonController.ShouldStop() {
			logrus.Debugln("stop component test run received, aborting runner")
			return
		}

		testResult := r.executeTest(testcase, ruleCtx, ctxLogger)
		r.daemonController.Results() <- testResult
	}
}

func (r *TestCaseRunner) setNotRunning() {
	r.runningLock.Lock()
	r.running = false
	r.runningLock.Unlock()
}

func (r *TestCaseRunner) makeRuleCtx(ctx *model.Context) *model.Context {
	ruleCtx := &model.Context{}
	ruleCtx.Put("SigningCert", r.definition.SigningCert)
	for k, v := range *ctx {
		ruleCtx.Put(k, v)
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

// AppMsg - application level trace
func (r *TestCaseRunner) AppMsg(msg string) string {
	tracer.AppMsg("TestCaseRunner", fmt.Sprintf("%s", msg), r.String())
	return msg
}

// AppErr - application level trace error msg
func (r *TestCaseRunner) AppErr(msg string) error {
	tracer.AppErr("TestCaseRunner", fmt.Sprintf("%s", msg), r.String())
	return errors.New(msg)
}

// String - object represetation
func (r *TestCaseRunner) String() string {
	bites, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		// String() doesn't return error but still want to log as error to tracer ...
		return r.AppErr(fmt.Sprintf("error converting TestCaseRunner  %s", err.Error())).Error()
	}
	return string(bites)
}
