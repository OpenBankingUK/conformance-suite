package server

// Note: Do not run the server tests in parallel.

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// conditionalityCheckerMock - implements model.ConditionalityChecker interface for tests
type conditionalityCheckerMock struct {
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsOptional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsMandatory true for POST /account-access-consents, false for all other endpoint/methods.
func (c conditionalityCheckerMock) IsMandatory(method, endpoint string, specification string) (bool, error) {
	if method == "POST" && endpoint == "/account-access-consents" {
		return true, nil
	}
	return false, nil
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsConditional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsPresent true for valid GET/POST/DELETE endpoints.
func (c conditionalityCheckerMock) IsPresent(method, endpoint string, specification string) (bool, error) {
	if method == "GET" || method == "POST" || method == "DELETE" {
		return true, nil
	}
	return false, nil
}

func (c conditionalityCheckerMock) MissingMandatory(endpoints []model.Input, specification string) ([]model.Input, error) {
	return []model.Input{}, nil
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()

	// silence log output when running tests...
	logrus.SetLevel(logrus.WarnLevel)

	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &mocks.Version{})

	t.Run("NewServer() returns non-nil value", func(t *testing.T) {
		assert.NotNil(t, server)
	})

	t.Run("GET / returns index.html", func(t *testing.T) {
		code, body, _ := request(http.MethodGet, "/", nil, server)

		assert.Equal(t, true, strings.HasPrefix(body.String(), "<!DOCTYPE html>"))
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("GET /favicon.ico returns favicon.ico", func(t *testing.T) {
		code, body, _ := request(http.MethodGet, "/favicon.ico", nil, server)

		assert.NotEmpty(t, body.String())
		assert.Equal(t, http.StatusOK, code)
	})

	require.NoError(t, server.Shutdown(context.TODO()))
}

// TestServerConformanceSuiteCallback - Test that `/conformancesuite/callback` returns `./web/dist/index.html`.
func TestServerConformanceSuiteCallback(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// read the file we expect to be served.
	bytes, err := ioutil.ReadFile("./web/dist/index.html")
	require.NoError(err)
	bodyExpected := string(bytes)

	code, body, headers := request(
		http.MethodGet,
		`/conformancesuite/callback`,
		nil,
		server)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 5)
	require.Equal("text/html; charset=utf-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.Equal(bodyExpected, bodyActual)
}

func TestServerSkipper(t *testing.T) {
	require := require.New(t)

	echo := echo.New()
	context := echo.AcquireContext()
	defer func() {
		echo.ReleaseContext(context)
	}()

	paths := map[string]bool{
		"/index.html":                false,
		"/index.js":                  false,
		"/conformancesuite/callback": false,
		"/ipa":                       false,
		"/reggaws":                   false,
		"/api":                       true,
		"/swagger":                   true,
	}
	for path, shouldSkip := range paths {
		context.SetPath(path)         // set path on the Context
		isSkipped := skipper(context) // check if the `skipper` skips the path
		require.Equal(shouldSkip, isSkipped)
	}
}

// TestServerHTTPS - tests that TLS works.
func TestServerHTTPS(t *testing.T) {
	require := require.New(t)

	certFile := "../../certs/conformancesuite_cert.pem"
	keyFile := "../../certs/conformancesuite_key.pem"
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)
	server.HideBanner = true

	// Start HTTPS server
	go func() {
		require.EqualError(server.StartTLS(":443", certFile, keyFile), "http: Server closed")
	}()
	time.Sleep(100 * time.Millisecond)

	res, err := client.Get("https://localhost/")
	require.NoError(err)
	require.NotNil(res)
	require.Equal(http.StatusOK, res.StatusCode)

	require.NotNil(res.Body)
	defer func() {
		require.NoError(res.Body.Close())
	}()
	bodyActual, err := ioutil.ReadAll(res.Body)
	require.NoError(err)

	bodyExpected, err := ioutil.ReadFile("./web/dist/index.html")
	require.NoError(err)

	require.Equal(string(bodyExpected), string(bodyActual))
}

// Generic util function for making test requests.
func request(method, path string, body io.Reader, server *Server) (int, *bytes.Buffer, http.Header) {
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	server.ServeHTTP(rec, req)

	return rec.Code, rec.Body, rec.HeaderMap
}

// nullLogger - create a logger that discards output.
func nullLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger.WithField("app", "test")
}
