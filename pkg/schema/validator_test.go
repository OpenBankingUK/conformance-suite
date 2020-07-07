package schema

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidators_Validate_Transactions(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
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

const getTransactionsResponse = `
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
	require.NoError(t, err)
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

const getTransactionsResponseEmptyTransactionReference = `
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

func TestCheckRequestSchema(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/payment-initiation-swagger.flattened.json")
	require.NoError(t, err)

	spec := doc.Spec()

	for path, props := range spec.Paths.Paths {
		for meth, op := range getOperations(&props) {
			_ = meth
			if path == "/domestic-standing-order-consents" && meth == "POST" {
				for _, param := range op.Parameters {
					if param.ParamProps.In == "body" {
						schema := param.ParamProps.Schema
						found, _ := findPropertyInSchema(schema, "Data.Initiation.CreditorAccount.SecondaryIdentification", "")
						if found {
							t.Log("*** FOUND IT ******")
						} else {
							t.Fail()
						}
					}
				}
			}
		}
	}
}

func TestTraverseSchemaLookingforNonRequiredProperties(t *testing.T) { // Example traversal routine - no test fail
	doc, err := loads.Spec("spec/v3.1.0/payment-initiation-swagger.flattened.json")
	require.NoError(t, err)

	spec := doc.Spec()

	for path, props := range spec.Paths.Paths {
		for meth, op := range getOperations(&props) {
			t.Logf("%s %s %s\n", meth, path, op.ID)
			for _, param := range op.Parameters {
				if param.ParamProps.In == "body" {
					t.Logf("%s %s %t %s\n", param.Name, param.In, param.Required, param.Type)
					sc := param.ParamProps.Schema
					if sc != nil {
						dumpSchema(t, sc, "")
					}
				}
			}
		}
	}
}

func TestCheckPostalAddressFormat(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/payment-initiation-swagger.flattened.json")
	require.NoError(t, err)

	spec := doc.Spec()

	for path, props := range spec.Paths.Paths {
		for meth, op := range getOperations(&props) {
			_ = meth
			if path == "/international-payment-consents" && meth == "POST" {
				for _, param := range op.Parameters {
					if param.ParamProps.In == "body" {
						schema := param.ParamProps.Schema
						found, objtype := findPropertyInSchema(schema, "Data.Initiation.CreditorAgent.PostalAddress.AddressLine", "")
						if found {
							fmt.Printf("ObjectType: %s\n", objtype)
							fmt.Printf("*** FOUND IT ******")
						} else {
							t.Fail()
						}
					}
				}
			}
		}
	}
}

func getObjectType(sc *spec.Schema, path string) {
	for k, j := range sc.SchemaProps.Properties {
		var element string
		if len(path) == 0 {
			element = k
		} else {
			element = path + "." + k
		}

		if element == path {
			fmt.Printf("%s\n", element)
		}
		getObjectType(&j, element)
	}
}

func dumpSchema(t *testing.T, sc *spec.Schema, previousPath string) {
	for k, j := range sc.SchemaProps.Properties {
		var element string
		if len(previousPath) == 0 {
			element = k
		} else {
			element = previousPath + "." + k
		}
		if len(j.Required) > 0 {
			fmt.Printf("*** %s required:%s\n", element, j.Required)
		} else {
			fmt.Printf("%s\n", element)
			if element == "Data.Initiation.Creditor.PostalAddress.AddressLine" {
				fmt.Printf("%#v\n", j.SchemaProps.Type)

			}
		}
		dumpSchema(t, &j, element)
	}
}

func TestValidators_ValidateStandingOrderWithFreeformField(t *testing.T) {
	for _, spec := range []string{
		"spec/v3.1.0/account-info-swagger.flattened.json",
		"spec/v3.1.1/account-info-swagger.flattened.json",
		"spec/v3.1.2/account-info-swagger.flattened.json",
	} {
		doc, err := loads.Spec(spec)
		require.NoError(t, err)
		validator, err := newValidator(doc)
		require.NoError(t, err)
		for _, testCase := range getStandingOrdersWithFreeformFields {
			body := strings.NewReader(testCase)
			header := &http.Header{}
			header.Add("Content-type", "application/json; charset=utf-8")
			r := Response{
				Method:     "GET",
				Path:       "/accounts/500000000000000000000001/standing-orders",
				StatusCode: http.StatusOK,
				Body:       body,
				Header:     *header,
			}

			failures, err := validator.Validate(r)

			require.NoError(t, err)
			assert.Len(t, failures, 0, fmt.Sprintf("spec: %s", spec))
		}
	}
}

var getStandingOrdersWithFreeformFields = []string{`{
  "Data": {
    "StandingOrder": [
      {
        "AccountId": "string",
        "StandingOrderId": "string",
        "Frequency": "EvryDay",
        "Reference": "string",
        "FirstPaymentDateTime": "2019-12-02T09:34:48Z",
        "NextPaymentDateTime": "2019-12-02T09:34:48Z",
        "FinalPaymentDateTime": "2019-12-02T09:34:48Z",
        "StandingOrderStatusCode": "Active",
        "FirstPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "NextPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "FinalPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "SupplementaryData": {
          "FirstRecurringPaymentAmount": 123.45,
          "ExampleKey": "ExampleValue",
          "ExampleKey2": "ExampleValue2",
          "FirstPaymentAmount": {
            "Amount": "0.0"
          },
          "NestedKey": {
            "SubKey": {}
          }
        },
        "CreditorAgent": {
          "SchemeName": "UK.OBIE.BICFI",
          "Identification": "string"
        },
        "CreditorAccount": {
          "SchemeName": "UK.OBIE.SortCodeAccountNumber",
          "Identification": "string",
          "Name": "string",
          "SecondaryIdentification": "string"
        }
      }
    ]
  },
  "Links": {
    "Self": "http://example.com",
    "First": "http://example.com",
    "Prev": "http://example.com",
    "Next": "http://example.com",
    "Last": "http://example.com"
  },
  "Meta": {
    "TotalPages": 0,
    "FirstAvailableDateTime": "2019-12-02T09:34:48Z",
    "LastAvailableDateTime": "2019-12-02T09:34:48Z"
  }
}`,
	`{
  "Data": {
    "StandingOrder": [
      {
        "AccountId": "string",
        "StandingOrderId": "string",
        "Frequency": "EvryDay",
        "Reference": "string",
        "FirstPaymentDateTime": "2019-12-02T09:34:48Z",
        "NextPaymentDateTime": "2019-12-02T09:34:48Z",
        "FinalPaymentDateTime": "2019-12-02T09:34:48Z",
        "StandingOrderStatusCode": "Active",
        "FirstPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "NextPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "FinalPaymentAmount": {
          "Amount": "0.0",
          "Currency": "AAA"
        },
        "SupplementaryData": {},
        "CreditorAgent": {
          "SchemeName": "UK.OBIE.BICFI",
          "Identification": "string"
        },
        "CreditorAccount": {
          "SchemeName": "UK.OBIE.SortCodeAccountNumber",
          "Identification": "string",
          "Name": "string",
          "SecondaryIdentification": "string"
        }
      }
    ]
  },
  "Links": {
    "Self": "http://example.com",
    "First": "http://example.com",
    "Prev": "http://example.com",
    "Next": "http://example.com",
    "Last": "http://example.com"
  },
  "Meta": {
    "TotalPages": 0,
    "FirstAvailableDateTime": "2019-12-02T09:34:48Z",
    "LastAvailableDateTime": "2019-12-02T09:34:48Z"
  }
}`,
}
