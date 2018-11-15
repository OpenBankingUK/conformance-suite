package discovery

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// invalidTestCase
// name - name of the test case.
// config - the discovery model config.
// err - the expected err
type invalidTestCase struct {
	name        string
	config      string
	expectedErr string
}

func TestDiscovery_FromJSONString_Invalid_Cases(t *testing.T) {
	t.Parallel()

	testCases := []invalidTestCase{
		{
			name:        `json_needs_to_be_valid`,
			config:      ` `,
			expectedErr: `unexpected end of JSON input`,
		},
		{
			name:   `version_and_discoveryItems_array_needs_to_specified`,
			config: `{}`,
			expectedErr: `Key: 'Model.DiscoveryModel.Version' Error:Field validation for 'Version' failed on the 'required' tag
Key: 'Model.DiscoveryModel.DiscoveryItems' Error:Field validation for 'DiscoveryItems' failed on the 'required' tag`,
		},
		{
			name: `discoveryItems_array_needs_to_be_greater_than_one`,
			config: `
{
  "discoveryModel": {
	"version": "v0.0.1",
	"discoveryItems": [
	]
  }
}
			`,
			expectedErr: `Key: 'Model.DiscoveryModel.DiscoveryItems' Error:Field validation for 'DiscoveryItems' failed on the 'gt' tag`,
		},
		{
			name: `endpoints_needs_to_be_specified`,
			config: `
{
	"discoveryModel": {
		"version": "v0.0.1",
		"discoveryItems": [
			{
				"apiSpecification": {
					"name": "Account and Transaction API Specification",
					"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
					"version": "v3.0",
					"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
				},
				"openidConfigurationUri": "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
				"resourceBaseUri": "https://rs.aspsp.ob.forgerock.financial:443/",
				"endpoints": [
				]
			}
		]
	}
}
			`,
			expectedErr: `Key: 'Model.DiscoveryModel.DiscoveryItems[0].Endpoints' Error:Field validation for 'Endpoints' failed on the 'gt' tag`,
		},
		{
			name: `endpoints_path_and_method_need_to_be_valid`,
			config: `
{
	"discoveryModel": {
		"version": "v0.0.1",
		"discoveryItems": [
			{
				"apiSpecification": {
					"name": "Account and Transaction API Specification",
					"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
					"version": "v3.0",
					"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
				},
				"openidConfigurationUri": "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
				"resourceBaseUri": "https://rs.aspsp.ob.forgerock.financial:443/",
				"endpoints": [
					{
						"method": "FAKE-METHOD",
						"path": "/fake-path"
					},
					{
						"method": "FAKE-METHOD2",
						"path": "/fake-path2"
					}
				]
			}
		]
	}
}
			`,
			expectedErr: `discoveryItemIndex=0, invalid endpoint Method=FAKE-METHOD, Path=/fake-path
discoveryItemIndex=0, invalid endpoint Method=FAKE-METHOD2, Path=/fake-path2`,
		},
		{
			name: `endpoints_missing_mandatory_endpoints_accounts`,
			config: `
{
	"discoveryModel": {
		"version": "v0.0.1",
		"discoveryItems": [
			{
				"apiSpecification": {
					"name": "Account and Transaction API Specification",
					"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
					"version": "v3.0",
					"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
				},
				"openidConfigurationUri": "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
				"resourceBaseUri": "https://rs.aspsp.ob.forgerock.financial:443/",
				"endpoints": [
					{
						"method": "GET",
						"path": "/accounts/{AccountId}/balances"
					}
				]
			}
		]
	}
}
			`,
			expectedErr: `discoveryItemIndex=0, missing mandatory endpoint Method=POST, Path=/account-access-consents
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}
discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/accounts/{AccountId}/transactions`,
		},
	}

	for _, testCaseEntry := range testCases {
		// See: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		// for why we need this. Basically because we are running the tests in parallel using `t.Parallel`
		// we cannot bind to `testCaseEntry` as  there is a very good chance that when you run this code
		// you will see the last element being used all the time.
		func(testCase invalidTestCase) {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()
				assert := assert.New(t)

				model, err := FromJSONString(testCase.config)
				// fmt.Println()
				// fmt.Printf("%+v", err)
				// fmt.Println()

				assert.Nil(model)
				assert.EqualError(err, testCase.expectedErr)
			})
		}(testCaseEntry)
	}
}

