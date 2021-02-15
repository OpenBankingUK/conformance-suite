//go:generate mockery -name Journey -inpkg
package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schemaprops"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
)

var (
	errDiscoveryModelNotSet            = errors.New("error discovery model not set")
	errTestCasesNotGenerated           = errors.New("error test cases not generated")
	errTestCasesGenerated              = errors.New("error test cases already generated")
	errNotFinishedCollectingTokens     = errors.New("error not finished collecting tokens")
	errConsentIDAcquisitionFailed      = errors.New("ConsentId acquistion failed")
	errDynamicResourceAllocationFailed = errors.New("Dynamic Resource allocation failed")
	errNoTestCases                     = errors.New("No testcases were generated - please select a wider set of endpoints to test")
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
	TestCases() (generation.SpecRun, error)
	CollectToken(code, state, scope string) error
	AllTokenCollected() bool
	RunTests() error
	StopTestRun()
	NewDaemonController()
	Results() executors.DaemonController
	SetConfig(config JourneyConfig) error
	ConditionalProperties() []discovery.ConditionalAPIProperties
	Events() events.Events
	TLSVersionResult() map[string]*discovery.TLSValidationResult
}

type journey struct {
	generator             generation.Generator
	validator             discovery.Validator
	daemonController      executors.DaemonController
	journeyLock           *sync.Mutex
	specRun               generation.SpecRun
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
	tlsValidator          discovery.TLSValidator
	conditionalProperties []discovery.ConditionalAPIProperties
	dynamicResourceIDs    bool
}

// NewJourney creates an instance for a user journey
func NewJourney(logger *logrus.Entry, generator generation.Generator, validator discovery.Validator, tlsValidator discovery.TLSValidator, dynamicResourceIDs bool) *journey {
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
		tlsValidator:          tlsValidator,
		dynamicResourceIDs:    dynamicResourceIDs,
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

	if discoveryModel.DiscoveryModel.DiscoveryVersion == "v0.4.0" { // Conditional properties requires 0.4.0
		//TODO: remove this constraint once support for v0.3.0 discovery model is dropped
		conditionalApiProperties, hasProperties, err := discovery.GetConditionalProperties(discoveryModel)
		if err != nil {
			return nil, errors.Wrap(err, "journey.SetDiscoveryModel: error processing conditional properties")
		}
		if hasProperties {
			wj.conditionalProperties = conditionalApiProperties
			logrus.Tracef("conditionalProperties from discovery model: %#v", wj.conditionalProperties)
		} else {
			logrus.Trace("No Conditional Properties found")
		}
	}

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

func (wj *journey) TLSVersionResult() map[string]*discovery.TLSValidationResult {
	logger := wj.log.WithFields(logrus.Fields{
		"package":  "server",
		"module":   "journey",
		"function": "TLSVersionResult",
	})
	tlsValidationResult := make(map[string]*discovery.TLSValidationResult, len(wj.validDiscoveryModel.DiscoveryModel.DiscoveryItems))
	for _, discoveryItem := range wj.validDiscoveryModel.DiscoveryModel.DiscoveryItems {
		tlsVersionKey := wj.tlsVersionCtxKey(discoveryItem.APISpecification.Name)
		tlsVersion, err := wj.context.GetString(tlsVersionKey)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":                err,
				"discoveryItem[key]": tlsVersionKey,
				"discoveryItem.discoveryItem.APISpecification.Name": discoveryItem.APISpecification.Name,
			}).Errorf("Error getting %s from context ...", tlsVersionKey)
			continue
		}
		tlsValidKey := wj.tlsValidCtxKey(discoveryItem.APISpecification.Name)
		tlsValid, ok := wj.context.Get(tlsValidKey)
		if !ok {
			logger.WithFields(logrus.Fields{
				"discoveryItem[key]": tlsValidKey,
				"discoveryItem.discoveryItem.APISpecification.Name": discoveryItem.APISpecification.Name,
			}).Errorf("Error getting %s from context ...", tlsValidKey)
			continue
		}
		tlsValidationResult[strings.ReplaceAll(discoveryItem.APISpecification.Name, " ", "-")] = &discovery.TLSValidationResult{TLSVersion: tlsVersion, Valid: tlsValid.(bool)}
	}

	return tlsValidationResult
}

