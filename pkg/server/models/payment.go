package models

import (
	"regexp"

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

// InstructedAmount represents global details for the payment test cases
// As in the Payment struct, structure was deduced from this specification:
// https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json
type InstructedAmount struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value,string"`
}

// Validate validates value and currency of the instructed amount
// provided in input
func (a InstructedAmount) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Value, validation.Max(float64(1))),
		validation.Field(&a.Currency, validation.Match(regexp.MustCompile("/^AED|AFN|ALL|AMD|ANG|AOA|ARS|AUD|AWG|AZN|BAM|BBD|BDT|BGN|BHD|BIF|BMD|BND|BOB|BRL|BSD|BTN|BWP|BYR|BZD|CAD|CDF|CHF|CLP|CNY|COP|CRC|CUC|CUP|CVE|CZK|DJF|DKK|DOP|DZD|EGP|ERN|ETB|EUR|FJD|FKP|GBP|GEL|GGP|GHS|GIP|GMD|GNF|GTQ|GYD|HKD|HNL|HRK|HTG|HUF|IDR|ILS|IMP|INR|IQD|IRR|ISK|JEP|JMD|JOD|JPY|KES|KGS|KHR|KMF|KPW|KRW|KWD|KYD|KZT|LAK|LBP|LKR|LRD|LSL|LYD|MAD|MDL|MGA|MKD|MMK|MNT|MOP|MRO|MUR|MVR|MWK|MXN|MYR|MZN|NAD|NGN|NIO|NOK|NPR|NZD|OMR|PAB|PEN|PGK|PHP|PKR|PLN|PYG|QAR|RON|RSD|RUB|RWF|SAR|SBD|SCR|SDG|SEK|SGD|SHP|SLL|SOS|SPL|SRD|STD|SVC|SYP|SZL|THB|TJS|TMT|TND|TOP|TRY|TTD|TVD|TWD|TZS|UAH|UGX|USD|UYU|UZS|VEF|VND|VUV|WST|XAF|XCD|XDR|XOF|XPF|YER|ZAR|ZMW|ZWD$/"))),
	)
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
