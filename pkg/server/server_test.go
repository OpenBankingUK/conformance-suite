package server

// Note: Do not run the server tests in parallel.

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/OpenBankingUK/conformance-suite/pkg/client"

	"github.com/OpenBankingUK/conformance-suite/pkg/test"
	"github.com/OpenBankingUK/conformance-suite/pkg/version/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

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
	require.Len(headers, 9)
	require.Equal("text/html; charset=utf-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.Equal(bodyExpected, bodyActual)
}

func TestServerSkipperSwagger(t *testing.T) {
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
		context.SetPath(path)                // set path on the Context
		isSkipped := skipperSwagger(context) // check if the `skipper` skips the path
		require.Equal(shouldSkip, isSkipped)
	}
}

func TestServerSkipperGzip(t *testing.T) {
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
		"/api":                       false,
		"/swagger":                   false,
		"/api/export":                true,
		"/api/import":                true,
	}
	for path, shouldSkip := range paths {
		context.SetPath(path)             // set path on the Context
		isSkipped := skipperGzip(context) // check if the `skipper` skips the path
		require.Equal(shouldSkip, isSkipped)
	}
}

// TestServerHTTPS - tests that TLS works.
func TestServerHTTPS(t *testing.T) {
	require := test.NewRequire(t)

	certFile := "../../certs/conformancesuite_cert.pem"
	keyFile := "../../certs/conformancesuite_key.pem"
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
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
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(err)
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	require.True(ok)
	address := net.JoinHostPort("127.0.0.1", fmt.Sprint(tcpAddr.Port))
	require.NoError(listener.Close())

	go func() {
		require.EqualError(server.StartTLS(address, certFile, keyFile), "http: Server closed")
	}()
	time.Sleep(100 * time.Millisecond)

	url := fmt.Sprintf("https://127.0.0.1:%d/", tcpAddr.Port)
	res, err := httpClient.Get(url)
	require.NoError(err)
	require.NotNil(res)
	require.Equal(http.StatusOK, res.StatusCode)

	for _, version := range []uint16{tls.VersionTLS10, tls.VersionTLS11} {
		conn, err := tls.Dial("tcp", address, &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         version,
			MaxVersion:         version,
		})
		if err == nil {
			_ = conn.Close()
		}
		require.Error(err)
	}

	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		CipherSuites:       []uint16{tls.TLS_RSA_WITH_AES_128_GCM_SHA256},
	})
	if err == nil {
		_ = conn.Close()
	}
	require.Error(err)

	conn, err = tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	})
	require.NoError(err)
	require.NoError(conn.Close())

	conn, err = tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
	})
	require.NoError(err)
	require.NoError(conn.Close())

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
