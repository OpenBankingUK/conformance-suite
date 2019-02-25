//go:generate mockery -name Journey -inpkg
package server

import (
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	errDiscoveryModelNotSet        = errors.New("error discovery model not set")
	errTestCasesNotGenerated       = errors.New("error test cases not generated")
	errNotFinishedCollectingTokens = errors.New("error not finished collecting tokens")
	errConsentIDAcquisitionFailed  = errors.New("ConsentId acquistion failed")
)

// Journey represents all possible steps for a user test conformance journey
//
// Happy path journey would look like:
// 1. SetCertificates - sets configuration to run test cases
// 2. SetDiscoveryModel - this validates and if successful set this as your discovery model
// 3. TestCases - Generates test cases, generates permission set requirements to run tests and starts a token collector
// 3.1 CollectToken - collects all tokens required to RunTest
// 4. RunTest - Runs triggers a background run on all generated test from previous steps, needs all token to be already collected
// 5. Results - returns a background process control, so we can monitor on finished tests
//
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() (generation.TestCasesRun, error)
	CollectToken(code, state, scope string) error
	AllTokenCollected() bool
	RunTests() error
	StopTestRun()
	Results() executors.DaemonController
	SetConfig(config JourneyConfig) error
}

type journey struct {
	generator             generation.Generator
	validator             discovery.Validator
	daemonController      executors.DaemonController
	journeyLock           *sync.Mutex
	testCasesRun          generation.TestCasesRun
	testCasesRunGenerated bool
	collector             executors.TokenCollector
	allCollected          bool
	validDiscoveryModel   *discovery.Model
	context               model.Context
	log                   *logrus.Entry
	config                JourneyConfig
}

// NewJourney creates an instance for a user journey
func NewJourney(logger *logrus.Entry, generator generation.Generator, validator discovery.Validator) *journey {
	return &journey{
		generator:             generator,
		validator:             validator,
		daemonController:      executors.NewBufferedDaemonController(),
		journeyLock:           &sync.Mutex{},
		allCollected:          false,
		testCasesRunGenerated: false,
		context:               model.Context{},
		log:                   logger.WithField("module", "Journey"),
	}
}

func (wj *journey) SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error) {
	failures, err := wj.validator.Validate(discoveryModel)
	if err != nil {
		return nil, errors.Wrap(err, "error setting discovery model")
	}

	if !failures.Empty() {
		return failures, nil
	}

	wj.journeyLock.Lock()
	wj.validDiscoveryModel = discoveryModel
	wj.testCasesRunGenerated = false
	wj.allCollected = false
	wj.journeyLock.Unlock()

	return discovery.NoValidationFailures, nil
}

func (wj *journey) TestCases() (generation.TestCasesRun, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if wj.validDiscoveryModel == nil {
		return generation.TestCasesRun{}, errDiscoveryModelNotSet
	}

	if !wj.testCasesRunGenerated {
		config := generation.GeneratorConfig{
			ClientID:              wj.config.clientID,
			Aud:                   wj.config.authorizationEndpoint,
			ResponseType:          "code id_token",
			Scope:                 "openid accounts",
			AuthorizationEndpoint: wj.config.authorizationEndpoint,
			RedirectURL:           wj.config.redirectURL,
		}
		wj.testCasesRun = wj.generator.GenerateSpecificationTestCases(config, wj.validDiscoveryModel.DiscoveryModel, &wj.context)
		runDefinition := executors.RunDefinition{
			DiscoModel:    wj.validDiscoveryModel,
			TestCaseRun:   wj.testCasesRun,
			SigningCert:   wj.config.certificateSigning,
			TransportCert: wj.config.certificateTransport,
		}
		if wj.validDiscoveryModel.DiscoveryModel.TokenAcquisition == "psu" {
			consentIds, err := executors.InitiationConsentAcquisition(wj.testCasesRun.SpecConsentRequirements, runDefinition, &wj.context)
			if err != nil {
				return generation.TestCasesRun{}, errConsentIDAcquisitionFailed
			}
			if len(consentIds) > 0 {
				wj.collector = executors.NewTokenCollector(consentIds, wj.doneCollectionCallback)
				consentIdsToTestCaseRun(consentIds, &wj.testCasesRun)
				wj.allCollected = false
			} else {
				wj.allCollected = true
			}
		} else {
			wj.allCollected = true
		}
		wj.testCasesRunGenerated = true
	}

	return wj.testCasesRun, nil
}

