# Swagger Schema Validation

This packages provides functionality for validate http responses, including body schema against a 
swagger definition using go-swagger library.


## Open Banking Specs


- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/account-info-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/account-info-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json


- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/payment-initiation-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/payment-initiation-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json


- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/confirmation-funds-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/confirmation-funds-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/confirmation-funds-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/confirmation-funds-swagger.json

- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/event-notifications-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/event-notifications-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/event-notifications-swagger.json
- https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/event-notifications-swagger.json


### Setup

This package uses flattened version of Open Banking API specs, you will the `swaggger` CLI tool from:

https://goswagger.io/install.html

To flatten a spec use command:

```bash
swagger flatten --with-expand -o account-info-swagger.flattened.json account-info-swagger.json
```


### Usage

This package is a wrapper around swagger library validator with adicional status code and content type check, 
also allowing implementing other custom validator and version variations.

```go 
func main() {
    v, err := NewSwaggerValidator("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json")
    if err != nil {
        panic(err)
    }
    
    failures, err := v.Validate(r)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Validation failures found:\n%v", failures)    
}
```
