# Functional Conformance Suite Discovery design

The Functional Conformance Suite provides a configurable discovery model that
allows an ASPSP to describe information on endpoint availability, and data
schema properties provided.

The suite uses this discovery information to configure which tests cases are run.

## Discovery Model

The discovery model information per specification includes:

* endpoint and methods implemented
* optional/conditional properties provided for online channel equivalence

The main purpose is to determine whether an ASPSP resource complies with
requirements of specifications and conditions, regulations, and standards.

### Model format

The discovery model defines in a JSON format endpoints implemented per
specification and optional payload schema properties provided for online channel
equivalence.

The discovery model consists of an `discoveryModel` root object with these
properties:

* `version` - version number of the discovery model format, e.g. "v0.0.1".
* `discoveryItems` - an array of discovery items, see below for details.

#### Discovery version

The version number is used to track changes to made to the discovery model.

The version number is formatted as MAJOR.MINOR.PATCH, following the
[Semantic Versioning](https://semver.org/) approach to convey meaning about what
has been modified from one version to the next. For details see: https://semver.org/

#### Discovery item

Each discovery item contains information related to a particular specification
version.

Properties in each discovery item are:

* `apiSpecification` - details of API specification
* `openidConfigurationUri` - URI of openid configuration well-known endpoint
* `resourceBaseUri` - Base of resource URI, i.e. the part before "/open-banking/v3.0".
* `endpoints` - Array of endpoint and method implementation details.

#### API Specification

The Discovery model records specification details in an unambiguous way:

* `apiSpecification`
  * `name` - the `info.title` field from the Swagger/OpenAPI specification file
  * `url` - URI identifier of the specification, i.e. link to specification document
  * `version` - API version number that appears in API paths, e.g. "v3.0"
  * `schemaVersion` - URI identifier of the Swagger/OpenAPI specification file patch version

The property names `url`, `version`, and `schemaVersion` are from the schema.org
[APIReference schema](https://schema.org/APIReference) defined here:
https://schema.org/APIReference

Example

```json
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
        "endpoints": []
      }
    ]
  }
}
```

#### Endpoint and method

A discovery item contains a list of endpoint and methods that have been
implemented by an ASPSP. This list includes:
  * all mandatory endpoints
  * conditional and optional endpoints implemented

Properties in each endpoint definition include (mandatory properties marked with *):
  - `method`* - HTTP method, e.g. "GET" or "POST"
  - `path`* - endpoint path, e.g. "/account-access-consents"
  - `conditionalProperties` - list of optional schema properties that an ASPSP attests it provides, more details in the next section.

Example

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

##### Required properties

The specification lists some resource schema properties may occur `0..1`, or `0..n` times.

When an ASPSP provides a `0..1`, `0..n` occurrence property via its online channel,
it must attest that it provides those properties in its API implementation. It adds
such properties to a `conditionalProperties` properties list in the relevant endpoint definition.

The `conditionalProperties` list contains items. Each item states:
 * `schema` - schema definition name from the Swagger/OpenAPI specification, e.g. "OBTransaction3Detail"
 * `property` - property name from schema, e.g. "Balance"
 * `path` - path to property expressed in [JSON Path](https://goessner.net/articles/JsonPath/) format, e.g. "Data.Transaction[* ].Balance"

Example: for online channel equivalence an ASPSP provides account
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
        "path": "Data.Transaction[*].Balance"
      },
      {
        "schema": "OBTransaction3Detail",
        "property": "MerchantDetails",
        "path": "Data.Transaction[*].MerchantDetails"
      },
      {
        "schema": "OBTransaction3Basic",
        "property": "TransactionReference",
        "path": "Data.Transaction[*].TransactionReference"
      },
      {
        "schema": "OBTransaction3Detail",
        "property": "TransactionReference",
        "path": "Data.Transaction[*].TransactionReference"
      }
    ]
  },
  ...
]
```

### Example file

See ./docs/discovery-example.json for a [longer example file](./discovery-example.json).
Note, this file is a nonnormative incomplete example of a discovery model.
