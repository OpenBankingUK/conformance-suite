//go:generate mockery -name Journey
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
	CollectToken(setName, token, scope string) error
	AllTokenCollected() bool
	RunTests() error
	StopTestRun()
	Results() executors.DaemonController
	SetConfig(signing, transport authentication.Certificate, clientID, clientSecret, tokenEndpoint, authorizationEndpoint, resourceBaseURL, xXFAPIFinancialID, redirectURL string)
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
	certificateSigning    authentication.Certificate
	certificateTransport  authentication.Certificate
	context               model.Context
	log                   *logrus.Entry
	clientID              string
	clientSecret          string
	tokenEndpoint         string
	authorizationEndpoint string
	resourceBaseURL       string
	xXFAPIFinancialID     string
	redirectURL           string
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
			ClientID:              wj.clientID,
			Aud:                   wj.authorizationEndpoint,
			ResponseType:          "code id_token",
			Scope:                 "openid accounts",
			AuthorizationEndpoint: wj.authorizationEndpoint,
			RedirectURL:           wj.redirectURL,
		}
		wj.testCasesRun = wj.generator.GenerateSpecificationTestCases(config, wj.validDiscoveryModel.DiscoveryModel, &wj.context)
		runDefinition := executors.RunDefinition{
			DiscoModel:    wj.validDiscoveryModel,
			TestCaseRun:   wj.testCasesRun,
			SigningCert:   wj.certificateSigning,
			TransportCert: wj.certificateTransport,
		}
		if wj.validDiscoveryModel.DiscoveryModel.TokenAcquisition == "psu" {
			consentIds, err := executors.InitiationConsentAcquisition(wj.testCasesRun.SpecConsentRequirements, runDefinition, &wj.context)
			if err != nil {
				return generation.TestCasesRun{}, errConsentIDAcquisitionFailed
			}
			if len(consentIds) > 0 {
				wj.collector = executors.NewTokenCollector(consentIds, wj.doneCollectionCallback)
				consentIdsToTestCaseRun(consentIds, &wj.testCasesRun)
			} else {
				wj.allCollected = true
			}
		}
		wj.testCasesRunGenerated = true
		wj.allCollected = false
	}

	return wj.testCasesRun, nil
}

func (wj *journey) CollectToken(tokenName, code, scope string) error {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		TestCaseRun:   wj.testCasesRun,
		SigningCert:   wj.certificateSigning,
		TransportCert: wj.certificateTransport,
	}
	accessToken, err := executors.ExchangeCodeForAccessToken(tokenName, code, scope, runDefinition, &wj.context)
	if err != nil {
		return err
	}

	return wj.collector.Collect(tokenName, accessToken)
}

func (wj *journey) AllTokenCollected() bool {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	return wj.allCollected
}

func (wj *journey) doneCollectionCallback() {
	wj.journeyLock.Lock()
	wj.allCollected = true
	wj.journeyLock.Unlock()
}

func (wj *journey) RunTests() error {
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
	}

	runner := executors.NewTestCaseRunner(runDefinition, wj.daemonController)
	return runner.RunTestCases(&wj.context)
}

func (wj *journey) Results() executors.DaemonController {
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.daemonController.Stop()
}

func (wj *journey) SetConfig(signing, transport authentication.Certificate, clientID, clientSecret, tokenEndpoint, authorizationEndpoint, resourceBaseURL, xXFAPIFinancialID, redirectURL string) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	wj.certificateSigning = signing
	wj.certificateTransport = transport
	wj.clientID = clientID
	wj.clientSecret = clientSecret
	wj.tokenEndpoint = tokenEndpoint
	wj.authorizationEndpoint = authorizationEndpoint
	wj.resourceBaseURL = resourceBaseURL
	wj.xXFAPIFinancialID = xXFAPIFinancialID
	wj.redirectURL = redirectURL
	wj.configParametersToJourneyContext()
	wj.customTestParametersToJourneyContext()
}

const ctxConstClientID = "client_id"
const ctxConstClientSecret = "client_secret"
const ctxConstTokenEndpoint = "token_endpoint"
const ctxConstFapiFinancialID = "fapi_financial_id"
const ctxConstRedirectURL = "redirect_url"
const ctxConstAuthorisationEndpoint = "authorisation_endpoint"
const ctxConstBasicAuthentication = "basic_authentication"
const ctxConstResourceBaseURL = "resource_base_url"

func (wj *journey) configParametersToJourneyContext() error {
	wj.context.PutString(ctxConstClientID, wj.clientID)
	wj.context.PutString(ctxConstClientSecret, wj.clientSecret)
	wj.context.PutString(ctxConstTokenEndpoint, wj.tokenEndpoint)
	wj.context.PutString(ctxConstFapiFinancialID, wj.xXFAPIFinancialID)
	wj.context.PutString(ctxConstFapiFinancialID, wj.resourceBaseURL) // tmp mapping fix
	wj.context.PutString(ctxConstRedirectURL, wj.redirectURL)
	wj.context.PutString(ctxConstAuthorisationEndpoint, wj.authorizationEndpoint)
	wj.context.PutString(ctxConstResourceBaseURL, wj.resourceBaseURL)
	wj.context.PutString(ctxConstResourceBaseURL, wj.xXFAPIFinancialID) // tmp mapping fix
	basicauth, err := authentication.CalculateClientSecretBasicToken(wj.clientID, wj.clientSecret)
	if err != nil {
		return err
	}
	wj.context.PutString(ctxConstBasicAuthentication, basicauth)
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
