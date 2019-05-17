package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	validation "github.com/go-ozzo/ozzo-validation"
)

func TestPaymentValidateSchemeName(t *testing.T) {
	require := test.NewRequire(t)

	for _, validSchemeName := range OBExternalAccountIdentification4Codes {
		data := fmt.Sprintf(`
{
    "scheme_name": "%s",
    "identification": "20202010981789"
}
		`, validSchemeName)
		payment := Payment{}
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

		require.NoError(payment.Validate())
	}

	for _, validSchemeName := range OBExternalAccountIdentification4Codes {
		invalidSchemeName := fmt.Sprintf("FAKE_%s", validSchemeName)
		data := fmt.Sprintf(`
{
    "scheme_name": "%s",
    "identification": "20202010981789"
}
		`, invalidSchemeName)
		payment := Payment{}
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

		require.EqualError(payment.Validate(), "scheme_name: the length must be between 1 and 40.")
	}
}

func TestPaymentValidateIdentification(t *testing.T) {
	require := test.NewRequire(t)
	schemaName, ok := OBExternalAccountIdentification4Codes[0].(string)
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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

		require.EqualError(payment.Validate(), "identification: the length must be between 1 and 256.")
	}
}

func TestPaymentValidateName(t *testing.T) {
	require := test.NewRequire(t)
	schemaName, ok := OBExternalAccountIdentification4Codes[0].(string)
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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

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
		err := json.Unmarshal([]byte(data), &payment)
		require.NoError(err)

		require.EqualError(payment.Validate(), "name: the length must be between 1 and 70.")
	}
}

func TestPaymentValidateInstructedAmount(t *testing.T) {
	require := test.NewRequire(t)
	a := InstructedAmount{Currency: "USD", Value: 1.0}
	err := validation.Validate(&a)
	require.Nil(err)
}

func TestPaymentValidateInstructedAmountFails(t *testing.T) {
	require := test.NewRequire(t)
	a := InstructedAmount{Currency: "not a valid currency", Value: 1.0}
	err := validation.Validate(&a)
	require.NotNil(err)
}
