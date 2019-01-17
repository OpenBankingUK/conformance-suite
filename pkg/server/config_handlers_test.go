package server

import (
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

// /api/config - POST - can POST config
func TestServerConfigPOSTCreatesProxy(t *testing.T) {
	require := require.New(t)

	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	mockedServer, serverURL := test.HTTPServer(http.StatusBadRequest, "body", nil)
	appConfig := appConfigJSONWithURL(serverURL)

	// assert server isn't started before call
	frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
	_, err := http.Get(frontendProxy.String())
	require.Error(err)
	require.Nil(server.proxy)

	// create the request to post the config
	// this should start the proxy
	req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfig))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// do the request
	server.ServeHTTP(rec, req)

	require.NotNil(rec.Body)
	require.JSONEq(appConfig, rec.Body.String())
	require.Equal(http.StatusOK, rec.Code)

	// check the proxy is up now, we should hit the forgerock server
	resp, err := http.Get(frontendProxy.String())
	require.NoError(err)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(err)
	require.Equal(http.StatusBadRequest, resp.StatusCode)
	require.Equal("body", string(body))
	require.NotNil(server.proxy)

	mockedServer.Close()
}

// /api/config - POST - cannot POST config twice without first deleting it
func TestServerConfigPOSTCannotPOSTConfigTwiceWithoutFirstDeletingIt(t *testing.T) {
	require := require.New(t)

	// Version helper
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	// assert server isn't started before call
	frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
	_, err := http.Get(frontendProxy.String())
	require.Error(err)

	// create the request to post the config
	// this should start the proxy
	req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// do the request
	server.ServeHTTP(rec, req)

	require.NotNil(rec.Body)
	require.JSONEq(appConfigJSON, rec.Body.String())
	require.Equal(http.StatusOK, rec.Code)

	// create another request to POST the config again
	// this should fail because a DELETE need to happen first.
	req = httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	// do the request
	server.ServeHTTP(rec, req)

	require.NotNil(rec.Body)
	require.JSONEq(
		`{"error":"listen tcp :8989: bind: address already in use"}`,
		rec.Body.String(),
	)
	require.Equal(http.StatusBadRequest, rec.Code)
}

// /api/config - DELETE - DELETE stops the proxy
func TestServerConfigDELETEStopsTheProxy(t *testing.T) {
	require := require.New(t)

	// Version helper
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()

	// assert server isn't started before call
	frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
	_, err := http.Get(frontendProxy.String())
	require.Error(err)

	// create the request to post the config
	// this should start the proxy
	req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// do the request
	server.ServeHTTP(rec, req)

	require.NotNil(rec.Body)
	require.JSONEq(appConfigJSON, rec.Body.String())
	require.Equal(http.StatusOK, rec.Code)

	// create request to delete config
	req = httptest.NewRequest(http.MethodDelete, "/api/config", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	// do the request
	server.ServeHTTP(rec, req)

	require.NotNil(rec.Body)
	require.Equal(
		"",
		rec.Body.String(),
	)
	require.Equal(http.StatusOK, rec.Code)

	// call proxy and assert it is no longer up
	// check the proxy is up now, we should hit the forgerock server
	resp, err := http.Get(frontendProxy.String())
	require.Equal(
		`Get http://0.0.0.0:8989/open-banking/v2.0/accounts: dial tcp 0.0.0.0:8989: connect: connection refused`,
		err.Error(),
	)
	require.Nil(resp)
}

// TestServerConfigGlobalPostValid - tests /api/config/global
func TestServerConfigGlobalPostValid(t *testing.T) {
	require := require.New(t)

	// Version helper
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate: `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		SigningPublic: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`,
		TransportPrivate: `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		TransportPublic: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`,
	}
	globalConfigurationJSON, err := json.MarshalIndent(globalConfiguration, ``, `  `)
	require.NoError(err)
	require.NotNil(globalConfigurationJSON)

	// make the request
	//
	// `?pretty` makes the JSON more readable in the event of a failure
	// see the example: https://echo.labstack.com/guide/response#json-pretty
	code, body, headers := request(
		http.MethodPost,
		"/api/config/global?pretty",
		strings.NewReader(string(globalConfigurationJSON)),
		server)

	// do assertions
	require.Equal(http.StatusCreated, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])

	require.NotNil(body)

	bodyExpected := string(globalConfigurationJSON)
	bodyActual := body.String()
	// do not use `require.Equal`.
	require.JSONEq(bodyExpected, bodyActual)
}

// TestServerConfigGlobalPostInvalidSigning - tests /api/config/global
func TestServerConfigGlobalPostInvalidSigning(t *testing.T) {
	require := require.New(t)

	// Version helper
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate: ``,
		SigningPublic:  ``,
		TransportPrivate: `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		TransportPublic: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`,
	}
	globalConfigurationJSON, err := json.MarshalIndent(globalConfiguration, ``, `  `)
	require.NoError(err)
	require.NotNil(globalConfigurationJSON)

	// make the request
	//
	// `?pretty` makes the JSON more readable in the event of a failure
	// see the example: https://echo.labstack.com/guide/response#json-pretty
	code, body, headers := request(
		http.MethodPost,
		"/api/config/global?pretty",
		strings.NewReader(string(globalConfigurationJSON)),
		server)

	// do assertions
	require.Equal(http.StatusBadRequest, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])

	require.NotNil(body)

	bodyExpected := `
	{
		"error": "error with signing certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"
	}
	`
	bodyActual := body.String()
	// do not use `require.Equal`.
	require.JSONEq(bodyExpected, bodyActual)
}

// TestServerConfigGlobalPostInvalidTransport - tests /api/config/global
func TestServerConfigGlobalPostInvalidTransport(t *testing.T) {
	require := require.New(t)

	// Version helper
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	globalConfiguration := &GlobalConfiguration{
		SigningPrivate: `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		SigningPublic: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`,
		TransportPrivate: ``,
		TransportPublic:  ``,
	}
	globalConfigurationJSON, err := json.MarshalIndent(globalConfiguration, ``, `  `)
	require.NoError(err)
	require.NotNil(globalConfigurationJSON)

	// make the request
	//
	// `?pretty` makes the JSON more readable in the event of a failure
	// see the example: https://echo.labstack.com/guide/response#json-pretty
	code, body, headers := request(
		http.MethodPost,
		"/api/config/global?pretty",
		strings.NewReader(string(globalConfigurationJSON)),
		server)

	// do assertions
	require.Equal(http.StatusBadRequest, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])

	require.NotNil(body)

	bodyExpected := `
	{
		"error": "error with transport certificate: error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key"
	}
	`
	bodyActual := body.String()
	// do not use `require.Equal`.
	require.JSONEq(bodyExpected, bodyActual)
}
