package model

import "github.com/pkg/errors"

// Code is a string representing a OB access permission
type Code string

// permission holds endpoint permission data
type permission struct {
	Code              Code
	Endpoints         []string
	Default           bool
	RequiredOneOrMore []Code
	Optional          []Code
}

type standardPermissions struct {
	permissions []permission
}

func newStandardPermissions() standardPermissions {
	return standardPermissions{
		permissions: staticApiPermissions,
	}
}

func newStandardPermissionsWithOptions(data []permission) standardPermissions {
	return standardPermissions{
		permissions: data,
	}
}

// defaultForEndpoint finds the default permission for and endpoint
// either marked as default or it only has one permission
func (sp standardPermissions) defaultForEndpoint(endpoint string) (Code, error) {
	perms := sp.permissionsForEndpoint(endpoint)

	if len(perms) == 0 {
		return Code(""), errors.New("no default permissions found")
	}

	// only one permission so always DEFAULT
	if len(perms) == 1 {
		return perms[0].Code, nil
	}

	for _, p := range perms {
		if p.Default == true {
			return p.Code, nil
		}
	}

	return Code(""), errors.New("no default permission found, but found more then one")
}

// permissionsForEndpoint returns a list of Permissions required by an endpoint
func (sp standardPermissions) permissionsForEndpoint(endpoint string) []permission {
	var endpointPermissions []permission
	for _, p := range sp.permissions {
		for _, e := range p.Endpoints {
			if e == endpoint {
				endpointPermissions = append(endpointPermissions, p)
			}
		}
	}
	return endpointPermissions
}

// staticApiPermission is the standard for OB permission
// accesses to account endpoints
var staticApiPermissions = []permission{
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
