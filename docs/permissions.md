# Permission Ingestion from TestCase JSON definitions

The way we've implemented Permission ingestion from test cases is to allow for specifiying the particualar permissions required and particular permissions to be excluded for a test case. Its also possible to rely on default permissions by not specifying permissions.

Within a testcase json definition, the application of permissions is controlled by the an entry in the testcase "context" section, under the "permissions" and "permissions_excluded" keys.

Each test case has three ways of declaratively specifiying permissions :

1. context.permissions: array of required permission strings
2. context.permissions_excluded: array of permissions that must not be present
3. No "permissions" or "permissions_excluded" lists, in which case the default permissions for the endpoint are used if the endpoint has any defaults defined in the specification. Where an endpoint has a number of possible defaults, the defaults that will be used are specificed in the conformance suite permission definition file.

Implications are, if you don't specify permissions you get the default permissions for the endpoint. If you want no allowed permissions, you explicity specifiy the permissions in the "permissions_excluded" list.

The example below specifies a set of permissions for the /transactions endpoint

- "ReadTransactionsBasic"
- "ReadTransactionDetail"
- "ReadTransactionsCredits"
- "ReadTransactionsDebits"

```json
{
        "@id": "#t0001",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/transactions"
        },
        "context": {
            "permissions":["ReadTransactionsBasic","ReadTransactionDetail","ReadTransactionsCredits","ReadTransactionsDebits"]
        },
        "expect": {
            "status-code": 200
        }
    }
```

The following example requires the following permissions

- "ReadTransactionsBasic"
- "ReadTransactionDetail"
- "ReadTransactionsCredits"

and requires that the "ReadTransactionsDebits" not be specified in the consent.

```json
{
        "@id": "#t0002",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/transactions"
        },
        "context": {
            "permissions":["ReadTransactionsBasic","ReadTransactionDetail","ReadTransactionsCredits"],
            "permissions_excluded":["ReadTransactionsDebits"]
        },
        "expect": {
            "status-code": 200
        }
    }
```

The following example will use the default permissions as specified in the Accounts and Transactions 3.0 specification:

```json
{
        "@id": "#t0002",
        "name": "Transaction Test with Permissions",
        "input": {
            "method": "GET",
            "endpoint": "/transactions"
        },
        "context": {
        },
        "expect": {
            "status-code": 200
        }
    }
```

## Permission Sets

PermissionSets are used to define the permission requirements for running testcases against available access_token permissions.

The initial release of the Functional Conformation test suite will use a small number of common predefined PermissionSets. The expectation is that a small number of access tokens will be supplied to the conformance suite as part of the configuration, that match the permission sets and therfore allow a predictable set of core tests to be run. Over time the flexibilty with which permissions can be expressed across test cases will be expanded.

As an example, the following JSON fragment defines two PermissionSets: "accounts-v3.0-all-basic-variant-without-read-pan", and "accounts-v3.0-all-detail-variant-without-read-pan".

```json
{
    "accounts-v3.0-all-basic-variant-without-read-pan": ["ReadAccountsBasic", "ReadBalances", "ReadBeneficiariesBasic",
         "ReadDirectDebits", "ReadOffers", "ReadParty", "ReadPartyPSU", "ReadProducts", "ReadScheduledPaymentsBasic",
          "ReadStandingOrdersBasic", "ReadStatementsBasic", "ReadTransactionsBasic", "ReadTransactionsCredits",
           "ReadTransactionsDebits"],
    "accounts-v3.0-all-detail-variant-without-read-pan": ["ReadAccountsDetail", "ReadBalances", "ReadBeneficiariesDetail",
         "ReadDirectDebits", "ReadOffers", "ReadParty", "ReadPartyPSU", "ReadProducts", "ReadScheduledPaymentsDetail",
         "ReadStandingOrdersDetail", "ReadStatementsDetail", "ReadTransactionsCredits", "ReadTransactionsDebits",
          "ReadTransactionsDetail"]
}
```

## Using Manifests, Rules, Test cases and Permissions Sets

A Manifest defines many Rules which can contain many test cases. Within a Rule, test cases exist in test case sequences. A Rule can have many test case sequences.  A test case sequence contains one or more test cases.

In its simplest form, a Rule contains a single test case. The single test case is contained within a single test sequence as the only member.

Each test case can have a list of required permissions, and a list of excluded permissions that must not be present when the test case is run.

In order to simplify things in the initial iteration, we're taking the view that a Rule presents two permission sets that must be satified in order to run all the test cases defined under that Rule. One permission set contains all the permission required to run all the test cases under a Rule. The second permission set contains all the permission that must not be present (excluded) in order for the test cases defined under to the rule to run.

When determining the Permission Sets required to run all the test cases in a rule, the rule traverses all associated test cases and accumulates two Permission Sets. One set containing all the required permissions, one set containing all the excluded permissions. Once we have the two permission sets, its a realtively straightfoward exercise to match these sets again the avaialble permission sets provided by any supplied access tokens.

