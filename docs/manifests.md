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

| Name              | Occurrence | Description                                             | Type             | Value(s)    |
|-------------------|------------|---------------------------------------------------------|------------------|-------------|
| id                | 1..1       | A unique identifier used to identify a test.            | UUID             |             |
| description       | 1..1       | A short description describing the and expected result. | String (max 256) |             |
| refURI            | 0..1       | A URI to identify regulatory or specification.          | String (max 256) |             |
| detail            | 0..1       | Long description describing the and expected result     | String (max 256) |             |
| parameters        | 1..1       | Maps context                                            | json             | see example |
| uri               | 1..1       | A resource to test.                                     | String           |             |
| asserts           | 1..1       | List of linked asserts.                                 | List             |             |
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

[WIP]

## Manifests

Open Banking Implementation Entity (OBIE) has created a number of manifests to help Implementers (Account Providers, Third Party Providers, Vendors and Technical Service Providers) test or provide evidence you have implemented each part of the OBIE Standard correctly. If required these manifests should be used or referenced in your discovery file. 

* Open Banking Implementation Entity Discovery File.
* Open Banking Implementation Manifests.

**Each test is only pickup if a corresponding endpoint is detected in your Discovery.**