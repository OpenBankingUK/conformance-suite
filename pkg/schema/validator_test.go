package schema

import (
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestValidators_Validate_AggregatesMultipleFailures(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/confirmation-funds-swagger.flattened.json")
	validator, err := newValidator(doc)
	require.NoError(t, err)
	body := strings.NewReader(getAccountsResponse)
	header := &http.Header{}
	header.Add("Content-type", "application/klingon")
	r := Response{
		Method:     "POST",
		Path:       "/funds-confirmation-consents",
		StatusCode: http.StatusTeapot,
		Body:       body,
		Header:     *header,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 3)
	expected := []Failure{
		{Message: "Content-Type Error: Should produce 'application/json', but got: 'application/klingon'"},
		{Message: "server Status 418 not defined by the spec"},
		{Message: "could't find a schema to validate for status code 418"},
	}
	assert.Equal(t, expected, failures)
}

func TestValidators_Validate_Transactions(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	validator, err := newValidator(doc)
	require.NoError(t, err)
	body := strings.NewReader(getTransactionsResponse)
	header := &http.Header{}
	header.Add("Content-type", "application/json; charset=utf-8")
	r := Response{
		Method:     "GET",
		Path:       "/accounts/500000000000000000000001/transactions",
		StatusCode: http.StatusOK,
		Body:       body,
		Header:     *header,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 0)
}

var getTransactionsResponse = `
		{
			"Data": {
				"Transaction": [
					{
						"AccountId": "500000000000000000000001",
						"Status": "Booked",
						"CreditDebitIndicator": "Credit",
						"BookingDateTime": "2017-04-05T10:43:07+00:00",
						"Amount": {
							"Amount": "100.10",
							"Currency": "GBP"
						}
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

func TestValidators_Validate_FailureEmptyOptionalProperty(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	validator, err := newValidator(doc)
	require.NoError(t, err)
	body := strings.NewReader(getTransactionsResponseEmptyTransactionReference)
	header := &http.Header{}
	header.Add("Content-type", "application/json; charset=utf-8")
	r := Response{
		Method:     "GET",
		Path:       "/accounts/500000000000000000000001/transactions",
		StatusCode: http.StatusOK,
		Body:       body,
		Header:     *header,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 1)
	assert.Equal(t, []Failure{{"Data.Transaction.TransactionReference in body should be at least 1 chars long"}}, failures)
}

var getTransactionsResponseEmptyTransactionReference = `
		{
			"Data": {
				"Transaction": [
					{
						"AccountId": "500000000000000000000001",
						"Status": "Booked",
						"CreditDebitIndicator": "Credit",
						"BookingDateTime": "2017-04-05T10:43:07+00:00",
						"Amount": {
							"Amount": "100.10",
							"Currency": "GBP"
						},
						"TransactionReference": ""
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
