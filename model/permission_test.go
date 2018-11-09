package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// check we have some permissions loaded from the configuration
func TestPermissionJsonRead(t *testing.T) {
	count := len(permissions)
	result := count > 10
	assert.True(t, true, result)

}

// Get the list of permissions associated with this an endpoint
func TestPermissionListReturned(t *testing.T) {
	list := GetPermissionsForEndpoint("/accounts/{AccountId}/transactions")
	count := len(list)
	assert.Equal(t, 4, count) // get 4 permissions return that refer to /accounts/{AccountId}/transactions
}

// For a specified permission name, get the permission object to which it refers
func TestSpecifiedPermissionName(t *testing.T) {
	perm := GetPermissionFromName("ReadTransactionsDetail")
	assert.Equal(t, "ReadTransactionsDetail", perm.Permission)
}
