package schema

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountsRoute(t *testing.T) {
	data := []struct {
		method string
		url    string
	}{
		{"POST", "/open-banking/v3.1/aisp/account-access-consents"},
		{"GET", "/open-banking/v3.1/aisp/account-access-consents/10001"},
		{"DELETE", "/open-banking/v3.1/aisp/account-access-consents/10002"},
		{"GET", "/open-banking/v3.1/aisp/accounts"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10003"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10004/balances"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10005/beneficiaries"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10006/direct-debits"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10007/offers"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10008/parties"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10009/party"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10010/product"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10011/scheduled-payments"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10012/standing-orders"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10013/statements"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10014/statements/20000"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10015/statements/20001/file"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10016/statements/20002/transactions"},
		{"GET", "/open-banking/v3.1/aisp/accounts/10017/transactions"},
		{"GET", "/open-banking/v3.1/aisp/balances"},
		{"GET", "/open-banking/v3.1/aisp/beneficiaries"},
		{"GET", "/open-banking/v3.1/aisp/direct-debits"},
		{"GET", "/open-banking/v3.1/aisp/offers"},
		{"GET", "/open-banking/v3.1/aisp/party"},
		{"GET", "/open-banking/v3.1/aisp/scheduled-payments"},
		{"GET", "/open-banking/v3.1/aisp/standing-orders"},
		{"GET", "/open-banking/v3.1/aisp/statements"},
		{"GET", "/open-banking/v3.1/aisp/transactions"},
	}

	validator, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v3.1.8")
	require.NoError(t, err)

	for _, row := range data {
		req, err := createHTTPReq(row.method, row.url)
		route, pathParams, err := validator.findTestRoute(req)
		for key, element := range pathParams {
			fmt.Printf("%s: %s\n", key, element)
		}
		require.NoError(t, err)
		_ = route
	}
}

func TestVrpRoutes(t *testing.T) {
	data := []struct {
		method string
		url    string
	}{
		{"POST", "/open-banking/v3.1/pisp/domestic-vrp-consents"},
		{"GET", "/open-banking/v3.1/pisp/domestic-vrp-consents/10001"},
		{"DELETE", "/open-banking/v3.1/pisp/domestic-vrp-consents/10002"},
		{"POST", "/open-banking/v3.1/pisp/domestic-vrp-consents/10003/funds-confirmation"},
		{"POST", "/open-banking/v3.1/pisp/domestic-vrps"},
		{"GET", "/open-banking/v3.1/pisp/domestic-vrps/10004"},
		{"GET", "/open-banking/v3.1/pisp/domestic-vrps/10005/payment-details"},
	}

	validator, err := NewRawOpenAPI3Validator("OBIE VRP Profile", "v3.1.8")
	require.NoError(t, err)

	for _, row := range data {
		req, err := createHTTPReq(row.method, row.url)
		route, pathParams, err := validator.findTestRoute(req)
		for key, element := range pathParams {
			fmt.Printf("%s: %s\n", key, element)
		}
		require.NoError(t, err)

		_ = route
	}
}

func createHTTPReqFromResponse(resp HTTPResponse) (*http.Request, error) {
	req, err := http.NewRequest(resp.Method, resp.Path, strings.NewReader(""))
	req.Header = http.Header{"Content-type": []string{"application/json; charset=utf-8"}}
	return req, err
}

func TestAcc10000TestResponse(t *testing.T) {
	validator, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v3.1.8")
	require.NoError(t, err)

	r := HTTPResponse{
		Method:     "GET",
		Path:       acc10000responseReqURL,
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(acc10000response),
		Header:     http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
	}

	_, err = validator.Validate(r)
	require.NoError(t, err)
}

func TestAcc10000TestResponseCapitalUtfNoSpace(t *testing.T) {
	validator, err := NewRawOpenAPI3Validator("Account and Transaction API Specification", "v3.1.8")
	require.NoError(t, err)

	r := HTTPResponse{
		Method:     "GET",
		Path:       acc10000responseReqURL,
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(acc10000response),
		Header:     http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}},
	}

	_, err = validator.Validate(r)
	require.NoError(t, err)
}

func TestCbpIITestResponseCapitalUtfNoSpace(t *testing.T) {
	validator, err := NewRawOpenAPI3Validator("Confirmation of Funds API Specification", "v3.1.8")
	require.NoError(t, err)

	r := HTTPResponse{
		Method:     "GET",
		Path:       cbpiiGoodResponseUrl,
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(cbpiiGoodResponse),
		Header:     http.Header{"Content-Type": []string{"application/json;charset=UTF-8"}},
	}

	_, err = validator.Validate(r)
	require.NoError(t, err)
}

func TestVrp100200Response(t *testing.T) {
	validator, err := NewRawOpenAPI3Validator("OBIE VRP Profile", "v3.1.8")
	require.NoError(t, err)

	r := HTTPResponse{
		Method:     "GET",
		Path:       vrp100200ReqURL,
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(vrp100200Response),
		Header:     http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
	}

	_, err = validator.Validate(r)
	require.NoError(t, err)
}

