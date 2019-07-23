# Manifest Specification - v0.1

**Warning**: The Manifest specification is in alpha and is subject to change without notice.

** Contents **

[TOC]

## Overview

This document is a specification (alpha) for describing a Manifest as a JSON document.  It is intended to be used by the Functional Conformance Suite with the goal of:

* Providing a standardised and structured approach to the presentation of a Manifest data used by the suite in the generation of test cases.
* Facilitating a cross-platform format that allows 3rd parties to view and contribute easily to tests.

## Manifest Dictionary

TABLE A:

| Name              | Occurrence | Description                                             | Type             | Value(s)    |
|-------------------|------------|---------------------------------------------------------|------------------|-------------|
| id                | 1..1       | A unique identifier used to identify a test.            | UUID             |             |
| description       | 1..1       | A short description describing the and expected result. | String (max 256) |             |
| refURI            | 0..1       | A URI to identify regulatory or specification.          | String (max 256) |             |
| detail            | 0..1       | Long description describing the and expected result     | String (max 256) |             |
| parameters        | 1..1       | Maps context                                            | json             | see example |
| uri               | 1..1       | A resource to test.                                     | String           |             |
| asserts           | 1..1       | List of linked asserts all of which must be true.       | List             |             |
| asserts_one_of    | 0..1       | List of linked asserts one of which must be true.       | List             |             |
| uriImplementation | 1..1       |                                                         |                  |             |
| resource          | 1..1       |                                                         |                  |             |
| keepContext       | 1..1       |                                                         |                  |             |
| method            | 1..1       |                                                         |                  |             |
| schemaCheck       | 1..1       |                                                         |                  |             |
| headers           | 0..1       |                                                         |                  |             |
| body              | 0..1       |                                                         |                  |             |

### Example Test in a Manifest

        {
			"description": "Domestic Payment consents succeeds with minimal data set with additional schema checks.",
            "id": "OB-301-DOP-206111",
            "refURI": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937984109/Domestic+Payments+v3.1#DomesticPaymentsv3.1-POST/domestic-payment-consents",
            "detail" : "Checks that the resource succeeds posting a domestic payment consents with a minimal data set and checks additional schema.",
			"parameters": {
                "tokenRequestScope": "payments",
                "paymentType": "domestic-payment-consents",
                "post" : "minimalDomesticPaymentConsent"    
            },
            "uri": "domestic-payment-consents",
            "uriImplementation": "mandatory",
            "resource": "DomesticPayment",
            "asserts": ["OB3DOPAssertOnSuccess", "OB3GLOAAssertConsentId"],
            "keepContext": ["OB3GLOAAssertConsentId"],
            "method":"post",
            "schemaCheck": true
        },

## Manifest Asserts

Re-usable assertions can be defined as standalone units in a JSON file named `assertions.json`. An assertion can be defined in
the context of a HTTP request in a test case. A HTTP status code, header field(s) and JSON path can be specified as items to be
tested against, along with a custom implementation of an "expectation", designated with the label `custom`.

A boilerplate assertions file is structured such that `assertion-ref` is a reference to an assertion, many references can be defined
alongside each other as long as the reference is unique:

    {
        "references": {
            "assertion-ref": {}
        }
    }

TABLE B - Structure of an assertion:

| Name                  | Occurrence | Description                                             | Type             | Value(s)    |
|-----------------------|------------|---------------------------------------------------------|------------------|-------------|
| expect                | 1..1       | Container for test expectations                         | JSON             |             |
| expect.status-code    | 0..1       | Expected HTTP status code                               | Integer          |             |
| expect.matches        | 0..N       | Array of "MatchType" checks                             | Array of JSON    | see example |
| expect.custom         | 0..N       | Reference to an implementation of a custom expectation. Can be defined multiple times.                          | String           |             |



_In the following example `assertion-ref` is the reference of the assertion._
    
    "assertion-ref": {
        "expect": {
            "status-code": 201,
            "matches": [{
                    "header": "x-fapi-interaction-id",
                    "detail": "An RFC4122 UID used as a correlation id. The ASPSP 'plays back' the value given. If a value is not given the ASPSP MUST play back their own UUID.",
                    "value": "ba70ddff-b9ea-43a2-8ade-567746f39fff"
                },
                {
                    "JSON": "Data.Status",
                    "detail: "Status label for request",
                    "Value": "AcceptedSettlementCompleted"
                }
            ],
            "custom": "bodyNotEmpty"
        }
    }

### Custom expectations

** WIP **

## Custom Data

** WIP ** 

Custom data can be defined, which can then be consumed by test cases. This data is stored in a JSON file name `data.json`.
The items are stored as referencable JSON objects similar to that of the asserts - this means we would start off with the same boilerplate
as assertions, seen previously.

The schema for data items is currently "free form" and specific to each test. With this in mind,
it would be useful to examine any associated notes for the each test.

## Manifest Functions

Manifests have the ability to call a function which is mapped to a Go function in `pkg/model/macro.go`. An example function is shown below, which generates a unique identifier, to be used in the
`instructionIdentication` parameter in some payment tests.

Manifest Functions also supports any number of parameters, which are passed and parsed as strings. It is worth noting that all function parameters will passed to the implementation as strings. If other types
are required, the specific function implementation is required to perform type assertions and casting of types.

Function implementations must return one value, of type `string`

Register and implement function in `pkg/model/macro.go`
```
var macroMap = map[string]interface{}{
	"instructionIdentificationID": instructionIdentificationID,
}
```

```
func instructionIdentificationID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
```

In manifest file, call the fuction. Note the pattern required here. The function should be followed by parentheses (`()`).

```
"parameters": {
        "tokenRequestScope": "payments",
        "instructedAmountValue": "$instructedAmountValue",
        "instructedAmountCurrency": "$instructedAmountCurrency",
        "currencyOfTransfer": "$currencyOfTransfer",
        "instructionIdentification": "$fn:instructionIdentificationID()",
        "endToEndIdentification": "e2e-internat-sched-pay",
        "postData": "$minimalInternationalScheduledPayment",
        "consentId": "$OB-301-DOP-102000-ConsentId"
      },
```

## Supplementary Manifests

Open Banking Implementation Entity (OBIE) has created a number of manifests to help Implementers (Account Providers, Third Party Providers, Vendors and Technical Service Providers) test or provide evidence you have implemented each part of the OBIE Standard correctly. If required these manifests should be used or referenced in your discovery file. 

* Open Banking Implementation Entity Discovery File.
* Open Banking Implementation Manifests.

**A test is only picked up if a corresponding endpoint is detected in your Discovery.**
