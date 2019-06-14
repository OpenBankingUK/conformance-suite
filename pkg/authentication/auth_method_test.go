package authentication

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenOneMatch(t *testing.T) {
	require := test.NewRequire(t)
	expected := SuiteSupportedAuthMethodsMostSecureFirst()[0]

	openIDConfigAuthMethods := []string{
		"client_secret_post",
		"client_secret_jwt",
		expected,
		"private_key_jwt",
		"tls_client_auth",
	}

	actual := DefaultAuthMethod(openIDConfigAuthMethods, test.NullLogger())
	require.Equal(expected, actual)
	require.Equal("tls_client_auth", actual)
}

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenNoMatch(t *testing.T) {
	require := test.NewRequire(t)
	openIDConfigAuthMethods := []string{
		"tls_client_auth",
	}
	actual := DefaultAuthMethod(openIDConfigAuthMethods, test.NullLogger())

	expected := SuiteSupportedAuthMethodsMostSecureFirst()[0]
	require.Equal(expected, actual)
	require.Equal("tls_client_auth", actual)
}

func TestDefaultAuthMethodReturnsFirstSuiteSupportedMethodWhenMultipleMatches(t *testing.T) {
	require := test.NewRequire(t)
	openIDConfigAuthMethods := []string{
		"client_secret_basic",
		"tls_client_auth",
	}
	mockSuiteSupported := []string{
		TlsClientAuth,
		PrivateKeyJwt,
		ClientSecretBasic,
	}

	actual := defaultAuthMethod(mockSuiteSupported, openIDConfigAuthMethods, test.NullLogger())
	expected := mockSuiteSupported[0]
	require.Equal(expected, actual)
	require.Equal("tls_client_auth", actual)
}
