package executors

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	resty "gopkg.in/resty.v1"
)

// MakeExecutor creates an executor
func MakeExecutor() *Executor {
	return &Executor{}
}

// Executor - passes request to system under test across an matls connection
type Executor struct {
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

// SetCertificates receives transport and signing certificates
func (e *Executor) SetCertificates(certificateSigning, certificationTransport authentication.Certificate) error {
	e.SigningCert = certificateSigning
	e.TransportCert = certificationTransport
	return e.setupTLSCertificate(e.TransportCert.TLSCert())
}

// ExecuteTestCase - makes this a generic executor
func (e *Executor) ExecuteTestCase(r *resty.Request, t *model.TestCase, ctx *model.Context) (*resty.Response, error) {
	e.appMsg("Execute Testcase")
	e.appMsg(fmt.Sprintf("Request: %#v", r))
	resp, err := r.Execute(r.Method, r.URL)
	if err != nil {
		if resp.StatusCode() == http.StatusFound { // catch status code 302 redirects and pass back as good response
			return resp, nil
		}
	}
	t.ResponseTime = resp.Time()
	t.ResponseSize = len(resp.Body())
	e.appMsg(fmt.Sprintf("Response: (%s)", resp.String()))
	return resp, err
}

func (e *Executor) setupTLSCertificate(tlsCert tls.Certificate) error {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return errors.New("setupTLSCertificate SystemCertPool:" + err.Error())
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionSSL30,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, // not available by default however used by OB
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	}
	tlsConfig.BuildNameToCertificate()
	resty.SetTLSClientConfig(tlsConfig)
	resty.SetDebug(false)
	return nil

}

func (e *Executor) appMsg(msg string) {
	tracer.AppMsg("OZONE", msg, "")
}

func (e *Executor) appErr(msg string) error {
	tracer.AppErr("OZONE", msg, "")
	return errors.New(msg)
}

func (e *Executor) appEntry(msg string) {
	tracer.AppEntry("OZONE", msg)
}

func (e *Executor) appExit(msg string) {
	tracer.AppExit("OZONE", msg)
}
