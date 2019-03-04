package authentication

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"testing"
)

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenOneMatch(t *testing.T) {
	expected := SuiteSupportedAuthMethodsMostSecureFirst[0]

	openIDConfigAuthMethods := []string{
		"client_secret_post",
		"client_secret_jwt",
		expected,
		"private_key_jwt",
		"tls_client_auth",
	}

	actual := DefaultAuthMethod(openIDConfigAuthMethods, test.NullLogger())
	test.NewRequire(t).Equal(expected, actual)
}

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenNoMatch(t *testing.T) {
	openIDConfigAuthMethods := []string{
		"tls_client_auth",
	}
	actual := DefaultAuthMethod(openIDConfigAuthMethods, test.NullLogger())

	expected := SuiteSupportedAuthMethodsMostSecureFirst[0]
	test.NewRequire(t).Equal(expected, actual)
}

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenMultipleMatches(t *testing.T) {
	openIDConfigAuthMethods := []string{
		"client_secret_basic",
		"tls_client_auth",
	}
	mockSuiteSupported := []string{
		TlsClientAuth,
		PrivateKeyJwt,
		ClientSecretJwt,
		ClientSecretPost,
		ClientSecretBasic,
	}

	actual := defaultAuthMethod(mockSuiteSupported, openIDConfigAuthMethods, test.NullLogger())
	expected := mockSuiteSupported[0]
	test.NewRequire(t).Equal(expected, actual)
}
