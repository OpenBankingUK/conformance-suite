package executors

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	resty "gopkg.in/resty.v1"
)

// MakeOzoneExecutor creates and ozone live executor
func MakeOzoneExecutor() (*OzoneExecutor, error) {
	err := ozoneCertSetup()
	if err != nil {
		return nil, err
	}
	return &OzoneExecutor{}, nil
}

func ozoneCertSetup() error {
	certConfig := CertConfig{}
	err := certConfig.getCertsFromEnvironment()
	if err != nil {
		return err
	}
	tlsConfig, err := certConfig.createMatlsTLSConfig()
	if err != nil {
		return err
	}
	resty.SetDebug(true)
	resty.SetTLSClientConfig(tlsConfig)
	return nil
}

// OzoneExecutor - passes request to ozone over matls
type OzoneExecutor struct {
	CertConfig    CertConfig
	SigningCert   authentication.Certificate
	TransportCert authentication.Certificate
}

// SetCertificates receives transport and signing certificates
func (e *OzoneExecutor) SetCertificates(certificateSigning, certificationTransport authentication.Certificate) {
	e.SigningCert = certificateSigning
	e.TransportCert = certificationTransport
}

// ExecuteTestCase - makes this a generic executor
func (e *OzoneExecutor) ExecuteTestCase(r *resty.Request, t *model.TestCase, ctx *model.Context) (*resty.Response, error) {
	e.appMsg("EXECUTE TEST CASE")
	e.appMsg(fmt.Sprintf("Request: %#v", r))
	resp, err := r.Execute(r.Method, r.URL)
	if err != nil {
		if resp.StatusCode() == 302 { // catch redirects and pass back as good response
			return resp, nil
		}
	}
	e.appMsg(fmt.Sprintf("Response: (%s)", resp.String()))
	return resp, err
}

func (e *OzoneExecutor) appMsg(msg string) {
	tracer.AppMsg("OZONE", msg, "")
}

func (e *OzoneExecutor) appErr(msg string) error {
	tracer.AppErr("OZONE", msg, "")
	return errors.New(msg)
}

func (e *OzoneExecutor) appEntry(msg string) {
	tracer.AppEntry("OZONE", msg)
}

func (e *OzoneExecutor) appExit(msg string) {
	tracer.AppExit("OZONE", msg)
}

// CertConfig Basic Certifcate Environment Configuration
type CertConfig struct {
	CertSigning   []byte
	CertTransport []byte
	KeySigning    []byte
	KeyTransport  []byte
}

func (c *CertConfig) getCertsFromEnvironment() error {
	signingPrivateString := os.Getenv("SIGNING_PRIVATE_KEY")
	if signingPrivateString == "" {
		return errors.New("Connot read Signing Private Key from environment : SIGNING_PRIVATE_KEY")
	}
	signingPublicString := os.Getenv("SIGNING_PUBLIC_CERT")
	if signingPublicString == "" {
		return errors.New("Connot read Signing Public Key from environment : SIGNING_PUBLIC_CERT")
	}

	transportPrivateString := os.Getenv("TRANSPORT_PRIVATE_KEY")
	if transportPrivateString == "" {
		return errors.New("Connot read Transport Private Key from environment : TRANSPORT_PRIVATE_KEY")
	}

	transportPublicString := os.Getenv("TRANSPORT_PUBLIC_CERT")
	if transportPublicString == "" {
		return errors.New("Connot read Transport Public Key from environment : TRANSPORT_PUBLIC_CERT")
	}

	c.CertSigning = []byte(signingPublicString)
	c.CertTransport = []byte(transportPublicString)
	c.KeySigning = []byte(signingPrivateString)
	c.KeyTransport = []byte(transportPrivateString)
	return nil

}

func (c *CertConfig) createMatlsTLSConfig() (*tls.Config, error) {
	transportCert, err := tls.X509KeyPair(c.CertTransport, c.KeyTransport)
	if err != nil {
		return nil, errors.New("createMATLSConfig X509Pair:" + err.Error())
	}
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.New("createMATLSConfig SystemCertPool:" + err.Error())
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{transportCert},
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
	return tlsConfig, nil

}
