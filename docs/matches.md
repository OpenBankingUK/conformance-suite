# Matching Model

The conformance suite matching model specifies how the suite and check the responses that are returned from API implementations.

The following methods are available for checking a API http response and its contents:

- Http Status Code field
- Http Header value comparision
- Http Header value regex
- Http Header present
- Http Body regex
- Http Body - Json field present
- Http Body - Json specific number of particualar fields present
- Http Body - Json field content
- Http Body length

The following json fragments show examples of each of the selection options :-

```json
    "expect": {
        "status-code": 200,
    }
```

```json
    "expect": {
        "matches": [{
            "description": "Example match an http header value",
            "header": "Content-Type",
            "value": "application/json"
        }],
    }
```

```json
    "expect": {
        "matches": [{
            "description": "match a header value using a regex",
            "header":"Proxied-via",
            "header-regex": "^mybox$",
        }],
    }
```

```json
    "expect": {
        "matches": [{
            "description": "check that a header is present",
            "header-present": "content-length"
         }]
    }
```

```json
    "expect": {
        "matches": [{
            "description": "body-regex",
            "body-regex": ".*London Bridge.*",
        }],
    }
```

```json
    "expect": {
        "matches": [{
            "description": "A json field present",
            "json": "Data.Account.Accountid",
        }],
    }
```

```json
    "expect": {
        "matches": [{
            "description": "A json field count present",
            "json": "Data.Account.[*]Accountid",
            "count: 4
        }],
    }
```




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

```json
    "matches": [{
        "description": "Content length??",
        "response-length": "Data.Account.Accountid",
        "value": "0"
    }],
```