func (wj *journey) SetFilteredManifests(fmfs manifest.Scripts) {
	wj.filteredManifests = fmfs
}

func (wj *journey) FilteredManifests() (manifest.Scripts, error) {
	return wj.filteredManifests, nil
}

func (wj *journey) TestCases() (generation.SpecRun, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	logger := wj.log.WithFields(logrus.Fields{
		"package":  "server",
		"module":   "journey",
		"function": "TestCases",
	})

	if wj.validDiscoveryModel == nil {
		return generation.SpecRun{}, errDiscoveryModelNotSet
	}

	if wj.testCasesRunGenerated {
		logger.WithFields(logrus.Fields{
			"err":                      errTestCasesGenerated,
			"wj.testCasesRunGenerated": wj.testCasesRunGenerated,
		}).Error("Error getting generation.TestCasesRun ...")
		return generation.SpecRun{}, errTestCasesGenerated
	}

	jwks_uri := authentication.GetJWKSUri()
	if jwks_uri != "" { // STORE jwks_uri from well known endpoint in journey context
		wj.context.PutString("jwks_uri", jwks_uri)
	} else {
		logrus.Warn("JWKS URI is empty")
	}

	if tlsCheck {
		for k, discoveryItem := range wj.validDiscoveryModel.DiscoveryModel.DiscoveryItems {
			tlsValidationResult, err := wj.tlsValidator.ValidateTLSVersion(discoveryItem.ResourceBaseURI)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err":                           errors.Wrapf(err, "unable to validate TLS version for uri %s", discoveryItem.ResourceBaseURI),
					"discoveryItem[key]":            k,
					"discoveryItem.ResourceBaseURI": discoveryItem.ResourceBaseURI,
				}).Error("Error validating TLS version for discovery item ResourceBaseURI")
			}
			wj.context.PutString(wj.tlsVersionCtxKey(discoveryItem.APISpecification.Name), tlsValidationResult.TLSVersion)
			wj.context.Put(wj.tlsValidCtxKey(discoveryItem.APISpecification.Name), tlsValidationResult.Valid)
		}
	} else {
		logrus.Warn("TLS Check disabled")
	}

	wj.context.PutString(CtxPhase, "generation")
	config := wj.makeGeneratorConfig()
	discovery := wj.validDiscoveryModel.DiscoveryModel
	if len(discovery.DiscoveryItems) > 0 { // default currently "v3.1" ... allow "v3.0"
		apiversions := DetermineAPIVersions(discovery.DiscoveryItems)
		if len(apiversions) > 0 {
			wj.context.PutStringSlice("apiversions", apiversions)
		}
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
	logrus.Tracef("conditionalProperties from journey config: %#v", wj.config.conditionalProperties)
	wj.specRun, wj.filteredManifests, wj.permissions = wj.generator.GenerateManifestTests(wj.log, config, discovery, &wj.context, wj.config.conditionalProperties)

	tests := 0
	for _, sp := range wj.specRun.SpecTestCases {
		tests += len(sp.TestCases)
	}
	if tests == 0 { // no tests to run
		logrus.Warn("No TestCases Generated!!!")
		return generation.SpecRun{}, errNoTestCases
	}

	for _, spec := range wj.permissions {
		for _, required := range spec {
			logger.WithFields(logrus.Fields{
				"permission": required.Name,
				"idlist":     required.IDs,
			}).Debug("We have a permission ([]manifest.RequiredTokens)")
		}
	}

	collector := schemaprops.GetPropertyCollector()
	collector.SetCollectorAPIDetails(schemaprops.ConsentGathering, "")

	if discovery.TokenAcquisition == "psu" { // Handle  PSU Consent
		logger.WithFields(logrus.Fields{
			"discovery.TokenAcquisition": discovery.TokenAcquisition,
		}).Debug("AcquirePSUTokens ...")
		definition := wj.makeRunDefinition()

		consentIds, tokenMap, err := executors.GetPsuConsent(definition, &wj.context, &wj.specRun, wj.permissions)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error on executors.GetPsuConsent ...")
			return generation.SpecRun{}, errors.WithMessage(errConsentIDAcquisitionFailed, err.Error())
		}

		for k := range wj.permissions {
			if k == "payments" {
				paymentpermissions := wj.permissions["payments"]
				if len(paymentpermissions) > 0 {
					for _, spec := range wj.specRun.SpecTestCases {
						manifest.MapTokensToPaymentTestCases(paymentpermissions, spec.TestCases, &wj.context)
					}
				}
			}
			if k == "cbpii" {
				cbpiiPerms := wj.permissions["cbpii"]
				if len(cbpiiPerms) > 0 {
					for _, spec := range wj.specRun.SpecTestCases {
						manifest.MapTokensToCBPIITestCases(cbpiiPerms, spec.TestCases, &wj.context)
					}
				}
			}
		}

		for k, v := range tokenMap {
			wj.context.PutString(k, v)
			logger.Tracef("processtokenMap %s:%s into context", k, v)
		}

		wj.createTokenCollector(consentIds)

	} else { // Handle headless token acquistion

		logger.WithFields(logrus.Fields{
			"discovery.TokenAcquisition": discovery.TokenAcquisition,
		}).Debug("AcquireHeadlessTokens ...")
		definition := wj.makeRunDefinition()

		tokenPermissionsMap, err := executors.GetHeadlessConsent(definition, &wj.context, &wj.specRun, wj.permissions)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error on executors.AcquireHeadlessTokens ...")
			return generation.SpecRun{}, errConsentIDAcquisitionFailed
		}

		tokenMap := map[string]string{} // Put access tokens into context
		for _, v := range wj.specRun.SpecTestCases {
			aMap := manifest.MapTokensToTestCases(tokenPermissionsMap, v.TestCases)
			for x, y := range aMap {
				tokenMap[x] = y
			}
		}

		for k, v := range tokenMap {
			wj.context.PutString(k, v)
			logger.Tracef("processtokenMap %s:%s into context", k, v)
		}

		for k := range wj.permissions {
			if k == "payments" {
				paymentpermissions := wj.permissions["payments"]
				logger.Tracef("We have %d Payment Permissions", len(paymentpermissions))
				if len(paymentpermissions) > 0 {
					for _, spec := range wj.specRun.SpecTestCases {
						logger.Tracef("Analysing %d test cases for token mapping", len(spec.TestCases))
						manifest.MapTokensToPaymentTestCases(paymentpermissions, spec.TestCases, &wj.context)
					}
				}
			}
			if k == "cbpii" {
				cbpiiPerms := wj.permissions["cbpii"]
				if len(cbpiiPerms) > 0 {
					for _, spec := range wj.specRun.SpecTestCases {
						manifest.MapTokensToCBPIITestCases(cbpiiPerms, spec.TestCases, &wj.context)
					}
				}
			}
		}

		wj.allCollected = true
	}
	wj.testCasesRunGenerated = true

	logger.Tracef("SpecRun.SpecConsentRequirements: %#v", wj.specRun.SpecConsentRequirements)
	for k := range wj.specRun.SpecTestCases {
		logger.Tracef("SpecRun-Specification: %#v", wj.specRun.SpecTestCases[k].Specification)
	}
	return wj.specRun, nil
}

