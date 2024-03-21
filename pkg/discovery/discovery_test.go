package discovery

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConditionalProperties(t *testing.T) {

	tt := []struct{ DiscoveryPayload []byte }{
		{
			//testing legacy validator (Swagger)
			DiscoveryPayload: testDiscoveryPayload,
		},
		{
			//testing OAS3 Validator
			DiscoveryPayload: testDiscoveryPayloadOAS3,
		},
	}

	for _, ti := range tt {
		disco, err := UnmarshalDiscoveryJSONBytes(ti.DiscoveryPayload)
		assert.Nil(t, err)

		conditionalProperties, _, _ := GetConditionalProperties(disco)
		for _, api := range conditionalProperties {
			log.Printf("API: %s", api.Name)
			for _, endpoint := range api.Endpoints {
				log.Printf("\tEndpoint: %s:%s", endpoint.Method, endpoint.Path)
				for _, p := range endpoint.ConditionalProperties {
					log.Printf("\t\t%s, %s, %s, %t", p.Schema, p.Name, p.Path, p.Required)
				}
			}
		}
		log.Print("---------------")
		assert.Equal(t, "SecondaryIdentification", conditionalProperties[0].Endpoints[4].ConditionalProperties[1].Name)
		assert.Equal(t, true, conditionalProperties[0].Endpoints[4].ConditionalProperties[1].Request)
	}
}

var testDiscoveryPayload = []byte(`
{
	"discoveryModel": {
	  "name": "ob-v3.1-ozone",
	  "description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data. PSU consent flow.",
	  "discoveryVersion": "v0.4.0",
	  "tokenAcquisition": "psu",
	  "discoveryItems": [{
		"apiSpecification": {
		  "name": "Account and Transaction API Specification",
		  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
		  "version": "v3.1.0",
		  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
		  "manifest": "file://manifests/ob_3.1_accounts_transactions_fca.json"
		},
		"openidConfigurationUri": "https://auth1.obie.uk.ozoneapi.io/.well-known/openid-configuration",
		"resourceBaseUri": "https://rs1.obie.uk.ozoneapi.io/open-banking/v3.1/aisp",
		"endpoints": [
		  {
			"method": "POST",
			"path": "/account-access-consents"
		  },
		  {
			"method": "GET",
			"path": "/account-access-consents/{ConsentId}"
		  },
		  {
			"method": "DELETE",
			"path": "/account-access-consents/{ConsentId}"
		  },
		  {
			"method": "GET",
			"path": "/accounts"
		  },
		  {
			"method": "GET",
			"path": "/accounts/{AccountId}"
		  },
		  {
			"method": "GET",
			"path": "/accounts/{AccountId}/balances"
		  }
		]
	  }, {
		"apiSpecification": {
		  "name": "Payment Initiation API",
		  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1",
		  "version": "v3.1.0",
		  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json",
		  "manifest": "file://manifests/ob_3.1_payment_fca.json"
		},
		"openidConfigurationUri": "https://auth1.obie.uk.ozoneapi.io/.well-known/openid-configuration",
		"resourceBaseUri": "https://rs1.obie.uk.ozoneapi.io/open-banking/v3.1/pisp",
		"endpoints": [{
			"method": "POST",
			"path": "/domestic-payment-consents",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticConsentResponse1",
			  "property": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payment-consents/{ConsentId}",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticConsentResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payment-consents/{ConsentId}/funds-confirmation"
		  },
		  {
			"method": "POST",
			"path": "/domestic-payments",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payments/{DomesticPaymentId}",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "POST",
			"path": "/domestic-scheduled-payment-consents"
		  },
		  {
			"method": "GET",
			"path": "/domestic-scheduled-payment-consents/{ConsentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-scheduled-payments"
		  },
		  {
			"method": "GET",
			"path": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-standing-order-consents",
			"conditionalProperties": [{
                "schema": "OBWriteDataDomesticStandingOrder2",
                "name":"Reference",
                "path": "Data.Initiation.Reference"                
            }, {
                "schema": "OBWriteDataDomesticStandingOrder2",
                "name":"SecondaryIdentification",
                "path": "Data.Initiation.CreditorAccount.SecondaryIdentification"                
            }]			
		  },
		  {
			"method": "GET",
			"path": "/domestic-standing-order-consents/{ConsentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-standing-orders"
		  },
		  {
			"method": "GET",
			"path": "/domestic-standing-orders/{DomesticStandingOrderId}"
		  }
		]
	  }]
  
	}
  }
  
`)