const acc10000responseReqURL = "/open-banking/v3.1/aisp/accounts/100004000000000000000001"
const acc10000response = `{
	"Data": {
		 "Account": [{
				"AccountId": "100004000000000000000001",
				"Currency": "GBP",
				"Account": [
					 {
							"Name": "Mario International",
							"SchemeName": "UK.OBIE.SortCodeAccountNumber",
							"Identification": "10000109010101"
					 },
					 {
							"Name": "Mario International",
							"SchemeName": "UK.OBIE.IBAN",
							"Identification": "10000109010101"
					 }
				]
		 }]
	},
	"Links": {
		 "Self": "http://localhost:4700/open-banking/v3.1/aisp/accounts/100004000000000000000001"
	},
	"Meta": {
		 "TotalPages": 1
	}
}`

const vrp100200ReqURL = "/open-banking/v3.1/pisp/domestic-vrp-consents/vrp-8ba1c1a1-6ffd-43fa-aac0-c1d1f8524f5d"

const vrp100200Response = `{
	"Data": {
		 "ControlParameters": {
				"MaximumIndividualAmount": {
					 "Amount": "1.00",
					 "Currency": "GBP"
				},
				"PSUAuthenticationMethods": [
					 "UK.OBIE.SCA"
				],
				"PeriodicLimits": [
					 {
							"Amount": "15.00",
							"Currency": "GBP",
							"PeriodAlignment": "Consent",
							"PeriodType": "Week"
					 }
				],
				"VRPType": [
					 "UK.OBIE.VRPType.Sweeping"
				],
				"ValidFromDateTime": "2017-06-05T15:15:13+00:00",
				"ValidToDateTime": "2020-06-05T15:15:13+00:00"
		 },
		 "Initiation": {
				"CreditorAccount": {
					 "Identification": "30949330000010",
					 "Name": "Marcus Sweepimus",
					 "SchemeName": "SortCodeAccountNumber",
					 "SecondaryIdentification": "Roll 90210"
				},
				"RemittanceInformation": {
					 "Reference": "Sweepco"
				}
		 },
		 "DebtorAccount": {
			 "SchemeName": "SortCodeAccountNumber",
			 "Identification": "Identification",
			 "Name": "Joe Smoe"
		 },
		 "ReadRefundAccount": "Yes",
		 "ConsentId": "vrp-8ba1c1a1-6ffd-43fa-aac0-c1d1f8524f5d",
		 "Status": "Authorised",
		 "CreationDateTime": "2021-08-01T18:28:23.523Z",
		 "StatusUpdateDateTime": "2021-08-01T18:28:33.252Z"
	},
	"Risk": {
		 "PaymentContextCode": "PartyToParty"
	},
	"Links": {
		 "Self": "http://localhost:4700/open-banking/v3.1/pisp/domestic-vrp-consents/vrp-8ba1c1a1-6ffd-43fa-aac0-c1d1f8524f5d"
	},
	"Meta": {}
}`
const cbpiiGoodResponseUrl = "/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-22a6e08c-d5fa-4159-9eed-c9f0c7398fff"
const cbpiiGoodResponse = `
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

func TestIsRequestPropertyOas3(t *testing.T) {
	validator, err := buildValidator("Payment Initiation API", "v3.1.8")
	require.NoError(t, err)

	/*
		Context:

		"name": "Payment Initiation API",
		"endpoints": [
		  {
		    "method": "POST",
		    "path": "/international-payment-consents",
		    "conditionalProperties": [
		      {
		        "schema": "OBCreditor1Name",
		        "name": "Name",
		        "path": "Data.Initiation.Creditor.Name",
		        "request": true,
		        "value": "Will B. Wealthy"
		      },
	*/
	isRequestProperty, objType, err := validator.IsRequestProperty("POST", "/international-payment-consents", "Data.Initiation.Creditor.Name")
	require.NoError(t, err)
	require.True(t, isRequestProperty)
	require.Equal(t, "[string]", objType)

	isRequestProperty, objType, err = validator.IsRequestProperty("POST", "/international-payment-consents", "Data.Initiation.Creditor.PostalAddress.AddressLine")
	require.Equal(t, "[array]", objType)
	require.NoError(t, err)
	require.True(t, isRequestProperty)

	isRequestProperty, objType, err = validator.IsRequestProperty("POST", "/international-payment-consents", "Data.Initiation.Creditor.PostalAddress")
	require.Equal(t, "[object]", objType)

	isRequestProperty, objType, err = validator.IsRequestProperty("POST", "/international-payment-consents", "Data.Initiation.Creditor.FakeProperty")
	require.NoError(t, err)
	require.False(t, isRequestProperty)

	isRequestProperty, objType, err = validator.IsRequestProperty("GET", "/international-payment-consents/{ConsentId}/funds-confirmation", "Data.Initiation.Creditor.Name")
	require.NoError(t, err)
	require.False(t, isRequestProperty)
}
