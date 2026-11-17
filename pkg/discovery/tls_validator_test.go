package discovery

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"sync/atomic"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStdTLSValidator(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	assert.Equal(t, uint16(tls.VersionTLS11), validator.minSupportedTLSVersion)
	assert.True(t, validator.tlsConfig.InsecureSkipVerify)
}

func TestValidateTLSVersionFailsOnInvalidURI(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	_, err := validator.ValidateTLSVersion("covfefe")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unable to parse the provided uri")
}

func TestValidateTLSVersionFailsOnInvalidHost(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	// using a non-tls server to trigger tls error
	srv, uri := test.HTTPServer(http.StatusServiceUnavailable, "", nil)
	defer srv.Close()
	_, err := validator.ValidateTLSVersion(uri)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unable to detect tls version for hostname")
}

func TestValidateTLSVersionFailsOnLowerVersion(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	srv, uri := test.HTTPSServer(&tls.Config{MinVersion: tls.VersionTLS10, MaxVersion: tls.VersionTLS10}, http.StatusServiceUnavailable, "", nil)
	defer srv.Close()
	r, err := validator.ValidateTLSVersion(uri)
	assert.False(t, r.Valid)
	assert.Empty(t, r.TLSVersion)
	assert.Contains(t, err.Error(), "unable to detect tls version for hostname 127.0.0.1")
	assert.Contains(t, err.Error(), "remote error: tls: protocol version not supported")
}

func TestValidateTLSVersionSucceeds(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	srv, uri := test.HTTPSServer(&tls.Config{MinVersion: tls.VersionTLS12}, http.StatusServiceUnavailable, "", nil)
	defer srv.Close()
	r, err := validator.ValidateTLSVersion(uri)
	result := r.TLSVersion
	assert.True(t, result == "TLS12" || result == "TLS13")
	assert.Nil(t, err)
}

func TestValidateTLSVersionSucceedsWhenPortIsntSpecified(t *testing.T) {
	validator := NewStdTLSValidator(tls.VersionTLS11)
	r, err := validator.ValidateTLSVersion("https://google.com")
	assert.True(t, r.Valid)
	assert.Nil(t, err)
}

// proxyProbes is how often each proxy scenario is probed, so that any nondeterministic
// handshake behaviour (e.g. a randomised ClientHello) shows up as a failure.
const proxyProbes = 50

// connectProxy is an in-process HTTP CONNECT proxy that counts the tunnels requested through it.
type connectProxy struct {
	*httptest.Server
	credentials string // "user:password" required via Proxy-Authorization, if not empty
	connects    atomic.Int32
}

func newConnectProxy(t *testing.T, useTLS bool, credentials string) *connectProxy {
	t.Helper()
	p := &connectProxy{credentials: credentials}
	p.Server = newQuietServer(http.HandlerFunc(p.serveConnect))
	if useTLS {
		p.StartTLS()
	} else {
		p.Start()
	}
	t.Cleanup(p.Close)
	return p
}

// newQuietServer returns an unstarted server that doesn't log connections closed mid-handshake by the probe.
func newQuietServer(handler http.Handler) *httptest.Server {
	srv := httptest.NewUnstartedServer(handler)
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	return srv
}

// startTLSTarget starts a TLS server using tlsConfig and returns its URL.
func startTLSTarget(t *testing.T, tlsConfig *tls.Config) string {
	t.Helper()
	srv := newQuietServer(http.NotFoundHandler())
	srv.TLS = tlsConfig
	srv.StartTLS()
	t.Cleanup(srv.Close)
	return srv.URL
}

func (p *connectProxy) serveConnect(w http.ResponseWriter, r *http.Request) {
	p.connects.Add(1)
	if p.credentials != "" && r.Header.Get("Proxy-Authorization") != "Basic "+base64.StdEncoding.EncodeToString([]byte(p.credentials)) {
		http.Error(w, "proxy authentication required", http.StatusProxyAuthRequired)
		return
	}
	target, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	client, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		_ = target.Close()
		return
	}
	_, _ = client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	go pipe(target, client)
	go pipe(client, target)
}

func pipe(dst, src net.Conn) {
	_, _ = io.Copy(dst, src)
	_ = dst.Close()
	_ = src.Close()
}

// useProxy points the validator at proxyURL, shielding the test from any proxy set in the environment.
func useProxy(t *testing.T, proxyURL string) {
	t.Helper()
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}
	t.Setenv("HTTPS_PROXY", proxyURL)
}

