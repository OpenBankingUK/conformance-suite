package manifest

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"encoding/json"
	"testing"
)

func TestDiscoveryEndpointsMapToManifestCorrectly(t *testing.T) {
	discoJSON := `
{
  "discoveryModel": {
    "name": "ob-v3.1-ozone",
    "description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data.",
    "discoveryVersion": "v0.3.0",
    "tokenAcquisition": "psu",
    "discoveryItems": [{
        "apiSpecification": {
          "name": "Account and Transaction API Specification",
          "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
          "version": "v3.1",
          "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
          "manifest": "file://manifests/ob_3.1__accounts_fca.json"
        },
        "openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
        "resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/aisp",
        "resourceIds": {
			 "AccountId": "500000000000000000000001",
     	     "ConsentId": "$consent_id",
			 "StatementId":"140000000000000000000001"},
        "endpoints": [
          {
            "method": "HEAD",
            "path": "/accounts/{AccountId}"
          },
          {
            "method": "GET",
            "path": "/accounts/{AccountId}"
          }
        ]
      }],
    "customTests": [
      {}
    ]
  }
}
`
	mfJSON := `
{
	"scripts": [
        {
			"description": "Minimal data returned for a given account using the ReadAccountsBasic permission, status and headers.",
            "id": "OB-301-ACC-120382",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource differs depending on the permissions (ReadAccountsBasic and ReadAccountsDetail) used to access the resource with additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "basicAccountAccessConsent",
				"tokenRequestScope": "accounts",
                "accountId": "$consentedAccountId"         
            },
            "uri": "accounts/$accountId",
            "uriImplementation": "mandatory",
            "resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"get",
            "schemaCheck": true
        },
        {
			"description": "All data returned for a given account with ReadAccountsDetail permission, status and headers.",
            "id": "OB-301-ACC-352203",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource returns the correct data depending on the permissions ReadAccountsDetail with additional additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "detailAccountAccessConsent",
				"tokenRequestScope": "accounts",
				"accountId": "$consentedAccountId"
            },
            "uri": "accounts/$accountId",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"head",
            "schemaCheck": true
        }
	]
}
`
	require := test.NewRequire(t)

	var mf Scripts
	err := json.Unmarshal([]byte(mfJSON), &mf)
	require.Nil(err)

	disco, err := discovery.UnmarshalDiscoveryJSON(discoJSON)
	require.Nil(err)

	mpParams := map[string]string{
		"$AccountID":"500000000000000000000004",
	}

	mpResults := MapDiscoveryEndpointsToManifestTestIDs(disco, mf, mpParams)


	exp := DiscoveryPathsTestIDs {
		"/accounts/500000000000000000000004": {
			"GET": {"OB-301-ACC-120382"},
			"HEAD": {"OB-301-ACC-352203"},
		},
	}

	require.Equal(exp, mpResults)
}

func TestUnMappedManifestItemsReportedCorrectly(t *testing.T) {
	discoJSON := `
{
  "discoveryModel": {
    "name": "ob-v3.1-ozone",
    "description": "An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data.",
    "discoveryVersion": "v0.3.0",
    "tokenAcquisition": "psu",
    "discoveryItems": [{
        "apiSpecification": {
          "name": "Account and Transaction API Specification",
          "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1",
          "version": "v3.1",
          "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
          "manifest": "file://manifests/ob_3.1__accounts_fca.json"
        },
        "openidConfigurationUri": "https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration",
        "resourceBaseUri": "https://modelobank2018.o3bank.co.uk:4501/open-banking/v3.1/aisp",
        "resourceIds": {
			 "AccountId": "500000000000000000000001",
     	     "ConsentId": "$consent_id",
			 "StatementId":"140000000000000000000001"},
        "endpoints": [
          {
            "method": "HEAD",
            "path": "/accounts/{AccountId}"
          },
          {
            "method": "GET",
            "path": "/accounts/{AccountId}"
          }
        ]
      }],
    "customTests": [
      {}
    ]
  }
}
`
	mfJSON := `
{
	"scripts": [
        {
			"description": "Minimal data returned for a given account using the ReadAccountsBasic permission, status and headers.",
            "id": "OB-301-ACC-120382",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource differs depending on the permissions (ReadAccountsBasic and ReadAccountsDetail) used to access the resource with additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "basicAccountAccessConsent",
				"tokenRequestScope": "accounts",
                "accountId": "$consentedAccountId"         
            },
            "uri": "accounts/$accountId",
            "uriImplementation": "mandatory",
            "resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"get",
            "schemaCheck": true
        },
        {
			"description": "All data returned for a given account with ReadAccountsDetail permission, status and headers.",
            "id": "OB-301-ACC-352203",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937623627/Accounts+v3.1#Accountsv3.1-PermissionCodes",
            "detail" : "Checks that the resource returns the correct data depending on the permissions ReadAccountsDetail with additional additional schema checks on status and headers.",
			"parameters": {
				"accountAccessConsent": "detailAccountAccessConsent",
				"tokenRequestScope": "accounts",
				"accountId": "$consentedAccountId"
            },
            "uri": "accounts/$accountId",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": ["OB3ACCAssertOnSuccess"],
            "method":"head",
            "schemaCheck": true
        },
		{
			"description": "",
            "id": "unmapped-test-id",
            "refURI": "",
            "detail" : "",
			"parameters": {},
            "uri": "FOO-BAR",
            "uriImplementation": "mandatory",
			"resource": "Account",
            "asserts": [],
            "method":"head",
            "schemaCheck": true
        }
	]
}
`
	require := test.NewRequire(t)

	var mf Scripts
	err := json.Unmarshal([]byte(mfJSON), &mf)
	require.Nil(err)

	disco, err := discovery.UnmarshalDiscoveryJSON(discoJSON)
	require.Nil(err)

	mpParams := map[string]string{
		"$AccountID":"500000000000000000000004",
	}

	mpResults := MapDiscoveryEndpointsToManifestTestIDs(disco, mf, mpParams)

	unmatched := FindUnmatchedManifestTests(mf, mpResults)

	exp := []string { "unmapped-test-id" }

	require.Equal(exp, unmatched)
}
