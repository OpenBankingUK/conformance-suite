package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"

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

// feature/refapp_466_add_permissions_to_testcase
var (
	transactionTestcase01 = []byte(`
	{
        "@id": "#t1008",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/accounts"
        },
        "context": {
			"permissions":["ReadTransactionsBasic","ReadTransactionDetail","ReadTransactionsCredits"],
			"permissions_excluded":["ReadTransactionsDebits"]
		},
        "expect": {
            "status-code": 200,
            "schema-validation": true
        }
    }
	`)
)

// Read a testcase with permissions and check that they are all retrieved
func TestIncludedAndExcludedPermissions(t *testing.T) {
	tc := TestCase{}
	err := json.Unmarshal(transactionTestcase01, &tc)
	assert.NoError(t, err)
	included, excluded := tc.GetPermissions()
	assert.Equal(t, len(included), 3)
	assert.Equal(t, len(excluded), 1)
}

// A conformance suite testcase needs to know how to emit the permissions that it contains
// A rule needs to know how to emit the permissionSet that a testcase sequence contains
// A collection of permissions not required should also be expressed
// - so two sets - in inclusion set and an exclusion set

var (
	transactionTestcase02 = []byte(`
	{
        "@id": "#t1008",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/transactions"
        },
        "context": {
		},
        "expect": {
            "status-code": 200,
            "schema-validation": true
        }
    }
	`)
)

// figure out the default permission for /transactions
func TestGetDefaultPermissionsForEndpoint(t *testing.T) {
	tc := TestCase{}
	err := json.Unmarshal(transactionTestcase02, &tc)
	assert.NoError(t, err)
	included, excluded := tc.GetPermissions()
	assert.Equal(t, len(included), 1)
	assert.Equal(t, len(excluded), 0)
	assert.Equal(t, included[0], "ReadTransactionsBasic")
}

//
func TestGetIncludedAndExcludedPermissionSetsFromTestcaseSequence(t *testing.T) {
	m, err := loadPermissionTestData() // from testdata/permissionTestData.json
	assert.Nil(t, err)
	rule := m.Rules[0]
	includedSet, excludedSet := rule.GetPermissionSets()
	assert.Equal(t, len(includedSet), 3)
	assert.Equal(t, len(excludedSet), 2)
}

func loadPermissionTestData() (Manifest, error) {
	plan, _ := ioutil.ReadFile("testdata/permissionTestData.json")
	var m Manifest
	err := json.Unmarshal(plan, &m)
	if err != nil {
		return Manifest{}, err
	}
	fmt.Printf(string(pkgutils.DumpJSON(m)))
	return m, nil
}