func TestValidateTLSVersionThroughProxy(t *testing.T) {
	tests := []struct {
		name    string
		server  *tls.Config
		version string
		valid   bool
	}{
		{"TLS 1.2 and 1.3", &tls.Config{MinVersion: tls.VersionTLS12}, "TLS13", true},
		{"TLS 1.2 only", &tls.Config{MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS12}, "TLS12", true},
		{"TLS 1.1 only", &tls.Config{MinVersion: tls.VersionTLS11, MaxVersion: tls.VersionTLS11}, "TLS11", true},
		{"TLS 1.0 only is reported as unsupported", &tls.Config{MinVersion: tls.VersionTLS10, MaxVersion: tls.VersionTLS10}, "TLS10", false},
		{"P-256 only", &tls.Config{CurvePreferences: []tls.CurveID{tls.CurveP256}}, "TLS13", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri := startTLSTarget(t, tc.server)
			proxy := newConnectProxy(t, false, "")
			useProxy(t, proxy.URL)

			validator := NewStdTLSValidator(tls.VersionTLS11)
			for i := 0; i < proxyProbes; i++ {
				r, err := validator.ValidateTLSVersion(uri)
				require.NoError(t, err, "probe %d", i)
				require.Equal(t, TLSValidationResult{Valid: tc.valid, TLSVersion: tc.version}, r, "probe %d", i)
			}
			assert.EqualValues(t, proxyProbes, proxy.connects.Load())
		})
	}
}

func TestValidateTLSVersionThroughProxyWithHybridPostQuantumServer(t *testing.T) {
	for _, mlkem := range []string{"0", "1"} {
		t.Run("tlsmlkem="+mlkem, func(t *testing.T) {
			t.Setenv("GODEBUG", "tlsmlkem="+mlkem)
			var offeredMLKEM atomic.Bool
			uri := startTLSTarget(t, &tls.Config{
				CurvePreferences: []tls.CurveID{tls.X25519MLKEM768, tls.X25519},
				GetConfigForClient: func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
					offeredMLKEM.Store(slices.Contains(hello.SupportedCurves, tls.X25519MLKEM768))
					return nil, nil
				},
			})
			proxy := newConnectProxy(t, false, "")
			useProxy(t, proxy.URL)

			validator := NewStdTLSValidator(tls.VersionTLS11)
			for i := 0; i < proxyProbes; i++ {
				r, err := validator.ValidateTLSVersion(uri)
				require.NoError(t, err, "probe %d", i)
				require.Equal(t, TLSValidationResult{Valid: true, TLSVersion: "TLS13"}, r, "probe %d", i)
			}
			// guards the scenario itself: the probe only offers the hybrid group when GODEBUG enables it
			assert.Equal(t, mlkem == "1", offeredMLKEM.Load())
		})
	}
}

func TestValidateTLSVersionThroughProxyWithCredentials(t *testing.T) {
	tests := []struct {
		name string
		user *url.Userinfo
		ok   bool
	}{
		{"valid credentials", url.UserPassword("alice", "s3cret"), true},
		{"wrong password", url.UserPassword("alice", "wrong"), false},
		{"no credentials", nil, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri := startTLSTarget(t, &tls.Config{MinVersion: tls.VersionTLS12})
			proxy := newConnectProxy(t, false, "alice:s3cret")
			proxyURL, err := url.Parse(proxy.URL)
			require.NoError(t, err)
			proxyURL.User = tc.user
			useProxy(t, proxyURL.String())

			r, err := NewStdTLSValidator(tls.VersionTLS11).ValidateTLSVersion(uri)

			assert.EqualValues(t, 1, proxy.connects.Load())
			if tc.ok {
				require.NoError(t, err)
				assert.Equal(t, TLSValidationResult{Valid: true, TLSVersion: "TLS13"}, r)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), "407")
			assert.Equal(t, TLSValidationResult{}, r)
		})
	}
}

func TestValidateTLSVersionThroughHTTPSProxy(t *testing.T) {
	uri := startTLSTarget(t, &tls.Config{MinVersion: tls.VersionTLS12})
	proxy := newConnectProxy(t, true, "")
	useProxy(t, proxy.URL)

	r, err := NewStdTLSValidator(tls.VersionTLS11).ValidateTLSVersion(uri)

	require.NoError(t, err)
	assert.Equal(t, TLSValidationResult{Valid: true, TLSVersion: "TLS13"}, r)
	assert.EqualValues(t, 1, proxy.connects.Load())
}
