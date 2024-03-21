package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var router routers.Router

func TestTesterResopd(t *testing.T) {
	filename := "spec/v3.1.8/account-info-openapi.json"
	sl := openapi3.NewLoader()
	doc, err := sl.LoadFromFile(filename)
	require.NoError(t, err)
	err = doc.Validate(context.Background())
	require.NoError(t, err)

	router, _ = legacyrouter.NewRouter(doc)
	require.NotNil(t, router)

	resp := ResponseWrapper{
		Status:      200,
		ContentType: "application/json",
		Body:        getAccounts,
	}

	req := RequestWrapper{
		Method: "GET",
		URL:    "/open-banking/v3.1/aisp/accounts",
	}
	err = validateTestResponse(t, req, resp, router)
	assert.NoError(t, err)
	if err != nil {
		t.Logf(err.Error())
	}
}

const getAccounts = `
		{
			"Data": {
			},
			"Links": {
				"Self": "https://rs1.obie.uk.ozoneapi.io/open-banking/v3.1/aisp/accounts"
			},
			"Meta": {
				"TotalPages": 1
			}
		}
	`

func TestVRPPost(t *testing.T) {
	filename := "spec/v3.1.8/variable-recurring-payments-openapi.json"
	sl := openapi3.NewLoader()
	doc, err := sl.LoadFromFile(filename)
	require.NoError(t, err)
	err = doc.Validate(context.Background())
	require.NoError(t, err)

	resp := ResponseWrapper{
		Status:      201,
		ContentType: "application/json",
		Body:        goodVrpConsents,
	}

	badresp := ResponseWrapper{
		Status:      201,
		ContentType: "application/json",
		Body:        badVrpConsents,
	}

	req := RequestWrapper{
		Method: "POST",
		URL:    "/open-banking/v3.1/pisp/domestic-vrp-consents",
	}

	router, _ = legacyrouter.NewRouter(doc)
	require.NotNil(t, router)

	err = validateTestResponse(t, req, resp, router)
	assert.NoError(t, err)

	err = validateTestResponse(t, req, badresp, router)
	assert.Error(t, err)

}

func validateTestResponse(t *testing.T, req RequestWrapper, resp ResponseWrapper, arouter routers.Router) error {
	httpReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(""))
	require.NoError(t, err)
	httpReq.Header.Set(headerCT, req.ContentType)

	route, pathParams, err := arouter.FindRoute(httpReq)
	if err != nil {
		fmt.Printf("Failed route params %s %s\n", req.Method, req.URL)
	}
	require.NoError(t, err)

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 resp.Status,
		Header: http.Header{
			headerCT: []string{
				resp.ContentType,
			},
		},
		Options: &openapi3filter.Options{
			ExcludeRequestBody:    true,
			IncludeResponseStatus: true,
			MultiError:            true,
		},
	}

	if resp.Body != "" {
		responseValidationInput.SetBodyBytes([]byte(resp.Body))
	}

	err = openapi3filter.ValidateResponse(context.Background(), responseValidationInput)
	return err
}

func marshalReader(value interface{}) io.ReadCloser {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(data))
}

const goodVrpConsents = `{
    "Data": {
        "ReadRefundAccount": "Yes",
        "ControlParameters": {
            "PSUAuthenticationMethods": [
                "UK.OBIE.SCANotRequired"
            ],
            "VRPType": [
                "UK.OBIE.VRPType.Sweeping"
            ],
            "ValidFromDateTime": "2017-06-05T15:15:13+00:00",
            "ValidToDateTime": "2020-06-05T15:15:13+00:00",
            "MaximumIndividualAmount": {
                "Amount": "100.00",
                "Currency": "GBP"
            },
            "PeriodicLimits": [
                {
                    "Amount": "200.00",
                    "Currency": "GBP",
                    "PeriodAlignment": "Consent",
                    "PeriodType": "Week"
                }
            ]
        },
        "Initiation": {
            "CreditorAccount": {
                "SchemeName": "SortCodeAccountNumber",
                "Identification": "30949330000010",
                "SecondaryIdentification": "Roll 90210",
                "Name": "Marcus Sweepimus"
            },
            "RemittanceInformation": {
                "Reference": "Sweepco"
            }
        },
        "DebtorAccount": {
            "SchemeName": "UK.OBIE.SortCodeAccountNumber",
            "Identification": "70000170000001",
            "Name": "Marcus Sweepimus"
        },
        "ConsentId": "vrp-7c55935d-3ff9-4210-a695-b4e7646afd0c",
        "Status": "AwaitingAuthorisation",
        "CreationDateTime": "2021-07-12T18:03:54.314Z",
        "StatusUpdateDateTime": "2021-07-12T18:03:54.314Z"
    },
    "Risk": {
        "PaymentContextCode": "PartyToParty"
    },
    "Links": {
        "Self": "http://localhost:4700/open-banking/v3.1/pisp/domestic-vrp-consents/vrp-7c55935d-3ff9-4210-a695-b4e7646afd0c"
    },
    "Meta": {}
}`

const badVrpConsents = `{
    "Data": {
        "ReadRefundAccount": "Yes2",
        "ControlParameters": {
            "PSUAuthenticationMethods": [
                "UK.OBIE.SCANotRequired"
            ],
            "VRPType": [
                "UK.OBIE.VRPType.Sweeping"
            ],
            "ValidFromDateTime": "2017-06-05T15:15:13+00:00",
            "ValidToDateTime": "2020-06-05T15:15:13+00:00",
            "MaximumIndividualAmount": {
                "Amount": "100.00",
                "Currency": "GBP"
            },
            "PeriodicLimits": [
                {
                    "Amount": "200.00",
                    "Currency": "GBP",
                    "PeriodAlignment": "Consent",
                    "PeriodType": "Week"
                }
            ]
        },
        "Initiation": {
            "CreditorAccount": {
                "SchemeName": "SortCodeAccountNumber",
                "Identification": "30949330000010",
                "SecondaryIdentification": "Roll 90210",
                "Name": "Marcus Sweepimus"
            },
            "RemittanceInformation": {
                "Reference": "Sweepco"
            }
        },
        "DebtorAccount": {
            "SchemeName": "UK.OBIE.SortCodeAccountNumber",
            "Identification": "70000170000001",
            "Name": "Marcus Sweepimus"
        },
        "ConsentId": "vrp-7c55935d-3ff9-4210-a695-b4e7646afd0c",
        "Status": "AwaitingAuthorisations",
        "CreationDateTime": "2021-07-12T18:03:54.314Z",
        "StatusUpdateDateTime": "2021-07-12T18:03:54.314Z"
    },
    "Risk": {
        "PaymentContextCode": "PartyToParty"
    },
    "Links": {
        "Self": "http://localhost:4700/open-banking/v3.1/pisp/domestic-vrp-consents/vrp-7c55935d-3ff9-4210-a695-b4e7646afd0c"
    },
    "Meta": {}
}`
