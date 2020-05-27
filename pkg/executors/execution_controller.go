package executors

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schemaprops"

	"gopkg.in/resty.v1"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
)

// RunDefinition captures all the information required to run the test cases
type RunDefinition struct {
	DiscoModel    *discovery.Model
	SpecRun       generation.SpecRun
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

type TestCaseRunner struct {
	executor         TestCaseExecutor
	definition       RunDefinition
	daemonController DaemonController
	logger           *logrus.Entry
	runningLock      *sync.Mutex
	running          bool
}

// NewTestCaseRunner -
func NewTestCaseRunner(logger *logrus.Entry, definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logger.WithField("module", "TestCaseRunner"),
		runningLock:      &sync.Mutex{},
		running:          false,
	}
}

// NewConsentAcquisitionRunner -
func NewConsentAcquisitionRunner(logger *logrus.Entry, definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logger.WithField("module", "ConsentAcquisitionRunner"),
		runningLock:      &sync.Mutex{},
		running:          false,
	}
}

// NewExchangeComponentRunner -
func NewExchangeComponentRunner(definition RunDefinition, daemonController DaemonController) *TestCaseRunner {
	return &TestCaseRunner{
		executor:         NewExecutor(),
		definition:       definition,
		daemonController: daemonController,
		logger:           logrus.StandardLogger().WithField("module", "ExchangeComponent"),
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
func (r *TestCaseRunner) RunConsentAcquisition(item TokenConsentIDItem, ctx *model.Context, consentType string, consentIDChannel chan<- TokenConsentIDItem) error {
	r.runningLock.Lock()
	defer r.runningLock.Unlock()
	if r.running {
		return errors.New("consent acquisition test cases runner already running")
	}
	r.running = true
	logrus.Tracef("runConsentAquisition with %s, %s, %s", item.TokenName, item.ConsentURL, item.Permissions)
	go r.runConsentAcquisitionAsync(item, ctx, consentType, consentIDChannel)

	return nil
}

func (r *TestCaseRunner) runTestCasesAsync(ctx *model.Context) {
	err := r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	if err != nil {
		r.logger.WithError(err).Error("running test cases async")
	}

	ruleCtx := r.makeRuleCtx(ctx)

	ctxLogger := r.logger.WithField("id", uuid.New())
	for _, spec := range r.definition.SpecRun.SpecTestCases {
		r.executeSpecTests(spec, ruleCtx, ctxLogger) // Run Tests for each spec
	}

	collector := schemaprops.GetPropertyCollector()
	r.daemonController.AddResponseFields(collector.OutputJSON())

	r.daemonController.SetCompleted()

	r.setNotRunning()
}

func (r *TestCaseRunner) runConsentAcquisitionAsync(item TokenConsentIDItem, ctx *model.Context, consentType string, consentIDChannel chan<- TokenConsentIDItem) {
	err := r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	if err != nil {
		r.logger.WithError(err).Error("running consent acquisition async")
	}

	ruleCtx := r.makeRuleCtx(ctx)
	ruleCtx.PutString("consent_id", item.TokenName)
	ruleCtx.PutString("token_name", item.TokenName)
	ruleCtx.PutString("permission_list", item.Permissions)

	ctxLogger := r.logger.WithField("id", uuid.New())
	var comp model.Component

	// Check for MTLS vs client basic authentication
	authMethod, err := ctx.GetString("token_endpoint_auth_method")
	if err != nil {
		authMethod = authentication.ClientSecretBasic
	}

	if consentType == "psu" {
		comp, err = model.LoadComponent("PSUConsentProviderComponent.json")
		if err != nil {
			r.AppMsg("Load PSU Component Failed: " + err.Error())
			r.setNotRunning()
			return
		}

	} else {
		comp, err = model.LoadComponent("headlessTokenProviderProviderComponent.json")
		if err != nil {
			r.AppMsg("Load HeadlessConsent Component Failed: " + err.Error())
			r.setNotRunning()
			return
		}
	}

	err = comp.ValidateParameters(ruleCtx) // correct parameters for component exist in context
	if err != nil {
		msg := fmt.Sprintf("component execution error: component (%s) cannot ValidateParameters: %s", comp.Name, err.Error())
		r.AppMsg(msg)
		r.setNotRunning()
		return
	}

	for k, v := range comp.GetTests() {
		v.ProcessReplacementFields(ruleCtx, true)
		v.Validator = schema.NewNullValidator()
		comp.Tests[k] = v
	}

	r.executeComponentTests(&comp, ruleCtx, ctxLogger, item, consentIDChannel, authMethod)
	clientGrantToken, err := ruleCtx.GetString("client_access_token")
	if err == nil {
		logrus.StandardLogger().WithFields(logrus.Fields{
			"clientGrantToken": clientGrantToken,
		}).Debugf("Setting client_access_token")
		ctx.PutString("client_access_token", clientGrantToken)
	}

	r.setNotRunning()
}

func (r *TestCaseRunner) executeComponentTests(comp *model.Component, ruleCtx *model.Context, logger *logrus.Entry, item TokenConsentIDItem, consentIDChannel chan<- TokenConsentIDItem, authMethod string) {
	ctxLogger := logger.WithFields(logrus.Fields{
		"component": comp.Name,
		"module":    "TestCaseRunner",
		"function":  "executeComponentTests",
	})

	for _, testcase := range comp.Tests {
		if r.daemonController.ShouldStop() {
			ctxLogger.Debug("stop component test run received, aborting runner")
			return
		}

		if testcase.ID == "#compPsuConsent01" {
			switch authMethod {
			case authentication.ClientSecretBasic:
				testcase.Input.SetHeader("authorization", "Basic $basic_authentication")
			case authentication.TlsClientAuth:
				clientid, err := ruleCtx.GetString("client_id")
				if err != nil {
					ctxLogger.WithFields(logrus.Fields{
						"authMethod": authMethod,
						"err":        err,
					}).Error("cannot locate client_id to populate form field")
					continue
				}

				testcase.Input.SetFormField("client_id", clientid)
			case authentication.PrivateKeyJwt:
				clientID, err := ruleCtx.GetString("client_id")
				if err != nil {
					ctxLogger.WithFields(logrus.Fields{
						"authMethod": authMethod,
						"err":        err,
					}).Error("cannot locate client_id to populate form field")
					continue
				}

				tokenEndpoint, err := ruleCtx.GetString("token_endpoint")
				if err != nil {
					ctxLogger.WithFields(logrus.Fields{
						"authMethod": authMethod,
						"err":        err,
					}).Error("cannot locate token_endpoint to populate form field")
				}

				if testcase.Input.Claims == nil {
					testcase.Input.Claims = map[string]string{}
				}

				// https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication
				// iss
				// REQUIRED. Issuer. This MUST contain the client_id of the OAuth Client.
				// sub
				// REQUIRED. Subject. This MUST contain the client_id of the OAuth Client.
				// aud
				// REQUIRED. Audience. The aud (audience) Claim. Value that identifies the Authorization Server as an intended audience. The Authorization Server MUST verify that it is an intended audience for the token. The Audience SHOULD be the URL of the Authorization Server's Token Endpoint.
				testcase.Input.Claims["iss"] = clientID
				testcase.Input.Claims["sub"] = clientID
				testcase.Input.Claims["aud"] = tokenEndpoint
				clientAssertion, err := testcase.Input.GenerateRequestToken(ruleCtx)
				if err != nil {
					ctxLogger.WithFields(logrus.Fields{
						"testcase": testcase,
						"err":      err,
					}).Error("failed on testcase.Input.GenerateRequestToken")
					continue
				}

				testcase.Input.SetFormField(authentication.ClientAssertionType, authentication.ClientAssertionTypeValue)
				testcase.Input.SetFormField(authentication.ClientAssertion, clientAssertion)
			default:
				ctxLogger.WithFields(logrus.Fields{
					"authMethod": authMethod,
				}).Error("Unsupported token_endpoint_auth_method")
				continue
			}
		}

		testResult := r.executeTest(testcase, ruleCtx, logger)
		r.daemonController.AddResult(testResult)

		if testResult.Pass {
			ctxLogger.WithFields(logrus.Fields{
				"item": fmt.Sprintf("%#v", item),
			}).Debug("hanging around for token (TokenConsentIDItem)")
			consentURL, err := ruleCtx.GetString("consent_url")
			if err == model.ErrNotFound {
				continue
			}

			item.ConsentURL = consentURL
			ruleCtx.DumpContext()
			consentID, err := ruleCtx.GetString(item.TokenName)
			if err == model.ErrNotFound {
				ctxLogger.WithFields(logrus.Fields{
					"item.TokenName": fmt.Sprintf("%+v", item.TokenName),
					"err":            err,
				}).Warn("Did not find consentID in context for item.TokenName")
			}
			item.ConsentID = consentID

			ctxLogger.WithFields(logrus.Fields{
				"item": fmt.Sprintf("%#v", item),
			}).Debug("Sending item (TokenConsentIDItem) to consentIDChannel")
			consentIDChannel <- item
		} else if len(testResult.Fail) > 0 {
			item.Error = testResult.Fail[0]
			consentIDChannel <- item
		}
	}
}

func (r *TestCaseRunner) setNotRunning() {
	logger := logrus.StandardLogger().WithFields(logrus.Fields{
		"function": "setNotRunning",
		"module":   "TestCaseRunner",
	})

	logger.Debug("acquiring runningLock")
	r.runningLock.Lock()
	logger.Debug("acquired runningLock")
	defer func() {
		logger.Debug("releasing runningLock")
		r.runningLock.Unlock()

	}()
	r.running = false
}

func (r *TestCaseRunner) makeRuleCtx(ctx *model.Context) *model.Context {
	ruleCtx := &model.Context{}
	ruleCtx.PutContext(ctx)
	return ruleCtx
}

func (r *TestCaseRunner) executeSpecTests(spec generation.SpecificationTestCases, ruleCtx *model.Context, ctxLogger *logrus.Entry) {
	ctxLogger = ctxLogger.WithField("spec", spec.Specification.Name)
	collector := schemaprops.GetPropertyCollector()
	collector.SetCollectorAPIDetails(spec.Specification.Name, spec.Specification.Version)

	for _, testcase := range spec.TestCases {
		if r.daemonController.ShouldStop() {
			ctxLogger.Info("stop test run received, aborting runner")
			return
		}
		ctxLogger = ctxLogger.WithField("ID", testcase.ID)
		ruleCtx.DumpContext("ruleCtx before: " + testcase.ID)
		testResult := r.executeTest(testcase, ruleCtx, ctxLogger)
		r.daemonController.AddResult(testResult)
	}
}

func (r *TestCaseRunner) executeTest(tc model.TestCase, ruleCtx *model.Context, logger *logrus.Entry) results.TestCase {
	ctxLogger := logWithTestCase(logger, tc)
	req, err := tc.Prepare(ruleCtx)
	if err != nil {
		ctxLogger.WithError(err).Error("preparing executing test")
		return results.NewTestCaseFail(tc.ID, results.NoMetrics(), []error{err}, tc.Input.Endpoint, tc.APIName, tc.APIVersion, tc.Detail, tc.RefURI, tc.StatusCode)
	}
	resp, metrics, err := r.executor.ExecuteTestCase(req, &tc, ruleCtx)
	ctxLogger = logWithMetrics(ctxLogger, metrics)
	if err != nil {
		ctxLogger.WithError(err).WithFields(logrus.Fields{"result": "FAIL", "ID": tc.ID}).Error("test result")
		return results.NewTestCaseFail(tc.ID, metrics, []error{err}, tc.Input.Endpoint, tc.APIName, tc.APIVersion, tc.Detail, tc.RefURI, tc.StatusCode)
	}
	tc.StatusCode = resp.Status()
	result, errs := tc.Validate(resp, ruleCtx)
	if errs != nil {
		detailedErrors := detailedErrors(errs, resp)
		ctxLogger.WithField("errs", detailedErrors).WithFields(logrus.Fields{"result": passText()[result], "ID": tc.ID}).Error("test result validate")
		return results.NewTestCaseFail(tc.ID, metrics, detailedErrors, tc.Input.Endpoint, tc.APIName, tc.APIVersion, tc.Detail, tc.RefURI, tc.StatusCode)
	}

	if !result {
		ctxLogger.WithError(err).WithFields(logrus.Fields{"result": passText()[result], "ID": tc.ID}).Error("test result blank")
	} else {
		ctxLogger.WithError(err).WithFields(logrus.Fields{"result": passText()[result], "ID": tc.ID}).Info("test result")
	}

	return results.NewTestCaseResult(tc.ID, result, metrics, []error{}, tc.Input.Endpoint, tc.APIName, tc.APIVersion, tc.Detail, tc.RefURI, tc.StatusCode)
}

type DetailError struct {
	EndpointResponse string `json:"endpointResponse"`
	TestCaseMessage  string `json:"testCaseMessage"`
}

func (de DetailError) Error() string {
	j, _ := json.Marshal(de)

	return string(j)
}

func detailedErrors(errs []error, resp *resty.Response) []error {
	detailedErrors := []error{}
	for _, err := range errs {
		detailedError := DetailError{
			EndpointResponse: string(resp.Body()),
			TestCaseMessage:  err.Error(),
		}
		detailedErrors = append(detailedErrors, detailedError)
	}
	return detailedErrors
}

func passText() map[bool]string {
	return map[bool]string{
		true:  "PASS",
		false: "FAIL",
	}
}

func logWithTestCase(logger *logrus.Entry, tc model.TestCase) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"TestCase.Name":              tc.Name,
		"TestCase.Input.Method":      tc.Input.Method,
		"TestCase.Input.Endpoint":    tc.Input.Endpoint,
		"TestCase.Expect.StatusCode": tc.Expect.StatusCode,
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
	tracer.AppMsg("TestCaseRunner", msg, r.String())
	return msg
}

// AppErr - application level trace error msg
func (r *TestCaseRunner) AppErr(msg string) error {
	tracer.AppErr("TestCaseRunner", msg, r.String())
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
