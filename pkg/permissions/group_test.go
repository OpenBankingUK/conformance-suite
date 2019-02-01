package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupIsCompatible(t *testing.T) {
	groupUnderTest := &group{
		included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
		excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
	}

	groups := &group{
		included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
		excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
	}

	assert.False(t, groupUnderTest.isCompatible(groups))
}

func TestGroupIsCompatibleLess(t *testing.T) {
	groupUnderTest := &group{
		included: []Code{"ReadAccountsBasic", "ReadAccountsDetail"},
		excluded: []Code{},
	}

	group := &group{
		included: CodeSet{"ReadAccountsBasic"},
		excluded: CodeSet{},
	}

	assert.True(t, groupUnderTest.isCompatible(group))
}

func TestGroupIsCompatibleMoreCodes(t *testing.T) {
	groupUnderTest := &group{
		included: []Code{"ReadAccountsBasic"},
		excluded: []Code{},
	}

	group := &group{
		included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
		excluded: CodeSet{},
	}

	assert.True(t, groupUnderTest.isCompatible(group))
}

func TestGroupIsCompatibleByExcluded(t *testing.T) {
	groupUnderTest := &group{
		included: []Code{"ReadAccountsBasic"},
		excluded: []Code{"ReadAccountsDetail"},
	}

	group := &group{
		included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
		excluded: CodeSet{},
	}

	assert.False(t, groupUnderTest.isCompatible(group))
}

func TestGroupsIsSatisfiedByAnyOf(t *testing.T) {
	groupUnderTest := &group{
		included: []Code{"ReadAccountsBasic"},
		excluded: []Code{"ReadAccountsDetail"},
	}

	groups := []*group{
		{
			included: CodeSet{"ReadAccountsBasic"},
			excluded: CodeSet{},
		},
	}

	assert.True(t, groupUnderTest.isSatisfiedByAnyOf(groups))
}

func TestGroupIsSatisfiedByAnyOfFindsFirstOne(t *testing.T) {
	groupUnderTest := &group{
		included: []Code{"ReadAccountsBasic"},
		excluded: []Code{"ReadAccountsDetail"},
	}

	groups := []*group{
		{
			included: CodeSet{"ReadAccountsBasic", "ReadAccountsDetail"},
			excluded: CodeSet{},
		},
		{
			included: CodeSet{"ReadAccountsBasic"},
			excluded: CodeSet{},
		},
	}

	assert.True(t, groupUnderTest.isSatisfiedByAnyOf(groups))
}

func TestGroup_Union(t *testing.T) {
	group1 := &group{
		included: []Code{"ReadAccountsBasic-1"},
		excluded: []Code{"ReadAccountsDetail-1"},
	}

	group2 := &group{
		included: []Code{"ReadAccountsBasic-2"},
		excluded: []Code{"ReadAccountsDetail-2"},
	}

	group1.Union(group2)

	assert.Equal(t, CodeSet{"ReadAccountsBasic-1", "ReadAccountsBasic-2"}, group1.included)
	assert.Equal(t, CodeSet{"ReadAccountsDetail-1", "ReadAccountsDetail-2"}, group1.excluded)
}
