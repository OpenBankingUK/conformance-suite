package executors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	resty "gopkg.in/resty.v1"
)

func TestSimulatedChainedOzoneRequest(t *testing.T) {
	tracer.Silent = true
	executor := &ozoneResponder{}
	chainedOzoneHeadlessAccounts(t, executor)
}

func chainedOzoneHeadlessAccounts(t *testing.T, executor model.TestCaseExecutor) {
	manifest, err := loadManifest("testdata/ozoneconnect.json")
	require.NoError(t, err)
	for _, rule := range manifest.Rules {
		rule.Executor = executor
		rulectx := &model.Context{}
		for _, sequence := range rule.Tests {
			for _, testcase := range sequence {
				testcase.AppEntry("Sequence Loop")
				req, err := testcase.Prepare(rulectx)
				assert.Nil(t, err)
				assert.NotNil(t, req)
				if err == nil {
					resp, err := rule.Execute(req, &testcase)
					require.Nil(t, err)
					require.NotNil(t, resp)
					if resp != nil {
						result, err := testcase.Validate(resp, rulectx)
						require.Nil(t, err)
						require.True(t, result)
					}
				}
				testcase.AppMsg(fmt.Sprintf("Context=%v", *rulectx))
				testcase.AppExit("End Sequence Loop")
			}
		}
	}
}

type ozoneResponder struct {
}

// ExecuteTestCase signature makes this an instance of  TestCaseExecutor
func (o *ozoneResponder) ExecuteTestCase(r *resty.Request, t *model.TestCase, ctx *model.Context) (*resty.Response, error) {
	appEntry("|--- OZONE RESPONDER ---|")
	defer appExit("OZONE-RESPONDER")
	appMsg("This should be the place")
	responseKey := t.Input.Method + " " + t.Input.Endpoint
	if strings.HasPrefix(responseKey, "GET https://modelobankauth2018.o3bank.co.uk:4101/auth?client_id") {
		responseKey = "GET https://modelobankauth2018.o3bank.co.uk:4101/auth?client_id"
	}
	appMsg(fmt.Sprintf("responsekey: %s", responseKey))
	fn := chainTest[responseKey]
	if fn != nil {
		if r.Body != nil {
			appMsg(fmt.Sprintf("%s %s", responseKey, r.Body.(string)))
		} else {
			appMsg(fmt.Sprintf("%s %v", responseKey, r.RawRequest))
		}
		return fn(r), nil
	}
	return nil, appErr("Cannot find handler function for:" + responseKey)
}

var chainTest = map[string]func(*resty.Request) *resty.Response{
	"POST https://modelobank2018.o3bank.co.uk:4201/token":                                          tokenEndpoint(),
	"POST https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.0/aisp/account-access-consents": accountAccessConsents(),
	"GET https://modelobankauth2018.o3bank.co.uk:4101/auth?client_id":                              consentIDFlow(),
	"GET https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.0/aisp/accounts":                 getAccounts(),
}

func appMsg(msg string) {
	tracer.AppMsg("OZONE", msg, "")
}

func appErr(msg string) error {
	tracer.AppErr("OZONE", msg, "")
	return errors.New(msg)
}

func appEntry(msg string) {
	tracer.AppEntry("OZONE", msg)
}

func appExit(msg string) {
	tracer.AppExit("OZONE", msg)
}

func consentIDFlow() func(r *resty.Request) *resty.Response {
	return func(r *resty.Request) *resty.Response {
		appMsg("GetConsentID")
		raw := http.Response{}
		raw.Header = make(http.Header)
		raw.StatusCode = 302
		raw.Header.Add("Location", "https://test.mybank.co.uk/redir?code=93101909-ea24-44ce-b237-d8245e729c9d&test=big")
		response := &resty.Response{
			RawResponse: &raw,
		}
		//appMsg(fmt.Sprintf("RESPONSE: %#v", response.RawResponse))
		//appMsg(fmt.Sprintf("Response: (%s)", response.String()))
		return response
	}
}

func tokenEndpoint() func(r *resty.Request) *resty.Response {
	return func(r *resty.Request) *resty.Response {
		appMsg("-------------")
		appMsg("tokenEndpoint:POST")

		if r.RawRequest != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.RawRequest.Body)
			newStr := buf.String()
			appMsg("rawBody: " + newStr)
		} else {
			appMsg(fmt.Sprintf(`REQUEST: %#v`, r))
		}

		grantType := r.FormData.Get("grant_type")
		var response *resty.Response
		appMsg("Grant type :" + grantType)
		if grantType == "client_credentials" {
			appMsg("Serve client_credentials")
			response = test.CreateHTTPResponse(200, "OK", string(clientCredentialsResponse))
		}

		if grantType == "authorization_code" {
			appMsg("Serve authorization_code")
			response = test.CreateHTTPResponse(200, "OK", string(authorizationCodeResponse))
		}
		if response == nil {
			appMsg("OMG !!!! Response == nil ")
		}
		appMsg(fmt.Sprintf("RESPONSE: %#v", response.RawResponse))
		return response
	}
}

