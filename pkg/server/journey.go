//go:generate mockery -name Journey -inpkg
package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
)

var (
	errDiscoveryModelNotSet            = errors.New("error discovery model not set")
	errTestCasesNotGenerated           = errors.New("error test cases not generated")
	errTestCasesGenerated              = errors.New("error test cases already generated")
	errNotFinishedCollectingTokens     = errors.New("error not finished collecting tokens")
	errConsentIDAcquisitionFailed      = errors.New("ConsentId acquistion failed")
	errDynamicResourceAllocationFailed = errors.New("Dynamic Resource allocation failed")

	dynamicResourceIDs = false
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
	DiscoveryModel() (discovery.Model, error)
	SetFilteredManifests(manifest.Scripts)
	FilteredManifests() (manifest.Scripts, error)
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
	permissions           map[string][]manifest.RequiredTokens
	manifests             []manifest.Scripts
	filteredManifests     manifest.Scripts
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
		log:                   logger.WithField("module", "journey"),
		events:                events.NewEvents(),
		permissions:           make(map[string][]manifest.RequiredTokens),
		manifests:             make([]manifest.Scripts, 0),
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
		return nil, errors.Wrap(err, "journey.SetDiscoveryModel: error setting discovery model")
	}

	if !failures.Empty() {
		return failures, nil
	}

	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	wj.validDiscoveryModel = discoveryModel
	wj.testCasesRunGenerated = false
	wj.allCollected = false

	return discovery.NoValidationFailures(), nil
}

func (wj *journey) DiscoveryModel() (discovery.Model, error) {
	wj.journeyLock.Lock()
	discoveryModel := wj.validDiscoveryModel
	wj.journeyLock.Unlock()

	if discoveryModel == nil {
		return discovery.Model{}, errors.New("journey.DiscoveryModel: discovery model not set yet")
	}
	return *discoveryModel, nil
}

func (wj *journey) SetFilteredManifests(fmfs manifest.Scripts) {
	wj.filteredManifests = fmfs
}

func (wj *journey) FilteredManifests() (manifest.Scripts, error) {
	return wj.filteredManifests, nil
}

