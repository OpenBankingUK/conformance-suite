package model

var permissions = []Permission{
	{
		Code: "ReadAccountsBasic",
		Endpoints: []string{
			"/accounts",
			"/accounts/{AccountId}",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadAccountsDetail",
		Endpoints: []string{
			"/accounts",
			"/accounts/{AccountId}",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadBalances",
		Endpoints: []string{
			"/balances",
			"/accounts/{AccountId}/balances",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadBeneficiariesBasic",
		Endpoints: []string{
			"/beneficiaries",
			"/accounts/{AccountId}/beneficiaries",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadBeneficiariesDetail",
		Endpoints: []string{
			"/beneficiaries",
			"/accounts/{AccountId}/beneficiaries",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadDirectDebits",
		Endpoints: []string{
			"/direct-debits",
			"/accounts/{AccountId}/direct-debits",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadStandingOrdersBasic",
		Endpoints: []string{
			"/standing-orders",
			"/accounts/{AccountId}/standing-orders",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadStandingOrdersDetail",
		Endpoints: []string{
			"/standing-orders",
			"/accounts/{AccountId}/standing-orders",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadTransactionsBasic",
		Endpoints: []string{
			"/transactions",
			"/accounts/{AccountId}/transactions",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
		},
		Default: true,
		RequiredOneOrMore: []Code{
			"ReadTransactionsCredits",
			"ReadTransactionsDebits",
		},
		Optional: []Code{"ReadPAN"},
	},
	{
		Code: "ReadTransactionsDetail",
		Endpoints: []string{
			"/transactions",
			"/accounts/{AccountId}/transactions",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
		},
		Default: false,
		RequiredOneOrMore: []Code{
			"ReadTransactionsCredits",
			"ReadTransactionsDebits",
		},
		Optional: []Code{"ReadPAN"},
	},
	{
		Code: "ReadTransactionsCredits",
		Endpoints: []string{
			"/transactions",
			"/accounts/{AccountId}/transactions",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
		},
		Default: false,
		RequiredOneOrMore: []Code{
			"ReadTransactionsBasic",
			"ReadTransactionsDetail",
		},
		Optional: []Code{"ReadPAN"},
	},
	{
		Code: "ReadTransactionsDebits",
		Endpoints: []string{
			"/transactions",
			"/accounts/{AccountId}/transactions",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
		},
		Default: false,
		RequiredOneOrMore: []Code{
			"ReadTransactionsBasic",
			"ReadTransactionsDetail",
		},
		Optional: []Code{"ReadPAN"},
	},
	{
		Code: "ReadStatementsBasic",
		Endpoints: []string{
			"/statements",
			"/accounts/{AccountId}/statements",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadStatementsDetail",
		Endpoints: []string{
			"/statements",
			"/accounts/{AccountId}/statements",
			"/accounts/{AccountId}/statements/{StatementId}/file",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadProducts",
		Endpoints: []string{
			"/products",
			"/accounts/{AccountId}/product",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadOffers",
		Endpoints: []string{
			"/offers",
			"/accounts/{AccountId}/offers",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadParty",
		Endpoints: []string{
			"/accounts/{AccountId}/party",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadPartyPSU",
		Endpoints: []string{
			"/party",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadScheduledPaymentsBasic",
		Endpoints: []string{
			"/scheduled-payments",
			"/accounts/{AccountId}/scheduled-payments",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadScheduledPaymentsDetail",
		Endpoints: []string{
			"/scheduled-payments",
			"/accounts/{AccountId}/scheduled-payments",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	}, {
		Code: "ReadPAN",
		Endpoints: []string{
			"",
			"",
		},
		Default:           false,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
}
