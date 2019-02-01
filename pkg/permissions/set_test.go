package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestPermissionSetIntersection(t *testing.T) {
	set1 := NewPermissionSet("set1", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})

	t.Run("given mutually excusive sets", func(t *testing.T) {
		set2 := NewPermissionSet("set2", []string{"ReadProducts", "ReadOffers", "ReadPartyPSU"})
		assert.False(t, set1.IsSubset(set2))
		assert.False(t, set2.IsSubset(set1))
		inter := set1.Intersection(set2)
		assert.Equal(t, len(inter.GetPermissions()), 0)
	})

	t.Run("given intersecting sets", func(t *testing.T) {
		set3 := NewPermissionSet("set3", []string{"ReadTransactionsBasic", "ReadProducts", "ReadOffers", "ReadPartyPSU"})
		inter := set1.Intersection(set3)
		assert.Equal(t, len(inter.GetPermissions()), 1)
		assert.Equal(t, inter.GetPermissions()[0], "ReadTransactionsBasic")
	})
}

func TestPermissionSetEqual(t *testing.T) {
	set1 := NewPermissionSet("set1", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})

	t.Run("given mutually excusive sets", func(t *testing.T) {
		set2 := NewPermissionSet("set2", []string{"ReadProducts", "ReadOffers", "ReadPartyPSU"})
		assert.False(t, set1.IsSubset(set2))
		assert.False(t, set2.IsSubset(set1))
		assert.False(t, set1.Equal(set2))
		assert.False(t, set2.Equal(set1))
	})

	t.Run("given set and subset", func(t *testing.T) {
		sub := NewPermissionSet("sub", []string{"ReadTransactionsBasic", "ReadTransactionsCredits"})
		assert.True(t, set1.IsSubset(sub))
		assert.False(t, sub.IsSubset(set1))
		assert.False(t, set1.Equal(sub))
		assert.False(t, sub.Equal(set1))
	})

	t.Run("given matching sets", func(t *testing.T) {
		other := NewPermissionSet("other", []string{"ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"})
		assert.True(t, set1.IsSubset(other))
		assert.True(t, other.IsSubset(set1))
		assert.True(t, set1.Equal(other))
		assert.True(t, other.Equal(set1))
	})
}
