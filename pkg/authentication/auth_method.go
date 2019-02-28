package authentication

import (
	"github.com/sirupsen/logrus"
)

const tls_client_auth = "tls_client_auth"
const private_key_jwt = "private_key_jwt"
const client_secret_jwt = "client_secret_jwt"
const client_secret_post = "client_secret_post"
const client_secret_basic = "client_secret_basic"

// SUITE_SUPPORTED_AUTH_METHODS_MOST_SECURE_FIRST -
// We have made our own determination of security offered by each auth method.
// It is not from a formal definition.
var SUITE_SUPPORTED_AUTH_METHODS_MOST_SECURE_FIRST = []string{
	client_secret_basic,
}

func DefaultAuthMethod(openIDConfigAuthMethods []string, logger *logrus.Entry) string {
	return defaultAuthMethod(SUITE_SUPPORTED_AUTH_METHODS_MOST_SECURE_FIRST, openIDConfigAuthMethods, logger)
}

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
