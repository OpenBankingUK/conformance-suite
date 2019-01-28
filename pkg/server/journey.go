package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"github.com/pkg/errors"
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
	SetCertificateSigning(authentication.Certificate)
	CertificateSigning() authentication.Certificate
	SetCertificateTransport(authentication.Certificate)
	CertificateTransport() authentication.Certificate
}

var errTestCasesNotSet = errors.New("error running test cases, test cases not set")

type journey struct {
	generator            generation.Generator
	testCases            []generation.SpecificationTestCases
	validator            discovery.Validator
	validDiscoveryModel  *discovery.Model
	certificateSigning   authentication.Certificate
	daemonController     executors.DaemonController
	certificateTransport authentication.Certificate
}

// NewJourney creates an instance for a user journey
func NewJourney(generator generation.Generator, validator discovery.Validator) Journey {
	daemonController := executors.NewDaemonController(
		make(chan results.TestCase, 100),
		make(chan error, 100),
	)
	return &journey{
		generator:        generator,
		validator:        validator,
		daemonController: daemonController,
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

	wj.validDiscoveryModel = discoveryModel
	wj.testCases = nil

	return discovery.NoValidationFailures, nil
}

func (wj *journey) DiscoveryModel() (*discovery.Model, error) {
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	return wj.validDiscoveryModel, nil
}

func (wj *journey) TestCases() ([]generation.SpecificationTestCases, error) {
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	if wj.testCases == nil {
		wj.testCases = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
	}
	return wj.testCases, nil
}

func (wj *journey) RunTests() error {
	if wj.validDiscoveryModel == nil {
		return errDiscoveryModelNotSet
	}

	if wj.testCases == nil {
		return errTestCasesNotSet
	}

	specTestCases, err := wj.TestCases()
	if err != nil {
		return err
	}

	runDefinition := executors.RunDefinition{
		DiscoModel:    wj.validDiscoveryModel,
		SpecTests:     specTestCases,
		SigningCert:   wj.CertificateSigning(),
		TransportCert: wj.CertificateTransport(),
	}

	runner := executors.NewTestCaseRunner(runDefinition, wj.daemonController)
	err = runner.RunTestCases()
	if err != nil {
		return err
	}

	return nil
}

func (wj *journey) Results() executors.DaemonController {
	return wj.daemonController
}

func (wj *journey) StopTestRun() {
	wj.daemonController.Stop()
}

func (wj *journey) SetCertificateSigning(certificateSigning authentication.Certificate) {
	wj.certificateSigning = certificateSigning
}

func (wj *journey) CertificateSigning() authentication.Certificate {
	return wj.certificateSigning
}

func (wj *journey) SetCertificateTransport(certificateTransport authentication.Certificate) {
	wj.certificateTransport = certificateTransport
}

func (wj *journey) CertificateTransport() authentication.Certificate {
	return wj.certificateTransport
}
