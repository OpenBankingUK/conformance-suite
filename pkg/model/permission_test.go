package model

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDefaultForEndpointErrNoDefaults(t *testing.T) {
	var data []permission
	r := newStandardPermissionsWithOptions(data)

	_, err := r.defaultForEndpoint("/home")

	assert.EqualError(t, err, "no default permissions found")
}

func TestDefaultForEndpointOneFoundDefaults(t *testing.T) {
	data := []permission{
		{
			Code:      "read",
			Default:   false,
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	perms, err := r.defaultForEndpoint("/home")

	require.NoError(t, err)
	assert.Equal(t, []Code{"read"}, perms)
}

func TestDefaultForEndpointUsesDefault(t *testing.T) {
	data := []permission{
		{
			Code:      "read",
			Default:   false,
			Endpoints: []string{"/home"},
		},
		{
			Code:      "write",
			Default:   true,
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	perms, err := r.defaultForEndpoint("/home")

	require.NoError(t, err)
	assert.Equal(t, []Code{"write"}, perms)
}

func TestDefaultForErrIfNoDefaultAndMoreThenOne(t *testing.T) {
	data := []permission{
		{
			Code:      "read",
			Default:   false,
			Endpoints: []string{"/home"},
		},
		{
			Code:      "write",
			Default:   false,
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	perms, err := r.defaultForEndpoint("/home")

	assert.EqualError(t, err, "no default permissions found, but found more than one")
	assert.Equal(t, []Code{}, perms)
}

func TestDefaultForTwoDefaultsReturnBoth(t *testing.T) {
	data := []permission{
		{
			Code:      "write",
			Default:   true,
			Endpoints: []string{"/home"},
		},
		{
			Code:      "read",
			Default:   true,
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	perms, err := r.defaultForEndpoint("/home")

	require.NoError(t, err)
	assert.Equal(t, []Code{"write", "read"}, perms)
}

func TestPermissionsForEndpoint(t *testing.T) {
	data := []permission{
		{
			Code:      "a",
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	search := r.permissionsForEndpoint("/home")

	assert.Len(t, search, 1)
	assert.Equal(t, Code("a"), search[0].Code)
}

func TestPermissionsForEndpointMultiplePerms(t *testing.T) {
	data := []permission{
		{
			Code:      "a",
			Endpoints: []string{"/home"},
		},
		{
			Code:      "b",
			Endpoints: []string{"/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	search := r.permissionsForEndpoint("/home")

	assert.Len(t, search, 2)
	assert.Equal(t, Code("a"), search[0].Code)
	assert.Equal(t, Code("b"), search[1].Code)
}

func TestPermissionsForEndpointDifferentEndpoints(t *testing.T) {
	data := []permission{
		{
			Code:      "a",
			Endpoints: []string{"/home"},
		},
		{
			Code:      "b",
			Endpoints: []string{"/home2"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	search := r.permissionsForEndpoint("/home")

	assert.Len(t, search, 1)
	assert.Equal(t, Code("a"), search[0].Code)
}

func TestPermissionsForEndpointMultipleEndpointsURL(t *testing.T) {
	data := []permission{
		{
			Code:      "a",
			Endpoints: []string{"/home"},
		},
		{
			Code:      "b",
			Endpoints: []string{"/home2", "/home"},
		},
	}
	r := newStandardPermissionsWithOptions(data)

	search := r.permissionsForEndpoint("/home")

	assert.Len(t, search, 2)
	assert.Equal(t, Code("a"), search[0].Code)
	assert.Equal(t, Code("b"), search[1].Code)
}

func TestStaticPermissionsHaveNotChanged(t *testing.T) {
	expected, err := json.MarshalIndent(staticApiPermissions, "", "    ")
	require.NoError(t, err)

	goldenFile := filepath.Join("testdata", "permissions.golden")
	if *update {
		t.Log("update golden file")
		require.NoError(t, ioutil.WriteFile(goldenFile, expected, 0644), "failed to update golden file")
	}

	perms, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err, "failed reading .golden")

	assert.JSONEq(t, string(expected), string(perms))
}

// This is a test of our matching function logic operating
// over our static permissions configuration.
//
// The examples below are intended to be realistic examples
// returning permissions based on the rules defined in the
// Accounts API specification. See:
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1#AccountandTransactionAPISpecification-v3.1-Permissions
func TestStaticPermissionsDefaultEnpointMatchingIntegration(t *testing.T) {
	config := newStandardPermissions()

	t.Run("when single default permission code", func(t *testing.T) {
		permissions, err := config.defaultForEndpoint("/accounts")
		assert.NoError(t, err)

		assert.Len(t, permissions, 1)
		// At a minimum to access the "/accounts" endpoint you need
		// either "ReadAccountsBasic" or "ReadAccountsDetail".
		// So this default permission makes sense:
		assert.Equal(t, Code("ReadAccountsBasic"), permissions[0])
	})

	t.Run("when multiple default permission codes", func(t *testing.T) {
		permissions, err := config.defaultForEndpoint("/accounts/{AccountId}/statements/{StatementId}/transactions")
		assert.NoError(t, err)

		assert.Len(t, permissions, 5)
		// At a minimum to access the "/accounts/{AccountId}/statements/{StatementId}/transactions"
		// endpoint you need
		// either "ReadAccountsBasic" or "ReadAccountsDetail" AND
		// either "ReadStatementsBasic" or "ReadStatementsDetail" AND
		// either "ReadTransactionsBasic" or "ReadTransactionsDetail" AND
		// one or more of "ReadTransactionsDebits" or "ReadTransactionsCredits".
		// So these default permissions make sense:
		assert.Equal(t, Code("ReadAccountsBasic"), permissions[0])
		assert.Equal(t, Code("ReadTransactionsBasic"), permissions[1])
		assert.Equal(t, Code("ReadTransactionsCredits"), permissions[2])
		assert.Equal(t, Code("ReadTransactionsDebits"), permissions[3])
		assert.Equal(t, Code("ReadStatementsBasic"), permissions[4])
	})
}
