# Discovery Specification - v0.4.0

**Warning**: The Discovery Specification is in alpha and is subject to change without notice.

** Contents **

[TOC]

## Overview

The Functional Conformance Suite generates a set of tests that can be used to verify conformance to a set of specifications. To facilitate this, an implementer must describe their system using a Discovery file. This discovery file allows the Functional Conformance Suite to generate and run 'tests cases' to ensure conformance to a specification by running test cases and analysing the results.

Currently, the suite supports the following standards:

* [Open Banking UK](https://www.openbanking.org.uk/customers/what-is-open-banking/)- Read/Write Data API Specifications v3.0/3.1 (alpha)

## Discovery Templates

The Functional Conformance Suite provides several discovery templates that can be used to help implementers describe a system.

The following discovery templates are available:

* [Open Banking](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937656404/Read+Write+Data+API+Specification+-+v3.1) - Read/Write Data API Specifications v3.0/v3.1 templates:
 * Generic - a customizable template for implementers of the Open Banking v3.0/v3.1 to describe their API endpoints.
 * Ozone -  a customizable template that is pre-populated with Ozone endpoints and data.

## Data Model

The discovery data model defines in a JSON format the endpoints implemented per
specification and optional payload schema properties provided for online channel
equivalence.

### Data dictionary

Name             | Occurrence | Path | Description
-----------------|:----------:|------|----------------
discoveryModel   | 1..1       | discoveryModel |
name             | 1..1       | discoveryModel.name | Name of the model, e.g. "ob-v3.0-ozone".
description      | 1..1       | discoveryModel.description | Description of the model, e.g. "An Open Banking UK discovery template for v3.0 of Accounts and Payments with pre-populated model Bank (Ozone) data."
discoveryVersion | 1..1       | discoveryModel.discoveryVersion | Version of the discovery model format, e.g. "v0.4.0"
tokenAcquisition | 1..1       | discoveryModel.tokenAcquisition | Define how access tokens will be acquired, e.g. "headless", "psu", "store"
discoveryItems   | 1..n       | discoveryModel.discoveryItems.* | List of items. Each item contains information related to a particular specification version.
apiSpecification | 1..1       | discoveryModel.discoveryItems.*.apiSpecification | Details of API specification
name             | 1..1       | discoveryModel.discoveryItems.*.apiSpecification.name | The `info.title` field from the Swagger/OpenAPI specification file
url              | 1..1       | discoveryModel.discoveryItems.*.apiSpecification.url | URI identifier of the specification, i.e. link to specification document
version          | 1..1       | discoveryModel.discoveryItems.*.apiSpecification.version | API version number that appears in API paths, e.g. "v3.0"
schemaVersion    | 1..1       | discoveryModel.discoveryItems.*.apiSpecification.schemaVersion | URI identifier of the Swagger/OpenAPI specification file patch version
manifest         | 1..1       | discoveryModel.discoveryItems.*.apiSpecification.manifest | Path to manifest file for custom tests. Can be `http://` or `https://` or `file://`.
openidConfigurationUri | 1..1 | discoveryModel.discoveryItems.*.openidConfigurationUri | URI of the openid configuration well-known endpoint
resourceBaseUri  | 1..1       | discoveryModel.discoveryItems.*.resourceBaseUri | Base of resource URI, i.e. the part before "/open-banking/v3.0".
endpoints        | 1..n       | discoveryModel.discoveryItems.*.endpoints | List of endpoint and methods that have been implemented.
method           | 1..1       | discoveryModel.discoveryItems.\*.endpoints.\*.method | HTTP method, e.g. "GET" or "POST"
path             | 1..1       | discoveryModel.discoveryItems.\*.endpoints.\*.path | Endpoint path, e.g. "/account-access-consents"
conditionalProperties | 0..n  | discoveryModel.discoveryItems.\*.endpoints.\*.conditionalProperties | List of optional schema properties that an ASPSP attests it provides.
schema           | 1..1       | discoveryModel.discoveryItems.\*.endpoints.\*.conditionalProperties.*.schema | Schema definition name from the Swagger/OpenAPI specification, e.g. "OBTransaction3Detail"
property         | 1..1       | discoveryModel.discoveryItems.\*.endpoints.\*.conditionalProperties.*.property | Property name from schema, e.g. "Balance"
path             | 1..1       | discoveryModel.discoveryItems.\*.endpoints.\*.conditionalProperties.*.path | Path to property expressed in [JSON dot notation](https://github.com/tidwall/gjson#path-syntax) format, e.g. Data.Transaction.*.Balance

### Discovery version

The version number is used to track changes to made to the discovery model format.

The version number is formatted as MAJOR.MINOR.PATCH, following the
[Semantic Versioning](https://semver.org/) approach to convey meaning about what
has been modified from one version to the next. For details see: https://semver.org/

### Token Acquisition

The Token Acquisition field informs the application how it shall acquire access tokens to access endpoints under test.
The following values are valid `psu`, `headless`, `store`

* `psu` - tokens are acquired following the "[Hybrid Flow](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937656404/Read+Write+Data+API+Specification+-+v3.1#Read/WriteDataAPISpecification-v3.1-GrantTypesforidentifyingtheTPPandPSU)" authorisation method.
This process involves directing the PSU to the ASPSP's authorisation pages, requiring manual effort from the PSU.
* `headless` - Similar to `psu`, except the PSU is not required to intervene and perform any actions. This mode of operation enables developers to integrate
the operation of this suite into their build tooling e.g. continuous integration/deployment (CI/CD), thus removing the manual element from `psu`.
* `store` - As a final step to to the `psu` and `headless` methods, an access token is generated and used to access the protected endpoints.
The access tokens for use in this method would typically be generated in a developer/application management portal hosted by the ASPS.

### Discovery item

Each discovery item contains information related to a particular specification
version.


### API Specification

The discovery model records specification details in an unambiguous way.

The `apiSpecification` property names `name`, `url`, `version`, and `schemaVersion` are from the schema.org
[APIReference schema](https://schema.org/APIReference)

The `manifest` property is proprietary and contains a reference to a manifest file for custom tests. This location can
be a `https://` endpoint or a relative location on the local filesystem defined using the `file://` scheme.

Non-normative example

```json
{
  "discoveryModel": {
    "discoveryVersion": "v0.4.0",
    "discoveryItems": [
      {
        "apiSpecification": {
          "name": "Account and Transaction API Specification",
          "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
          "version": "v3.0",
          "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json",
          "manifest": "file://../pkg/discovery/templates/ob_3.1_accounts_fca.json"
        },
        "openidConfigurationUri": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
        "resourceBaseUri": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp",
        "endpoints": []
      }
    ]
  }
}
```

### Endpoint and method

A discovery item contains a list of endpoint and methods that have been
implemented by an ASPSP. This list includes:

  * all mandatory endpoints
  * conditional and optional endpoints implemented

Non-normative example fragment

```json
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
    "path": "/accounts/{AccountId}/product"
  },
  ...
]
```

#### Required properties

The API specifications define that some resource schema properties that may occur `0..1`, or `0..*` times.

When an ASPSP provides a `0..1`, `0..*` occurrence property via its online channel,
it must attest that it provides those properties in its API implementation. An ASPSP must add
such properties to a `conditionalProperties` list in the relevant endpoint definition.

Non-normative example:

For online channel equivalence an ASPSP provides account
transaction data including `Balance`, `MerchantDetails`, `TransactionReference`.
The ASPSP attests to that in an endpoint definition, via a `conditionalProperties` list
as follows:

```json
"endpoints": [
  {
    "method": "GET",
    "path": "/accounts/{AccountId}/transactions",
    "conditionalProperties": [
      {
        "schema": "OBTransaction3Detail",
        "property": "Balance",
        "path": "Data.Transaction.*.Balance"
      },
      {
        "schema": "OBTransaction3Detail",
        "property": "MerchantDetails",
        "path": "Data.Transaction.*.MerchantDetails"
      },
      {
        "schema": "OBTransaction3Basic",
        "property": "TransactionReference",
        "path": "Data.Transaction.*.TransactionReference"
      },
      {
        "schema": "OBTransaction3Detail",
        "property": "TransactionReference",
        "path": "Data.Transaction.*.TransactionReference"
      }
    ]
  },
  ...
]
```

## Resource IDs

We've introduced a "resourceId" section to the discovery model which allows a tester to provide resource ids to be used when swagger/openapi calls are made.

If we use Accounts as an example; to get a list of accounts the following endpoint is called:-

```json
GET /accounts
```

Once we have the list of accounts we can query an individual account using the following swagger definition

```json
GET /accounts/{AccountId}
```

Notice that the swagger definition includes the term `{AccountId}` in the endpoint path.

The purpose of the Discovery resourceId section is to allow the replacement of resourceIds like `{AccountId}` with a resource ID that exists in the system under test.

For example, if the account id `5000000001` exists in the system under test, and was to be used in the account detail call, the following `resourceId` section would be added next to the endpoint definitions in the discovery file.

```json
  "resourceIds":  {
              "AccountId": "5000000001"
          },
  "endpoints":[
            {
              "method": "GET",
              "path": "/accounts"
            },
            {
              "method": "GET",
              "path": "/accounts/{AccountId}"
            }
  ]

```

The above json Discovery file fragment would result in the following endpoint path being used in the generated test case:

```json
GET /accounts/5000000001
```

So in summary, using the `resouceIds` section of the discovery file is one way of specificying which resource id to use for detailed calls when a list of candidate resource calls are available.



## Example file

Discovery templates can be found in the [templates directory here](../pkg/discovery/templates).
