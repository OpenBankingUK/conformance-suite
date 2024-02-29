package discovery

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"github.com/magisterquis/connectproxy"
	"github.com/pkg/errors"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http/httpproxy"
	"golang.org/x/net/proxy"
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
	}, minSupportedTLSVersion}
}

func (v StdTLSValidator) ValidateTLSVersion(uri string) (TLSValidationResult, error) {
	var version uint16
	var strVersion string

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

	proxyState, err := getProxyConnectionState(addr, v.tlsConfig.InsecureSkipVerify)
	if err != nil {
		return TLSValidationResult{}, errors.Wrapf(err, "unable to get the proxy connection state for %s", addr)
	}

	if proxyState.Version != 0 {
		version = proxyState.Version
		strVersion, err = tlsVersionToString(proxyState.Version)
		if err != nil {
			return TLSValidationResult{}, errors.Wrapf(err, "unable to parse proxy proxy tls version `%d` to string", proxyState.Version)
		}
	} else {
		state, err := v.getConnectionState(addr)
		if err != nil {
			return TLSValidationResult{}, errors.Wrapf(err, "unable to get the connection state for %s", addr)
		}

		version = state.Version
		strVersion, err = tlsVersionToString(state.Version)
		if err != nil {
			return TLSValidationResult{}, errors.Wrapf(err, "unable to parse tls version `%d` to string", state.Version)
		}
	}

	return TLSValidationResult{
		version >= v.minSupportedTLSVersion,
		strVersion,
	}, nil
}

func (v StdTLSValidator) getConnectionState(addr string) (tls.ConnectionState, error) {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%s:%d", addr, 443)
	}

	conn, err := tls.Dial("tcp", addr, v.tlsConfig)
	if err != nil {
		return tls.ConnectionState{}, errors.Wrapf(err, "unable to detect tls version for hostname %s", addr)
	}

	defer conn.Close()

	state := conn.ConnectionState()
	return state, nil
}

func getProxyConnectionState(addr string, InsecureSkipVerify bool) (utls.ConnectionState, error) {
	var proxyStr string

	if httpsProxy := httpproxy.FromEnvironment().HTTPSProxy; httpsProxy != "" {
		proxyStr = httpsProxy
	}

	if httpProxy := httpproxy.FromEnvironment().HTTPProxy; httpProxy != "" {
		proxyStr = httpProxy
	}

	if proxyStr != "" {
		uris, err := url.Parse(proxyStr)
		if err != nil {
			return utls.ConnectionState{}, errors.Wrap(err, "unable to parse proxy URL")
		}

		proxyTLSConfig := &connectproxy.Config{
			InsecureSkipVerify: InsecureSkipVerify,
		}

		proxyDialer, err := connectproxy.NewWithConfig(uris, proxy.Direct, proxyTLSConfig)
		if err != nil {
			return utls.ConnectionState{}, errors.Wrapf(err, "unable to detect tls version for hostname %s using proxy", addr)
		}

		conn, err := proxyDialer.Dial("tcp", addr)
		if err != nil {
			return utls.ConnectionState{}, errors.Wrapf(err, "unable to detect tls version for hostname %s using proxy", addr)
		}

		defer conn.Close()

		config := utls.Config{InsecureSkipVerify: true, Renegotiation: utls.RenegotiateFreelyAsClient}
		uconn := utls.UClient(conn, &config, utls.HelloRandomizedALPN)

		err = uconn.Handshake()
		if err != nil {
			return utls.ConnectionState{}, errors.Wrapf(err, "unable to Handshake utls")
		}

		state := uconn.ConnectionState()

		return state, nil
	}

	return utls.ConnectionState{}, nil
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
