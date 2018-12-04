# Functional Conformance Suite Test case design

The functional conformance suite provides a cross-platform, language agnostic representation of test cases. 
The main purpose is to determine whether an ASPSP resource complies with requirements of specifications and conditions, regulations and standards.

The use case this addresses is described **as an ASPSP, I want Test Case documentation so that I can understand and contribute.**  
_N.B. For the rest of this file functional conformance is referred to as 'conformance'._

```json

      +-------------+           +-------------+           +--------------+
      |             |           |             |           |              |
      |             |  1..n     |             |  1..n     |              |
      |  Manifest   +---------->|   Rule      |+--------->|   TestCase   |
      |             |           |             |           |              |
      |             |           |             |           |              |
      +-------------+           +-------------+           +--------------+
```

## Manifest

A Manifest defines all the rules and test cases to be run that cover the conformance to a discrete set of specifications. 
The Manifest provides details of the specifications covered and links to all tests run as part of the conformance assessment.

Example

```json
{
    "@context": "https://json-ld.org/test-suite/context.jsonld",
    "@id": "",
    "@type": "mf:Manifest",
    "name": "Basic Swagger 2.0 test run",
    "description": "Tests appropriate behaviour of the Open Banking Limited 2.0 Read/Write APIs",
    "baseIri": "https://json-ld.org/test-suite/tests/",
    "rules": [{
```

## Rules

A conformance Manifest consists of one or more rules. Each rule identifies a specific item in a specification which is under test. 
Each rule also contains one or more test cases used to determine adherence to the specification item under consideration.

Example

```json
    "rules": [{
        "@id": "#r0001",
        "@type": ["jld:testsuiteRule"],
        "name": "Get Accounts Basic Rule",
        "description": "Accesses the Accounts endpoint and retreives a list of PSU accounts",
        "specref": "Read Write 2.0 section subsection 1 point 1a",
        "speclocation": "https://openbanking.org.uk/rw2.0spec/errata#1.1a",
        "tests": [{
```

## Test Case

A test case encapsulates a test against a specific endpoint in order to determine an expected outcome. 
The test case defines the expected outcomes using HTTP response status codes and an optional set of javacript and regular expression matching tools.

```json

                                         +--------------+
                                         |              |    Method:
                                         |              |
                        +--------------->|   Input      |    Endpoint:
                        |                |              |
                        |                |              |
     +-----------------+|                +--------------+
     |                 ||                +--------------+
     |                 ||                |              |
     |   Testcase      ||                |              |    Access_token:
     |                 |+--------------->|  Context     |
     |                 ||                |              |    ASPSP Endpoints:
     +-----------------+|                |              |
                        |                +--------------+
                        |                +--------------+
                        |                |              |    Status Code:
                        |                |              |
                        +--------------> |  Expect      |    JSON Field Value:
                                         |              |
                                         |              |    Regular Expression:
                                         +--------------+
```

A test case consists of the follow 3 sections: input, context and expect, used to define the activities of the test case

Example

```json
   {
        "@id": "#t1025",
        "name": "Get transactions related to an account",
        "input": {
            "method": "GET",
            "endpoint": "/accounts/{AccountId}/transactions"
        },
        "context": {
            "accountID":"ABC12300000007"
        },
        "expect": {
            "status-code": 200,
            "schema-validation": true,
            "matches": [{
                    "description": "Example Regex",
                    "regex": "^TheRainInSpainFallsMainlyOnThePlain$"
                }
            ]
        }
    }
```

## Context

Contexts provide pools of information (properties) for test cases. There are 3 main context pools from which a testcase can draw information.

## Global Context

The global context provides access to the Discovery Information provided by the ASPSP in order to define endpoints, certificates, tokens, etc. 
Information likely to be found in the global context includes :-

- ASPSP endpoints
- Access tokens for specific permission mixes
- Transport/Signing certificates
- Conditional endpoint information
  - what's implemented/ what isn't
- Optional endpoint information
  - what's implemented / what isn't

## Section Context

A section context provides access to a collection of common information that is applicable to a certain type of test. 
For example, an access token with a specific permission mix, required by a particular set of test cases across multiple rules.

## Local Context

A rule can have multiple sequentially executed test cases. 
A local context provides a mechanism for passing the results of one test case to the input of the next test case as part of a rule test sequence. 
For example GET /Accounts would provide a list of AccountIDs, the first accountID returned could be passed via the local context as a parameter 
to the next test case which calls GET /Accounts/{accountID}/balances.

Items that would typically be put in the local context include :-

- Account Numbers
- Statement IDs
- DomesticPaymentID
- ConsentID

## Response Matching

In order to check conformance to a particular specification, the primary mechanism available is checking the response data returned 
from an ASPSP API implementation.

There are a number of matching options available via the test case framework in this space :-

- Response status code
- Body
  - Regex matching
  - Json Value Extraction
  - Regex on Json Value
- Header
  - Header value matching
  - Regex on header value

## Declarative Matching

There are four declarative matching options available within test cases:-