func (wj *journey) tlsVersionCtxKey(discoveryItemName string) string {
	return fmt.Sprintf("tlsVersionForDiscoveryItem-%s", strings.ReplaceAll(discoveryItemName, " ", "-"))
}

func (wj *journey) tlsValidCtxKey(discoveryItemName string) string {
	return fmt.Sprintf("tlsIsValidForDiscoveryItem-%s", strings.ReplaceAll(discoveryItemName, " ", "-"))
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

	if wj.config.useDynamicResourceID {
		err := executors.GetDynamicResourceIds(state, accessToken, &wj.context, wj.permissions["accounts"])
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

	if wj.config.useDynamicResourceID {
		for _, accountPermissions := range wj.permissions["accounts"] {
			// cycle over all test case ids for this account permission/token set
			for _, tcID := range accountPermissions.IDs {
				for i := range wj.specRun.SpecTestCases {
					specType := wj.specRun.SpecTestCases[i].Specification.SpecType
					// isolate all testcases to be run that are from and 'account' spec type
					if specType == "accounts" {
						tc := wj.specRun.SpecTestCases[i].TestCases
						// look for test cases matching the permission set test case list
						for j, test := range tc {
							if test.ID == tcID {
								resourceCtx := model.Context{}
								resourceCtx.PutString(CtxConsentedAccountID, accountPermissions.AccountID)
								// perform the dynamic resource id replacement
								test.ProcessReplacementFields(&resourceCtx, false)
								wj.specRun.SpecTestCases[i].TestCases[j] = test
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

	requiredTokens := wj.permissions

	for k := range wj.specRun.SpecTestCases {
		specType := wj.specRun.SpecTestCases[k].Specification.SpecType
		manifest.MapTokensToTestCases(requiredTokens[specType], wj.specRun.SpecTestCases[k].TestCases)
		wj.dumpJSON(wj.specRun.SpecTestCases[k].TestCases)
	}

	runDefinition := wj.makeRunDefinition()
	runner := executors.NewTestCaseRunner(wj.log, runDefinition, wj.daemonController)
	wj.context.PutString(CtxPhase, "run")
	err := runner.RunTestCases(&wj.context)
	return err
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
		consentIdsToTestCaseRun(wj.log, consentIds, &wj.specRun)

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
		SpecRun:       wj.specRun,
		SigningCert:   wj.config.certificateSigning,
		TransportCert: wj.config.certificateTransport,
	}
}

// JourneyConfig main configuration variables
type JourneyConfig struct {
	certificateSigning            authentication.Certificate
	certificateTransport          authentication.Certificate
	tppSignatureKID               string
	tppSignatureIssuer            string
	tppSignatureTAN               string
	clientID                      string
	clientSecret                  string
	tokenEndpoint                 string
	ResponseType                  string
	tokenEndpointAuthMethod       string
	authorizationEndpoint         string
	resourceBaseURL               string
	xXFAPIFinancialID             string
	xXFAPICustomerIPAddress       string
	redirectURL                   string
	resourceIDs                   model.ResourceIDs
	creditorAccount               models.Payment
	internationalCreditorAccount  models.Payment
	instructedAmount              models.InstructedAmount
	paymentFrequency              models.PaymentFrequency
	firstPaymentDateTime          string
	requestedExecutionDateTime    string
	currencyOfTransfer            string
	apiVersion                    string
	transactionFromDate           string
	transactionToDate             string
	requestObjectSigningAlgorithm string
	signingPrivate                string
	signingPublic                 string
	useDynamicResourceID          bool
	AcrValuesSupported            []string
	conditionalProperties         []discovery.ConditionalAPIProperties
	cbpiiDebtorAccount            discovery.CBPIIDebtorAccount
	issuer                        string
}

func (wj *journey) SetConfig(config JourneyConfig) error {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()

	wj.config = config
	wj.config.useDynamicResourceID = wj.dynamicResourceIDs // fed from environment variable 'dynres'=true/false
	err := PutParametersToJourneyContext(wj.config, wj.context)
	if err != nil {
		return err
	}

	wj.customTestParametersToJourneyContext()
	return nil
}

// ConditionalProperties retrieve conditional properties right after
// they have been set from the discovery model to the webJourney.ConditionalProperties
func (wj *journey) ConditionalProperties() []discovery.ConditionalAPIProperties {
	return wj.conditionalProperties
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

func consentIdsToTestCaseRun(log *logrus.Entry, consentIds []executors.TokenConsentIDItem, testCasesRun *generation.SpecRun) {
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

// Utility to Dump Json
func (wj *journey) dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	wj.log.Traceln(string(model))
}

// EnableDynamicResourceIDs is triggered by and environment variable dynids=true
func (wj *journey) EnableDynamicResourceIDs() {
	wj.dynamicResourceIDs = true
}

// DetermineAPIVersions
func DetermineAPIVersions(apis []discovery.ModelDiscoveryItem) []string {
	apiversions := []string{}
	for _, v := range apis {
		v.APISpecification.SpecType, _ = manifest.GetSpecType(v.APISpecification.SchemaVersion)
		apiversions = append(apiversions, v.APISpecification.SpecType+"_"+v.APISpecification.Version)
		logrus.Warnf("spectype %s, specversion %s", v.APISpecification.SpecType, v.APISpecification.Version)
	}
	return apiversions
}

var tlsCheck = true

func EnableTLSCheck(state bool) {
	tlsCheck = state
}
