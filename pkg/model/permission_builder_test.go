package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPermissionGroup(t *testing.T) {
	tc := TestCase{ID: "1"}

	group := NewPermissionGroup(tc)

	assert.Equal(t, permissions.TestId("1"), group.TestId)
}

func TestIncludedPermissionsFromEmptyContext(t *testing.T) {
	ctx := Context{}
	b := newPermissionBuilder()

	codeSet := b.includedPermission(ctx, "/endpoint")

	assert.Equal(t, permissions.NoCodeSet(), codeSet)
}

func TestIncludedPermissionsFromContext(t *testing.T) {
	ctx := Context{}
	b := newPermissionBuilder()
	ctx.PutStringSlice(permissionIncludedKey, []string{"read", "write"})

	codeSet := b.includedPermission(ctx, "")

	assert.Equal(t, permissions.CodeSet{"read", "write"}, codeSet)
}

func TestExcludePermissionsFromEmptyContext(t *testing.T) {
	ctx := Context{}
	b := newPermissionBuilder()

	codeSet := b.excludedPermissions(ctx)

	assert.Equal(t, permissions.NoCodeSet(), codeSet)
}

func TestExcludePermissions(t *testing.T) {
	ctx := Context{}
	b := newPermissionBuilder()
	ctx.PutStringSlice(permissionExcludedKey, []string{"read", "write"})

	codeSet := b.excludedPermissions(ctx)

	assert.Equal(t, permissions.CodeSet{"read", "write"}, codeSet)
}

func TestExcludePermissionsHandlesCastPanic(t *testing.T) {
	b := newPermissionBuilder()
	var values []interface{}
	values = append(values, 0)
	ctx := Context{permissionExcludedKey: values}

	codeSet := b.excludedPermissions(ctx)

	assert.Equal(t, permissions.NoCodeSet(), codeSet)
}

func TestIncludePermissionsFromDefault(t *testing.T) {
	ctx := Context{}
	data := []permission{
		{
			Code:      "read",
			Default:   true,
			Endpoints: []string{"/accounts/{AccountId}/statements/{StatementId}/transactions"},
		},
	}
	stdPermissions := newStandardPermissionsWithOptions(data)
	b := newPermissionBuilderWithOptions(stdPermissions)

	codeSet := b.includedPermission(ctx, "/accounts/{AccountId}/statements/{StatementId}/transactions")

	assert.Equal(t, permissions.CodeSet{"read"}, codeSet)
}
