//go:generate mockery -name Journey -inpkg
package server

import (
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	errDiscoveryModelNotSet        = errors.New("error discovery model not set")
	errTestCasesNotGenerated       = errors.New("error test cases not generated")
	errNotFinishedCollectingTokens = errors.New("error not finished collecting tokens")
	errConsentIDAcquisitionFailed  = errors.New("ConsentId acquistion failed")
	errTokenMappingFailed          = errors.New("token mapping to testcases ailed")
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
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() (generation.TestCasesRun, error)
	CollectToken(code, state, scope string) error
	AllTokenCollected() bool
	RunTests() error
	StopTestRun()
	NewDaemonController()
	Results() executors.DaemonController
	SetConfig(config JourneyConfig) error
	Events() events.Events
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
	events                events.Events
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
		events:                events.NewEvents(),
	}
}

// NewDaemonController - calls StopTestRun and then sets new daemonController
// and new events on journey.
// This is a solution to prevent events being sent to a disconnected
// websocket instead of new websocket after the client reconnects.
func (wj *journey) NewDaemonController() {
	wj.StopTestRun()

	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	wj.daemonController = executors.NewBufferedDaemonController()
	wj.events = events.NewEvents()
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
	defer wj.journeyLock.Unlock()
	wj.validDiscoveryModel = discoveryModel
	wj.testCasesRunGenerated = false
	wj.allCollected = false

	return discovery.NoValidationFailures, nil
}

func (wj *journey) TestCases() (generation.TestCasesRun, error) {
	wj.journeyLock.Lock()
	defer func() {
		wj.journeyLock.Unlock()
	}()

	if wj.validDiscoveryModel == nil {
		return generation.TestCasesRun{}, errDiscoveryModelNotSet
	}

	if !wj.testCasesRunGenerated {
		config := wj.makeGeneratorConfig()
		discovery := wj.validDiscoveryModel.DiscoveryModel
		//TODO: Remove standard testcase generation
		//wj.testCasesRun = wj.generator.GenerateSpecificationTestCases(wj.log, config, discovery, &wj.context)

		wj.log.Debugln("Journey:GenerationManifestTests")
		//tcrun2 := wj.generator.GenerateManifestTests(wj.log, config, discovery, &wj.context) // Integration work in progress
		wj.testCasesRun = wj.generator.GenerateManifestTests(wj.log, config, discovery, &wj.context) // Integration work in progress

		// wj.testCasesRun.SpecConsentRequirements = append(wj.testCasesRun.SpecConsentRequirements, tcrun2.SpecConsentRequirements...)
		// wj.testCasesRun.TestCases = append(wj.testCasesRun.TestCases, tcrun2.TestCases...)

		//wj.testCasesRun.SpecConsentRequirements = tcrun2.SpecConsentRequirements
		//wj.testCasesRun.TestCases = tcrun2.TestCases

		if discovery.TokenAcquisition == "psu" {
			wj.log.Debugln("Journey:AcquirePSUTokens")
			definition := wj.makeRunDefinition()

			consentIds, tokenMap, err := executors.InitiationConsentAcquisition(wj.testCasesRun.SpecConsentRequirements, definition, &wj.context, &wj.testCasesRun)
			if err != nil {
				return generation.TestCasesRun{}, errors.WithMessage(errConsentIDAcquisitionFailed, err.Error())
			}
			for k, v := range tokenMap {
				wj.context.PutString(k, v)
			}
			wj.createTokenCollector(consentIds)
		} else {
			wj.log.Debugln("Journey:AcquireHeadlessTokens")
			runDefinition := wj.makeRunDefinition()
			// TODO:Process multiple specs ... don't restrict to element [0]!!
			tokenPermissionsMap, err := executors.AcquireHeadlessTokens(wj.testCasesRun.TestCases[0].TestCases, &wj.context, runDefinition)
			if err != nil {
				return generation.TestCasesRun{}, errConsentIDAcquisitionFailed
			}
			// TODO:Process multipe specs
			tokenMap := manifest.MapTokensToTestCases(tokenPermissionsMap, wj.testCasesRun.TestCases[0].TestCases)
			for k, v := range tokenMap {
				wj.context.PutString(k, v)
			}

			wj.allCollected = true
		}
		wj.testCasesRunGenerated = true
	}

	wj.log.Tracef("TestCaseRun.SpecConsentRequirements: %#v\n", wj.testCasesRun.SpecConsentRequirements)
	for k := range wj.testCasesRun.TestCases {
		wj.log.Tracef("TestCaseRun-Specificatino: %#v\n", wj.testCasesRun.TestCases[k].Specification)
	}
	wj.log.Tracef("Dumping Consents:---------------------------\n")
	for _, v := range wj.testCasesRun.SpecConsentRequirements {
		wj.log.Tracef("%s", v.Identifier)
		for _, x := range v.NamedPermissions {
			wj.log.Tracef("\tname: %s codeset: %#v\n\tconsent Url: %s", x.Name, x.CodeSet.CodeSet, x.ConsentUrl)
		}
	}
	return wj.testCasesRun, nil
}

func (wj *journey) CollectToken(code, state, scope string) error {
	wj.journeyLock.Lock()
	defer func() {
		wj.journeyLock.Unlock()
	}()

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	runDefinition := wj.makeRunDefinition()
	accessToken, err := executors.ExchangeCodeForAccessToken(state, code, scope, runDefinition, &wj.context)
	if err != nil {
		return err
	}

	wj.context.PutString(state, accessToken)
	if state == "Token001" {
		wj.log.Warnf("Setting 'access_token' to %s", accessToken)
		wj.context.PutString("access_token", accessToken) // tmp measure to get testcases running
	}

	return wj.collector.Collect(state, accessToken)
}