func TestDiscovery_FromJSONString_Valid(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	discoveryExample, err := ioutil.ReadFile("../../docs/discovery-example.json")
	assert.NoError(err)
	assert.NotNil(discoveryExample)
	config := string(discoveryExample)

	modelExpected := &Model{
		DiscoveryModel: ModelDiscovery{
			Version: "v0.0.1",
			DiscoveryItems: []ModelDiscoveryItem{
				ModelDiscoveryItem{
					APISpecification: ModelAPISpecification{
						Name:          "Account and Transaction API Specification",
						URL:           "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
						Version:       "v3.0",
						SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json",
					},
					OpenidConfigurationURI: "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
					ResourceBaseURI:        "https://rs.aspsp.ob.forgerock.financial:443/",
					Endpoints: []ModelEndpoint{
						ModelEndpoint{
							Method:                "POST",
							Path:                  "/account-access-consents",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{
							Method:                "GET",
							Path:                  "/account-access-consents/{ConsentId}",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{Method: "DELETE",
							Path:                  "/account-access-consents/{ConsentId}",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{Method: "GET",
							Path:                  "/accounts/{AccountId}/product",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{Method: "GET",
							Path: "/accounts/{AccountId}/transactions",
							ConditionalProperties: []ModelConditionalProperties{
								ModelConditionalProperties{
									Schema:   "OBTransaction3Detail",
									Property: "Balance",
									Path:     "Data.Transaction[*].Balance",
								},
								ModelConditionalProperties{
									Schema:   "OBTransaction3Detail",
									Property: "MerchantDetails",
									Path:     "Data.Transaction[*].MerchantDetails",
								},
								ModelConditionalProperties{
									Schema:   "OBTransaction3Basic",
									Property: "TransactionReference",
									Path:     "Data.Transaction[*].TransactionReference",
								},
								ModelConditionalProperties{
									Schema:   "OBTransaction3Detail",
									Property: "TransactionReference",
									Path:     "Data.Transaction[*].TransactionReference",
								},
							},
						},
						ModelEndpoint{
							Method:                "GET",
							Path:                  "/accounts",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{
							Method:                "GET",
							Path:                  "/accounts/{AccountId}",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
						ModelEndpoint{
							Method:                "GET",
							Path:                  "/accounts/{AccountId}/balances",
							ConditionalProperties: []ModelConditionalProperties(nil),
						},
					},
				},
				ModelDiscoveryItem{
					APISpecification: ModelAPISpecification{
						Name:          "Payment Initiation API",
						URL:           "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367011/Payment+Initiation+API+Specification+-+v3.0",
						Version:       "v3.0",
						SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json",
					},
					OpenidConfigurationURI: "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
					ResourceBaseURI:        "https://rs.aspsp.ob.forgerock.financial:443/",
					Endpoints: []ModelEndpoint{ModelEndpoint{
						Method: "POST",
						Path:   "/domestic-payment-consents",
						ConditionalProperties: []ModelConditionalProperties{
							ModelConditionalProperties{
								Schema:   "OBWriteDataDomesticConsentResponse1",
								Property: "Charges",
								Path:     "Data.Charges",
							},
						},
					},
						ModelEndpoint{Method: "GET",
							Path: "/domestic-payment-consents/{ConsentId}",
							ConditionalProperties: []ModelConditionalProperties{
								ModelConditionalProperties{
									Schema:   "OBWriteDataDomesticConsentResponse1",
									Property: "Charges",
									Path:     "Data.Charges",
								},
							},
						},
						ModelEndpoint{Method: "POST",
							Path: "/domestic-payments",
							ConditionalProperties: []ModelConditionalProperties{
								ModelConditionalProperties{
									Schema:   "OBWriteDataDomesticResponse1",
									Property: "Charges",
									Path:     "Data.Charges",
								},
							},
						},
						ModelEndpoint{Method: "GET",
							Path: "/domestic-payments/{DomesticPaymentId}",
							ConditionalProperties: []ModelConditionalProperties{
								ModelConditionalProperties{
									Schema:   "OBWriteDataDomesticResponse1",
									Property: "Charges",
									Path:     "Data.Charges",
								},
							},
						},
					},
				},
			},
		},
	}
	modelActual, err := FromJSONString(config)

	assert.NoError(err)
	assert.Exactly(modelExpected, modelActual)
}

func TestDiscovery_Version(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal(Version(), "v0.0.1")
}
