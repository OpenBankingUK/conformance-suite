# Reporting Specification - v0.1

**Warning**: The Reporting Specification is in alpha and is subject to change without notice.

** Contents **

[TOC]

## Overview

This document is a specification (alpha) for parsing and processing test case results as a JSON document.
It is intended to be used by the Functional Conformance Suite with the goal of:

* Providing a standardise and structured approach to the presentation of data throughout the suite.
* Providing a safe and alternative format to export data over conventional log files.
* Facilitating the importing and re-run tests.
* Facilitating a cross-platform format to validate results with other 3rd parties.

## UML

TBD

## Report Dictionary

### `Report`

| Name           | Occurrence | Description                                                    | Class                  | Example                                | Value(s)                                                                      | Notes                                                                       |
|----------------|------------|----------------------------------------------------------------|------------------------|----------------------------------------|-------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| id             | 1..1       | A unique and immutable identifier used to identify the report. | v4 UUID                | `f47ac10b-58cc-4372-8567-0e02b2c3d479` | Regex `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$` | The v4 UUIDs generated conform to RFC 4122                                  |
| created        | 1..1       | Date and time when the report was created.                     | timestamp              | `2006-01-02T15:04:05Z07:00`            | Formatted accorrding to RFC3339 (<https://tools.ietf.org/html/rfc3339>)       | RFC3339 is derived from ISO 8601 (<https://en.wikipedia.org/wiki/ISO_8601>) |
| expiration     | 0..1       | Date and time when the report should not longer be accepted.   | timestamp              | `2006-01-02T15:04:05Z07:00`            | Formatted accorrding to RFC3339 (<https://tools.ietf.org/html/rfc3339>)       | RFC3339 is derived from ISO 8601 (<https://en.wikipedia.org/wiki/ISO_8601>) |
| version        | 1..1       | The current version of the report model used.                  | string                 |                                        |                                                                               |                                                                             |
| status         | 1..1       | A status describing overall condition of the report.           | string(8)              | `Complete`                             | One of [`Pending`, `Complete`, `Error`]                                       |                                                                             |
| signatureChain | 0..1       | TBD                                                            | `SignatureChain`       |                                        |                                                                               |                                                                             |
| certifiedBy    | 1..1       | The certifier of the report.                                   | `CertifiedBy`          |                                        |                                                                               |                                                                             |
| apiSpecification|0..n       | The name of API being specified, version and tests that were run.| Array of `APISpecification`   | See class definition.                  |                                                                               |                                                                             |

### `CertifiedBy`

| Name         | Occurrence | Description                     | Class      | Value(s)                         |
|--------------|------------|---------------------------------|------------|----------------------------------|
| environment  | 1..1       | Name of the environment tested. | string(60) | One of [`testing`, `production`] |
| brand        | 1..1       | Name of the brand tested.       | string(60) |                                  |
| authorisedBy | 1..1       | Full name of the authoriser.    | string(60) |                                  |
| jobTitle     | 1..1       | Job title of the authoriser.    | string(60) |                                  |

### `SignatureChain`

TDB

## `Report` Example

TDB

## `APISpecification`

| Name      | Occurrence | Description          | Class     | Value(s)                          |
|-----------|------------|----------------------|-----------|-----------------------------------|
| name      | 1..1       | Name of the API      | string    |
| version   | 1..1       | Version of the API   | string    |
| results   | 0..n       | Results of tests     | string    |

### `APISpecification.Result`

| Name      | Occurrence | Description          | Class     | Value(s)                          |
|-----------|------------|----------------------|-----------|-----------------------------------|
| id        | 1..1       | Test case ID         | string    ||
| pass      | 1..1       | Test passed (true/false) | boolean ||
| metrics   | 0..n       | Metrics (response time/size) | `Metrics` | See example |
| endpoint  | 1..1       | Endpoint under test | string | ||

### Example

```json
    {
     "name": "Account and Transaction API Specification",
     "version": "v3.1",
     "results": [
       {
         "id": "OB-301-ACC-001000",
         "pass": true,
         "metrics": {
           "response_time": 17.000526,
           "response_size": 168
         },
         "endpoint": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/foobar"
       }
     ]
    }
```
