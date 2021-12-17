package discovery

import (
	"crypto/tls"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "TLS10", r.TLSVersion)
	assert.NoError(t, err)
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
