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
	return StdTLSValidator{&tls.Config{InsecureSkipVerify: true}, minSupportedTLSVersion}
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
	conn, err := tls.Dial("tcp", parsedURI.Host, v.tlsConfig)
	if err != nil {
		return TLSValidationResult{}, errors.Wrapf(err, "unable to detect tls version for hostname %s", parsedURI.Host)
	}
	state := conn.ConnectionState()

	return TLSValidationResult{
		state.Version >= v.minSupportedTLSVersion,
		string(tls.VersionTLS12),
	}, nil
}
