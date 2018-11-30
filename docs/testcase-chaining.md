# Test case parameter chaining

Test case chaining is the ability to use part of the results of one test case as input to another test case. 
Within the test case model, a rule has one or more test sequences, and each test sequence has one or more test cases. 
Parameter chaining occurs between two test cases within the same test sequence. The source test case is run before the receiving test case. 
Any part of a response that can be selected for comparison with a 'Matches' clause, can be used to extract data from the source test case, 
and make it available to the receiving test case.
Data is transferred between test cases using the context. Test cases have access to four levels of context scope:

1. test case scope - only visible to the test case
2. test sequence scope - visible to any test case within the test sequence
3. rule scope - visible to any test sequence within a rule, and all the test sequence's test cases
4. global/manifest scope - visible to all testing running - discovery information would typically be accessible via global scope

For test case chaining, the test sequence scope is involved as this exists for the duration of a test sequence
and therefore enables parameters to be available after a test case has completed and before a subsequent test case has started.

In order to move a parameter between test cases, two directives have been introduced, **contextGet** and **contextPut**. 
**contextGet** is used in the situation where you **get** a variable from the context and insert the variable into the request payload. 
**contextPut** is used in the situation where you **put** a variable into the context for later use.

Here is a worked example of test case parameter chaining to explain how it works:

Consider the following json fragment which defines two test cases:

- t0001
- t0002

```json
{
    "@id": "t0001",
    "name": "Get Accounts Basic",
    "purpose": "Accesses the Accounts endpoint and retrieves a list of PSU accounts",
    "input": {
        "method": "GET",
        "endpoint": "/accounts/"
    },
    "context": {
        "baseurl": "http://myaspsp"
    },
    "expect": {
        "status-code": 200,
        "matches": [{
            "description": "A json match on response body",
            "json": "Data.Account.0.AccountId",
            "value": "500000000000000000000001"
        }],
        "contextPut": {
            "matches": [{
                "name": "AccountId",
                "description": "A json match to extract variable to context",
                "json": "Data.Account.1.AccountId"
            }]
        }
    }
}, {
    "@id": "t0002",
    "name": "Get Accounts using AccountId",
    "description": "Retrieve account information for given AccountId",
    "input": {
        "method": "GET",
        "endpoint": "/accounts/{AccountId}",
        "contextGet": {
            "matches": [{
                "name": "AccountId",
                "description": "supplies the account number to be queried",
                "replaceInEndpoint": "{AccountId}"
            }]
        }
    },
    "context": {
        "baseurl": "http://myaspsp"
    },
    "expect": {
        "status-code": 200,
        "matches": [{
            "description": "A json match on response body",
            "json": "Data.Account.0.Account.0.Identification",
            "value": "GB29PAPA20000390210099"
        }]
    }
}]
```

Test case t0001 does the following:-

- Makes an http GET call to resource endpoint /accounts/
- Checks that the call response status code is 200
- Examines the response body for the first occurrence of the JSON field AccountId using the JSON match string "Data.Account.0.AccountId".  
In the response body, the matched field would appear as follows:-

```json
{
 "Data": {
    "Account": [{
                "AccountId": "500000000000000000000001",

```

- The AccountId field is checked to see if it matches the value "500000000000000000000001"
- If the matches are successful then the value of the JSON field expression "Data.Account.1.AccountId" is extracted from the response
and put into the **context**

```json
{
 "Data": {
    "Account": [{
                    "AccountId": "500000000000000000000001",
                    "etc":"..."
                }, {
                    "AccountId": "500000000000000000000007",
                    "etc":"..."
                }]
 }

```

- The JSON Expression "Data.Account.1.AccountId results in the AccountId value "500000000000000000000007" being inserted into the **context**
- The AccountId value is then added to the **Context** under the name "AccountId".
- The "AccountId" variable with value "500000000000000000000007" becomes available to other test cases that use the same **context**.

Test case t0002 then does the following:-

- Accesses the value of "AccountId" in the context and retrieves the value
- Uses the retrieved accountId of "500000000000000000000007" to perform an endpoint string replacement
- Finds the string "{AccountId} in the test case endpoint and replaces the string with the value "500000000000000000000007"
- Makes an HTTP GET call to resource endpoint /accounts/500000000000000000000007
- Checks the HTTP response code is 200
- Performs a JSON field match on the call response body
- Uses the JSON pattern "Data.Account.0.Account.0.Identification" to check the value of response field "Identification", 
as shown in the following JSON response fragment:

```json
{
"Data": {
    "Account": [
        {
            "AccountId": "500000000000000000000007",
            "Currency": "GBP",
            "Nickname": "xxxx0001",
            "AccountType": "Business",
            "AccountSubType": "CurrentAccount",
            "Account": [
                {
                    "SchemeName": "IBAN",
                    "Identification": "GB29PAPA20000390210099",
```

- Checks that the Identification field value matches "GB29PAPA20000390210099" as specified in the test case

In summary, this simple example, uses declarative text to create two test cases which :-

- Call an initial resource endpoint to get a list of resource id's
- Check that a specific AccountId exists in the returned list
- Extract a second AccountId from the returned list and puts the AccountId value in the context
- Run a second test case which modifies its resource endpoint based on the AccountId retrieved from the previous call
- Check the value of a response field returned for the second AccountId