func (wj *journey) TestCases() (generation.TestCasesRun, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	logger := wj.log.WithFields(logrus.Fields{
		"package":  "server",
		"module":   "journey",
		"function": "TestCases",
	})

	if wj.validDiscoveryModel == nil {
		return generation.TestCasesRun{}, errDiscoveryModelNotSet
	}

	if wj.testCasesRunGenerated {
		logger.WithFields(logrus.Fields{
			"err":                      errTestCasesGenerated,
			"wj.testCasesRunGenerated": wj.testCasesRunGenerated,
		}).Error("Error getting generation.TestCasesRun ...")
		return generation.TestCasesRun{}, errTestCasesGenerated
	}

	if !wj.testCasesRunGenerated {
		wj.context.PutString(CtxPhase, "generation")
		config := wj.makeGeneratorConfig()
		discovery := wj.validDiscoveryModel.DiscoveryModel
		if len(discovery.DiscoveryItems) > 0 { // default currently "v3.1" ... allow "v3.0"
			// version string gets replaced in URLS like  "endpoint": "/open-banking/$api-version/aisp/account-access-consents",
			version, err := semver.ParseTolerant(discovery.DiscoveryItems[0].APISpecification.Version)
			if err != nil {
				logger.WithError(err).Error("parsing spec version")
			} else {
				wj.config.apiVersion = fmt.Sprintf("v%d.%d", version.Major, version.Minor)
				wj.context.PutString(CtxAPIVersion, wj.config.apiVersion)
			}
			logger.WithField("version", wj.config.apiVersion).Info("API url version")
		}

		logger.Debug("generator.GenerateManifestTests ...")

		wj.testCasesRun, wj.filteredManifests, wj.permissions = wj.generator.GenerateManifestTests(wj.log, config, discovery, &wj.context)

		logger.WithFields(logrus.Fields{
			"len(wj.permissions)": len(wj.permissions),
		}).Debug("manifest.RequiredTokens")
		for _, permission := range wj.permissions {
			logger.WithFields(logrus.Fields{
				"permission": permission,
			}).Debug("We have a permission ([]manifest.RequiredTokens)")
		}
		if discovery.TokenAcquisition == "psu" {
			logger.WithFields(logrus.Fields{
				"discovery.TokenAcquisition": discovery.TokenAcquisition,
			}).Debug("AcquirePSUTokens ...")
			definition := wj.makeRunDefinition()

			consentIds, tokenMap, err := executors.GetPsuConsent(definition, &wj.context, &wj.testCasesRun, wj.permissions)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error on executors.GetPsuConsent ...")
				return generation.TestCasesRun{}, errors.WithMessage(errConsentIDAcquisitionFailed, err.Error())
			}

			for k := range wj.permissions {
				if k == "payments" {
					paymentpermissions := wj.permissions["payments"]
					if len(paymentpermissions) > 0 {
						for _, spec := range wj.testCasesRun.TestCases {
							manifest.MapTokensToPaymentTestCases(paymentpermissions, spec.TestCases, &wj.context)
						}
					}
				}
			}

			for k, v := range tokenMap {
				wj.context.PutString(k, v)
			}

			wj.createTokenCollector(consentIds)
		} else {
			logger.WithFields(logrus.Fields{
				"discovery.TokenAcquisition": discovery.TokenAcquisition,
			}).Debug("AcquireHeadlessTokens ...")
			runDefinition := wj.makeRunDefinition()
			// TODO:Process multiple specs ... don't restrict to element [0]!!
			tokenPermissionsMap, err := executors.AcquireHeadlessTokens(wj.testCasesRun.TestCases[0].TestCases, &wj.context, runDefinition)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error on executors.AcquireHeadlessTokens ...")
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

	logger.Tracef("TestCaseRun.SpecConsentRequirements: %#v\n", wj.testCasesRun.SpecConsentRequirements)
	for k := range wj.testCasesRun.TestCases {
		logger.Tracef("TestCaseRun-Specificatino: %#v\n", wj.testCasesRun.TestCases[k].Specification)
	}
	logger.Tracef("Dumping Consents:---------------------------\n")
	for _, v := range wj.testCasesRun.SpecConsentRequirements {
		logger.Tracef("%s", v.Identifier)
		for _, x := range v.NamedPermissions {
			logger.Tracef("\tname: %s codeset: %#v\n\tconsent Url: %s", x.Name, x.CodeSet.CodeSet, x.ConsentUrl)
		}
	}
	return wj.testCasesRun, nil
}

func (wj *journey) CollectToken(code, state, scope string) error {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	logger := wj.log.WithFields(logrus.Fields{
		"package":  "server",
		"module":   "journey",
		"function": "CollectToken",
	})

	if !wj.testCasesRunGenerated {
		logger.WithFields(logrus.Fields{
			"err":   errTestCasesNotGenerated,
			"code":  code,
			"state": state,
			"scope": scope,
		}).Error("Error collecting token")
		return errTestCasesNotGenerated
	}

	accessToken, err := executors.ExchangeCodeForAccessToken(state, code, &wj.context)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":         err,
			"code":        code,
			"state":       state,
			"scope":       scope,
			"accessToken": accessToken,
		}).Error("Error collecting token due to error in executors.ExchangeCodeForAccessToken")
		return err
	}

	wj.context.PutString(state, accessToken)
	if state == "Token001" {
		logger.WithFields(logrus.Fields{
			"err":         err,
			"code":        code,
			"state":       state,
			"scope":       scope,
			"accessToken": accessToken,
		}).Warn(`Setting 'access_token' because state == "Token001"`)
		wj.context.PutString("access_token", accessToken) // tmp measure to get testcases running
	}

	accountPermissions := wj.permissions["accounts"]

	if wj.config.useDynamicResourceID {
		err := executors.GetDynamicResourceIds(state, accessToken, &wj.context, accountPermissions)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("Dynamic resource allocation failure")
			return errDynamicResourceAllocationFailed
		}
	}

	for _, v := range wj.permissions["accounts"] {
		logger.Tracef("journey perms: %#v", v)
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
	logger := wj.log.WithField("function", "RunTests")

	if !wj.testCasesRunGenerated {
		logger.WithFields(logrus.Fields{
			"err": errTestCasesNotGenerated,
		}).Error("Error on starting run")
		return errTestCasesNotGenerated
	}

	if !wj.allCollected {
		logger.WithFields(logrus.Fields{
			"err": errNotFinishedCollectingTokens,
		}).Error("Error on starting run")
		return errNotFinishedCollectingTokens
	}

	requiredTokens := wj.permissions

	if wj.config.useDynamicResourceID {
		for _, accountPermissions := range wj.permissions["accounts"] {
			// cycle over all test case ids for this account permission/token set
			for _, tcID := range accountPermissions.IDs {
				for i := range wj.testCasesRun.TestCases {
					specType := wj.testCasesRun.TestCases[i].Specification.SpecType
					// isolate all testcases to be run that are from and 'account' spec type
					if specType == "accounts" {
						tc := wj.testCasesRun.TestCases[i].TestCases
						// look for test cases matching the permission set test case list
						for j, test := range tc {
							if test.ID == tcID {
								resourceCtx := model.Context{}
								resourceCtx.PutString(CtxConsentedAccountID, accountPermissions.AccountID)
								resourceCtx.PutString(CtxStatementID, accountPermissions.StatementID)
								// perform the dynamic resource id replacement
								test.ProcessReplacementFields(&resourceCtx, false)
								wj.testCasesRun.TestCases[i].TestCases[j] = test
							}
						}
					}
				}
			}
		}
		// put a default accountid and statement id in the journey context for those tests that haven't got a token that can call /accounts
		wj.context.PutString(CtxConsentedAccountID, wj.config.resourceIDs.AccountIDs[0].AccountID)
		wj.context.PutString(CtxStatementID, wj.config.resourceIDs.StatementIDs[0].StatementID)
	}

	for k := range wj.testCasesRun.TestCases {
		specType := wj.testCasesRun.TestCases[k].Specification.SpecType
		manifest.MapTokensToTestCases(requiredTokens[specType], wj.testCasesRun.TestCases[k].TestCases)
		wj.dumpJSON(wj.testCasesRun.TestCases[k].TestCases)
	}

	runDefinition := wj.makeRunDefinition()
	runner := executors.NewTestCaseRunner(wj.log, runDefinition, wj.daemonController)
	wj.context.PutString(CtxPhase, "run")
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

	wj.log.WithFields(logrus.Fields{
		"package":      "server",
		"module":       "journey",
		"function":     "createTokenCollector",
		"consentIds":   fmt.Sprintf("%#v", consentIds),
		"allCollected": wj.allCollected,
	}).Debug("TokenCollector status ...")
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
	certificateSigning            authentication.Certificate
	certificateTransport          authentication.Certificate
	clientID                      string
	clientSecret                  string
	tokenEndpoint                 string
	ResponseType                  string
	tokenEndpointAuthMethod       string
	authorizationEndpoint         string
	resourceBaseURL               string
	xXFAPIFinancialID             string
	issuer                        string
	redirectURL                   string
	resourceIDs                   model.ResourceIDs
	creditorAccount               models.Payment
	instructedAmount              models.InstructedAmount
	currencyOfTransfer            string
	apiVersion                    string
	transactionFromDate           string
	transactionToDate             string
	requestObjectSigningAlgorithm string
	signingPrivate                string
	signingPublic                 string
	useNonOBDirectory             bool
	signingKid                    string
	signatureTrustAnchor          string
	useDynamicResourceID          bool
}