var testDiscoveryPayloadOAS3 = []byte(`
{
	"discoveryModel": {
	  "name": "ob-v3.1-ozone",
	  "description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data. PSU consent flow.",
	  "discoveryVersion": "v0.4.0",
	  "tokenAcquisition": "psu",
	  "discoveryItems": [{
		"apiSpecification": {
		  "name": "Account and Transaction API Specification",
		  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
		  "version": "v3.1.8",
		  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.8/dist/account-info-swagger.json",
		  "manifest": "file://manifests/ob_3.1_accounts_transactions_fca.json"
		},
		"openidConfigurationUri": "https://auth1.obie.uk.ozoneapi.io/.well-known/openid-configuration",
		"resourceBaseUri": "https://rs1.obie.uk.ozoneapi.io/open-banking/v3.1/aisp",
		"endpoints": [
		  {
			"method": "POST",
			"path": "/account-access-consents"
		  },
		  {
			"method": "GET",
			"path": "/account-access-consents/{ConsentId}"
		  },
		  {
			"method": "DELETE",
			"path": "/account-access-consents/{ConsentId}"
		  },
		  {
			"method": "GET",
			"path": "/accounts"
		  },
		  {
			"method": "GET",
			"path": "/accounts/{AccountId}"
		  },
		  {
			"method": "GET",
			"path": "/accounts/{AccountId}/balances"
		  }
		]
	  }, {
		"apiSpecification": {
		  "name": "Payment Initiation API",
		  "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1",
		  "version": "v3.1.8",
		  "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json",
		  "manifest": "file://manifests/ob_3.1_payment_fca.json"
		},
		"openidConfigurationUri": "https://auth1.obie.uk.ozoneapi.io/.well-known/openid-configuration",
		"resourceBaseUri": "https://rs1.obie.uk.ozoneapi.io/open-banking/v3.1/pisp",
		"endpoints": [{
			"method": "POST",
			"path": "/domestic-payment-consents",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticConsentResponse1",
			  "property": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payment-consents/{ConsentId}",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticConsentResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payment-consents/{ConsentId}/funds-confirmation"
		  },
		  {
			"method": "POST",
			"path": "/domestic-payments",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "GET",
			"path": "/domestic-payments/{DomesticPaymentId}",
			"conditionalProperties": [{
			  "schema": "OBWriteDataDomesticResponse1",
			  "name": "Charges",
			  "path": "Data.Charges"
			}]
		  },
		  {
			"method": "POST",
			"path": "/domestic-scheduled-payment-consents"
		  },
		  {
			"method": "GET",
			"path": "/domestic-scheduled-payment-consents/{ConsentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-scheduled-payments"
		  },
		  {
			"method": "GET",
			"path": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-standing-order-consents",
			"conditionalProperties": [{
                "schema": "OBWriteDataDomesticStandingOrder2",
                "name":"Reference",
                "path": "Data.Initiation.Reference"                
            }, {
                "schema": "OBWriteDataDomesticStandingOrder2",
                "name":"SecondaryIdentification",
                "path": "Data.Initiation.CreditorAccount.SecondaryIdentification"                
            }]			
		  },
		  {
			"method": "GET",
			"path": "/domestic-standing-order-consents/{ConsentId}"
		  },
		  {
			"method": "POST",
			"path": "/domestic-standing-orders"
		  },
		  {
			"method": "GET",
			"path": "/domestic-standing-orders/{DomesticStandingOrderId}"
		  }
		]
	  }]
  
	}
  }
  
`)
