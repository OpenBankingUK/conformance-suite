package model

import (
	"github.com/pkg/errors"
)

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

// defaultForEndpoint finds the default permissions for an endpoint
// i.e. slice of permission codes marked as default true
func (sp standardPermissions) defaultForEndpoint(endpoint string) ([]Code, error) {
	perms := sp.permissionsForEndpoint(endpoint)

	if len(perms) == 0 {
		return []Code{}, errors.New("no default permissions found")
	}

	// only one permission so always DEFAULT
	if len(perms) == 1 {
		code := []Code{perms[0].Code}
		return code, nil
	}

	codes := []Code{}
	for _, p := range perms {
		if p.Default {
			codes = append(codes, p.Code)
		}
	}
	if len(codes) > 0 {
		return codes, nil
	}

	return []Code{}, errors.New("no default permissions found, but found more than one")
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
			"/accounts/{AccountId}/balances",
			"/accounts/{AccountId}/beneficiaries",
			"/accounts/{AccountId}/direct-debits",
			"/accounts/{AccountId}/offers",
			"/accounts/{AccountId}/party",
			"/accounts/{AccountId}/product",
			"/accounts/{AccountId}/scheduled-payments",
			"/accounts/{AccountId}/standing-orders",
			"/accounts/{AccountId}/statements",
			"/accounts/{AccountId}/statements/{StatementId}/file",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
			"/accounts/{AccountId}/transactions",
			"/balances",
			"/beneficiaries",
			"/direct-debits",
			"/offers",
			"/party",
			"/products",
			"/scheduled-payments",
			"/standing-orders",
			"/statements",
			"/transactions",
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
		Default:           true,
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
		Default:           true,
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
		Default: true,
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
		Default: true,
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
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadStatementsDetail",
		Endpoints: []string{
			"/statements",
			"/accounts/{AccountId}/statements",
			"/accounts/{AccountId}/statements/{StatementId}/file",
			"/accounts/{AccountId}/statements/{StatementId}/transactions",
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
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadOffers",
		Endpoints: []string{
			"/offers",
			"/accounts/{AccountId}/offers",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadParty",
		Endpoints: []string{
			"/accounts/{AccountId}/party",
		},
		Default:           true,
		RequiredOneOrMore: []Code{},
		Optional:          []Code{},
	},
	{
		Code: "ReadPartyPSU",
		Endpoints: []string{
			"/party",
		},
		Default:           true,
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