- status-code - simple matching of HTTP status codes
- openapi-schema validation - simple validation of API responses against the standard OpenAPI/Swagger schema
- json field value matching - custom JSON field value matching
- regex matching - custom regular expression matching

## Assumptions

- Token acquisition is decoupled
- Tokens and discovery will be provided by via JSON objects
- Conditionality and optionality will be captured in discovery
- Endpoints will be captured in discovery

## Additional Notes

OpenAPI-Schema validation allows enabling and disabling validation for specific test cases. 
It also needs a swagger validation off mode, where swagger json validation checks are not performed, 
which allows for defining custom json response validation.

## Auto-generated Test

A number of test cases will be auto generated for the OpenAPI/Swagger json files. 
These test cases will have the OpenAPI-Schema validation feature enabled by default. 
For specific custom test cases, its possible to disable default schema-validation and provide a custom regex to perform the desired checks.

## Key Features

- Checking JSON variable values against a regex
- Regex header checking
- Regex body checking
- Json body field extraction, then application of regex
- Extensibility - create your own tests
- Using a result field from one test as an input to the next via the local context

### Example Test Case

The following example JSON show an example test case which demonstrates the following:-

```json
{
    "@context": "https://json-ld.org/test-suite/context.jsonld",
    "@id": "",
    "@type": "mf:Manifest",
    "name": "Basic Swagger 2.0 test run",
    "description": "Tests appropriate behaviour of the Open Banking Limited 2.0 Read/Write APIs",
    "baseIri": "https://json-ld.org/test-suite/tests/",
    "rules": [{
        "@id": "#r0001",
        "@type": ["jld:testsuiteRule"],
        "name": "Get Accounts Basic Rule",
        "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
        "specref": "Read Write 2.0 section subsection 1 point 1a",
        "speclocation": "https://openbanking.org.uk/rw2.0spec/errata#1.1a",
        "tests": [{
                "@id": "#t0001",
                "@type": ["jld:PositiveEvaluationTest"],
                "name": "Get Accounts Basic - Positive",
                "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
                "input": {
                    "method": "GET",
                    "endpoint": "/accounts/"
                },
                "context": {
                    "token": {
                        "@type": ["jld:BasicAccountsToken"]
                    },
                    "account": "XLY12300010202"
                },
                "expect": {
                    "status-code": 200,
                    "matches": [{
                        "description": "A json match on response body",
                        "json": "Data.Account.Accountid",
                        "value": "XYZ1231231231231"
                    }, {
                        "description": "a regex match - on response body",
                        "regex": ".*"
                    }, {
                        "description": "a header match - using context reference",
                        "header": "x-fapi-id",
                        "value": "@ref:ctx-aspspid"
                    }, {
                        "description": "a header regex match",
                        "header": "x-fapi-id",
                        "regex": "***"
                    }]
                }
            },
            {
                "@id": "#t0002",
                "@type": ["jld:NegativeEvaluationTest"],
                "name": "Get Accounts Basic - Negative Test",
                "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
                "input": {
                    "method": "POST",
                    "endpoint": "/accounts/"
                },
                "context": {},
                "expect": {
                    "status-code": 201
                }
            },
            {
                "@id": "#t0003",
                "@type": ["jld:WarningEvaluationTest"],
                "name": "Get Accounts Basic - WarningTest",
                "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
                "input": {
                    "method": "POST",
                    "endpoint": "/accounts/"
                },
                "expect": {
                    "status-code": 201
                }
            }
        ]
    }, {
        "@id": "#r0002",
        "name": "Rule with shared context across tests",
        "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
        "specref": "Read Write 2.0 section subsection 2 point 2",
        "speclocation": "https://openbanking.org.uk/rw2.0spec/errata#2.2a",
        "tests": [{
            "@id": "#t0033",
            "@type": ["jld:PositiveEvaluationTest"],
            "name": "Simple SwaggerSchema check",
            "purpose": "Simple test to check basic swagger schema validation for call",
            "input": {
                "method": "GET",
                "endpoint": "/accounts"
            },
            "context": {},
            "expects": {
                "status-code": 200,
                "schema-validation": true
            }
        }]
    }, {
        "@id": "#r0003",
        "name": "Rule to capture OBSD-5408 Data/StandingOrder/Frequency regex",
        "description": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
        "specref": "Read Write 2.0 section subsection 3 point 4a",
        "speclocation": "https://openbanking.org.uk/rw2.0spec/errata#3.4a",
        "tests": [{
            "@id": "#t0099",
            "@type": ["jld:PositiveEvaluationTest"],
            "name": "Get Standing Order Frequency Field Validation",
            "purpose": "Validates the Standing order frequency field using complex regex",
            "input": {
                "method": "GET",
                "endpoint": "/accounts/{AccountId}/standing-orders"
            },
            "context": {
                "token": {
                    "@type": ["jld:BasicAccountsToken"]
                },
                "AccountId": "XLY12300010202"
            },
            "expect": {
                "matches": [{
                    "description": "Standing order frequency validated against regex",
                    "json-selector": "Data.StandingOrder.Frequency",
                    "regex": "^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$"
                }]
            }
        }]
    }]
}

```