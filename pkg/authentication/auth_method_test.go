package authentication

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"testing"
)

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenOneMatch(t *testing.T) {
	expected := SUITE_SUPPORTED_AUTH_METHODS_MOST_SECURE_FIRST[0]

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

	expected := SUITE_SUPPORTED_AUTH_METHODS_MOST_SECURE_FIRST[0]
	test.NewRequire(t).Equal(expected, actual)
}

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenMultipleMatches(t *testing.T) {
	openIDConfigAuthMethods := []string{
		"client_secret_basic",
		"tls_client_auth",
	}
	mockSuiteSupported := []string{
		tls_client_auth,
		private_key_jwt,
		client_secret_jwt,
		client_secret_post,
		client_secret_basic,
	}

	actual := defaultAuthMethod(mockSuiteSupported, openIDConfigAuthMethods, test.NullLogger())
	expected := mockSuiteSupported[0]
	test.NewRequire(t).Equal(expected, actual)
}
