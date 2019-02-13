package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"

	"github.com/pkg/errors"
)

var (
	privateKey = readFile("./testdata/certs/key.pem")
	publicKey  = readFile("./testdata/certs/cert.pem")
)

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

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate:   privateKey,
		SigningPublic:    publicKey,
		TransportPrivate: privateKey,
		TransportPublic:  publicKey,
		ClientID:         `8672384e-9a33-439f-8924-67bb14340d71`,
		ClientSecret:     `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
		TokenEndpoint:    `https://modelobank2018.o3bank.co.uk:4201/token`,
		XFAPIFinancialID: `0015800001041RHAAY`,
		RedirectURL:      `https://0.0.0.0:8443/conformancesuite/callback`,
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
		expectedBody       error
		expectedStatusCode int
		config             GlobalConfiguration
	}{
		{
			name:               `InvalidSigning`,
			expectedBody:       errors.New(`error with signing certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key`),
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   ``,
				SigningPublic:    ``,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
			},
		},
		{
			name:               `InvalidTransport`,
			expectedBody:       errors.New(`error with transport certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key`),
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: ``,
				TransportPublic:  ``,
			},
		},

		{
			name:               `MissingClientID`,
			expectedBody:       ErrEmptyClientID,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         ``,
				ClientSecret:     `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
				TokenEndpoint:    `https://modelobank2018.o3bank.co.uk:4201/token`,
				XFAPIFinancialID: `0015800001041RHAAY`,
				RedirectURL:      `https://0.0.0.0:8443/conformancesuite/callback`,
			},
		},
		{
			name:               `MissingClientSecret`,
			expectedBody:       ErrEmptyClientSecret,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         `8672384e-9a33-439f-8924-67bb14340d71`,
				ClientSecret:     ``,
				TokenEndpoint:    `https://modelobank2018.o3bank.co.uk:4201/token`,
				XFAPIFinancialID: `0015800001041RHAAY`,
				RedirectURL:      `https://0.0.0.0:8443/conformancesuite/callback`,
			},
		},
		{
			name:               `MissingTokenEndpoint`,
			expectedBody:       ErrEmptyTokenEndpoint,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         `8672384e-9a33-439f-8924-67bb14340d71`,
				ClientSecret:     `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
				TokenEndpoint:    ``,
				XFAPIFinancialID: `0015800001041RHAAY`,
				RedirectURL:      `https://0.0.0.0:8443/conformancesuite/callback`,
			},
		},
		{
			name:               `MissingXFAPIFinancialID`,
			expectedBody:       ErrEmptyXFAPIFinancialID,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         `8672384e-9a33-439f-8924-67bb14340d71`,
				ClientSecret:     `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
				TokenEndpoint:    `https://modelobank2018.o3bank.co.uk:4201/token`,
				XFAPIFinancialID: ``,
				RedirectURL:      `https://0.0.0.0:8443/conformancesuite/callback`,
			},
		},
		{
			name:               `MissingRedirectURL`,
			expectedBody:       ErrEmptyRedirectURL,
			expectedStatusCode: http.StatusBadRequest,
			config: GlobalConfiguration{
				SigningPrivate:   privateKey,
				SigningPublic:    publicKey,
				TransportPrivate: privateKey,
				TransportPublic:  publicKey,
				ClientID:         `8672384e-9a33-439f-8924-67bb14340d71`,
				ClientSecret:     `2cfb31a3-5443-4e65-b2bc-ef8e00266a77`,
				TokenEndpoint:    `https://modelobank2018.o3bank.co.uk:4201/token`,
				XFAPIFinancialID: `0015800001041RHAAY`,
				RedirectURL:      ``,
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			require := test.NewRequire(t)

			server := NewServer(nullLogger(), conditionalityCheckerMock{}, &mocks.Version{})
			defer func() {
				require.NoError(server.Shutdown(context.TODO()))
			}()
			require.NotNil(server)

			configJson, err := json.MarshalIndent(testCase.config, ``, `  `)
			require.NoError(err)
			require.NotNil(configJson)

			// make the request
			//
			// `?pretty` makes the JSON more readable in the event of a failure
			// see the example: https://echo.labstack.com/guide/response#json-pretty
			code, body, headers := request(http.MethodPost, "/api/config/global?pretty", bytes.NewReader(configJson), server)

			// do assertions
			require.NotNil(body)
			bodyExpected, err := json.MarshalIndent(NewErrorResponse(testCase.expectedBody), ``, `  `)
			require.NoError(err)
			bodyActual := body.String()
			require.JSONEq(string(bodyExpected), bodyActual)

			require.Equal(testCase.expectedStatusCode, code)
			require.Equal(http.Header{
				"Vary":         []string{"Accept-Encoding"},
				"Content-Type": []string{"application/json; charset=UTF-8"},
			}, headers)
		})
	}
}
