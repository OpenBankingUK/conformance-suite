package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermissionResolver(t *testing.T) {
	t.Run("given empty Group expect empty slice result", func(t *testing.T) {
		var configs []Group

		result := Resolver(configs)

		var expected CodeSetResultSet
		assert.Equal(t, result, expected)
	})

	t.Run("given one Group expect single permission set result", func(t *testing.T) {
		configs := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
		}

		result := Resolver(configs)

		expected := CodeSet{"ReadAccountsBasic"}
		assert.Len(t, result, 1)
		assert.True(t, result[0].CodeSet.HasAll(expected))
	})

	t.Run("with duplicated Group expect single permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
		}

		result := Resolver(set)

		expected := CodeSet{"ReadAccountsBasic"}
		assert.Len(t, result, 1)
		assert.Equal(t, expected, result[0].CodeSet)
	})

	t.Run("with two mutually exclusive configs expect two permission sets result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsDetail"},
				Excluded: CodeSet{},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic"},
			{"ReadAccountsDetail"},
		}

		assert.Len(t, result, 2)
		assert.Equal(t, expected[0], result[0].CodeSet)
		assert.Equal(t, expected[1], result[1].CodeSet)
	})

	t.Run("with two non-mutually exclusive configs expect single permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
			},
		}

		result := Resolver(set)
		expected := CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"}

		assert.Len(t, result, 1)
		assert.Equal(t, expected, result[0].CodeSet)
	})

	t.Run("with two non-mutually exclusive configs and one mutually exclusive expect two permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				Excluded: CodeSet{},
			},
		}

		result := Resolver(set)

		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
			{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
		}
		assert.Len(t, result, 2)
		assert.Equal(t, expected[0], result[0].CodeSet)
		assert.Equal(t, expected[1], result[1].CodeSet)
	})

	t.Run("with four mutually exclusive configs and one non-mutually exclusive expect four permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
			},
		}

		result := Resolver(set)

		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
			{"ReadAccountsBasic", "ReadTransactionsBasic"},
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
		}
		assert.Len(t, result, 4)
		assert.Equal(t, expected[0], result[0].CodeSet)
		assert.Equal(t, expected[1], result[1].CodeSet)
		assert.Equal(t, expected[2], result[2].CodeSet)
		assert.Equal(t, expected[3], result[3].CodeSet)
	})

	t.Run("with four mutually exclusive configs and four non-mutually exclusive expect four permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsDetail"},
			},

			{
				Included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				Included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic"},
				Excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				Included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsDetail"},
			},
		}

		result := Resolver(set)

		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits"},
			{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsDebits"},
			{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic"},
			{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
		}
		assert.Equal(t, len(result), 4)
		assert.Len(t, expected, 4)
		assert.True(t, result[0].CodeSet.HasAll(expected[0]))
		assert.True(t, result[1].CodeSet.HasAll(expected[1]))
		assert.True(t, result[2].CodeSet.HasAll(expected[2]))
		assert.True(t, result[3].CodeSet.HasAll(expected[3]))
	})

	t.Run("with large set of Accounts API configs expect four permission set result", func(t *testing.T) {
		set := []Group{
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts
			},
			{
				Included: CodeSet{"ReadAccountsDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadPAN"},
				Excluded: CodeSet{},
				// endpoint: /accounts
			},
			{
				Included: CodeSet{"ReadAccountsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}
			},
			{
				Included: CodeSet{"ReadAccountsDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadPAN"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBalances"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /balances
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBalances"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/balances
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				// endpoint: /beneficiaries
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				Excluded: CodeSet{},
				// endpoint: /beneficiaries
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				// endpoint: /accounts/{AccountId}/beneficiaries
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/beneficiaries
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadDirectDebits"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /direct-debits
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadDirectDebits"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/direct-debits
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStandingOrdersBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				// endpoint: /standing-orders
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				Excluded: CodeSet{},
				// endpoint: /standing-orders
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStandingOrdersBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				// endpoint: /accounts/{AccountId}/standing-orders
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/standing-orders
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{},
				// endpoint: /transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				// endpoint: /statements
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				Excluded: CodeSet{},
				// endpoint: /statements
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				// endpoint: /accounts/{AccountId}/statements
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/statements
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/file
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadProducts"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /products
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadProducts"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/product
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadOffers"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /offers
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadOffers"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/offers
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadParty"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/party
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadPartyPSU"},
				Excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /party
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadScheduledPaymentsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				// endpoint: /scheduled-payments
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				Excluded: CodeSet{},
				// endpoint: /scheduled-payments
			},
			{
				Included: CodeSet{"ReadAccountsBasic", "ReadScheduledPaymentsBasic"},
				Excluded: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				// endpoint: /accounts/{AccountId}/scheduled-payments
			},
			{
				Included: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				Excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/scheduled-payments
			},
		}

		result := Resolver(set)

		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadBalances", "ReadBeneficiariesBasic", "ReadDirectDebits", "ReadOffers", "ReadPAN", "ReadParty", "ReadPartyPSU", "ReadProducts", "ReadScheduledPaymentsBasic", "ReadStandingOrdersBasic", "ReadStatementsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
			{"ReadAccountsDetail", "ReadBeneficiariesDetail", "ReadPAN", "ReadScheduledPaymentsDetail", "ReadStandingOrdersDetail", "ReadStatementsDetail", "ReadTransactionsCredits", "ReadTransactionsDetail"},
		}
		assert.True(t, result[0].CodeSet.HasAll(expected[0]))
		assert.True(t, result[1].CodeSet.HasAll(expected[1]))
	})
}

