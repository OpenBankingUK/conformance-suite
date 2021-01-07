package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"
)

const (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvxbCEsZsrweSIQMpXQluO1anf8RYqEoVbdi0q09886Chl6N0
2f1UTQEiVDivBrXCawz00MoLSbrSiI05h7K6DH0+gUjdFfO4pVfbHsTktNbSy/Qs
L08KCBkxlOKHUCJi3AkhmX5orqzh1Nv89Q4sN1xFNlMXZ6CR0CJxtSVBPgqf4DeM
eIT4BMKHcyVXYEdquECFZvQWs6DZsS0WszIUNotscjvFuP68H0KQYH0vm05s/Yz2
VALEfx0SVaQmPGDu92SEnn6xY1pjzS602KlR1U7zjm0gRHBFdFG+IVik5y41+clg
L2TBGO8JFqeHd6laEbRwVBVXl3jmwd/q3kX6OwIDAQABAoIBAQCR69EcAUZxinh+
mSl3EIKK8atLGCcTrC8dCQU+ZJ7odFuxrnLHHHrJqvoKEpclqprioKw63G8uSGoJ
OL8b7tHAQ8v9ciTSZKE2Mhb0MirsJbgnYzhykAr7EDIanbny6a9Qk/CChFNwQDjc
EXnjsIT3aZC44U7YJXfz1rm6OM7Pjn6z8H4vYGRDOsYkhXvPfnPW8C2LFJVr9nvE
0gIAOVoGejEJrsJVK3Uj/nPcqSQYXmwEmtjtzOw7u6yp1b2VZEK7tR47HwJt6ltG
Z9zhpwhpvdOuXNMqMOYRf9bLBWnSqIlTHOO0UlAnyRCY1HxluZB7ZSg9VnoJDrD7
w+JqAGnBAoGBAO5qyIzjldwR004YjepmZfuX3PnGLZhzhmTTC7Pl9gqv1TvxfxvD
6yBFL2GrN1IcnrX9Qk2xncUAbpM989MF+EC7I4++1t1I6akUKFEDkfvQwQjCXfPS
Jv2rkwIVSkt8F0X/tOb13OeIiHuFVI/Bb9VoJSP/k4DfPV+/HnwBxvzLAoGBAM0u
b/rYfm5rb20/PKClUs154s0eKSokVogqiJkf+5qLsV+TD50JVZBVw8s4XM79iwQI
PyGY9nI1AvqG7yIzxSy5/Qk1+ZVdVYpmWIO5PnJ8TVraDVhCQ3fVz1uWtcyaqPVr
3QzdyvsEgFUGFItmRdhSvA8RGrpVCHTBzrDj3jpRAoGBAKNaSLS3jkstb3D3w+yR
YliisYX1cfIdXTyhmUgWTKD/3oLmsSdt8iC3JoKt1AaPk3Kv5ojjJG0BIcIC1ZeF
ZJW9Yt0vbXpKZcYyCHmRj6lQW6JLwiG3oH133A62VaQojq2oSONiG4wL8S9oqAqj
B6PZanEiwIaw7hU3FoTylstHAoGAFYvE0pCdZjb98njrgusZcN5VxLhgFj7On2no
AjxrjWUR8TleMF1kkM2Qy+xVQp85U+kRyBNp/cA3WduFjQ/mqrW1LpxuYxL0Ap6Q
uPRg7GDFNr8jG5uJvjHDnpiK6rtq9qqnAczgnc9xMnx699B7kSXO/b4MEnkPdENN
0yF6mqECgYA88UELxbhqMSdG24DX0zHXvkXLIml2JNVb54glFByIIem+acff9oG9
X5GajlBroPoKk7FgA9ouqcQMH66UnFi6qh07l0J2xb0aXP8yzLAGauVGTTNIQCR4
VpqyDpjlc1ZqfZWOrvwSrUH1mEkxbeVvQsOUja2Jvu+lc3Zo099ILw==
-----END RSA PRIVATE KEY-----

`
	publicKey = `-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQe/dw9alKTWAPhsHoLdkn+TANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTE2MDkyNTAwNDcxN1oXDTE3MDkyNTAwNDcx
N1owEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAL8WwhLGbK8HkiEDKV0JbjtWp3/EWKhKFW3YtKtPfPOgoZejdNn9VE0B
IlQ4rwa1wmsM9NDKC0m60oiNOYeyugx9PoFI3RXzuKVX2x7E5LTW0sv0LC9PCggZ
MZTih1AiYtwJIZl+aK6s4dTb/PUOLDdcRTZTF2egkdAicbUlQT4Kn+A3jHiE+ATC
h3MlV2BHarhAhWb0FrOg2bEtFrMyFDaLbHI7xbj+vB9CkGB9L5tObP2M9lQCxH8d
ElWkJjxg7vdkhJ5+sWNaY80utNipUdVO845tIERwRXRRviFYpOcuNfnJYC9kwRjv
CRanh3epWhG0cFQVV5d45sHf6t5F+jsCAwEAAaNLMEkwDgYDVR0PAQH/BAQDAgWg
MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwFAYDVR0RBA0wC4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQAdd3ZW6R4cImmxIzfoz7Ttq862
oOiyzFnisCxgNdA78epit49zg0CgF7q9guTEArXJLI+/qnjPPObPOlTlsEyomb2F
UOS+2hn/ZyU5/tUxhkeOBYqdEaryk6zF6vPLUJ5IphJgOg00uIQGL0UvupBLEyIG
Rsa/lKEtW5Z9PbIi9GeVn51U+9VMCYft/T7SDziKl7OcE/qoVh1G0/tTRkAqOqpZ
bzc8ssEhJVNZ/DO+uYHNYf/waB6NjfXQuTegU/SyxnawvQ4oBHIzyuWplGCcTlfT
IXsOQdJo2xuu8807d+rO1FpN8yWi5OF/0sif0RrocSskLAIL/PI1qfWuuPck
-----END CERTIFICATE-----

`

	timeFormat = "2006-01-02T15:04:05-07:00"
)

var (
	executionDateTime = time.Now().Add(24 * time.Hour).Format(timeFormat)
	paymentDateTime   = time.Now().Add(24 * time.Hour).Format(timeFormat)
)

func TestValidateConfig(t *testing.T) {
	config := &GlobalConfiguration{
		SigningPrivate:                "---",
		SigningPublic:                 "---",
		TransportPrivate:              "---",
		TransportPublic:               "--",
		ClientSecret:                  "secret",
		AuthorizationEndpoint:         "https://server/auth",
		TransactionFromDate:           defaultTxnFrom,
		TransactionToDate:             defaultTxnTo,
		TokenEndpoint:                 "https://server/token",
		TokenEndpointAuthMethod:       "client_secret_basic",
		ResponseType:                  "code id_token",
		ResourceBaseURL:               "https://server",
		XFAPIFinancialID:              "2cfb31a3-5443-4e65-b2bc-ef8e00266a77",
		RedirectURL:                   "https://localhost",
		Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
		ClientID:                      "8672384e-9a33-439f-8924-67bb14340d71",
		RequestObjectSigningAlgorithm: "PS256",
		ResourceIDs: model.ResourceIDs{
			AccountIDs:   []model.ResourceAccountID{{AccountID: "account-id"}},
			StatementIDs: []model.ResourceStatementID{{StatementID: "statement-id"}},
		},
	}

	ok, msg := validateConfig(config)

	assert.True(t, ok)
	assert.Empty(t, msg)
}

func configStubMissing(missingField string) GlobalConfiguration {
	creditorAccount := models.Payment{
		SchemeName:     "UK.OBIE.SortCodeAccountNumber",
		Identification: "20202010981789",
	}
	c := &GlobalConfiguration{
		SigningPublic:                 "------------",
		SigningPrivate:                "------------",
		TransportPrivate:              privateKey,
		TransportPublic:               publicKey,
		ClientID:                      "client_id",
		ClientSecret:                  "client_secret",
		ResourceBaseURL:               "http://server",
		ResponseType:                  "code id_token",
		TokenEndpoint:                 "http://server",
		TokenEndpointAuthMethod:       "client_secret_basic",
		TransactionFromDate:           defaultTxnFrom,
		TransactionToDate:             defaultTxnTo,
		AuthorizationEndpoint:         "http://server",
		RedirectURL:                   "http://server",
		XFAPIFinancialID:              "123",
		RequestObjectSigningAlgorithm: "PS256",
		Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
		CreditorAccount:               creditorAccount,
		InternationalCreditorAccount:  creditorAccount,
		ResourceIDs: model.ResourceIDs{
			AccountIDs:   []model.ResourceAccountID{{AccountID: "account-id"}},
			StatementIDs: []model.ResourceStatementID{{StatementID: "statement-id"}},
		},
	}
	v := reflect.ValueOf(c).Elem()
	f := v.FieldByName(missingField)
	f.SetString("")
	return *c
}

func TestValidateConfigTestsEmpty(t *testing.T) {
	testCases := []struct {
		name        string
		config      GlobalConfiguration
		expectedMsg string
	}{
		{
			name:        "missing signing private",
			config:      configStubMissing("SigningPrivate"),
			expectedMsg: "signing_private is empty",
		},
		{
			name:        "missing signing public",
			config:      configStubMissing("SigningPublic"),
			expectedMsg: "signing_public is empty",
		},
		{
			name:        "missing transport private",
			config:      configStubMissing("TransportPrivate"),
			expectedMsg: "transport_private is empty",
		},
		{
			name:        "missing transport public",
			config:      configStubMissing("TransportPublic"),
			expectedMsg: "transport_public is empty",
		},
		{
			name:        "missing client id",
			config:      configStubMissing("ClientID"),
			expectedMsg: "client_id is empty",
		},
		{
			name:        "missing client secret",
			config:      configStubMissing("ClientSecret"),
			expectedMsg: "client_secret is empty",
		},
		{
			name:        "missing token endpoint",
			config:      configStubMissing("TokenEndpoint"),
			expectedMsg: "token_endpoint is empty",
		},
		{
			name:        "response_type_missing",
			config:      configStubMissing("ResponseType"),
			expectedMsg: "response_type is empty",
		},
		{
			name:        "missing token endpoint auth method",
			config:      configStubMissing("TokenEndpointAuthMethod"),
			expectedMsg: "token_endpoint_auth_method is empty",
		},
		{
			name:        "missing client authorization_endpoint",
			config:      configStubMissing("AuthorizationEndpoint"),
			expectedMsg: "authorization_endpoint is empty",
		},
		{
			name:        "missing resource base URL",
			config:      configStubMissing("ResourceBaseURL"),
			expectedMsg: "resource_base_url is empty",
		},
		{
			name:        "missing client redirect_url",
			config:      configStubMissing("RedirectURL"),
			expectedMsg: "redirect_url is empty",
		},
		{
			name:        "missing client issuer",
			config:      configStubMissing("Issuer"),
			expectedMsg: "issuer is empty",
		},
		{
			name:        "missing x_fapi_financial_id id",
			config:      configStubMissing("XFAPIFinancialID"),
			expectedMsg: "x_fapi_financial_id is empty",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert := test.NewAssert(t)
			ok, msg := validateConfig(&testCase.config)
			assert.Equal(false, ok)
			assert.Equal(testCase.expectedMsg, msg)
		})
	}
}

// TestServerConfigGlobalPostValid - tests /api/config/global
func TestServerConfigGlobalPostValid(t *testing.T) {
	require := test.NewRequire(t)
	logger := logrus.New()
	server := NewServer(testJourney(), logger.WithField("component", "test-server"), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate:                privateKey,
		SigningPublic:                 publicKey,
		TransportPrivate:              privateKey,
		TransportPublic:               publicKey,
		ClientID:                      `8672384e-9a33-439f-8924-67bb14340d71`,
		ClientSecret:                  `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
		TokenEndpoint:                 `https://modelobank2018.o3bank.co.uk:4201/token`,
		TokenEndpointAuthMethod:       "client_secret_basic",
		TransactionFromDate:           defaultTxnFrom,
		TransactionToDate:             defaultTxnTo,
		ResponseType:                  "code id_token",
		XFAPIFinancialID:              `0015800001041RHAAY`,
		RedirectURL:                   fmt.Sprintf(`https://%s:8443/conformancesuite/callback`, ListenHost),
		AuthorizationEndpoint:         `https://modelobank2018.o3bank.co.uk:4201/token`,
		ResourceBaseURL:               `https://ob19-rs1.o3bank.co.uk:4501`,
		Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
		RequestObjectSigningAlgorithm: "PS256",
		ResourceIDs: model.ResourceIDs{
			AccountIDs:   []model.ResourceAccountID{{AccountID: "account-id"}},
			StatementIDs: []model.ResourceStatementID{{StatementID: "statement-id"}},
		},
		CreditorAccount: models.Payment{
			SchemeName:     "UK.OBIE.SortCodeAccountNumber",
			Identification: "20202010981789",
		},
		InternationalCreditorAccount: models.Payment{
			SchemeName:     "UK.OBIE.SortCodeAccountNumber",
			Identification: "20202010981789",
		},
		RequestedExecutionDateTime: executionDateTime,
		FirstPaymentDateTime:       paymentDateTime,
		PaymentFrequency:           models.PaymentFrequency("EvryDay"),
		CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
			SchemeName:     "UK.OBIE.SortCodeAccountNumber",
			Identification: "20202010981789",
			Name:           "Bob Stone",
		},
	}
	globalConfigurationJSON, err := json.MarshalIndent(globalConfiguration, ``, `  `)
	require.NoError(err)
	require.NotNil(globalConfigurationJSON)

	// make the request
	//
	// `?pretty` makes the JSON more readable in the event of a failure
	// see the example: https://echo.labstack.com/guide/response#json-pretty
	code, body, headers := request(http.MethodPost, "/api/config/global?pretty", bytes.NewReader(globalConfigurationJSON), server)
	// t.Log(body.String())

	// do assertions
	require.NotNil(body)
	bodyExpected := string(globalConfigurationJSON)
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)

	require.Equal(http.StatusCreated, code)
	require.Equal(expectedJsonHeaders(), headers)
}