func (wj *journey) CollectToken(code, state, scope string) error {
	logrus.Debugf("state: %s, code: %s", state, code)
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		TestCaseRun:   wj.testCasesRun,
		SigningCert:   wj.config.certificateSigning,
		TransportCert: wj.config.certificateTransport,
	}
	accessToken, err := executors.ExchangeCodeForAccessToken(state, code, scope, runDefinition, &wj.context)
	if err != nil {
		return err
	}
	wj.context.PutString(state, accessToken)
	if state == "to1001" {
		logrus.Warnf("Setting 'access_token' to %s", accessToken)
		wj.context.PutString("access_token", accessToken) // tmp measure to get testcases running
	}
	return wj.collector.Collect(state, accessToken)
}

func (wj *journey) AllTokenCollected() bool {
	// wj.journeyLock.Lock()
	// defer wj.journeyLock.Unlock()
	logrus.Debugf("All tokens collected %t", wj.allCollected)
	return wj.allCollected
}

func (wj *journey) doneCollectionCallback() {
	//wj.journeyLock.Lock()
	logrus.Debug("Setting wj.allCollection=true")
	wj.allCollected = true
	//wj.journeyLock.Unlock()
}

func (wj *journey) RunTests() error {
	logrus.Debug("RunTests ...")
	//wj.journeyLock.Lock()
	//defer wj.journeyLock.Unlock()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	if !wj.allCollected {
		return errNotFinishedCollectingTokens
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		TestCaseRun:   wj.testCasesRun,
		SigningCert:   wj.config.certificateSigning,
		TransportCert: wj.config.certificateTransport,
	}

	runner := executors.NewTestCaseRunner(runDefinition, wj.daemonController)
	logrus.Debug("runTestCases with context ...")
	return runner.RunTestCases(&wj.context)
}

func (wj *journey) Results() executors.DaemonController {
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.daemonController.Stop()
}

type JourneyConfig struct {
	certificateSigning    authentication.Certificate
	certificateTransport  authentication.Certificate
	clientID              string
	clientSecret          string
	tokenEndpoint         string
	authorizationEndpoint string
	resourceBaseURL       string
	xXFAPIFinancialID     string
	issuer                string
	redirectURL           string
}

func (wj *journey) SetConfig(config JourneyConfig) error {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	wj.config = config
	err := wj.configParametersToJourneyContext()
	if err != nil {
		return err
	}
	wj.customTestParametersToJourneyContext()
	return nil
}

const ctxConstClientID = "client_id"
const ctxConstClientSecret = "client_secret"
const ctxConstTokenEndpoint = "token_endpoint"
const ctxConstFapiFinancialID = "fapi_financial_id"
const ctxConstRedirectURL = "redirect_url"
const ctxConstAuthorisationEndpoint = "authorisation_endpoint"
const ctxConstBasicAuthentication = "basic_authentication"
const ctxConstResourceBaseURL = "resource_server"
const ctxConstIssuer = "issuer"

func (wj *journey) configParametersToJourneyContext() error {
	wj.context.PutString(ctxConstClientID, wj.config.clientID)
	wj.context.PutString(ctxConstClientSecret, wj.config.clientSecret)
	wj.context.PutString(ctxConstTokenEndpoint, wj.config.tokenEndpoint)
	wj.context.PutString(ctxConstFapiFinancialID, wj.config.xXFAPIFinancialID)
	wj.context.PutString(ctxConstRedirectURL, wj.config.redirectURL)
	wj.context.PutString(ctxConstAuthorisationEndpoint, wj.config.authorizationEndpoint)
	wj.context.PutString(ctxConstResourceBaseURL, wj.config.resourceBaseURL)

	basicauth, err := authentication.CalculateClientSecretBasicToken(wj.config.clientID, wj.config.clientSecret)
	if err != nil {
		return err
	}
	wj.context.PutString(ctxConstBasicAuthentication, basicauth)

	wj.context.PutString(ctxConstIssuer, wj.config.issuer)

	wj.context.DumpContext("configParameters - dumpcontext")
	return nil
}

func (wj *journey) customTestParametersToJourneyContext() {
	if wj.validDiscoveryModel == nil {
		return
	}
	for _, customTest := range wj.validDiscoveryModel.DiscoveryModel.CustomTests { // assume ordering is prerun i.e. customtest run before other tests
		for k, v := range customTest.Replacements {
			wj.context.PutString(k, v)
		}
	}
}

func consentIdsToTestCaseRun(consentIds []executors.TokenConsentIDItem, testCasesRun *generation.TestCasesRun) {
	for _, v := range testCasesRun.SpecConsentRequirements {
		for x, permission := range v.NamedPermissions {
			for _, item := range consentIds {
				if item.TokenName == permission.Name {
					permission.ConsentUrl = item.ConsentURL
					logrus.Debugf("Setting consent url for token %s to %s", permission.Name, permission.ConsentUrl)
					v.NamedPermissions[x] = permission
				}
			}
		}
	}
}
