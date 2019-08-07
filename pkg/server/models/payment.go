package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
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
		validation.Field(&p.SchemeName, validation.Required, validation.Length(1, 40)),
		validation.Field(&p.Identification, validation.Required, validation.Length(1, 256)),
		validation.Field(&p.Name, validation.Length(1, 70)),
	)
}

// InstructedAmount - Represents global details for the payment test cases
// As in the Payment struct, structure was deduced from this specification:
// https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json
//
// `Value` is of the format specified below:
//
//     "OBActiveCurrencyAndAmount_SimpleType": {
//         "description": "A number of monetary units specified in an active currency where the unit of currency is explicit and compliant with ISO 4217.",
//         "type": "string",
//         "pattern": "^\\d{1,13}\\.\\d{1,5}$"
//     },
//
// See: https://github.com/OpenBankingUK/read-write-api-specs/blob/master/dist/account-info-swagger.json#L2964.
type InstructedAmount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

const (
	regexInstructedAmountCurrencyErr = `must be in a valid format (^[A-Z]{3,3}$)`
	regexInstructedAmountValueErr    = `must be in a valid format (^\d{1,13}\.\d{1,5}$)`
)

var (
	// nolint:gochecknoglobals
	regexInstructedAmountCurrency = regexp.MustCompile("^[A-Z]{3,3}$")
	// nolint:gochecknoglobals
	regexInstructedAmountValue = regexp.MustCompile(`^\d{1,13}\.\d{1,5}$`)
)

// Validate - validates value and currency of the instructed amount provided in input
func (a InstructedAmount) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Currency, validation.Match(regexInstructedAmountCurrency).Error(regexInstructedAmountCurrencyErr)),
		validation.Field(&a.Value, validation.Match(regexInstructedAmountValue).Error(regexInstructedAmountValueErr)),
	)
}

const (
	regexPaymentFrequencyErr = `must be in a valid format (^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$)`
)

var (
	// nolint:gochecknoglobals
	regexPaymentFrequency = regexp.MustCompile(`^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$`)
)

type PaymentFrequency string

// Validate - ensures
func (p PaymentFrequency) Validate() error {
	return validation.Match(regexPaymentFrequency).Error(regexPaymentFrequencyErr).Validate(p)
}