func (wj *journey) AllTokenCollected() bool {
	wj.log.Debugf("All tokens collected %t", wj.allCollected)
	return wj.allCollected
}

func (wj *journey) doneCollectionCallback() {
	wj.log.Debug("Setting wj.allCollection=true")
	wj.allCollected = true
}

func (wj *journey) RunTests() error {
	wj.log.Debug("RunTests ...")

	if !wj.testCasesRunGenerated {
		return errTestCasesNotGenerated
	}

	if !wj.allCollected {
		return errNotFinishedCollectingTokens
	}

	requiredTokens := []manifest.RequiredTokens{}
	// map tokens to Testcases
	var err error
	for _, tests := range wj.testCasesRun.TestCases {
		requiredTokens, err = manifest.GetRequiredTokensFromTests(tests.TestCases)
		if err != nil {
			wj.log.Warn("Testcase Token setup failed")
			return errTokenMappingFailed
		}
	}

	//TODO Extend to cover more that one set of testcases
	manifest.MapTokensToTestCases(requiredTokens, wj.testCasesRun.TestCases[0].TestCases)

	runDefinition := wj.makeRunDefinition()
	runner := executors.NewTestCaseRunner(wj.log, runDefinition, wj.daemonController)
	wj.log.Debug("runTestCases with context ...")
	return runner.RunTestCases(&wj.context)
}

func (wj *journey) Results() executors.DaemonController {
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.daemonController.Stop()
}

func (wj *journey) createTokenCollector(consentIds executors.TokenConsentIDs) {
	if len(consentIds) > 0 {
		wj.collector = executors.NewTokenCollector(wj.log, consentIds, wj.doneCollectionCallback, wj.events)
		consentIdsToTestCaseRun(wj.log, consentIds, &wj.testCasesRun)

		wj.allCollected = false
	} else {
		wj.allCollected = true
	}
}

func (wj *journey) makeGeneratorConfig() generation.GeneratorConfig {
	return generation.GeneratorConfig{
		ClientID:              wj.config.clientID,
		Aud:                   wj.config.authorizationEndpoint,
		ResponseType:          "code id_token",
		Scope:                 "openid accounts",
		AuthorizationEndpoint: wj.config.authorizationEndpoint,
		RedirectURL:           wj.config.redirectURL,
		ResourceIDs:           wj.config.resourceIDs,
	}
}

func (wj *journey) makeRunDefinition() executors.RunDefinition {
	return executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		TestCaseRun:   wj.testCasesRun,
		SigningCert:   wj.config.certificateSigning,
		TransportCert: wj.config.certificateTransport,
	}
}

type JourneyConfig struct {
	certificateSigning      authentication.Certificate
	certificateTransport    authentication.Certificate
	clientID                string
	clientSecret            string
	tokenEndpoint           string
	tokenEndpointAuthMethod string
	authorizationEndpoint   string
	resourceBaseURL         string
	xXFAPIFinancialID       string
	issuer                  string
	redirectURL             string
	resourceIDs             model.ResourceIDs
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

func (wj *journey) Events() events.Events {
	return wj.events
}

const ctxConstClientID = "client_id"
const ctxConstClientSecret = "client_secret"
const ctxConstTokenEndpoint = "token_endpoint"
const ctxConstTokenEndpointAuthMethod = "token_endpoint_auth_method"
const ctxConstFapiFinancialID = "x-fapi-financial-id"
const ctxConstRedirectURL = "redirect_url"
const ctxConstAuthorisationEndpoint = "authorisation_endpoint"
const ctxConstBasicAuthentication = "basic_authentication"
const ctxConstResourceBaseURL = "resource_server"
const ctxConstIssuer = "issuer"

func (wj *journey) configParametersToJourneyContext() error {
	wj.context.PutString(ctxConstClientID, wj.config.clientID)
	wj.context.PutString(ctxConstClientSecret, wj.config.clientSecret)
	wj.context.PutString(ctxConstTokenEndpoint, wj.config.tokenEndpoint)
	wj.context.PutString(ctxConstTokenEndpointAuthMethod, wj.config.tokenEndpointAuthMethod)
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
	return nil
}

func (wj *journey) customTestParametersToJourneyContext() {
	if wj.validDiscoveryModel == nil {
		return
	}

	// assume ordering is prerun i.e. customtest run before other tests
	for _, customTest := range wj.validDiscoveryModel.DiscoveryModel.CustomTests {
		for k, v := range customTest.Replacements {
			wj.context.PutString(k, v)
		}
	}
}

func consentIdsToTestCaseRun(log *logrus.Entry, consentIds []executors.TokenConsentIDItem, testCasesRun *generation.TestCasesRun) {
	for _, v := range testCasesRun.SpecConsentRequirements {
		for x, permission := range v.NamedPermissions {
			for _, item := range consentIds {
				if item.TokenName == permission.Name {
					permission.ConsentUrl = item.ConsentURL
					log.Debugf("Setting consent url for token %s to %s", permission.Name, permission.ConsentUrl)
					v.NamedPermissions[x] = permission
				}
			}
		}
	}
}
