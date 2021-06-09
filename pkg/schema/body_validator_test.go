package schema

import (
	"fmt"
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

func TestSchemaBodyChecks(t *testing.T) {
	ais316, _ := loads.Spec("spec/v3.1.6/account-info-swagger-flattened.json")
	ais317, _ := loads.Spec("spec/v3.1.7/account-info-swagger-flattened.json")
	pis316, _ := loads.Spec("spec/v3.1.6/payment-initiation-swagger-flattened.json")
	pis317, _ := loads.Spec("spec/v3.1.7/payment-initiation-swagger-flattened.json")
	cbpii316, _ := loads.Spec("spec/v3.1.6/confirmation-funds-swagger-flattened.json")
	cbpii317, _ := loads.Spec("spec/v3.1.7/confirmation-funds-swagger-flattened.json")

	var tests = []struct {
		name         string
		doc          *loads.Document
		path         string
		responseBody string
		want         int
	}{
		{"316-ais-good", ais316, "/accounts", getAccountsResponse, 0},
		{"317-ais-good", ais317, "/accounts", getAccountsResponse, 0},
		{"316-ais-bad", ais316, "/accounts", getBadAccountsResponse, 4},
		// account Data.Account.SubType not required in 3.1.7
		{"317-ais-bad", ais317, "/accounts", getBadAccountsResponse, 3},
		{"316-pis-good", pis316, "/domestic-payments/pv349a7", aGetDomesticPayment, 0},
		{"317-pis-good", pis317, "/domestic-payments/pv349a7", aGetDomesticPayment, 0},
		{"316-pis-bad", pis316, "/domestic-payments/pv349a7", aBadGetDomesticPayment, 5},
		// domestic-payment Data.Debtor response does not contain addtionalProperties=false in 3.1.7 pis
		{"317-pis-bad", pis317, "/domestic-payments/pv349a7", aBadGetDomesticPayment, 4},
		{"316-cbpii-good", cbpii316, "/funds-confirmation-consents/fccx", cbpiiResponse, 0},
		{"317-cbpii-good", cbpii317, "/funds-confirmation-consents/fccx", cbpiiResponse, 0},
		{"316-cbpii-bad", cbpii316, "/funds-confirmation-consents/fccx", badCbpiiResponse, 3},
		{"317-cbpii-bad", cbpii317, "/funds-confirmation-consents/fccx", badCbpiiResponse, 3},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("SchemaTest:%s", tt.name)
		t.Run(testname, func(t *testing.T) {
			f := newFinder(tt.doc)
			validator := newBodyValidator(f)
			body := strings.NewReader(tt.responseBody)
			r := Response{
				Method:     "GET",
				Path:       tt.path,
				StatusCode: http.StatusOK,
				Body:       body,
			}

			failures, err := validator.Validate(r)
			require.NoError(t, err)
			numberFailures := len(failures)

			if numberFailures != tt.want {
				t.Errorf("got %d, want %d", numberFailures, tt.want)
			}
		})
	}
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

const getBadAccountsResponse = `
	{
		"Data": {
			"Account": [
				{
					"AccounxtId": "500000000000000000000001",
					"Currency": "GBP",
					"Nickname": "xxxx0101",
					"AccountType": "Personal",
					"AccounxtSubType": "CurrentAccount",
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
const aGetDomesticPayment = `
{
	"Data": {
		 "ConsentId": "sdp-1-a6626130-16a5-40ea-a112-49334487a204",
		 "Initiation": {
				"CreditorAccount": {
					 "Identification": "70000170000002",
					 "Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
					 "SchemeName": "UK.OBIE.SortCodeAccountNumber"
				},
				"EndToEndIdentification": "e2e-domestic-pay",
				"InstructedAmount": {
					 "Amount": "1.00",
					 "Currency": "GBP"
				},
				"InstructionIdentification": "a1a71ce50ec7490ea53a50c9baa92564"
		 },
		 "DomesticPaymentId": "pv3-49a7a93a-1d54-4b78-bf14-a5e12cb2b78d",
		 "CreationDateTime": "2020-05-19T10:16:26.893Z",
		 "Status": "Pending",
		 "StatusUpdateDateTime": "2020-05-19T10:16:26.893Z",
		 "ExpectedExecutionDateTime": "2020-05-19T10:16:26.893Z",
		 "ExpectedSettlementDateTime": "2020-05-19T10:16:26.893Z",
		 "Debtor": {
				"Name": "Octon Inc"
		 },
		 "MultiAuthorisation": {
				"Status": "AwaitingFurtherAuthorisation",
				"NumberRequired": 3,
				"NumberReceived": 1,
				"LastUpdateDateTime": "2020-05-19T10:16:26.914Z"
		 }
	},
	"Links": {
		 "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/pisp/domestic-payments/pv3-49a7a93a-1d54-4b78-bf14-a5e12cb2b78d"
	},
	"Meta": {}
}
`

const aBadGetDomesticPayment = `
{
	"Data": {
		 "ConsentId": "sdp-1-a6626130-16a5-40ea-a112-49334487a204",
		 "Initiation": {
				"CreditxorAccount": {
					 "Identification": "70000170000002",
					 "Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
					 "SchemeName": "UK.OBIE.SortCodeAccountNumber"
				},
				"EndToEndIdentification": "e2e-domestic-pay",
				"InstructedAmount": {
					 "Amount": "1.00",
					 "Currency": "GBP"
				},
				"InstructionIdentification": "a1a71ce50ec7490ea53a50c9baa92564"
		 },
		 "DomesticPaymentId": "pv3-49a7a93a-1d54-4b78-bf14-a5e12cb2b78d",
		 "CreatxionDateTime": "2020-05-19T10:16:26.893Z",
		 "Status": "Pending",
		 "StatusUpdateDateTime": "2020-05-19T10:16:26.893Z",
		 "ExpectedExecutionDateTime": "2020-05-19T10:16:26.893Z",
		 "ExpectedSettlementDateTime": "2020-05-19T10:16:26.893Z",
		 "Debtor": {
				"Naxme": "Octon Inc"
		 },
		 "MultiAuthorisation": {
				"Status": "AwaitingFurtherAuthorisation",
				"NumberRequired": 3,
				"NumberReceived": 1,
				"LastUpdateDateTime": "2020-05-19T10:16:26.914Z"
		 }
	},
	"Links": {
		 "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/pisp/domestic-payments/pv3-49a7a93a-1d54-4b78-bf14-a5e12cb2b78d"
	},
	"Meta": {}
}
`
const cbpiiResponse = `
{
	"Data": {
		 "DebtorAccount": {
				"Identification": "70000170000002",
				"Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
				"SchemeName": "UK.OBIE.SortCodeAccountNumber"
		 },
		 "ExpirationDateTime": "2021-01-01T00:00:00+01:00",
		 "ConsentId": "fcc-22a6e08c-d5fa-4159-9eed-c9f0c7398fff",
		 "CreationDateTime": "2020-05-21T12:13:22.269Z",
		 "Status": "Authorised",
		 "StatusUpdateDateTime": "2020-05-21T12:13:37.323Z"
	},
	"Links": {
		 "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-22a6e08c-d5fa-4159-9eed-c9f0c7398fff"
	},
	"Meta": {}
}`

const badCbpiiResponse = `
{
	"Data": {
		 "DebtorAccount": {
				"IdentifiXcation": "70000170000002",
				"Name": "Mr. Roberto Rastapopoulos & Ivan Sakharine & mits",
				"SchemeName": "UK.OBIE.SortCodeAccountNumber"
		 },
		 "ExpirationDateTime": "2021-01-01T00:00:00+01:00",
		 "ConsentXId": "fcc-22a6e08c-d5fa-4159-9eed-c9f0c7398fff",
		 "CreationDateTime": "2020-05-21T12:13:22.269Z",
		 "StXatus": "Authorised",
		 "StatusUpdateDateTime": "2020-05-21T12:13:37.323Z"
	},
	"Links": {
		 "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-22a6e08c-d5fa-4159-9eed-c9f0c7398fff"
	},
	"Meta": {}
}`
