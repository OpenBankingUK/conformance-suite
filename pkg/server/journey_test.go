package server

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	gmocks "bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"github.com/pkg/errors"
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
	assert := test.NewAssert(t)

	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures(), nil)
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(discovery.NoValidationFailures(), failures)
	validator.AssertExpectations(t)
	generator.AssertExpectations(t)
}

func TestJourneySetDiscoveryModelHandlesErrorFromValidator(t *testing.T) {
	assert := test.NewAssert(t)

	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	expectedFailures := discovery.ValidationFailures{}
	validator.On("Validate", discoveryModel).Return(expectedFailures, errors.New("validator error"))
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.Error(t, err)
	assert.Equal("journey.SetDiscoveryModel: error setting discovery model: validator error", err.Error())
	assert.Nil(failures)
}

func TestJourneySetDiscoveryModelReturnsFailuresFromValidator(t *testing.T) {
	assert := test.NewAssert(t)

	discoveryModel := &discovery.Model{}
	validator := &mocks.Validator{}
	failure := discovery.ValidationFailure{
		Key:   "DiscoveryModel.Name",
		Error: "Field 'Name' is required",
	}
	expectedFailures := discovery.ValidationFailures{failure}
	validator.On("Validate", discoveryModel).Return(expectedFailures, nil)
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	failures, err := journey.SetDiscoveryModel(discoveryModel)

	require.NoError(t, err)
	assert.Equal(expectedFailures, failures)
}

func TestJourneyTestCasesCantGenerateIfDiscoveryNotSet(t *testing.T) {
	assert := test.NewAssert(t)

	validator := &mocks.Validator{}
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	testCases, err := journey.TestCases()

	assert.Error(err)
	assert.Equal(generation.SpecRun{}, testCases)
}

func TestJourneyRunTestCasesCantRunIfNoTestCases(t *testing.T) {
	assert := test.NewAssert(t)

	validator := &mocks.Validator{}
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	err := journey.RunTests()

	assert.EqualError(err, "error test cases not generated")
}

func TestJourneySetConfig(t *testing.T) {
	require := test.NewRequire(t)

	validator := &mocks.Validator{}
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	require.Equal(JourneyConfig{}, journey.config)

	certificateSigning, err := authentication.NewCertificate(publicCertValid, privateCertValid)
	require.NoError(err)
	require.NotNil(certificateSigning)
	certificateTransport, err := authentication.NewCertificate(publicCertValid, privateCertValid)
	require.NoError(err)
	require.NotNil(certificateTransport)

	resourceIDs := model.ResourceIDs{
		AccountIDs:   []model.ResourceAccountID{{AccountID: "account-id"}},
		StatementIDs: []model.ResourceStatementID{{StatementID: "statement-id"}},
	}
	config := JourneyConfig{
		certificateSigning:    certificateSigning,
		certificateTransport:  certificateTransport,
		clientID:              "8672384e-9a33-439f-8924-67bb14340d71",
		clientSecret:          "2cfb31a3-5443-4e65-b2bc-ef8e00266a77",
		tokenEndpoint:         "https://modelobank2018.o3bank.co.uk:4201/token",
		authorizationEndpoint: "https://modelobankauth2018.o3bank.co.uk:4101/auth",
		resourceBaseURL:       "https://ob19-rs1.o3bank.co.uk:4501",
		xXFAPIFinancialID:     "0015800001041RHAAY",
		issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
		redirectURL:           fmt.Sprintf("https://%s:8443/conformancesuite/callback", ListenHost),
		resourceIDs:           resourceIDs,
		apiVersion:            "v3.1",
		instructedAmount: models.InstructedAmount{
			Currency: "USD",
			Value:    "0.1",
		},
	}
	require.NoError(journey.SetConfig(config))
	require.Equal(config, journey.config)
}
