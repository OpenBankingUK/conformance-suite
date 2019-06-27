package model

import "bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"

// NewPermissionGroup returns a list of Code objects associated with a testcase
func NewPermissionGroup(tc TestCase) permissions.Group {
	b := newPermissionBuilder()
	return newPermissionGroupWithOptions(tc, b)
}

// newPermissionGroupWithOptions returns a list of Code objects associated with a testcase
// allows you to pass a diff static data access
func newPermissionGroupWithOptions(tc TestCase, builder permissionBuilder) permissions.Group {
	return permissions.NewGroup(
		tc.ID,
		builder.includedPermission(tc.Context, tc.Input.Endpoint),
		builder.excludedPermissions(tc.Context),
	)
}

func NewDefaultPermissionGroup(tc TestCase) permissions.Group {
	builder := newPermissionBuilder()
	return permissions.NewGroup(
		tc.ID,
		builder.includedPermission(Context{}, tc.Input.Endpoint),
		builder.excludedPermissions(Context{}),
	)
}

const (
	permissionIncludedKey = "permissions"
	permissionExcludedKey = "permissions_excluded"
)

// permissionBuilder helper to calculate permission from a testcase
// need api static data to fetch default set of permissions
type permissionBuilder struct {
	standardPermissions standardPermissions
}

func newPermissionBuilder() permissionBuilder {
	return newPermissionBuilderWithOptions(newStandardPermissions())
}

func newPermissionBuilderWithOptions(standardPermissions standardPermissions) permissionBuilder {
	return permissionBuilder{
		standardPermissions: standardPermissions,
	}
}

// includedPermission returns the list of permission names that need to be included
// in the access token for this testcase. See permission model docs for more information
func (b permissionBuilder) includedPermission(ctx Context, endpoint string) permissions.CodeSet {
	values, err := ctx.GetStringSlice(permissionIncludedKey)
	if err == nil {
		return mapStringToCodeSet(values)
	}

	defaultPerms, err := b.standardPermissions.defaultForEndpoint(endpoint)
	if err != nil {
		return permissions.NoCodeSet()
	}
	codes := []string{}
	for _, code := range defaultPerms {
		codes = append(codes, string(code))
	}
	return mapStringToCodeSet(codes)
}

// excludedPermissions return a list of excluded permissions from context
func (b permissionBuilder) excludedPermissions(ctx Context) permissions.CodeSet {
	values, err := ctx.GetStringSlice(permissionExcludedKey)
	if err != nil {
		return permissions.NoCodeSet()
	}
	return mapStringToCodeSet(values)
}

func mapStringToCodeSet(values []string) permissions.CodeSet {
	var result permissions.CodeSet
	for _, value := range values {
		result = append(result, permissions.Code(value))
	}
	return result
}
