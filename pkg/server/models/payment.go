package models

import (
	"encoding/json"
	"fmt"
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
//
//	{
//	    "SchemeName": "UK.OBIE.SortCodeAccountNumber",
//	    "Identification": "20202010981789",
//	    "Name": "Dr Foo"
//	}
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
//	"OBActiveCurrencyAndAmount_SimpleType": {
//	    "description": "A number of monetary units specified in an active currency where the unit of currency is explicit and compliant with ISO 4217.",
//	    "type": "string",
//	    "pattern": "^\\d{1,13}\\.\\d{1,5}$"
//	},
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
	regexPaymentFrequencyPattern = `^(NotKnown)$|^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlDay:((0[2-9])|([1-2][0-9])|3[0-1]))$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$`
	regexPaymentFrequencyErr     = `must be in a valid Frequency_1 format (` + regexPaymentFrequencyPattern + `)`
)

var (
	// nolint:gochecknoglobals
	regexPaymentFrequency = regexp.MustCompile(regexPaymentFrequencyPattern)
)

type PaymentFrequency string

// Validate - ensures
func (p PaymentFrequency) Validate() error {
	return validation.Match(regexPaymentFrequency).Error(regexPaymentFrequencyErr).Validate(p)
}

// V4StandingOrderFrequency models OBFrequency6 for v4 standing-order payloads.
type V4StandingOrderFrequency struct {
	Type           string `json:"Type"`
	PointInTime    string `json:"PointInTime,omitempty"`
	CountPerPeriod *int   `json:"CountPerPeriod,omitempty"`
}

const (
	defaultV4StandingOrderFrequencyType        = "WEEK"
	defaultV4StandingOrderFrequencyPointInTime = "03"
	regexV4StandingOrderFrequencyTypeErr       = `must be an OBFrequency6Code or Frequency_1 value`
	regexV4StandingOrderPointInTimeErr         = `must be numeric text up to two characters, including negative single-digit values`
	legacyV4StandingOrderFrequencyFieldsErr    = `legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod`
	maxV4StandingOrderCountPerPeriod           = 2147483647
)

var (
	// nolint:gochecknoglobals
	regexV4StandingOrderFrequencyCode = regexp.MustCompile(`^(ADHO|YEAR|DAIL|FRTN|INDA|MNTH|QURT|MIAN|WEEK|WODL|FOWK|TWMH|FOMH|FIMH|ALMH|NONE|LWMH|LXMH|TWYR)$`)
	// nolint:gochecknoglobals
	regexV4StandingOrderPointInTime = regexp.MustCompile(`^(-[0-9]|[0-9]{1,2})$`)
)

// DefaultV4StandingOrderFrequency returns the minimal viable v4 standing-order frequency.
func DefaultV4StandingOrderFrequency() V4StandingOrderFrequency {
	return V4StandingOrderFrequency{
		Type:        defaultV4StandingOrderFrequencyType,
		PointInTime: defaultV4StandingOrderFrequencyPointInTime,
	}
}

// Validate ensures the v4 standing-order frequency is a valid OBFrequency6 subset.
func (f V4StandingOrderFrequency) Validate() error {
	if f.Type == "" {
		return fmt.Errorf("Type: cannot be blank")
	}
	isFrequencyCode := regexV4StandingOrderFrequencyCode.MatchString(f.Type)
	isLegacyFrequency := regexPaymentFrequency.MatchString(f.Type)
	if !isFrequencyCode && !isLegacyFrequency {
		return fmt.Errorf("Type: %s", regexV4StandingOrderFrequencyTypeErr)
	}
	if isLegacyFrequency && !isFrequencyCode && (f.PointInTime != "" || f.CountPerPeriod != nil) {
		return fmt.Errorf(legacyV4StandingOrderFrequencyFieldsErr)
	}
	if f.PointInTime != "" && f.CountPerPeriod != nil {
		return fmt.Errorf("PointInTime and CountPerPeriod are mutually exclusive")
	}
	if f.PointInTime != "" && !regexV4StandingOrderPointInTime.MatchString(f.PointInTime) {
		return fmt.Errorf("PointInTime: %s", regexV4StandingOrderPointInTimeErr)
	}
	if f.CountPerPeriod != nil && (*f.CountPerPeriod <= 0 || *f.CountPerPeriod > maxV4StandingOrderCountPerPeriod) {
		return fmt.Errorf("CountPerPeriod: must be a positive Int32")
	}
	return nil
}

// JSON returns the payload-shaped OBFrequency6 object as a raw JSON fragment.
func (f V4StandingOrderFrequency) JSON() (string, error) {
	body, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
