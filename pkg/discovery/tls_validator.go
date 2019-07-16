package discovery

import (
	"crypto/tls"
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

func NewStdTLSValidator(minSupportedTLSVersion uint16) StdTLSValidator {
	return StdTLSValidator{&tls.Config{InsecureSkipVerify: true}, minSupportedTLSVersion}
}

func (s StdTLSValidator) ValidateTLSVersion(uri string) (TLSValidationResult, error) {
	conn, err := tls.Dial("tcp", uri, s.tlsConfig)
	if err != nil {
		return TLSValidationResult{}, err
	}
	state := conn.ConnectionState()

	return TLSValidationResult{
		state.Version >= s.minSupportedTLSVersion,
		string(tls.VersionTLS12),
	}, nil
}
