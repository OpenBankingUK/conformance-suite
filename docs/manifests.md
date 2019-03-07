# Manifest Specification - v0.1

**Warning**: The Manifest specification is in alpha and is subject to change without notice.

** Contents **

[TOC]

## Overview

This document is a specification (alpha) for describing a Manifest as a JSON document.  It is intended to be used by the Functional Conformance Suite with the goal of:

* Providing a standardise and structured approach to the presentation of a Manifest data used by the suite in the generation of test cases.
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
        "name": "A Test that fails when no token is provided.",
        "id": "OB-301-ACC-019281",
        "parameters": {
                        "accountAccessConsent": "basicAccountAccessConsent",
                        "tokenRequestScope": "accounts",
                        "consentedAccounts": "$consentedAccountIds",
                        "accountId": "$consentedAccountId"
                    },
        "resource": "Account",
        "asserts": ["OB3assertNoToken"]
    }

## Manifest Asserts

[WIP]

## Manifests

Open Banking Implementation Entity (OBIE) has created a number of manifests to help Implementers (Account Providers, Third Party Providers, Vendors and Technical Service Providers) test or provide evidence you have implemented each part of the OBIE Standard correctly. If required these manifests should be used or referenced in your discovery file. 

* Open Banking Implementation Entity Discovery File.
* Open Banking Implementation Manifests.