package authentication

import (
	"github.com/sirupsen/logrus"
)

const TlsClientAuth = "tls_client_auth"
const PrivateKeyJwt = "private_key_jwt"
const ClientSecretJwt = "client_secret_jwt"
const ClientSecretPost = "client_secret_post"
const ClientSecretBasic = "client_secret_basic"

// SuiteSupportedAuthMethodsMostSecureFirst -
// We have made our own determination of security offered by each auth method.
// It is not from a formal definition.
var SuiteSupportedAuthMethodsMostSecureFirst = []string{
	ClientSecretBasic,
}

func DefaultAuthMethod(openIDConfigAuthMethods []string, logger *logrus.Entry) string {
	return defaultAuthMethod(SuiteSupportedAuthMethodsMostSecureFirst, openIDConfigAuthMethods, logger)
}

// defaultAuthMethod - return first match from openIDConfigAuthMethods in
// suiteSupportedAuthMethods, else when no match return first method in
// suiteSupportedAuthMethods.
func defaultAuthMethod(suiteSupportedAuthMethods []string, openIDConfigAuthMethods []string, logger *logrus.Entry) string {
	intersection := make([]string, len(openIDConfigAuthMethods))
	i := 0
	for _, a := range suiteSupportedAuthMethods {
		for index, m := range openIDConfigAuthMethods {
			if a == m {
				intersection[i] = a
				openIDConfigAuthMethods[index] = ""
				i = i + 1
			}
		}
	}
	for _, m := range openIDConfigAuthMethods {
		if m != "" {
			logger.Infof("Invalid token endpoint auth method in OpenID config: %s", m)
		}
	}
	for _, method := range intersection {
		if method != "" {
			return method
		}
	}
	return suiteSupportedAuthMethods[0]
}
