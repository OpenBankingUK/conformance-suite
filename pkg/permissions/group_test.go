package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGroup(t *testing.T) {
	included := CodeSet{"a"}
	excluded := CodeSet{"b"}
	testId := "1"

	group := NewGroup(testId, included, excluded)

	assert.Equal(t, testId, string(group.TestId))
	assert.Equal(t, included, group.Included)
	assert.Equal(t, excluded, group.Excluded)
}

func TestGroupIsCompatible(t *testing.T) {
	groupUnderTest := &Group{
		Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
		Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
	}

	groups := &Group{
		Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
		Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
	}

	assert.False(t, groupUnderTest.isCompatible(groups))
}

func TestGroupIsCompatibleLess(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic", "ReadAccountsDetail"},
		Excluded: []Code{},
	}

	group := &Group{
		Included: CodeSet{"ReadAccountsBasic"},
		Excluded: CodeSet{},
	}

	assert.True(t, groupUnderTest.isCompatible(group))
}

func TestGroupIsCompatibleMoreCodes(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic"},
		Excluded: []Code{},
	}

	group := &Group{
		Included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
		Excluded: CodeSet{},
	}

	assert.True(t, groupUnderTest.isCompatible(group))
}

func TestGroupIsCompatibleByExcluded(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic"},
		Excluded: []Code{"ReadAccountsDetail"},
	}

	group := &Group{
		Included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
		Excluded: CodeSet{},
	}

	assert.False(t, groupUnderTest.isCompatible(group))
}

func TestGroupIsCompatibleByExcludedOnSource(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic"},
		Excluded: []Code{"ReadAccountsDetail"},
	}

	group := &Group{
		Included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
		Excluded: CodeSet{},
	}

	assert.False(t, groupUnderTest.isCompatible(group))
}

func TestGroupsIsSatisfiedByAnyOf(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic"},
		Excluded: []Code{"ReadAccountsDetail"},
	}

	groups := []*Group{
		{
			Included: CodeSet{"ReadAccountsBasic"},
			Excluded: CodeSet{},
		},
	}

	assert.True(t, groupUnderTest.isSatisfiedByAnyOf(groups))
}

func TestGroupIsSatisfiedByAnyOfFindsFirstOne(t *testing.T) {
	groupUnderTest := &Group{
		Included: []Code{"ReadAccountsBasic"},
		Excluded: []Code{"ReadAccountsDetail"},
	}

	groups := []*Group{
		{
			Included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
			Excluded: CodeSet{},
		},
		{
			Included: CodeSet{"ReadAccountsBasic"},
			Excluded: CodeSet{},
		},
	}

	assert.True(t, groupUnderTest.isSatisfiedByAnyOf(groups))
}

func TestGroup_Add(t *testing.T) {
	group1 := &Group{
		Included: []Code{"ReadAccountsBasic-1"},
		Excluded: []Code{"ReadAccountsDetail-1"},
	}

	group2 := &Group{
		Included: []Code{"ReadAccountsBasic-2"},
		Excluded: []Code{"ReadAccountsDetail-2"},
	}

	group1.add(group2)

	assert.Equal(t, CodeSet{"ReadAccountsBasic-1", "ReadAccountsBasic-2"}, group1.Included)
	assert.Equal(t, CodeSet{"ReadAccountsDetail-1", "ReadAccountsDetail-2"}, group1.Excluded)
}
