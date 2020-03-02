package discovery

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type TLSValidator interface {
	ValidateTLSVersion(uri string) (TLSValidationResult, error)
}

type TLSValidationResult struct {
	Valid      bool
	TLSVersion string
}

type StdTLSValidator struct {
	tlsConfig              *tls.Config
	minSupportedTLSVersion uint16
}

type NullTLSValidator struct{}

func NewNullTLSValidator() NullTLSValidator {
	return NullTLSValidator{}
}

func (v NullTLSValidator) ValidateTLSVersion(uri string) (TLSValidationResult, error) {
	return TLSValidationResult{}, nil
}

func NewStdTLSValidator(minSupportedTLSVersion uint16) StdTLSValidator {
	return StdTLSValidator{&tls.Config{
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	}, minSupportedTLSVersion}
}

func (v StdTLSValidator) ValidateTLSVersion(uri string) (TLSValidationResult, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return TLSValidationResult{}, errors.Wrapf(err, "unable to parse the provided uri %s", uri)
	}
	// url.Parse only returns error for uri containing ASCII CTL bytes
	// in this case checking for blank URI will suffice
	if strings.TrimSpace(parsedURI.Host) == "" {
		return TLSValidationResult{}, fmt.Errorf("unable to parse the provided uri %s", uri)
	}
	addr := parsedURI.Host
	if !strings.Contains(parsedURI.Host, ":") {
		addr = fmt.Sprintf("%s:%d", parsedURI.Host, 443)
	}
	conn, err := tls.Dial("tcp", addr, v.tlsConfig)
	if err != nil {
		return TLSValidationResult{}, errors.Wrapf(err, "unable to detect tls version for hostname %s", parsedURI.Host)
	}
	defer conn.Close()
	state := conn.ConnectionState()
	strVersion, err := tlsVersionToString(state.Version)
	if err != nil {
		return TLSValidationResult{}, errors.Wrapf(err, "unable to parse tls version `%d` to string", state.Version)
	}

	return TLSValidationResult{
		state.Version >= v.minSupportedTLSVersion,
		strVersion,
	}, nil
}

func tlsVersionToString(v uint16) (string, error) {
	switch int(v) {
	case tls.VersionSSL30:
		return "SSL30", nil
	case tls.VersionTLS10:
		return "TLS10", nil
	case tls.VersionTLS11:
		return "TLS11", nil
	case tls.VersionTLS12:
		return "TLS12", nil
	case tls.VersionTLS13:
		return "TLS13", nil
	}
	return "", fmt.Errorf("unable to cast tls version `%d` to string", v)
}
