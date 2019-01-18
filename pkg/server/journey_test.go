package server

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	gmocks "bitbucket.org/openbankingteam/conformance-suite/pkg/generation/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	publicCertValid = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`
	privateCertValid = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`
)

func TestJourneySetDiscoveryModelValidatesModel(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(t, discovery.NoValidationFailures, failures)
	validator.AssertExpectations(t)
	generator.AssertExpectations(t)
}

func TestJourneySetDiscoveryModelHandlesErrorFromValidator(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, errors.New("validator error"))
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.Error(t, err)
	assert.Equal(t, "error setting discovery model: validator error", err.Error())
	assert.Nil(t, failures)
}

func TestJourneySetDiscoveryModelReturnsFailuresFromValidator(t *testing.T) {
	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	failure := discovery.ValidationFailure{
		Key:   "DiscoveryModel.Name",
		Error: "Field 'Name' is required",
	}
	expectedFailures := discovery.ValidationFailures{failure}
	validator.On("Validate", discoveryModel).Return(expectedFailures, nil)
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(t, expectedFailures, failures)
}

func TestJourneyTestCasesCantGenerateIfDiscoveryNotSet(t *testing.T) {
	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	testCases, err := journey.TestCases()

	assert.Error(t, err)
	assert.Nil(t, testCases)
}

func TestJourneyTestCasesGenerate(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	expectedTestCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).Return(expectedTestCases)
	journey := NewJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	testCases, err := journey.TestCases()

	assert.NoError(t, err)
	assert.Equal(t, expectedTestCases, testCases)
}

func TestJourneyTestCasesDoesntREGenerate(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	expectedTestCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).
		Return(expectedTestCases).Times(1)

	journey := NewJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)
	firstRunTestCases, err := journey.TestCases()
	require.NoError(t, err)

	testCases, err := journey.TestCases()

	assert.NoError(t, err)
	assert.Equal(t, expectedTestCases, testCases)
	assert.Equal(t, firstRunTestCases, testCases)
	generator.AssertExpectations(t)
}

func TestJourneyRunTestCasesCantRunIfNoTestCases(t *testing.T) {
	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	result, err := journey.RunTests()

	assert.EqualError(t, err, "error running test cases, test cases not set")
	assert.Equal(t, reporting.Result{}, result)
}

func TestJourneyRunTestCases(t *testing.T) {
	validator := &mocks.Validator{}
	discoveryModel := &discovery.Model{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures, nil)
	testCases := []generation.SpecificationTestCases{}
	generator := &gmocks.Generator{}
	generator.On("GenerateSpecificationTestCases", discoveryModel.DiscoveryModel).
		Return(testCases).Times(1)

	journey := NewJourney(generator, validator)
	_, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(t, err)

	_, err = journey.TestCases()
	require.NoError(t, err)

	//result, err := journey.RunTests()
	//_ = result
	//assert.NoError(t, err)
	noResult := []reporting.Specification([]reporting.Specification{})
	//assert.Equal(t, noResult, result.Specifications)
	_ = noResult
	generator.AssertExpectations(t)
}

func TestJourneySetCertificateSigning(t *testing.T) {
	require := require.New(t)

	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	require.Nil(journey.CertificateSigning())

	certificateSigning, err := authentication.NewCertificate(publicCertValid, privateCertValid)
	require.NoError(err)
	require.NotNil(certificateSigning)

	journey.SetCertificateSigning(certificateSigning)

	require.Equal(certificateSigning, journey.CertificateSigning())
}

func TestJourneySetCertificateTransport(t *testing.T) {
	require := require.New(t)

	validator := &mocks.Validator{}
	generator := &gmocks.Generator{}
	journey := NewJourney(generator, validator)

	require.Nil(journey.CertificateTransport())

	certificateTransport, err := authentication.NewCertificate(publicCertValid, privateCertValid)
	require.NoError(err)
	require.NotNil(certificateTransport)

	journey.SetCertificateTransport(certificateTransport)

	require.Equal(certificateTransport, journey.CertificateTransport())
}
