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
	assert.Equal(t, Code("read"), perms)
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
	assert.Equal(t, Code("write"), perms)
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

	assert.EqualError(t, err, "no default permission found, but found more then one")
	assert.Equal(t, Code(""), perms)
}

func TestDefaultForTwoDefaultsReturnFirst(t *testing.T) {
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
	assert.Equal(t, Code("write"), perms)
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
