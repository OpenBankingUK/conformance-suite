//go:generate mockery -name Journey
package server

import (
	"os"
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
	CollectToken(setName, token string) error
	AllTokenCollected() bool
	RunTests() error
	StopTestRun()
	Results() executors.DaemonController
	SetConfig(signing, transport authentication.Certificate, clientID, clientSecret, tokenEndpoint, authorizationEndpoint, xXFAPIFinancialID, redirectURL string)
}

type journey struct {
	generator             generation.Generator
	validator             discovery.Validator
	daemonController      executors.DaemonController
	journeyLock           *sync.Mutex
	testCasesRun          generation.TestCasesRun
	testCasesRunGenerated bool
	collector             executors.Collector
	allCollected          bool
	validDiscoveryModel   *discovery.Model
	certificateSigning    authentication.Certificate
	certificateTransport  authentication.Certificate
	context               model.Context
	log                   *logrus.Entry
	clientID              string
	clientSecret          string
	tokenEndpoint         string
	authorizationEndpoint string
	xXFAPIFinancialID     string
	redirectURL           string
}

// NewJourney creates an instance for a user journey
func NewJourney(generator generation.Generator, validator discovery.Validator) *journey {
	return &journey{
		generator:             generator,
		validator:             validator,
		daemonController:      executors.NewBufferedDaemonController(),
		journeyLock:           &sync.Mutex{},
		allCollected:          false,
		testCasesRunGenerated: false,
		context:               model.Context{},
		log:                   (&logrus.Logger{Out: os.Stderr, Formatter: new(logrus.TextFormatter), Hooks: make(logrus.LevelHooks), Level: logrus.DebugLevel}).WithField("module", "Journey"),
	}
}

func (wj *journey) SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error) {
	wj.log.Debug("wj.SetDiscoveryModel")
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
	wj.log.Debug("wj.TestCases - Generate Test Cases")
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if wj.validDiscoveryModel == nil {
		return generation.TestCasesRun{}, errDiscoveryModelNotSet
	}

	if !wj.testCasesRunGenerated {
		wj.testCasesRun = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
		// replace this with a NewCollector to a real implementation
		runDefinition := executors.RunDefinition{
			DiscoModel:    wj.validDiscoveryModel,
			TestCaseRun:   wj.testCasesRun,
			SigningCert:   wj.certificateSigning,
			TransportCert: wj.certificateTransport,
			Context:       &wj.context,
		}
		wj.customTestParametersToJourneyContext()
		if wj.validDiscoveryModel.DiscoveryModel.TokenAcquisition == "psu" {
			executors.InitiationConsentAcquisition(wj.testCasesRun.SpecConsentRequirements, runDefinition)
			wj.collector = executors.NewNullCollector(wj.doneCollectionCallback)
		}
		wj.testCasesRunGenerated = true
		wj.allCollected = false
	}

	return wj.testCasesRun, nil
}

func (wj *journey) CollectToken(setName, token string) error {
	wj.log.Debug("wj.CollectToken")
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	return wj.collector.Collect(setName, token)
}

func (wj *journey) AllTokenCollected() bool {
	wj.log.Debug("wj.AllTokensCollected")
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	return wj.allCollected
}

func (wj *journey) doneCollectionCallback() {
	wj.log.Debug("wj.doneCollection Callback")
	wj.journeyLock.Lock()
	wj.allCollected = true
	wj.journeyLock.Unlock()
}

func (wj *journey) RunTests() error {
	wj.log.Debug("wj.RunTests")
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	if !wj.allCollected {
		return errNotFinishedCollectingTokens
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		TestCaseRun:   wj.testCasesRun,
		SigningCert:   wj.certificateSigning,
		TransportCert: wj.certificateTransport,
		Context:       &wj.context,
	}

	runner := executors.NewTestCaseRunner(runDefinition, wj.daemonController)
	return runner.RunTestCases()
}

func (wj *journey) Results() executors.DaemonController {
	wj.log.Debug("wj.Results")
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.log.Debug("wj.StopTestRun")
	wj.daemonController.Stop()
}

func (wj *journey) SetConfig(signing, transport authentication.Certificate, clientID, clientSecret, tokenEndpoint, authorizationEndpoint, xXFAPIFinancialID, redirectURL string) {
	wj.log.Debug("wj.SetConfig")
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	wj.certificateSigning = signing
	wj.certificateTransport = transport
	wj.clientID = clientID
	wj.clientSecret = clientSecret
	wj.tokenEndpoint = tokenEndpoint
	wj.authorizationEndpoint = authorizationEndpoint
	wj.xXFAPIFinancialID = xXFAPIFinancialID
	wj.redirectURL = redirectURL

	wj.configParametersToJourneyContext()
}

const ctxConstClientID = "client_id"
const ctxConstClientSecret = "client_secret"
const ctxConstTokenEndpoint = "token_endpoint"
const ctxConstFapiFinancialID = "fapi_financial_id"
const ctxConstRedirectURL = "redirect_url"
const ctxConstAuthorisationEndpoint = "authorisation_endpoint"
const ctxConstBasicAuthentication = "basic_authentication"

func (wj *journey) configParametersToJourneyContext() error {
	wj.context.PutString(ctxConstClientID, wj.clientID)
	wj.context.PutString(ctxConstClientSecret, wj.clientSecret)
	wj.context.PutString(ctxConstTokenEndpoint, wj.tokenEndpoint)
	wj.context.PutString(ctxConstFapiFinancialID, wj.xXFAPIFinancialID)
	wj.context.PutString(ctxConstRedirectURL, wj.redirectURL)
	wj.context.PutString(ctxConstAuthorisationEndpoint, wj.authorizationEndpoint)

	basicauth, err := authentication.CalculateClientSecretBasicToken(wj.clientID, wj.clientSecret)
	if err != nil {
		return err
	}
	wj.context.PutString(ctxConstBasicAuthentication, basicauth)
	return nil
}

func (wj *journey) customTestParametersToJourneyContext() {
	for _, customTest := range wj.validDiscoveryModel.DiscoveryModel.CustomTests { // assume ordering is prerun i.e. customtest run before other tests
		for k, v := range customTest.Replacements {
			wj.context.PutString(k, v)
		}
	}
}
