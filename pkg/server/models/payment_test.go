package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/test"
	validation "github.com/go-ozzo/ozzo-validation"
)

func TestPaymentValidateIdentification(t *testing.T) {
	require := test.NewRequire(t)
	schemaName := "UK.OBIE.IBAN"

	// `Identification` specified
	{
		data := fmt.Sprintf(`
{
	"scheme_name": "%s",
	"identification": "20202010981789"
}
	`, schemaName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.NoError(payment.Validate())
	}
	// `Identification` not specified
	{
		data := fmt.Sprintf(`
{
	"scheme_name": "%s"
}
	`, schemaName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.EqualError(payment.Validate(), "identification: cannot be blank.")
	}
	// `Identification` should be between 1-256 characters
	{
		identification := strings.Repeat("i", 257)
		data := fmt.Sprintf(`
{
	"scheme_name": "%s",
	"identification": "%s"
}
	`, schemaName, identification)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.EqualError(payment.Validate(), "identification: the length must be between 1 and 256.")
	}
}

func TestPaymentValidateName(t *testing.T) {
	require := test.NewRequire(t)
	schemaName := "UK.OBIE.IBAN"

	// `Name` does not need to be present according to specification
	{
		data := fmt.Sprintf(`
{
	"scheme_name": "%s",
	"identification": "20202010981789"
}
		`, schemaName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.NoError(payment.Validate())
	}
	// If `Name` is present, it should be between 1-70 characters
	{
		name := strings.Repeat("n", 71)
		data := fmt.Sprintf(`
{
	"scheme_name": "%s",
	"identification": "20202010981789",
	"name": "%s"
}
		`, schemaName, name)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))

		require.EqualError(payment.Validate(), "name: the length must be between 1 and 70.")
	}
}

func TestPaymentValidateInstructedAmount(t *testing.T) {
	require := test.NewRequire(t)
	a := InstructedAmount{Currency: "USD", Value: "1.0"}
	require.NoError(validation.Validate(&a))
}

func TestPaymentValidateInstructedAmountFails(t *testing.T) {
	require := test.NewRequire(t)
	a := InstructedAmount{Currency: "not a valid currency", Value: "1.0"}
	require.EqualError(validation.Validate(&a), fmt.Sprintf("currency: %+v.", regexInstructedAmountCurrencyErr))
}

func TestServer_Payment_InstructedAmountValue_String(t *testing.T) {
	assert := test.NewAssert(t)

	tests := []struct {
		Value         string
		ExpectedError bool
	}{
		{
			Value:         "1.0",
			ExpectedError: false,
		},
		{
			Value:         "0.1",
			ExpectedError: false,
		},
		{
			Value:         "0.0001",
			ExpectedError: false,
		},
		{
			Value:         "0.00001",
			ExpectedError: false,
		},
		{
			Value:         "1111111111111.0",
			ExpectedError: false,
		},
		{
			Value:         "0.000001",
			ExpectedError: true,
		},
		{
			Value:         "0.0000001",
			ExpectedError: true,
		},
		{
			Value:         "0.00000001",
			ExpectedError: true,
		},
		{
			Value:         "0.000000001",
			ExpectedError: true,
		},
		{
			Value:         "0.0000000001",
			ExpectedError: true,
		},
		{
			Value:         "11111111111111.0",
			ExpectedError: true,
		},
		{
			Value:         "1111111111111.000001",
			ExpectedError: true,
		},
	}

	for _, test := range tests {
		i := InstructedAmount{Currency: "GBP", Value: test.Value}
		err := validation.Validate(&i)
		if test.ExpectedError {
			assert.EqualError(err, fmt.Sprintf("value: %+v.", regexInstructedAmountValueErr), fmt.Sprintf("Value=%+v", test.Value))
		} else {
			assert.NoError(err, fmt.Sprintf("Value=%+v", test.Value))
		}
	}
}

func TestPaymentFrequency(t *testing.T) {
	require := test.NewRequire(t)

	// Taken from
	// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623689/Domestic+Standing+Orders+v3.1#DomesticStandingOrdersv3.1-FrequencyExamples
	tests := []struct {
		Value         string
		ExpectedError bool
	}{
		{
			Value:         "EvryDay",
			ExpectedError: false,
		},
		{
			Value:         "EvryWorkgDay",
			ExpectedError: false,
		},
		{
			Value:         "NotKnown",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlDay:02",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlDay:31",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlWkDay:01:03",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlWkDay:02:03",
			ExpectedError: false,
		},
		{
			Value:         "WkInMnthDay:02:03",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlMnthDay:01:-01",
			ExpectedError: false,
		},
		{
			Value:         "IntrvlMnthDay:06:15",
			ExpectedError: false,
		},
		{
			Value:         "QtrDay:ENGLISH",
			ExpectedError: false,
		},
		{
			Value:         "WkInMnthDay:01:01",
			ExpectedError: false,
		},
		{
			Value:         "BadValue",
			ExpectedError: true,
		},
		{
			Value:         "IntrvlDay:01",
			ExpectedError: true,
		},
		{
			Value:         "IntrvlDay:32",
			ExpectedError: true,
		},
		{
			Value:         "WEEK",
			ExpectedError: true,
		},
	}

	for _, test := range tests {
		p := PaymentFrequency(test.Value)
		err := validation.Validate(&p)
		if test.ExpectedError {
			require.EqualError(err, regexPaymentFrequencyErr)
		} else {
			require.NoError(err, fmt.Sprintf("Value=%+v", test.Value))
		}
	}
}

func TestV4StandingOrderFrequency(t *testing.T) {
	require := test.NewRequire(t)
	one := 1
	zero := 0
	maxInt32 := 2147483647
	aboveMaxInt32 := 2147483648

	tests := []struct {
		name          string
		frequency     V4StandingOrderFrequency
		expectedError string
	}{
		{
			name:      "default",
			frequency: DefaultV4StandingOrderFrequency(),
		},
		{
			name: "type only",
			frequency: V4StandingOrderFrequency{
				Type: "MNTH",
			},
		},
		{
			name: "two digit point in time",
			frequency: V4StandingOrderFrequency{
				Type:        "MNTH",
				PointInTime: "01",
			},
		},
		{
			name: "single digit point in time",
			frequency: V4StandingOrderFrequency{
				Type:        "MNTH",
				PointInTime: "1",
			},
		},
		{
			name: "negative point in time",
			frequency: V4StandingOrderFrequency{
				Type:        "MNTH",
				PointInTime: "-1",
			},
		},
		{
			name: "maximum positive point in time",
			frequency: V4StandingOrderFrequency{
				Type:        "MNTH",
				PointInTime: "99",
			},
		},
		{
			name: "count per period",
			frequency: V4StandingOrderFrequency{
				Type:           "WEEK",
				CountPerPeriod: &one,
			},
		},
		{
			name: "maximum count per period",
			frequency: V4StandingOrderFrequency{
				Type:           "WEEK",
				CountPerPeriod: &maxInt32,
			},
		},
		{
			name: "missing type",
			frequency: V4StandingOrderFrequency{
				PointInTime: "03",
			},
			expectedError: "Type: cannot be blank",
		},
		{
			name: "new code list value LWMH",
			frequency: V4StandingOrderFrequency{
				Type: "LWMH",
			},
		},
		{
			name: "new code list value LXMH",
			frequency: V4StandingOrderFrequency{
				Type: "LXMH",
			},
		},
		{
			name: "new code list value TWYR",
			frequency: V4StandingOrderFrequency{
				Type: "TWYR",
			},
		},
		{
			name: "legacy type only",
			frequency: V4StandingOrderFrequency{
				Type: "EvryDay",
			},
		},
		{
			name: "legacy interval type only",
			frequency: V4StandingOrderFrequency{
				Type: "IntrvlWkDay:01:03",
			},
		},
		{
			name: "invalid type",
			frequency: V4StandingOrderFrequency{
				Type: "BadValue",
			},
			expectedError: "Type: must be an OBFrequency6Code or Frequency_1 value",
		},
		{
			name: "legacy type with point in time",
			frequency: V4StandingOrderFrequency{
				Type:        "EvryDay",
				PointInTime: "03",
			},
			expectedError: "legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod",
		},
		{
			name: "legacy type with count per period",
			frequency: V4StandingOrderFrequency{
				Type:           "IntrvlWkDay:01:03",
				CountPerPeriod: &one,
			},
			expectedError: "legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod",
		},
		{
			name: "invalid point in time too long",
			frequency: V4StandingOrderFrequency{
				Type:        "WEEK",
				PointInTime: "100",
			},
			expectedError: "PointInTime: must be numeric text up to two characters, including negative single-digit values",
		},
		{
			name: "invalid negative point in time too long",
			frequency: V4StandingOrderFrequency{
				Type:        "WEEK",
				PointInTime: "-10",
			},
			expectedError: "PointInTime: must be numeric text up to two characters, including negative single-digit values",
		},
		{
			name: "invalid point in time non numeric",
			frequency: V4StandingOrderFrequency{
				Type:        "WEEK",
				PointInTime: "AA",
			},
			expectedError: "PointInTime: must be numeric text up to two characters, including negative single-digit values",
		},
		{
			name: "invalid count per period zero",
			frequency: V4StandingOrderFrequency{
				Type:           "WEEK",
				CountPerPeriod: &zero,
			},
			expectedError: "CountPerPeriod: must be a positive Int32",
		},
		{
			name: "invalid count per period above Int32",
			frequency: V4StandingOrderFrequency{
				Type:           "WEEK",
				CountPerPeriod: &aboveMaxInt32,
			},
			expectedError: "CountPerPeriod: must be a positive Int32",
		},
		{
			name: "mutually exclusive point and count",
			frequency: V4StandingOrderFrequency{
				Type:           "WEEK",
				PointInTime:    "03",
				CountPerPeriod: &one,
			},
			expectedError: "PointInTime and CountPerPeriod are mutually exclusive",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.frequency.Validate()
			if tc.expectedError != "" {
				require.EqualError(err, tc.expectedError)
				return
			}
			require.NoError(err)
		})
	}
}
