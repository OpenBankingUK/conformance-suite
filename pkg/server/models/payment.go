package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Needs to be a interface{} slice, see the official test for an example
// https://github.com/go-ozzo/ozzo-validation/blob/master/in_test.go
type OBExternalAccountIdentification4Code = interface{}

var (
	// OBExternalAccountIdentification4Codes - valid SchemeName as per the specification.
	OBExternalAccountIdentification4Codes = [5]OBExternalAccountIdentification4Code{
		"UK.OBIE.BBAN",
		"UK.OBIE.IBAN",
		"UK.OBIE.PAN",
		"UK.OBIE.Paym",
		"UK.OBIE.SortCodeAccountNumber",
	}
)

// Payment - Provides the details to identify the beneficiary account.
// This is referred to `OBCashAccount5` (line 9488) in the specification linked to below.
//
// Structure was deduced from this specification:
// https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json
//
// Example value:
// {
//     "SchemeName": "UK.OBIE.SortCodeAccountNumber",
//     "Identification": "20202010981789",
//     "Name": "Dr Foo"
// }
type Payment struct {
	// Name of the identification scheme, in a coded form as published in an external list
	SchemeName string `json:"scheme_name" form:"scheme_name"`
	// Beneficiary account identification.
	Identification string `json:"identification" form:"identification"`
	// Name of the account, as assigned by the account servicing institution.
	// Usage: The account name is the name or names of the account owner(s) represented at an account level. The account name is not the product name or the nickname of the account.
	Name string `json:"name" form:"name"`
}

// Just an an alternate spelling to match the Account and Transaction API Specification.
type OBCashAccount5 = Payment

// Validate - used by https://github.com/go-ozzo/ozzo-validation to validate struct.
func (p Payment) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.SchemeName, validation.Required, validation.Length(1, 40), validation.In(OBExternalAccountIdentification4Codes[:]...)),
		validation.Field(&p.Identification, validation.Required, validation.Length(1, 256)),
		validation.Field(&p.Name, validation.Length(1, 70)),
	)
}
