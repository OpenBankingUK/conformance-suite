# Account and Transaction API endpoint permutations file

When a new version of the Account and Transaction API specification is
released we need to generate a new endpoint permissions permutations file.

## Generate endpoint permutations file

You can generate new endpoint permissions permutations file as follows.

### Create new endpoints CSV file

Find the list of endpoints in the specification and its swagger file. E.g. you can see [v2.0.0 endpoints in the specification here](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/127009546/Account+and+Transaction+API+Specification+-+v2.0.0#AccountandTransactionAPISpecification-v2.0.0-Endpoints).

Find permissions for endpoints in the specification. E.g. you can see [v2.0.0
permissions here](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/127009546/Account+and+Transaction+API+Specification+-+v2.0.0#AccountandTransactionAPISpecification-v2.0.0-Permissions).

Create a new endpoints CSV file as follows:

* Add a new `priv/endpoints/endpoints-*.*.csv` file.
* E.g. if spec version is `v3.1.*` create `priv/endpoints/endpoints-3.1.csv`.
* Copy the headers from an existing endpoints CSV.
* Manually populate fields for each endpoint in CSV.

Fields in the CSV file are:

* `endpoint` - resource path,  e.g. `/accounts`
* `version` - specification major minor, e.g. `2.0`
* `required` - either `optional`, `conditional`, or `mandatory`
* `read` - generic read permission, either blank or `ReadXxx`
* `readbasic` - basic read permission, either blank or `ReadXxxBasic`
* `readdetail` - detail read permission, either blank or `ReadXxxDetail`
* `readcredits` - credit transaction access, either blank or `ReadTransactionsCredits`
* `readdebits` - debit transaction access, either blank or `ReadTransactionsDebits`
* `readpan` - account number access, either blank or `ReadPAN` when resource can contain `PAN` numbers

E.g. here's some example endpoint CSV:
```csv
endpoint,version,required,read,readbasic,readdetail,readcredits,readdebits,readpan
/accounts,2.0,mandatory,,ReadAccountsBasic,ReadAccountsDetail,,,ReadPAN
...
/transactions,2.0,optional,,ReadTransactionsBasic,ReadTransactionsDetail,ReadTransactionsCredits,ReadTransactionsDebits,ReadPAN
...
/accounts/{AccountId}/statements/{StatementId}/transactions,2.0,conditional,ReadStatementsDetail,ReadTransactionsBasic,ReadTransactionsDetail,ReadTransactionsCredits,ReadTransactionsDebits,ReadPAN
...
/products,2.0,optional,ReadProducts,,,,,
...
```

#### Endpoint optionality

* `mandatory` endpoints *must* be implemented by an ASPSP.
* `optional` endpoints *may* be implemented by an ASPSP.
* `conditional` endpoints *must* be implemented by an ASPSP if these are made
available to the PSU in the ASPSP's existing Online Channel.

It is up to each ASPSP to make their own regulatory interpretation, based
on PSD2, as to which of the conditional endpoints and fields must be
implemented.

### Generate endpoint permutations JSON file

Update the `priv/endpoints/generate.exs` with the endpoints file to run.
E.g. add `"./priv/endpoints/endpoints-3.1.csv"`.

To generate the new permutations file, run:

```sh
mix compliance.permutations
```

This generates a new `priv/endpoints/permutations-*.*.csv` file.
