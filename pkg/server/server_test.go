package server

// Note: Do not run the server tests in parallel.

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()

	// silence log output when running tests...
	logrus.SetLevel(logrus.WarnLevel)

	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})

	t.Run("NewServer() returns non-nil value", func(t *testing.T) {
		assert := test.NewAssert(t)
		assert.NotNil(server)
	})

	t.Run("GET / returns index.html", func(t *testing.T) {
		assert := test.NewAssert(t)
		code, body, _ := request(http.MethodGet, "/", nil, server)

		assert.Equal(true, strings.HasPrefix(body.String(), "<!DOCTYPE html>"))
		assert.Equal(http.StatusOK, code)
	})

	t.Run("GET /favicon.ico returns favicon.ico", func(t *testing.T) {
		assert := test.NewAssert(t)
		code, body, _ := request(http.MethodGet, "/favicon.ico", nil, server)

		assert.NotEmpty(body.String())
		assert.Equal(http.StatusOK, code)
	})

	require.NoError(t, server.Shutdown(context.TODO()))
}

// TestServerConformanceSuiteCallback - Test that `/conformancesuite/callback` returns `./web/dist/index.html`.
func TestServerConformanceSuiteCallback(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
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
	require := test.NewRequire(t)

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
	require := test.NewRequire(t)

	certFile := "../../certs/conformancesuite_cert.pem"
	keyFile := "../../certs/conformancesuite_key.pem"
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient := client.NewHTTPClientWithTransport(client.DefaultTimeout, transport)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// For how random port works: https://github.com/labstack/echo/issues/1065#issuecomment-367961653

	// Start HTTPS server
	go func() {
		require.EqualError(server.StartTLS(":0", certFile, keyFile), "http: Server closed")
	}()
	time.Sleep(100 * time.Millisecond)

	tcpAddr, ok := server.TLSListener.Addr().(*net.TCPAddr)
	require.True(ok)
	url := fmt.Sprintf("https://localhost:%d/", tcpAddr.Port)
	res, err := httpClient.Get(url)
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

	return rec.Code, rec.Body, rec.Header()
}

// nullLogger - create a logger that discards output.
func nullLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger.WithField("app", "test")
}
