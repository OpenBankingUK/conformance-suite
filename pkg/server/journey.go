package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"github.com/pkg/errors"
	"sync"
)

var errDiscoveryModelNotSet = errors.New("error discovery model not set")

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	DiscoveryModel() (*discovery.Model, error)
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() ([]generation.SpecificationTestCases, error)
	RunTests() error
	StopTestRun()
	Results() executors.DaemonController
	SetCertificates(signing, transport authentication.Certificate)
}

type journey struct {
	generator        generation.Generator
	validator        discovery.Validator
	daemonController executors.DaemonController

	journeyLock          *sync.Mutex
	testCases            []generation.SpecificationTestCases
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
	wj.testCases = nil
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

func (wj *journey) TestCases() ([]generation.SpecificationTestCases, error) {
	wj.journeyLock.Lock()
	defer wj.journeyLock.Unlock()
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	if wj.testCases == nil {
		wj.testCases = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
	}
	return wj.testCases, nil
}

func (wj *journey) RunTests() error {
	specTestCases, err := wj.TestCases()
	if err != nil {
		return err
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		SpecTests:     specTestCases,
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
