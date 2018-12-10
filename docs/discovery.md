# Discovery - v0.1

The Functional Conformance Suite generates a set of tests that can be used to verify conformance to a set of specifications. To facilitate this, an implementer must describe their system using a Discovery file. This discovery file allows the Functional Conformance Suite to generate and run 'tests cases' to ensure conformance to a specification by running test cases and analysing the results.

Currently, the suite supports the following standards:

* [Open Banking UK](https://www.openbanking.org.uk/customers/what-is-open-banking/)- Read/Write Data API Specifications v3.0/3.1 (alpha)

## Discovery Templates

The Functional Conformance Suite provides several discovery templates that can be used to help implementers describe a system.

The following discovery templates are available:

* [Open Banking](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937656404/Read+Write+Data+API+Specification+-+v3.1) - Read/Write Data API Specifications v3.0/3.1 templates:
* * Generic - a customizable template for implementers of the Open Banking v3.0/v3 to describe their API endpoints.
* * Ozone -  a customizable template that is pre-populated with Ozone endpoints and data.
* * ForgeRock -  a customizable template that is pre-populated with ForgeRock endpoints and data.

### Model format

The discovery model defines in a JSON format endpoints implemented per
specification and optional payload schema properties provided for online channel
equivalence.

The discovery model consists of a `discoveryModel` root object with these
properties:

* `name` - the name of the model, e.g. "ob-v3.0-ozone".
* `description` - the description of the model, e.g. "An Open Banking UK discovery template for v3.0 of Accounts and Payments with pre-populated model Bank (Ozone) data.".
* `discoveryVersion` - version number of the discovery model format, e.g. "v0.1.0".
* `discoveryItems` - an array of discovery items (see below for details).

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
* `openidConfigurationUri` - URI of the openid configuration well-known endpoint
* `resourceBaseUri` - Base of resource URI, i.e. the part before "/open-banking/v3.0".
* `endpoints` - Array of endpoint and method implementation details.

#### API Specification

The discovery model records specification details in an unambiguous way:

* `apiSpecification`
  * `name` - the `info.title` field from the Swagger/OpenAPI specification file
  * `url` - URI identifier of the specification, i.e. link to specification document
  * `version` - API version number that appears in API paths, e.g. "v3.0"
  * `schemaVersion` - URI identifier of the Swagger/OpenAPI specification file patch version

The property names `url`, `version`, and `schemaVersion` are from the schema.org
[APIReference schema](https://schema.org/APIReference)

Example

```json
{
  "discoveryModel": {
    "discoveryVersion": "v0.1.0",
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

  * `method`* - HTTP method, e.g. "GET" or "POST"
  * `path`* - endpoint path, e.g. "/account-access-consents"
  * `conditionalProperties` - list of optional schema properties that an ASPSP attests it provides (more details in the next section).

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

The specification lists some resource schema properties that may occur `0..1`, or `0..*` times.

When an ASPSP provides a `0..1`, `0..*` occurrence property via its online channel,
it must attest that it provides those properties in its API implementation. An ASPSP must add
such properties to a `conditionalProperties` list in the relevant endpoint definition.

The `conditionalProperties` list contains items. Each item states:

 * `schema` - schema definition name from the Swagger/OpenAPI specification, e.g. "OBTransaction3Detail"
 * `property` - property name from schema, e.g. "Balance"
 * `path` - path to property expressed in [JSON dot notation](https://github.com/tidwall/gjson#path-syntax) format, e.g. Data.Transaction.*.Balance

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

### Resource IDs

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



### Example file

Discovery templates can be found in the [templates directory here](../pkg/discovery/templates).
