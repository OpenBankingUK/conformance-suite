package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"

	"github.com/pkg/errors"
)

var errDiscoveryModelNotSet = errors.New("error generation test cases, discovery model not set")

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
	TestCases() ([]generation.SpecificationTestCases, error)
	RunTests() (reporting.Result, error)
	SetCertificateSigning(certificateSigning authentication.Certificate) Journey
	CertificateSigning() authentication.Certificate
	SetCertificateTransport(certificateTransport authentication.Certificate) Journey
	CertificateTransport() authentication.Certificate
}

var errTestCasesNotSet = errors.New("error running test cases, test cases not set")

type journey struct {
	generator            generation.Generator
	testCases            []generation.SpecificationTestCases
	validator            discovery.Validator
	validDiscoveryModel  *discovery.Model
	reportService        reporting.Service
	certificateSigning   authentication.Certificate
	certificateTransport authentication.Certificate
}

// NewWebJourney creates an instance for a user journey
func NewWebJourney(generator generation.Generator, validator discovery.Validator) Journey {
	return &journey{
		generator:     generator,
		validator:     validator,
		reportService: reporting.NewMockedService(),
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

func (wj *journey) TestCases() ([]generation.SpecificationTestCases, error) {
	if wj.validDiscoveryModel == nil {
		return nil, errDiscoveryModelNotSet
	}
	if wj.testCases == nil {
		wj.testCases = wj.generator.GenerateSpecificationTestCases(wj.validDiscoveryModel.DiscoveryModel)
	}
	return wj.testCases, nil
}

func (wj *journey) RunTests() (reporting.Result, error) {
	if wj.testCases == nil {
		return reporting.Result{}, errTestCasesNotSet
	}
	return wj.reportService.Run(wj.testCases)
}

func (wj *journey) SetCertificateSigning(certificateSigning authentication.Certificate) Journey {
	wj.certificateSigning = certificateSigning
	return wj
}

func (wj *journey) CertificateSigning() authentication.Certificate {
	return wj.certificateSigning
}

func (wj *journey) SetCertificateTransport(certificateTransport authentication.Certificate) Journey {
	wj.certificateTransport = certificateTransport
	return wj
}

func (wj *journey) CertificateTransport() authentication.Certificate {
	return wj.certificateTransport
}
