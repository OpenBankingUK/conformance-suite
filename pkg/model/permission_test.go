package model

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

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
	perm = GetPermissionFromName("SugarCoatedApple")
	assert.Equal(t, len(perm.Endpoints), 0)
}

// feature/refapp_466_add_permissions_to_testcase
var (
	permissionTestcase01 = []byte(`
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
            "status-code": 200
        }
    }
	`)
)

// Read a testcase with permissions and check that they are all retrieved
func TestIncludedAndExcludedPermissions(t *testing.T) {
	tc := TestCase{}
	err := json.Unmarshal(permissionTestcase01, &tc)
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
	permissionTestcase02 = []byte(`
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
            "status-code": 200
        }
    }
	`)
)

// figure out the default permission for /transactions
func TestGetDefaultPermissionsForEndpoint(t *testing.T) {
	tc := TestCase{}
	err := json.Unmarshal(permissionTestcase02, &tc)
	assert.NoError(t, err)
	included, excluded := tc.GetPermissions()
	assert.Equal(t, len(included), 1)
	assert.Equal(t, len(excluded), 0)
	assert.Equal(t, included[0], "ReadTransactionsBasic")
}

var (
	excludedPermissionsOnlyTestcase = []byte(`
	{
        "@id": "#t1008",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/accounts"
        },
        "context": {
			"permissions_excluded":["ReadTransactionsDebits","DummyPermission"]
		},
        "expect": {
            "status-code": 200
        }
    }
	`)
)

// checks that when we just define
func TestExcludedPermissionsOnlyForTestcase(t *testing.T) {
	tc := TestCase{}
	err := json.Unmarshal(excludedPermissionsOnlyTestcase, &tc)
	assert.NoError(t, err)
	included, excluded := tc.GetPermissions()
	assert.Equal(t, len(included), 0)
	assert.Equal(t, len(excluded), 2)
	assert.Equal(t, excluded[1], "DummyPermission")
}

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
			"permissions":["ReadTransactionsBasic","ReadTransactionsCredits","ReadTransactionsDebits"]
		},
        "expect": {
            "status-code": 200
        }
    }
	`)
)

var (
	transactionTestcase02 = []byte(`
	{
        "@id": "#t1010",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/transactions"
        },
        "context": {
			"permissions_excluded":["ReadTransactionsBasic","ReadTransactionDetail"]
		},
        "expect": {
            "status-code": 403
        }
    }
	`)
)

// Checks that a rule can retrieve the permissionSets for both included and excluded permissions
// from a series of test cases defined in a manifest
func TestGetIncludedAndExcludedPermissionSetsFromTestcaseSequence(t *testing.T) {
	m, err := loadPermissionTestData() // from testdata/permissionTestData.json
	assert.Nil(t, err)
	rule := m.Rules[0]
	includedSet, excludedSet := rule.GetPermissionSets()
	assert.Equal(t, 3, len(includedSet))
	assert.Equal(t, 2, len(excludedSet))
}

// two testcases to show /transcates with relevant permission
// and without relevant permissions

func TestTransactionsWithCorrectPermissions(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/transactions").Reply(200).BodyString(string(getAccountResponse))
	tc := TestCase{}
	err := json.Unmarshal(transactionTestcase01, &tc)
	assert.NoError(t, err)

	included, excluded := tc.GetPermissions()
	assert.Equal(t, 3, len(included))
	assert.Equal(t, 0, len(excluded))
	assert.Equal(t, included[0], "ReadTransactionsBasic")

}

func TestTransctionWithoutCorrectPermissions(t *testing.T) {
	defer gock.Off()
	gock.New("http://myaspsp").Get("/transactions").Reply(403).BodyString(string(getAccountResponse))

	tc := TestCase{}
	err := json.Unmarshal(transactionTestcase02, &tc)
	assert.NoError(t, err)
	included, excluded := tc.GetPermissions()
	assert.Equal(t, 0, len(included))
	assert.Equal(t, 2, len(excluded))
	assert.Equal(t, excluded[0], "ReadTransactionsBasic")
}

func loadPermissionTestData() (Manifest, error) {
	plan, _ := ioutil.ReadFile("testdata/permissionTestData.json")
	var m Manifest
	err := json.Unmarshal(plan, &m)
	if err != nil {
		return Manifest{}, err
	}
	//fmt.Printf(string(pkgutils.DumpJSON(m)))
	return m, nil
}

// Permissionset Test cases

func TestGetSetPermissionSetNames(t *testing.T) {
	p := NewPermissionSet("test", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
	assert.Equal(t, "test", p.GetName())
	p.SetName("anothertest")
	assert.Equal(t, "anothertest", p.GetName())
}

func TestGetPermissionFromSet(t *testing.T) {
	p := NewPermissionSet("test", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
	permission := p.Get("ReadTransactionsDebits")
	assert.True(t, permission)
	permission = p.Get("nonexistent")
	assert.False(t, permission)
}

func TestRemovePermissionFromSet(t *testing.T) {
	p := NewPermissionSet("test", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
	permission := p.Get("ReadTransactionsDebits")
	assert.True(t, permission)
	p.Remove("ReadTransactionsDebits")
	assert.False(t, p.Get("ReadTransactionsDebit"))
}

func TestPermissionSetSubSet(t *testing.T) {
	superset := NewPermissionSet("super", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
	subset := NewPermissionSet("sub", []string{"ReadTransactionsDebits"})
	issubset := superset.IsSubset(subset)
	assert.True(t, issubset)
	subset2 := NewPermissionSet("notsub", []string{"ReadTransactionsDebits_1"})
	issubset = superset.IsSubset(subset2)
	assert.False(t, issubset)
}

func TestPermissionSetUnion(t *testing.T) {
	set1 := NewPermissionSet("set1", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
	set2 := NewPermissionSet("set2", []string{"ReadProducts", "ReadOffers", "ReadPartyPSU"})
	assert.False(t, set1.IsSubset(set2))
	assert.False(t, set2.IsSubset(set1))
	union := set1.Union(set2)
	assert.True(t, union.IsSubset(set1))
	assert.True(t, union.IsSubset(set2))
}
