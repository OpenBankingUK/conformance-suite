# Matching Model

The conformance suite matching model specifies how the suite and check the responses that are returned from API implementations.

The following methods are available for checking a API http response and its contents:

- Http Status Code field
- Http Header value comparison
- Http Header value regex
- Http Header present
- Http Body regex
- Http Body - Json field present
- Http Body - Json specific number of particular fields present
- Http Body - Json field content
- Http Body - Json field with Regex applied
- Http Body Length - Checks the expected response body length

The following json fragments show examples of each of the selection options :-

#### Status Code

Checks the HTTP Status code returned in the response

```json
    "expect": {
        "status-code": 200,
    }
```

#### Header Value

Check that the specified HTTP header is present in the response, and that it has the specified value.

**Note**: All header checks are case in-sensitive as specified in RFC7230 Hypertext Transfer Protocol (HTTP/1.1) - https://tools.ietf.org/html/rfc7230#section-3.2

```json
    "expect": {
        "matches": [{
            "description": "Example match an http header value",
            "header": "Content-Type",
            "value": "application/json"
        }],
    }
```

#### Header Regex

Check that the value of the specified http header matches the supplied regular expression

```json
    "expect": {
        "matches": [{
            "description": "match a header value using a regex",
            "header":"Proxied-via",
            "regex": "^mybox$",
        }],
    }
```

#### Header Present

Check that the specified http header is present in the response.

```json
    "expect": {
        "matches": [{
            "description": "check that a header is present",
            "header-present": "content-length"
         }]
    }
```

#### Body Regex

Check that the response body matches the specified regular expression

```json
    "expect": {
        "matches": [{
            "description": "body-regex",
            "regex": ".*London Bridge.*",
        }],
    }
```

#### Body JSON Present

Check that the JSON field specified exists in the response body

```json
    "expect": {
        "matches": [{
            "description": "A json field present",
            "json": "Data.Account.Accountid",
        }],
    }
```

#### Body JSON Count

Check that the specified JSON field - which must be within a JSON array structure - exists the specified number of times in the response body

```json
    "expect": {
        "matches": [{
            "description": "A json field count present",
            "json": "Data.Account.[*]Accountid",
            "count": 4
        }],
    }
```

#### Body JSON Value

Check that the specified JSON response body field has the specified value

```json
    "expect": {
        "status-code": 200,
        "matches": [{
            "description": "A json match on response body",
            "json": "Data.Account.Accountid",
            "value": "XYZ1231231231231"
        }],
    }
```

#### Body JSON Regex

Check that the specified JSON response body field matches the specified regular expression

```json
    "expect": {
        "status-code": 200,
        "matches": [{
            "description": "A json match on response body",
            "json": "Data.Account.Accountid",
            "regex": "$*.^"
        }],
    }
```

#### Body Length

Check that the response body has the expected length

```json
    "expect": {
        "matches": [{
            "description": "Check the length in bytes of the body is as specified",
            "body-length": 0
        }],
    }
```
