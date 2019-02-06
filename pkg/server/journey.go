package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/pkg/errors"
	"sync"
)

var errDiscoveryModelNotSet = errors.New("error discovery model not set")

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	DiscoveryModel() (*discovery.Model, error)
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() (TestCasesRun, error)
	RunTests() error
	StopTestRun()
	Results() executors.DaemonController
	SetCertificates(signing, transport authentication.Certificate)
}

type journey struct {
	generator            generation.Generator
	validator            discovery.Validator
	daemonController     executors.DaemonController
	resolver             func(groups []permissions.Group) permissions.CodeSetResultSet
	journeyLock          *sync.Mutex
	specTestCases        []generation.SpecificationTestCases
	validDiscoveryModel  *discovery.Model
	certificateSigning   authentication.Certificate
	certificateTransport authentication.Certificate
}

// NewJourney creates an instance for a user journey
func NewJourney(generator generation.Generator, validator discovery.Validator) *journey {
	return &journey{
		generator:        generator,
		validator:        validator,
		daemonController: executors.NewBufferedDaemonController(),
		resolver:         permissions.Resolver,
		journeyLock:      &sync.Mutex{},
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
	wj.specTestCases = nil
	wj.journeyLock.Unlock()

	return discovery.NoValidationFailures, nil
}

func (wj *journey) DiscoveryModel() (*discovery.Model, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	return wj.validDiscoveryModel, nil
}

// TestCasesRun represents all specs and their test and a list of tokens
// required to run those tests
type TestCasesRun struct {
	TestCases               []generation.SpecificationTestCases `json:"specCases"`
	SpecConsentRequirements []model.SpecConsentRequirements     `json:"specTokens"`
}

func (wj *journey) TestCases() (TestCasesRun, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	if wj.validDiscoveryModel == nil {
		return TestCasesRun{}, errDiscoveryModelNotSet
	}
	if wj.specTestCases == nil {
		wj.specTestCases = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
	}

	tokens := wj.permissionSpecTokens()

	return TestCasesRun{wj.specTestCases, tokens}, nil
}

// permissionSpecTokens calls resolver to get list of permission sets required to run all test cases
func (wj *journey) permissionSpecTokens() []model.SpecConsentRequirements {
	var tokens []model.SpecConsentRequirements
	for _, spec := range wj.specTestCases {
		var groups []permissions.Group
		for _, tc := range spec.TestCases {
			groups = append(groups, model.NewPermissionGroup(tc))
		}
		resultSet := wj.resolver(groups)
		tokens = append(tokens, model.NewSpecConsentRequirements(resultSet, spec.Specification.Name))
	}
	return tokens
}

func (wj *journey) RunTests() error {
	specTestCasesRun, err := wj.TestCases()
	if err != nil {
		return err
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		SpecTests:     specTestCasesRun.TestCases,
		SpecTokens:    []model.SpecConsentRequirements{},
		SigningCert:   wj.certificateSigning,
		TransportCert: wj.certificateTransport,
	}

	runner := executors.NewTestCaseRunner(runDefinition, wj.daemonController)
	return runner.RunTestCases()
}

func (wj *journey) Results() executors.DaemonController {
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.daemonController.Stop()
}

func (wj *journey) SetCertificates(signing, transport authentication.Certificate) {
	wj.journeyLock.Lock()
	wj.certificateSigning = signing
	wj.certificateTransport = transport
	wj.journeyLock.Unlock()
}
