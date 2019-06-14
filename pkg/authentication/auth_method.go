package authentication

// UK Open Banking OIDC Security Profile
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/83919096/Open+Banking+Security+Profile+-+Implementer+s+Draft+v1.1.2.
// https://bitbucket.org/openid/OBUK/raw/1.1.2/uk-openbanking-security-profile.md?_=1556543617032
//
//    1. shall authenticate the confidential client at the Token Endpoint using one of the following methods:
//     1. `tls_client_auth` as per [MTLS] (Recommended); or
//     2. `client_secret_basic` or `client_secret_post` provided the client identifier matches the client identifier bound to the underlying mutually authenticated TLS session (Allowed); or
//     3. `private_key_jwt` or 'client_secret_jwt` (Recommended);

import (
	"github.com/sirupsen/logrus"
)

// token_endpoint_auth_methods_supported
const (
	TlsClientAuth     = "tls_client_auth"
	PrivateKeyJwt     = "private_key_jwt"
	ClientSecretBasic = "client_secret_basic"
)

// SuiteSupportedAuthMethodsMostSecureFirst -
// We have made our own determination of security offered by each auth method.
// It is not from a formal definition.
func SuiteSupportedAuthMethodsMostSecureFirst() []string {
	return []string{
		TlsClientAuth,
		PrivateKeyJwt,
		ClientSecretBasic,
	}
}

func DefaultAuthMethod(openIDConfigAuthMethods []string, logger *logrus.Entry) string {
	return defaultAuthMethod(SuiteSupportedAuthMethodsMostSecureFirst(), openIDConfigAuthMethods, logger)
}

// defaultAuthMethod - return first match from openIDConfigAuthMethods in
// suiteSupportedAuthMethods, else when no match return first method in
// suiteSupportedAuthMethods.
func defaultAuthMethod(suiteSupportedAuthMethods []string, openIDConfigAuthMethods []string, logger *logrus.Entry) string {
	intersection, remaining := intersectionAndRemaining(openIDConfigAuthMethods, suiteSupportedAuthMethods)

	for _, method := range remaining {
		if method != "" {
			logger.WithFields(logrus.Fields{
				"method": method,
			}).Info("Invalid 'token_endpoint_auth_methods_supported' in OpenIDConfiguration")
		}
	}
	for _, methodMatch := range intersection {
		if methodMatch != "" {
			return methodMatch
		}
	}

	return suiteSupportedAuthMethods[0]
}

func intersectionAndRemaining(openIDConfigAuthMethods []string, suiteSupportedAuthMethods []string) ([]string, []string) {
	remaining := make([]string, len(openIDConfigAuthMethods))
	copy(remaining, openIDConfigAuthMethods)

	intersection := make([]string, len(openIDConfigAuthMethods))
	i := 0
	for _, supportedAuthMethod := range suiteSupportedAuthMethods {
		for index, suppliedAuthMethod := range openIDConfigAuthMethods {
			if supportedAuthMethod == suppliedAuthMethod {
				intersection[i] = supportedAuthMethod
				remaining[index] = ""
				i += 1
			}
		}
	}
	return intersection, remaining
}