func (wj *journey) SetConfig(config JourneyConfig) error {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	wj.config = config
	wj.config.useDynamicResourceID = dynamicResourceIDs // fed from environment variable 'dynres'=true/false
	err := PutParametersToJourneyContext(wj.config, wj.context)
	if err != nil {
		return err
	}

	wj.customTestParametersToJourneyContext()
	return nil
}

func (wj *journey) Events() events.Events {
	return wj.events
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
	log.WithFields(logrus.Fields{
		"package":    "server",
		"function":   "consentIdsToTestCaseRun",
		"consentIds": consentIds,
	}).Debug("...")
	for _, v := range testCasesRun.SpecConsentRequirements {
		for x, permission := range v.NamedPermissions {
			for _, consentID := range consentIds {
				if consentID.TokenName == permission.Name {
					permission.ConsentUrl = consentID.ConsentURL
					log.WithFields(logrus.Fields{
						"permission.Name":       permission.Name,
						"permission.ConsentUrl": permission.ConsentUrl,
						"consentID":             consentID,
					}).Debug("consentIdsToTestCaseRun ... Setting consent url for token")

					v.NamedPermissions[x] = permission
				}
			}
		}
	}
}

// func dumpPermissions(p map[string][]manifest.RequiredTokens, title string) {
// 	logrus.Tracef("Dump Permissions at %s \n", title)
// 	for _, v := range p {
// 		logrus.Tracef("%#v\n", v)
// 	}
// }

// Utility to Dump Json
func (wj *journey) dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	wj.log.Traceln(string(model))
}

// EnableDynamicResourceIDs is triggered by and environment variable dynids=true
func EnableDynamicResourceIDs() {
	dynamicResourceIDs = true
}
