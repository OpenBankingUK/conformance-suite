package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
)

var (
	privateKey = readFile("./testdata/certs/key.pem")
	publicKey  = readFile("./testdata/certs/cert.pem")
)

func TestValidateConfig(t *testing.T) {
	config := &GlobalConfiguration{
		SigningPrivate:        "---",
		SigningPublic:         "---",
		TransportPrivate:      "---",
		TransportPublic:       "--",
		ClientSecret:          "secret",
		AuthorizationEndpoint: "https://server/auth",
		TokenEndpoint:         "https://server/token",
		ResourceBaseURL:       "https://server",
		XFAPIFinancialID:      "2cfb31a3-5443-4e65-b2bc-ef8e00266a77",
		RedirectURL:           "https://localhost",
		Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
		ClientID:              "8672384e-9a33-439f-8924-67bb14340d71",
	}

	ok, msg := validateConfig(config)

	assert.True(t, ok)
	assert.Empty(t, msg)
}

func TestValidateConfigTestsEmpty(t *testing.T) {
	testCases := []struct {
		name        string
		config      GlobalConfiguration
		expectedOk  bool
		expectedMsg string
	}{
		{
			name: "missing signing private",
			config: GlobalConfiguration{
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "signing_private is empty",
		},
		{
			name: "missing signing public",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "signing_public is empty",
		},
		{
			name: "missing transport private",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "transport_private is empty",
		},
		{
			name: "missing transport public",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "transport_public is empty",
		},
		{
			name: "missing client id",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "client_id is empty",
		},
		{
			name: "missing client secret",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "client_secret is empty",
		},
		{
			name: "missing token endpoint",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "token_endpoint is empty",
		},
		{
			name: "missing client authorization_endpoint",
			config: GlobalConfiguration{
				SigningPrivate:   `------------`,
				SigningPublic:    `------------`,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         "client_id",
				ClientSecret:     "client_secret",
				TokenEndpoint:    "http://server",
				RedirectURL:      "http://server",
				XFAPIFinancialID: "123",
				Issuer:           "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "authorization_endpoint is empty",
		},
		{
			name: "missing client redirect_url",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				ResourceBaseURL:       "http://server",
				AuthorizationEndpoint: "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "redirect_url is empty",
		},
		{
			name: "missing client redirect_url",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				ResourceBaseURL:       "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "redirect_url is empty",
		},
		{
			name: "missing client issuer",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				ResourceBaseURL:       "http://server",
				XFAPIFinancialID:      "123",
			},
			expectedOk:  false,
			expectedMsg: "issuer is empty",
		},
		{
			name: "missing x_fapi_financial_id id",
			config: GlobalConfiguration{
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				ResourceBaseURL:       "http://server",
				AuthorizationEndpoint: "http://server",
				RedirectURL:           "http://server",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
			expectedOk:  false,
			expectedMsg: "x_fapi_financial_id is empty",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			assert := test.NewAssert(t)
			ok, msg := validateConfig(&testCase.config)
			assert.Equal(testCase.expectedOk, ok)
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
		SigningPrivate:        privateKey,
		SigningPublic:         publicKey,
		TransportPrivate:      privateKey,
		TransportPublic:       publicKey,
		ClientID:              `8672384e-9a33-439f-8924-67bb14340d71`,
		ClientSecret:          `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
		TokenEndpoint:         `https://modelobank2018.o3bank.co.uk:4201/token`,
		XFAPIFinancialID:      `0015800001041RHAAY`,
		RedirectURL:           fmt.Sprintf(`https://%s:8443/conformancesuite/callback`, ListenHost),
		AuthorizationEndpoint: `https://modelobank2018.o3bank.co.uk:4201/token`,
		ResourceBaseURL:       `https://modelobank2018.o3bank.co.uk:4501`,
		Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
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
	require.Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
	}, headers)
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
				SigningPrivate:        `------------`,
				SigningPublic:         `------------`,
				TransportPrivate:      privateKey,
				TransportPublic:       publicKey,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "http://server",
				AuthorizationEndpoint: "http://server",
				ResourceBaseURL:       "http://server",
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
			},
		},
		{
			name:               `InvalidTransport`,
			expectedBody:       `{"error": "error with transport certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"}`,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:        privateKey,
				SigningPublic:         publicKey,
				TransportPrivate:      `--------------`,
				TransportPublic:       `--------------`,
				ClientID:              "client_id",
				ClientSecret:          "client_secret",
				TokenEndpoint:         "token_endpoint",
				AuthorizationEndpoint: "http://server",
				ResourceBaseURL:       `https://server`,
				RedirectURL:           "http://server",
				XFAPIFinancialID:      "123",
				Issuer:                "https://modelobankauth2018.o3bank.co.uk:4101",
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
			require.Equal(http.Header{
				"Vary":         []string{"Accept-Encoding"},
				"Content-Type": []string{"application/json; charset=UTF-8"},
			}, headers)
		})
	}
}

func testJourney() Journey {
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	return NewJourney(nullLogger(), testGenerator, validatorEngine)
}
