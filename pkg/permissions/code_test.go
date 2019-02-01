package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeSetHas(t *testing.T) {
	codeSet := CodeSet{"a"}

	assert.True(t, codeSet.Has("a"))
	assert.False(t, codeSet.Has("b"))
	assert.False(t, codeSet.Has(""))
}

func TestCodeSetEmptyHas(t *testing.T) {
	codeSet := CodeSet{}

	assert.False(t, codeSet.Has("a"))
	assert.False(t, codeSet.Has("b"))
	assert.False(t, codeSet.Has(""))
}

func TestCodeSetEmptyUnion(t *testing.T) {
	codeSet1 := CodeSet{}
	codeSet2 := CodeSet{}

	expected := CodeSet{}
	assert.Equal(t, expected, codeSet1.Union(codeSet2))
}

func TestCodeSetUnionOneEmpty(t *testing.T) {
	codeSet1 := CodeSet{"a"}
	codeSet2 := CodeSet{}

	expected := CodeSet{"a"}
	assert.Equal(t, expected, codeSet1.Union(codeSet2))
}

func TestCodeSetUnionAll(t *testing.T) {
	codeSet1 := CodeSet{"a"}
	codeSet2 := CodeSet{"b"}

	expected := CodeSet{"a", "b"}
	assert.Equal(t, expected, codeSet1.Union(codeSet2))
}

func TestCodeSetUnionAllNonDuplicates(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"b", "c"}

	expected := CodeSet{"a", "b", "c"}
	assert.Equal(t, expected, codeSet1.Union(codeSet2))
}

func TestCodeSetHasAll(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"a", "b"}

	assert.True(t, codeSet1.HasAll(codeSet2))
}

func TestCodeSetHasAllLargerSet(t *testing.T) {
	codeSet1 := CodeSet{"a", "b", "c"}
	codeSet2 := CodeSet{"a", "b"}

	assert.True(t, codeSet1.HasAll(codeSet2))
}

func TestCodeSetHasAllSmallerSet(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"a", "b", "c"}

	assert.False(t, codeSet1.HasAll(codeSet2))
}

func TestCodeSetHasAny(t *testing.T) {
	codeSet1 := CodeSet{"a"}
	codeSet2 := CodeSet{"a", "b", "c"}

	assert.True(t, codeSet1.HasAny(codeSet2))
}

func TestCodeSetHasAnyNot(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"c"}

	assert.False(t, codeSet1.HasAny(codeSet2))
}

func TestCodeEquals(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"b", "a"}

	assert.True(t, codeSet1.Equals(codeSet2))
}

func TestCodeEqualsFalse(t *testing.T) {
	codeSet1 := CodeSet{"a", "b"}
	codeSet2 := CodeSet{"a"}

	assert.False(t, codeSet1.Equals(codeSet2))
}

func TestCodeEqualsFalseInverted(t *testing.T) {
	codeSet1 := CodeSet{"a"}
	codeSet2 := CodeSet{"a", "b"}

	assert.False(t, codeSet1.Equals(codeSet2))
}
