package schema

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyValidator_Validate(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)
	validator := newBodyValidator(f)
	body := strings.NewReader(getAccountsResponse)
	r := Response{
		Method:     "GET",
		Path:       "/accounts",
		StatusCode: http.StatusOK,
		Body:       body,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 0)
}

const getAccountsResponse = `
		{
			"Data": {
				"Account": [
					{
						"AccountId": "500000000000000000000001",
						"Currency": "GBP",
						"Nickname": "xxxx0101",
						"AccountType": "Personal",
						"AccountSubType": "CurrentAccount",
						"Account": [
						{
							"SchemeName": "UK.OBIE.SortCodeAccountNumber",
							"Identification": "10000119820101",
							"SecondaryIdentification": "Roll No. 001"
						}
						]
					}
				]
			},
			"Links": {
				"Self": "http://modelobank2018.o3bank.co.uk/open-banking/v3.1/aisp/accounts"
			},
			"Meta": {
				"TotalPages": 1
			}
		}
	`

func TestBodyValidator_Validate_ReturnsFailures(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)
	validator := newBodyValidator(f)
	body := strings.NewReader(`{}`)
	r := Response{
		Method:     "GET",
		Path:       "/accounts",
		StatusCode: http.StatusOK,
		Body:       body,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 3)
	expected := []Failure{
		{Message: ".Data in body is required"},
		{Message: ".Links in body is required"},
		{Message: ".Meta in body is required"},
	}
	assert.Equal(t, expected, failures)
}

func TestBodyValidator_Validate_ReturnsFailuresFor317(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.7/account-info-swagger-flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)
	validator := newBodyValidator(f)
	body := strings.NewReader(`{ "Data": {
		"Account": [
			 {
					"AccountId": "700004000000000000000005",
					"Currency": "GBP",
					"Nickname": "xxxx0005",
					"AccountType": "Business",
					"AccountSubType": "CurrentAccount",
					"Account": [
						 {
								"SchemeName": "UK.OBIE.SortCodeAccountNumber",
								"Identification": "70000170000005",
								"Name": "Octon Inc"
						 }
					]
			 }
		]
 },
 "Links": {
		"Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/accounts/700004000000000000000005"
 }
 }`)
	r := Response{
		Method:     "GET",
		Path:       "/accounts",
		StatusCode: http.StatusOK,
		Body:       body,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 1)
	expected := []Failure{
		{Message: "Data.Account.AccountSubType in body should be one of [ChargeCard CreditCard CurrentAccount EMoney Loan Mortgage PrePaidCard Savings]"},
	}
	assert.Equal(t, expected, failures)
}