// TestServerConfigGlobalPostInvalid - tests /api/config/global invalid cases.
func TestServerConfigGlobalPostInvalid(t *testing.T) {
	testCases := []struct {
		name               string
		expectedBody       string
		expectedStatusCode int
		config             GlobalConfiguration
	}{
		{
			name:               `InvalidSigning`,
			expectedBody:       `{"error": "error with signing certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                "------------",
				SigningPublic:                 "------------",
				TransportPrivate:              privateKey,
				TransportPublic:               publicKey,
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "http://server",
				TokenEndpointAuthMethod:       "client_secret_basic",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "http://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				RequestObjectSigningAlgorithm: "PS256",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `InvalidTransport`,
			expectedBody:       `{"error": "error with transport certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpoint:                 "token_endpoint",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				RequestObjectSigningAlgorithm: "PS256",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `EmptyResourceIDs`,
			expectedBody:       `{"error": "resource_ids is empty"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TokenEndpointAuthMethod:       "client_secret_basic",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				RequestObjectSigningAlgorithm: "PS256",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `MissingResourcesAccountID`,
			expectedBody:       `{"error": "resource_ids.AccountIDs is empty"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				RequestObjectSigningAlgorithm: "PS256",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `HasEmptyAccountID`,
			expectedBody:       `{"error": "resource_ids.AccountIDs contains an empty value at index 1"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				RequestObjectSigningAlgorithm: "PS256",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
						{AccountID: ""},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `MissingResourcesStatementID`,
			expectedBody:       `{"error": "resource_ids.StatementIDs is empty"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				RequestObjectSigningAlgorithm: "PS256",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `HasEmptyStatementID`,
			expectedBody:       `{"error": "resource_ids.StatementIDs contains an empty value at index 0"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				RequestObjectSigningAlgorithm: "PS256",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: ""},
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `invalid_credit_account_1`,
			expectedBody:       `{"error":"creditor_account: (identification: cannot be blank; scheme_name: cannot be blank.); international_creditor_account: (identification: cannot be blank; scheme_name: cannot be blank.)."}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              "--------------",
				TransportPublic:               "--------------",
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				RequestObjectSigningAlgorithm: "PS256",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: ""},
						{StatementID: "statement-id"},
					},
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `invalid_credit_account_2`,
			expectedBody:       `{"error":"creditor_account: (identification: cannot be blank.)."}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:                privateKey,
				SigningPublic:                 publicKey,
				TransportPrivate:              privateKey,
				TransportPublic:               publicKey,
				ClientID:                      "client_id",
				ClientSecret:                  "client_secret",
				TokenEndpoint:                 "token_endpoint",
				ResponseType:                  "code id_token",
				TokenEndpointAuthMethod:       "client_secret_basic",
				TransactionFromDate:           defaultTxnFrom,
				TransactionToDate:             defaultTxnTo,
				AuthorizationEndpoint:         "http://server",
				ResourceBaseURL:               "https://server",
				RedirectURL:                   "http://server",
				XFAPIFinancialID:              "123",
				Issuer:                        "https://modelobankauth2018.o3bank.co.uk:4101",
				RequestObjectSigningAlgorithm: "PS256",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: ""},
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName: "UK.OBIE.SortCodeAccountNumber",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `response_type_missing`,
			expectedBody:       `{"error":"response_type: cannot be blank."}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:          privateKey,
				SigningPublic:           publicKey,
				TransportPrivate:        "--------------",
				TransportPublic:         "--------------",
				ClientID:                "client_id",
				ClientSecret:            "client_secret",
				TokenEndpoint:           "token_endpoint",
				TokenEndpointAuthMethod: "client_secret_basic",
				AuthorizationEndpoint:   "http://server",
				ResourceBaseURL:         "https://server",
				RedirectURL:             "http://server",
				XFAPIFinancialID:        "123",
				Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `response_type_invalid`,
			expectedBody:       `{"error":"response_type: must be a valid value."}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:          privateKey,
				SigningPublic:           publicKey,
				TransportPrivate:        "--------------",
				TransportPublic:         "--------------",
				ClientID:                "client_id",
				ClientSecret:            "client_secret",
				TokenEndpoint:           "token_endpoint",
				ResponseType:            "INVALID",
				TokenEndpointAuthMethod: "client_secret_basic",
				AuthorizationEndpoint:   "http://server",
				ResourceBaseURL:         "https://server",
				RedirectURL:             "http://server",
				XFAPIFinancialID:        "123",
				Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				PaymentFrequency:           models.PaymentFrequency("EvryDay"),
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
		{
			name:               `payment_frequency_invalid`,
			expectedBody:       `{"error":"payment_frequency: must be in a valid format (^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$)."}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:          privateKey,
				SigningPublic:           publicKey,
				TransportPrivate:        "--------------",
				TransportPublic:         "--------------",
				ClientID:                "client_id",
				ClientSecret:            "client_secret",
				TokenEndpoint:           "token_endpoint",
				ResponseType:            "code id_token",
				TokenEndpointAuthMethod: "client_secret_basic",
				AuthorizationEndpoint:   "http://server",
				ResourceBaseURL:         "https://server",
				RedirectURL:             "http://server",
				XFAPIFinancialID:        "123",
				Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
				ResourceIDs: model.ResourceIDs{
					AccountIDs: []model.ResourceAccountID{
						{AccountID: "account-id"},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
				CreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				InternationalCreditorAccount: models.Payment{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
				},
				PaymentFrequency:           models.PaymentFrequency("INVALID"),
				RequestedExecutionDateTime: executionDateTime,
				FirstPaymentDateTime:       paymentDateTime,
				CBPIIDebtorAccount: discovery.CBPIIDebtorAccount{
					SchemeName:     "UK.OBIE.SortCodeAccountNumber",
					Identification: "20202010981789",
					Name:           "Bob Stone",
				},
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			require := test.NewRequire(t)

			server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
			defer func() {
				require.NoError(server.Shutdown(context.TODO()))
			}()
			require.NotNil(server)

			configJSON, err := json.MarshalIndent(testCase.config, ``, `  `)
			require.NoError(err)
			require.NotNil(configJSON)

			code, body, headers := request(http.MethodPost, "/api/config/global", bytes.NewReader(configJSON), server)

			require.NotNil(body)
			require.NoError(err)
			bodyActual := body.String()
			require.JSONEq(testCase.expectedBody, bodyActual)

			require.Equal(testCase.expectedStatusCode, code)
			require.Equal(expectedJsonHeaders(), headers)
		})
	}
}

func testJourney() Journey {
	logger := nullLogger()
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	return NewJourney(logger, testGenerator, validatorEngine, discovery.NewNullTLSValidator(), false)
}
