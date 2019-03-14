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

| Name        | Occurrence | Description                                             | Type             | Value(s)    |
|-------------|------------|---------------------------------------------------------|------------------|-------------|
| id          | 1..1       | A unique identifier used to identify a test.            | UUID             |             |
| description | 1..1       | A short description describing the and expected result. | String (max 256) |             |
| reference   | 0..1       | A URI to identify regulatory or specification.          | String (max 256) |             |
| parameters  | 1..1       |                                                         | json             | see example |
| resource    | 1..1       | A resource to test.                                     | String           |             |
| asserts     | 1..1       | List of linked asserts.                                 | List             |             |

### Example Test in a Manifest

    {
        "description": "A Test that fails when no token is provided.",
        "id": "OB-301-ACC-019281",
        "parameters": {
                        "accountAccessConsent": "basicAccountAccessConsent",
                        "tokenRequestScope": "accounts",
                        "consentedAccounts": "$consentedAccountIds",
                        "accountId": "$consentedAccountId"
                    },
        "resource": "Account",
        "asserts": ["OB3DOPAssertOnSuccess"]
    }

## Manifest Asserts

Re-usable assertions can be defined as standalone units in a JSON file named `assertions.json`. An assertion can be defined in
the context of a HTTP request in a test case. A HTTP status code, header field(s) and JSON path can be specified as items to be
tested against.

A boilerplate assertions file is structured such that `assertion-ref` is a reference to an assertion:

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
| expect.matches        | 0..1       | Matches on JSON or HTTP header fields                   | Array of JSON    | see example |
| expect.custom         | 0..N       | Reference to an implementation of a custom expectation. Can be defined multiple times.  | String           |             |

_In the following example `assertion-ref` is the reference of the assertion._
    
    "assertion-ref": {
        "expect": {
            "status-code": 201,
            "matches": [{
                    "header": "x-fapi-interaction-id",
                    "detail": "An RFC4122 UID used as a correlation id. The ASPSP 'plays back' the value given. If a value is not given the ASPSP MUST play back their own UUID."
    
                },
                {
                    "JSON": "json.Data.Status",
                    "Value": "AcceptedSettlementCompleted"
                }
            ],
            "custom": "customExpectation"
        }
    }

### Custom expectations

** WIP **

## Manifests

Open Banking Implementation Entity (OBIE) has created a number of manifests to help Implementers (Account Providers, Third Party Providers, Vendors and Technical Service Providers) test or provide evidence you have implemented each part of the OBIE Standard correctly. If required these manifests should be used or referenced in your discovery file. 

* Open Banking Implementation Entity Discovery File.
* Open Banking Implementation Manifests.

**A test is only picked up if a corresponding endpoint is detected in your Discovery.**