func TestAddTestToGroupFound(t *testing.T) {
	var codeSets CodeSetResultSet
	groupInput := Group{TestId: "123"}
	groupFound := &Group{Included: CodeSet{"a"}}

	codeSets.addTestToGroupFound(groupInput, groupFound)

	assert.Len(t, codeSets, 1)
	assert.Len(t, codeSets[0].TestIds, 1)
	assert.Equal(t, TestId("123"), codeSets[0].TestIds[0])
	assert.Len(t, codeSets[0].CodeSet, 1)
	assert.Equal(t, Code("a"), codeSets[0].CodeSet[0])
}

func TestAddTestToGroupFoundAppendsTestId(t *testing.T) {
	var codeSets CodeSetResultSet
	groupInput := Group{TestId: "1"}
	groupFound := &Group{Included: CodeSet{"a"}}
	codeSets.addTestToGroupFound(groupInput, groupFound)

	groupTestIdAppend := Group{TestId: "2"}
	codeSets.addTestToGroupFound(groupTestIdAppend, groupFound)

	assert.Len(t, codeSets, 1)
	assert.Len(t, codeSets[0].TestIds, 2)
	assert.Equal(t, TestId("1"), codeSets[0].TestIds[0])
	assert.Equal(t, TestId("2"), codeSets[0].TestIds[1])
}

func TestMapCodeSets(t *testing.T) {
	groups := []Group{
		{TestId: "1", Included: CodeSet{"a"}},
		{TestId: "2", Included: CodeSet{"a"}},
		{TestId: "3", Included: CodeSet{"b"}},
		{TestId: "4", Included: CodeSet{"a"}},
	}
	groupsFound := []*Group{
		{Included: CodeSet{"a"}},
		{Included: CodeSet{"b"}},
	}

	results := mapToCodeSets(groups, groupsFound)

	expected := []CodeSetResult{
		{
			CodeSet: CodeSet{"a"},
			TestIds: []TestId{"1", "2", "4"},
		},
		{
			CodeSet: CodeSet{"b"},
			TestIds: []TestId{"3"},
		},
	}
	assert.Equal(t, expected, results)
}

func TestMapCodeSetsWrongMatch(t *testing.T) {
	groups := []Group{
		{TestId: "1", Included: CodeSet{"a"}},
		{TestId: "2", Included: CodeSet{"a", "b"}},
		{TestId: "3", Included: CodeSet{"b"}, Excluded: CodeSet{"a"}},
		{TestId: "4", Included: CodeSet{"a"}, Excluded: CodeSet{"b"}},
	}
	groupsFound := []*Group{
		{Included: CodeSet{"a"}, Excluded: CodeSet{"b"}},
		{Included: CodeSet{"b"}, Excluded: CodeSet{"a"}},
		{Included: CodeSet{"a", "b"}},
	}

	results := mapToCodeSets(groups, groupsFound)

	expected := []CodeSetResult{
		{
			CodeSet: CodeSet{"a"},
			TestIds: []TestId{"1", "4"},
		},
		{
			CodeSet: CodeSet{"b"},
			TestIds: []TestId{"3"},
		},
		{
			CodeSet: CodeSet{"a", "b"},
			TestIds: []TestId{"1", "2"},
		},
	}
	assert.Equal(t, expected, results)
}