func accountAccessConsents() func(r *resty.Request) *resty.Response {
	return func(r *resty.Request) *resty.Response {
		return test.CreateHTTPResponse(201, "OK", string(accountAccessConsentsResponse), "content-type", "klingon/text")
	}
}

func getAccounts() func(r *resty.Request) *resty.Response {
	return func(r *resty.Request) *resty.Response {
		appMsg("GetAccounts")
		appMsg(fmt.Sprintf("Request: %#v", r))
		return test.CreateHTTPResponse(200, "OK", string(getAccountsResponse), "content-type", "klingon/text")
	}
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadManifest(filename string) (model.Manifest, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return model.Manifest{}, err
	}
	var i model.Manifest
	err = json.Unmarshal(plan, &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

var (
	clientCredentialsResponse = []byte(`
	{
		"access_token": "4e407447-c1b3-4212-b355-71c425cbefb7",
		"token_type": "Bearer",
		"expires_in": 3600
	}`)

	accountAccessConsentsResponse = []byte(`
	{
		"Data": {
		   "Permissions": [
			  "ReadAccountsBasic",
			  "ReadAccountsDetail",
			  "ReadBalances",
			  "ReadBeneficiariesBasic",
			  "ReadBeneficiariesDetail",
			  "ReadDirectDebits",
			  "ReadTransactionsBasic",
			  "ReadTransactionsCredits",
			  "ReadTransactionsDebits",
			  "ReadTransactionsDetail",
			  "ReadProducts",
			  "ReadStandingOrdersDetail"
		   ],
		   "TransactionFromDateTime": "2016-01-01T10:40:00+02:00",
		   "TransactionToDateTime": "2025-12-31T10:40:00+02:00",
		   "ConsentId": "aac-ca19c2a5-f1fc-411a-828c-4e96390e9c95",
		   "CreationDateTime": "2018-12-14T10:55:18.748Z",
		   "Status": "AwaitingAuthorisation",
		   "StatusUpdateDateTime": "2018-12-14T10:55:18.748Z"
		},
		"Risk": {},
		"Links": {
		   "Self": "http://modelobank2018.o3bank.co.uk/open-banking/v3.0/aisp/account-access-consents/aac-ca19c2a5-f1fc-411a-828c-4e96390e9c95"
		},
		"Meta": {}
	 } `)

	authorizationCodeResponse = []byte(`
	 {
		"access_token": "1ddc15bb-94ad-4c15-a482-653ee9ec024d",
		"token_type": "Bearer",
		"expires_in": 3600,
		"scope": "openid accounts",
		"id_token": "eyJhbGciOiJQUzI1NiIsImtpZCI6ImlpRzh6UEJ4ZGJKQTFqdzN6ejMxcmJ4RldTayJ9.eyJzdWIiOiJhYWMtYzAyZmY2NWUtOTY3Yy00OTAyLTliZWItNjI4YWRlNTFlZmY4Iiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiYWFjLWMwMmZmNjVlLTk2N2MtNDkwMi05YmViLTYyOGFkZTUxZWZmOCIsImlzcyI6Imh0dHBzOi8vbW9kZWxvYmFua2F1dGgyMDE4Lm8zYmFuay5jby51azo0MTAxIiwiYXVkIjoiODY3MjM4NGUtOWEzMy00MzlmLTg5MjQtNjdiYjE0MzQwZDcxIiwiaWF0IjoxNTQ1MTYwMTgzLCJleHAiOjE1NDUxNjM3ODMsImNfaGFzaCI6IjFMNC1oYkVfemVZa3FLeU5hLTZwd1EiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EifQ.KaAbalGxaUDgFQjMtE8oC9pwT8N0DcwpSo08QnB8Y4wRu9FIXcq6RoWVbp6977370JCR9RNabzxLEQUDpYc_Sw_s812ur97_UD3hhkJcjFZf1kfZP3-Yw0qhcgkwIFJT0l1BbJBitFyDcpVOWghQu1vAjJ6x897LqxZiABBm7MFnvWxboQnnvMJlO16AJG7MKwi4n2QaCvQmIfRN_bmi1KfqYDPDfayi89A8O8e8mM-B4jgX8eeBLuvivq6TyVQ3Kc28Zh7vqDiwnpUuWoNYHaqPnWrHNpfT2QTFBaT5CUj4YUSg93wLgNBchiCjQz_V-lPEy4KCqJdeOxreZBptcg"
	 }	 
	 `)

	getAccountsResponse = []byte(`
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
			"Self": "http://modelobank2018.o3bank.co.uk/open-banking/v3.0/aisp/accounts"
		},
		"Meta": {
			"TotalPages": 1
		}
	}
	 `)
)
