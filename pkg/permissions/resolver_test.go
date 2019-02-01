package permissions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermissionResolver(t *testing.T) {
	t.Run("given empty group expect empty slice result", func(t *testing.T) {
		var configs []group
		result := Resolver(configs)
		var expected []CodeSet
		assert.Equal(t, result, expected)
	})

	t.Run("given one group expect single permission set result", func(t *testing.T) {
		configs := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
		}
		result := Resolver(configs)
		expected := []CodeSet{
			{"ReadAccountsBasic"},
		}
		assert.Equal(t, result, expected)
	})

	t.Run("with duplicated group expect single permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic"},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("with two mutually exclusive configs expect two permission sets result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				included: CodeSet{"ReadAccountsDetail"},
				excluded: CodeSet{},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic"},
			{"ReadAccountsDetail"},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("with two non-mutually exclusive configs expect single permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("with two non-mutually exclusive configs and one mutually exclusive expect two permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				excluded: CodeSet{},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
			{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
		}

		assert.Equal(t, result, expected)
	})

	t.Run("with four mutually exclusive configs and one non-mutually exclusive expect four permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
			{"ReadAccountsBasic", "ReadTransactionsBasic"},
			{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("with four mutually exclusive configs and four non-mutually exclusive expect four permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsDetail"},
			},

			{
				included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits"},
			},
			{
				included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic"},
				excluded: CodeSet{"ReadTransactionsDetail", "ReadTransactionsDebits", "ReadTransactionsCredits"},
			},
			{
				included: CodeSet{"ReadStatementsDetail", "ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsDetail"},
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
		assert.True(t, result[0].HasAll(expected[0]))
		assert.True(t, result[1].HasAll(expected[1]))
		assert.True(t, result[2].HasAll(expected[2]))
		assert.True(t, result[3].HasAll(expected[3]))
	})

	t.Run("with large set of Accounts API configs expect four permission set result", func(t *testing.T) {
		set := []group{
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts
			},
			{
				included: CodeSet{"ReadAccountsDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadPAN"},
				excluded: CodeSet{},
				// endpoint: /accounts
			},
			{
				included: CodeSet{"ReadAccountsBasic"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}
			},
			{
				included: CodeSet{"ReadAccountsDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadPAN"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBalances"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /balances
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBalances"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/balances
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				// endpoint: /beneficiaries
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				excluded: CodeSet{},
				// endpoint: /beneficiaries
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadBeneficiariesBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				// endpoint: /accounts/{AccountId}/beneficiaries
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadBeneficiariesDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/beneficiaries
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadDirectDebits"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /direct-debits
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadDirectDebits"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/direct-debits
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStandingOrdersBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				// endpoint: /standing-orders
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				excluded: CodeSet{},
				// endpoint: /standing-orders
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStandingOrdersBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				// endpoint: /accounts/{AccountId}/standing-orders
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStandingOrdersDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/standing-orders
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{},
				// endpoint: /transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				// endpoint: /statements
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				excluded: CodeSet{},
				// endpoint: /statements
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				// endpoint: /accounts/{AccountId}/statements
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/statements
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/file
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsDebits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadTransactionsCredits"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadStatementsDetail", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadTransactionsDetail"},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadStatementsDetail", "ReadTransactionsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadPAN"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/statements/{StatementId}/transactions
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadProducts"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /products
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadProducts"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/product
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadOffers"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /offers
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadOffers"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/offers
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadParty"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /accounts/{AccountId}/party
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadPartyPSU"},
				excluded: CodeSet{"ReadAccountsDetail"},
				// endpoint: /party
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadScheduledPaymentsBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				// endpoint: /scheduled-payments
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				excluded: CodeSet{},
				// endpoint: /scheduled-payments
			},
			{
				included: CodeSet{"ReadAccountsBasic", "ReadScheduledPaymentsBasic"},
				excluded: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				// endpoint: /accounts/{AccountId}/scheduled-payments
			},
			{
				included: CodeSet{"ReadAccountsDetail", "ReadScheduledPaymentsDetail"},
				excluded: CodeSet{},
				// endpoint: /accounts/{AccountId}/scheduled-payments
			},
		}

		result := Resolver(set)
		expected := []CodeSet{
			{"ReadAccountsBasic", "ReadBalances", "ReadBeneficiariesBasic", "ReadDirectDebits", "ReadOffers", "ReadPAN", "ReadParty", "ReadPartyPSU", "ReadProducts", "ReadScheduledPaymentsBasic", "ReadStandingOrdersBasic", "ReadStatementsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits"},
			{"ReadAccountsDetail", "ReadBeneficiariesDetail", "ReadPAN", "ReadScheduledPaymentsDetail", "ReadStandingOrdersDetail", "ReadStatementsDetail", "ReadTransactionsCredits", "ReadTransactionsDetail"},
		}
		assert.Len(t, result, 9)
		assert.True(t, result[0].HasAll(expected[0]))
		assert.True(t, result[1].HasAll(expected[1]))
	})

}
