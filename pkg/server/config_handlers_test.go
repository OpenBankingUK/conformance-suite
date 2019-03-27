package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
)

var (
	privateKey = readFile("./testdata/certs/key.pem")
	publicKey  = readFile("./testdata/certs/cert.pem")
)

func TestValidateConfig(t *testing.T) {
	config := &GlobalConfiguration{
		SigningPrivate:          "---",
		SigningPublic:           "---",
		TransportPrivate:        "---",
		TransportPublic:         "--",
		ClientSecret:            "secret",
		AuthorizationEndpoint:   "https://server/auth",
		TokenEndpoint:           "https://server/token",
		TokenEndpointAuthMethod: "client_secret_basic",
		ResourceBaseURL:         "https://server",
		XFAPIFinancialID:        "2cfb31a3-5443-4e65-b2bc-ef8e00266a77",
		RedirectURL:             "https://localhost",
		Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
		ClientID:                "8672384e-9a33-439f-8924-67bb14340d71",
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
	c := &GlobalConfiguration{
		SigningPublic:           "------------",
		SigningPrivate:          "------------",
		TransportPrivate:        privateKey,
		TransportPublic:         publicKey,
		ClientID:                "client_id",
		ClientSecret:            "client_secret",
		ResourceBaseURL:         "http://server",
		TokenEndpoint:           "http://server",
		TokenEndpointAuthMethod: "client_secret_basic",
		AuthorizationEndpoint:   "http://server",
		RedirectURL:             "http://server",
		XFAPIFinancialID:        "123",
		Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
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

func readFile(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
	return string(file)
}

// TestServerConfigGlobalPostValid - tests /api/config/global
func TestServerConfigGlobalPostValid(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate:          privateKey,
		SigningPublic:           publicKey,
		TransportPrivate:        privateKey,
		TransportPublic:         publicKey,
		ClientID:                `8672384e-9a33-439f-8924-67bb14340d71`,
		ClientSecret:            `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
		TokenEndpoint:           `https://modelobank2018.o3bank.co.uk:4201/token`,
		TokenEndpointAuthMethod: "client_secret_basic",
		XFAPIFinancialID:        `0015800001041RHAAY`,
		RedirectURL:             fmt.Sprintf(`https://%s:8443/conformancesuite/callback`, ListenHost),
		AuthorizationEndpoint:   `https://modelobank2018.o3bank.co.uk:4201/token`,
		ResourceBaseURL:         `https://modelobank2018.o3bank.co.uk:4501`,
		Issuer:                  "https://modelobankauth2018.o3bank.co.uk:4101",
		ResourceIDs: model.ResourceIDs{
			AccountIDs:   []model.ResourceAccountID{{AccountID: "account-id"}},
			StatementIDs: []model.ResourceStatementID{{StatementID: "statement-id"}},
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

	// do assertions
	require.NotNil(body)
	bodyExpected := string(globalConfigurationJSON)
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)

	require.Equal(http.StatusCreated, code)
	require.Equal(expectedJsonHeaders, headers)
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
				SigningPrivate:          "------------",
				SigningPublic:           "------------",
				TransportPrivate:        privateKey,
				TransportPublic:         publicKey,
				ClientID:                "client_id",
				ClientSecret:            "client_secret",
				TokenEndpoint:           "http://server",
				TokenEndpointAuthMethod: "client_secret_basic",
				AuthorizationEndpoint:   "http://server",
				ResourceBaseURL:         "http://server",
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
			},
		},
		{
			name:               `InvalidTransport`,
			expectedBody:       `{"error": "error with transport certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"}`,
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
			},
		},
		{
			name:               `EmptyResourceIDs`,
			expectedBody:       `{"error": "resource_ids is empty"}`,
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
			},
		},
		{
			name:               `MissingResourcesAccountID`,
			expectedBody:       `{"error": "resource_ids.AccountIDs is empty"}`,
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
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
			},
		},
		{
			name:               `HasEmptyAccountID`,
			expectedBody:       `{"error": "resource_ids.AccountIDs contains an empty value at index 1"}`,
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
						{AccountID: ""},
					},
					StatementIDs: []model.ResourceStatementID{
						{StatementID: "statement-id"},
					},
				},
			},
		},
		{
			name:               `MissingResourcesStatementID`,
			expectedBody:       `{"error": "resource_ids.StatementIDs is empty"}`,
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
				},
			},
		},
		{
			name:               `HasEmptyStatementID`,
			expectedBody:       `{"error": "resource_ids.StatementIDs contains an empty value at index 0"}`,
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
						{StatementID: ""},
						{StatementID: "statement-id"},
					},
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

			configJson, err := json.MarshalIndent(testCase.config, ``, `  `)
			require.NoError(err)
			require.NotNil(configJson)

			code, body, headers := request(http.MethodPost, "/api/config/global", bytes.NewReader(configJson), server)

			require.NotNil(body)
			require.NoError(err)
			bodyActual := body.String()
			require.JSONEq(testCase.expectedBody, bodyActual)

			require.Equal(testCase.expectedStatusCode, code)
			require.Equal(expectedJsonHeaders, headers)
		})
	}
}

func testJourney() Journey {
	logger := nullLogger()
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	return NewJourney(logger, testGenerator, validatorEngine)
}
