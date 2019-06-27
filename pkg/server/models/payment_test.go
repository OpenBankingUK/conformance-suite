package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	validation "github.com/go-ozzo/ozzo-validation"
)

func TestPaymentValidateSchemeName(t *testing.T) {
	require := test.NewRequire(t)

	for _, validSchemeName := range OBExternalAccountIdentification4Codes() {
		data := fmt.Sprintf(`
{
    "scheme_name": "%s",
    "identification": "20202010981789"
}
		`, validSchemeName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.NoError(payment.Validate())
	}

	for _, validSchemeName := range OBExternalAccountIdentification4Codes() {
		invalidSchemeName := fmt.Sprintf("FAKE_%s", validSchemeName)
		data := fmt.Sprintf(`
{
    "scheme_name": "%s",
    "identification": "20202010981789"
}
		`, invalidSchemeName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.EqualError(payment.Validate(), "scheme_name: must be a valid value.")
	}

	// `SchemaName` should be between 1-40 characters
	{
		invalidSchemeName := strings.Repeat("s", 41)
		data := fmt.Sprintf(`
{
    "scheme_name": "%s",
    "identification": "20202010981789"
}
		`, invalidSchemeName)
		payment := Payment{}
		require.NoError(json.Unmarshal([]byte(data), &payment))
		require.EqualError(payment.Validate(), "scheme_name: the length must be between 1 and 40.")
	}
}

func TestPaymentValidateIdentification(t *testing.T) {
	require := test.NewRequire(t)
	schemaName, ok := OBExternalAccountIdentification4Codes()[0].(string)
	require.True(ok)

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
	schemaName, ok := OBExternalAccountIdentification4Codes()[0].(string)
	require.True(ok)

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
