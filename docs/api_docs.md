
# API Documentation

## Base URL
All endpoints are prefixed with `/api`

## Health Check
### Ping
```
GET /ping
```

Simple health check endpoint.

**Request Payload**

_EMPTY BODY_

**Response Payload**

_EMPTY BODY_

## Import Operations
### Review Import
```
POST /import/review
```
Endpoint for reviewing import operations.

**Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["report"],
  "additionalProperties": false,
  "properties": {
    "report": {
      "type": "string",
      "description": "The exported report ZIP archive (base64 encoded)",
      "contentEncoding": "base64"
    }
  },
  "examples": [
    {
      "report": "UEsDBBQACAAIAJVLCVYAAAAA..."
    }
  ]
}
```

**Response Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["discoveryModel"],
  "additionalProperties": false,
  "properties": {
    "discoveryModel": {
      "type": "object",
      "required": ["name", "description", "discoveryVersion", "tokenAcquisition", "discoveryItems"],
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the discovery model"
        },
        "description": {
          "type": "string",
          "description": "Description of the model"
        },
        "discoveryVersion": {
          "type": "string",
          "description": "Version of the discovery protocol"
        },
        "tokenAcquisition": {
          "type": "string",
          "description": "Method of token acquisition"
        },
        "callbackProxyUrl": {
          "type": "string"
        },
        "discoveryItems": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["apiSpecification", "openidConfigurationUri", "resourceBaseUri", "endpoints"],
            "properties": {
              "apiSpecification": {
                "type": "object",
                "required": ["name", "url", "version", "schemaVersion", "manifest"],
                "properties": {
                  "name": {
                    "type": "string"
                  },
                  "url": {
                    "type": "string",
                    "format": "uri"
                  },
                  "version": {
                    "type": "string"
                  },
                  "schemaVersion": {
                    "type": "string",
                    "format": "uri"
                  },
                  "manifest": {
                    "type": "string",
                    "pattern": "^(file://|https://).*$"
                  }
                }
              },
              "openidConfigurationUri": {
                "type": "string",
                "format": "uri"
              },
              "resourceBaseUri": {
                "type": "string",
                "format": "uri"
              },
              "resourceIds": {
                "type": "object",
                "additionalProperties": {
                  "type": "string"
                }
              },
              "endpoints": {
                "type": "array",
                "minItems": 1,
                "items": {
                  "type": "object",
                  "required": ["method", "path"],
                  "properties": {
                    "method": {
                      "type": "string"
                    },
                    "path": {
                      "type": "string",
                      "format": "uri-reference"
                    },
                    "conditionalProperties": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "required": ["schema", "path"],
                        "properties": {
                          "schema": {
                            "type": "string"
                          },
                          "name": {
                            "type": "string"
                          },
                          "property": {
                            "type": "string",
                            "deprecated": true
                          },
                          "path": {
                            "type": "string"
                          },
                          "required": {
                            "type": "boolean"
                          },
                          "request": {
                            "type": "boolean"
                          },
                          "value": {
                            "type": "string"
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "customTests": {
          "type": "array",
          "items": {
            "type": "object"
          }
        }
      }
    }
  },
  "definitions": {
    "cbpiiDebtorAccount": {
      "type": "object",
      "required": ["identification", "scheme_name"],
      "properties": {
        "identification": {
          "type": "string",
          "minLength": 1,
          "maxLength": 256
        },
        "scheme_name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 40
        },
        "name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 70
        }
      }
    }
  }
}
```

### Rerun Import
```
POST /import/rerun
```
Endpoint for rerunning import operations.


**Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["report"],
  "additionalProperties": false,
  "properties": {
    "report": {
      "type": "string",
      "description": "The exported report ZIP archive (base64 encoded)",
      "contentEncoding": "base64"
    }
  },
  "examples": [
    {
      "report": "UEsDBBQACAAIAJVLCVYAAAAA..."
    }
  ]
}
```

**Response Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["discoveryModel"],
  "additionalProperties": false,
  "properties": {
    "discoveryModel": {
      "type": "object",
      "required": ["name", "description", "discoveryVersion", "tokenAcquisition", "discoveryItems"],
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the discovery model"
        },
        "description": {
          "type": "string",
          "description": "Description of the model"
        },
        "discoveryVersion": {
          "type": "string",
          "description": "Version of the discovery protocol"
        },
        "tokenAcquisition": {
          "type": "string",
          "description": "Method of token acquisition"
        },
        "callbackProxyUrl": {
          "type": "string"
        },
        "discoveryItems": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["apiSpecification", "openidConfigurationUri", "resourceBaseUri", "endpoints"],
            "properties": {
              "apiSpecification": {
                "type": "object",
                "required": ["name", "url", "version", "schemaVersion", "manifest"],
                "properties": {
                  "name": {
                    "type": "string"
                  },
                  "url": {
                    "type": "string",
                    "format": "uri"
                  },
                  "version": {
                    "type": "string"
                  },
                  "schemaVersion": {
                    "type": "string",
                    "format": "uri"
                  },
                  "manifest": {
                    "type": "string",
                    "pattern": "^(file://|https://).*$"
                  }
                }
              },
              "openidConfigurationUri": {
                "type": "string",
                "format": "uri"
              },
              "resourceBaseUri": {
                "type": "string",
                "format": "uri"
              },
              "resourceIds": {
                "type": "object",
                "additionalProperties": {
                  "type": "string"
                }
              },
              "endpoints": {
                "type": "array",
                "minItems": 1,
                "items": {
                  "type": "object",
                  "required": ["method", "path"],
                  "properties": {
                    "method": {
                      "type": "string"
                    },
                    "path": {
                      "type": "string",
                      "format": "uri-reference"
                    },
                    "conditionalProperties": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "required": ["schema", "path"],
                        "properties": {
                          "schema": {
                            "type": "string"
                          },
                          "name": {
                            "type": "string"
                          },
                          "property": {
                            "type": "string",
                            "deprecated": true
                          },
                          "path": {
                            "type": "string"
                          },
                          "required": {
                            "type": "boolean"
                          },
                          "request": {
                            "type": "boolean"
                          },
                          "value": {
                            "type": "string"
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "customTests": {
          "type": "array",
          "items": {
            "type": "object"
          }
        }
      }
    }
  },
  "definitions": {
    "cbpiiDebtorAccount": {
      "type": "object",
      "required": ["identification", "scheme_name"],
      "properties": {
        "identification": {
          "type": "string",
          "minLength": 1,
          "maxLength": 256
        },
        "scheme_name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 40
        },
        "name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 70
        }
      }
    }
  }
}
```

## Configuration Management
### Set Global Configuration
```
POST /config/global
```

**Request Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "signing_private",
    "signing_public",
    "transport_private",
    "transport_public",
    "client_id",
    "client_secret",
    "token_endpoint",
    "response_type",
    "token_endpoint_auth_method",
    "authorization_endpoint",
    "resource_base_url",
    "x_fapi_financial_id",
    "issuer",
    "redirect_url",
    "resource_ids",
    "transaction_from_date",
    "transaction_to_date"
  ],
  "properties": {
    "signing_private": {
      "type": "string",
      "pattern": "^-----BEGIN PRIVATE KEY-----\\n[A-Za-z0-9+/=\\n]+-----END PRIVATE KEY-----$"
    },
    "signing_public": {
      "type": "string",
      "pattern": "^-----BEGIN CERTIFICATE-----\\n[A-Za-z0-9+/=\\n]+-----END CERTIFICATE-----$"
    },
    "transport_private": {
      "type": "string",
      "pattern": "^-----BEGIN PRIVATE KEY-----\\n[A-Za-z0-9+/=\\n]+-----END PRIVATE KEY-----$"
    },
    "transport_public": {
      "type": "string",
      "pattern": "^-----BEGIN CERTIFICATE-----\\n[A-Za-z0-9+/=\\n]+-----END CERTIFICATE-----$"
    },
    "tpp_signature_kid": {
      "type": "string"
    },
    "tpp_signature_issuer": {
      "type": "string"
    },
    "tpp_signature_tan": {
      "type": "string"
    },
    "transaction_from_date": {
      "type": "string",
      "format": "date-time"
    },
    "transaction_to_date": {
      "type": "string",
      "format": "date-time"
    },
    "client_id": {
      "type": "string",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
    },
    "client_secret": {
      "type": "string",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
    },
    "token_endpoint": {
      "type": "string",
      "format": "uri"
    },
    "response_type": {
      "type": "string",
      "enum": ["code", "code id_token"]
    },
    "token_endpoint_auth_method": {
      "type": "string",
      "enum": ["private_key_jwt", "tls_client_auth", "client_secret_basic"]
    },
    "request_object_signing_alg": {
      "type": "string",
      "enum": ["PS256", "RS256", "none"]
    },
    "authorization_endpoint": {
      "type": "string",
      "format": "uri"
    },
    "resource_base_url": {
      "type": "string",
      "format": "uri"
    },
    "x_fapi_financial_id": {
      "type": "string"
    },
    "send_x_fapi_customer_ip_address": {
      "type": "boolean"
    },
    "x_fapi_customer_ip_address": {
      "type": "string"
    },
    "issuer": {
      "type": "string",
      "format": "uri"
    },
    "redirect_url": {
      "type": "string",
      "format": "uri"
    },
    "resource_ids": {
      "type": "object",
      "required": ["account_ids", "statement_ids"],
      "properties": {
        "account_ids": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["account_id"],
            "properties": {
              "account_id": {
                "type": "string"
              }
            }
          }
        },
        "statement_ids": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["statement_id"],
            "properties": {
              "statement_id": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "creditor_account": {
      "type": "object",
      "required": ["scheme_name", "identification", "name"],
      "properties": {
        "scheme_name": {
          "type": "string",
          "enum": ["UK.OBIE.SortCodeAccountNumber"]
        },
        "identification": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "international_creditor_account": {
      "type": "object",
      "required": ["scheme_name", "identification", "name"],
      "properties": {
        "scheme_name": {
          "type": "string",
          "enum": ["UK.OBIE.SortCodeAccountNumber"]
        },
        "identification": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "instructed_amount": {
      "type": "object",
      "required": ["value", "currency"],
      "properties": {
        "value": {
          "type": "string",
          "pattern": "^\\d+\\.\\d{2}$"
        },
        "currency": {
          "type": "string",
          "pattern": "^[A-Z]{3}$"
        }
      }
    },
    "currency_of_transfer": {
      "type": "string",
      "pattern": "^[A-Z]{3}$"
    },
    "acr_values_supported": {
      "type": "array",
      "items": {
        "type": "string",
        "pattern": "^urn:openbanking:psd2:(sca|ca)$"
      }
    },
    "payment_frequency": {
      "type": "string",
      "enum": ["EvryDay"]
    },
    "first_payment_date_time": {
      "type": "string",
      "format": "date-time"
    },
    "requested_execution_date_time": {
      "type": "string",
      "format": "date-time"
    },
    "conditional_properties": {
      "type": "array"
    },
    "cbpii_debtor_account": {
      "type": "object",
      "required": ["scheme_name", "identification"],
      "properties": {
        "scheme_name": {
          "type": "string",
          "enum": ["UK.OBIE.SortCodeAccountNumber"]
        },
        "identification": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    }
  }
}
```

**Response Payload**
_Same as Request Payload_

Endpoint to post global configuration settings.

### Get Conditional Property
```
GET /config/conditional-property
```
Retrieve conditional property configuration.

## Discovery Model
### Set Discovery Model
```
POST /discovery-model
```
Endpoint for setting the discovery model.

**Request Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["discoveryModel"],
  "properties": {
    "discoveryModel": {
      "type": "object",
      "required": [
        "name",
        "description",
        "discoveryVersion",
        "tokenAcquisition",
        "discoveryItems"
      ],
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the discovery model"
        },
        "description": {
          "type": "string",
          "description": "Detailed description of the discovery model"
        },
        "discoveryVersion": {
          "type": "string",
          "pattern": "^v\\d+\\.\\d+\\.\\d+$",
          "description": "Version of the discovery model in semver format"
        },
        "tokenAcquisition": {
          "type": "string",
          "description": "Token acquisition method"
        },
        "discoveryItems": {
          "type": "array",
          "minItems": 1,
          "items": {
            "type": "object",
            "required": [
              "apiSpecification",
              "openidConfigurationUri",
              "resourceBaseUri",
              "endpoints"
            ],
            "properties": {
              "apiSpecification": {
                "type": "object",
                "required": [
                  "name",
                  "url",
                  "version",
                  "schemaVersion",
                  "manifest"
                ],
                "properties": {
                  "name": {
                    "type": "string",
                    "description": "Name of the API specification"
                  },
                  "url": {
                    "type": "string",
                    "format": "uri",
                    "description": "URL to the API specification documentation"
                  },
                  "version": {
                    "type": "string",
                    "pattern": "^v\\d+\\.\\d+\\.\\d+$",
                    "description": "Version of the API specification"
                  },
                  "schemaVersion": {
                    "type": "string",
                    "format": "uri",
                    "description": "URL to the OpenAPI schema JSON"
                  },
                  "manifest": {
                    "type": "string",
                    "pattern": "^(file://|https://).+$",
                    "description": "Path or URL to the manifest file"
                  }
                }
              },
              "openidConfigurationUri": {
                "type": "string",
                "format": "uri",
                "description": "URI for OpenID configuration"
              },
              "resourceBaseUri": {
                "type": "string",
                "format": "uri",
                "description": "Base URI for API resources"
              },
              "endpoints": {
                "type": "array",
                "minItems": 1,
                "items": {
                  "type": "object",
                  "required": ["method", "path"],
                  "properties": {
                    "method": {
                      "type": "string",
                      "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"],
                      "description": "HTTP method for the endpoint"
                    },
                    "path": {
                      "type": "string",
                      "pattern": "^/",
                      "description": "Endpoint path, must start with /"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

**Response Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "token_endpoints",
    "token_endpoint_auth_methods",
    "default_token_endpoint_auth_method",
    "request_object_signing_alg_values_supported",
    "default_request_object_signing_alg_values_supported",
    "authorization_endpoints",
    "issuers",
    "default_transaction_from_date",
    "default_transaction_to_date",
    "response_types_supported",
    "acr_values_supported"
  ],
  "properties": {
    "token_endpoints": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "string",
          "format": "uri",
          "description": "Token endpoint URL for the specific schema version"
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their token endpoints"
    },
    "token_endpoint_auth_methods": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "array",
          "items": {
            "type": "string",
            "enum": ["tls_client_auth", "private_key_jwt", "client_secret_basic"]
          },
          "minItems": 1,
          "uniqueItems": true
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their supported authentication methods"
    },
    "default_token_endpoint_auth_method": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "string",
          "enum": ["tls_client_auth", "private_key_jwt", "client_secret_basic"]
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their default authentication method"
    },
    "request_object_signing_alg_values_supported": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "array",
          "items": {
            "type": "string",
            "enum": ["none", "PS256", "RS256"]
          },
          "minItems": 1,
          "uniqueItems": true
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their supported signing algorithms"
    },
    "default_request_object_signing_alg_values_supported": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "string",
          "enum": ["none", "PS256", "RS256"]
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their default signing algorithm"
    },
    "authorization_endpoints": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "string",
          "format": "uri"
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their authorization endpoints"
    },
    "issuers": {
      "type": "object",
      "patternProperties": {
        "^schema_version=https://raw\\.githubusercontent\\.com/OpenBankingUK/read-write-api-specs/v\\d+\\.\\d+\\.\\d+/dist/openapi/.*\\.json$": {
          "type": "string",
          "format": "uri"
        }
      },
      "additionalProperties": false,
      "description": "Map of schema versions to their issuer URLs"
    },
    "default_transaction_from_date": {
      "type": "string",
      "format": "date-time",
      "description": "Default start date for transaction queries"
    },
    "default_transaction_to_date": {
      "type": "string",
      "format": "date-time",
      "description": "Default end date for transaction queries"
    },
    "response_types_supported": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": ["code", "code id_token"]
      },
      "minItems": 1,
      "uniqueItems": true,
      "description": "List of supported OAuth response types"
    },
    "acr_values_supported": {
      "type": "array",
      "items": {
        "type": "string",
        "pattern": "^urn:openbanking:psd2:(sca|ca)$"
      },
      "minItems": 1,
      "uniqueItems": true,
      "description": "List of supported Authentication Context Class References"
    }
  }
}
```

## Test Case Management
### List Test Cases
```
GET /test-cases
```
Retrieve all test cases.

**Request Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["specCases", "specTokens"],
  "properties": {
    "specCases": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["apiSpecification", "testCases"],
        "properties": {
          "apiSpecification": {
            "type": "object",
            "required": ["name", "url", "version", "schemaVersion", "manifest"],
            "properties": {
              "name": {
                "type": "string"
              },
              "url": {
                "type": "string",
                "format": "uri"
              },
              "version": {
                "type": "string",
                "pattern": "^v\\d+\\.\\d+\\.\\d+$"
              },
              "schemaVersion": {
                "type": "string",
                "format": "uri"
              },
              "manifest": {
                "type": "string",
                "pattern": "^file://.*\\.json$"
              }
            }
          },
          "testCases": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["@id", "name", "input", "context", "expect"],
              "properties": {
                "@id": {
                  "type": "string",
                  "pattern": "^OB-\\d{3}-[A-Z]+-\\d{6}$"
                },
                "name": {
                  "type": "string"
                },
                "detail": {
                  "type": "string"
                },
                "purpose": {
                  "type": "string"
                },
                "refURI": {
                  "type": "string",
                  "format": "uri"
                },
                "input": {
                  "type": "object",
                  "required": ["method", "endpoint", "headers"],
                  "properties": {
                    "method": {
                      "type": "string",
                      "enum": ["GET", "POST", "PUT", "DELETE"]
                    },
                    "endpoint": {
                      "type": "string",
                      "pattern": "^/"
                    },
                    "headers": {
                      "type": "object",
                      "properties": {
                        "x-fapi-financial-id": { "type": "string" },
                        "x-fapi-interaction-id": { "type": "string" },
                        "x-fcs-testcase-id": { "type": "string" },
                        "Authorization": { "type": "string" }
                      }
                    },
                    "queryParameters": {
                      "type": "object",
                      "additionalProperties": {
                        "type": "string"
                      }
                    },
                    "removeheaders": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    }
                  }
                },
                "context": {
                  "type": "object",
                  "properties": {
                    "accountId": { "type": "string" },
                    "baseurl": { 
                      "type": "string",
                      "format": "uri"
                    },
                    "permissions": {
                      "type": "array",
                      "items": { "type": "string" }
                    },
                    "permissions-excluded": {
                      "type": "array",
                      "items": { "type": "string" }
                    },
                    "tokenScope": { "type": "string" },
                    "x-fapi-financial-id": { "type": "string" },
                    "x-fapi-interaction-id": { "type": "string" }
                  }
                },
                "expect": {
                  "type": "object",
                  "properties": {
                    "status-code": { "type": "integer" },
                    "schema-validation": { "type": "boolean" },
                    "matches": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "header": { "type": "string" },
                          "value": { "type": "string" },
                          "header-present": { "type": "string" },
                          "json": { "type": "string" }
                        }
                      }
                    },
                    "contextPut": { "type": "object" }
                  }
                },
                "expect_one_of": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "status-code": { "type": "integer" },
                      "matches": {
                        "type": "array",
                        "items": {
                          "type": "object",
                          "properties": {
                            "json": { "type": "string" },
                            "value": { "type": "string" }
                          }
                        }
                      },
                      "contextPut": { "type": "object" }
                    }
                  }
                },
                "apiName": { "type": "string" },
                "apiVersion": { 
                  "type": "string",
                  "pattern": "^v\\d+\\.\\d+\\.\\d+$"
                }
              }
            }
          }
        }
      }
    },
    "specTokens": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["specIdentifier", "namedPermissions"],
        "properties": {
          "specIdentifier": { "type": "string" },
          "namedPermissions": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["name", "codeSet", "consentUrl"],
              "properties": {
                "name": { "type": "string" },
                "codeSet": {
                  "type": "object",
                  "required": ["codes", "testIds"],
                  "properties": {
                    "codes": {
                      "type": "array",
                      "items": { "type": "string" }
                    },
                    "testIds": {
                      "type": "array",
                      "items": { 
                        "type": "string",
                        "pattern": "^OB-\\d{3}-[A-Z]+-\\d{6}$"
                      }
                    }
                  }
                },
                "consentUrl": {
                  "type": "string",
                  "format": "uri"
                }
              }
            }
          }
        }
      }
    }
  }
}
```

**Response Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["specCases", "specTokens"],
  "additionalProperties": false,
  "properties": {
    "specCases": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["apiSpecification", "testCases"],
        "additionalProperties": false,
        "properties": {
          "apiSpecification": {
            "$ref": "#/definitions/apiSpecification"
          },
          "testCases": {
            "type": "array",
            "items": {
              "$ref": "#/definitions/testCase"
            }
          }
        }
      }
    },
    "specTokens": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["specIdentifier", "namedPermissions"],
        "additionalProperties": false,
        "properties": {
          "specIdentifier": {
            "type": "string"
          },
          "namedPermissions": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["name", "codeSet", "consentUrl"],
              "additionalProperties": false,
              "properties": {
                "name": {
                  "type": "string",
                  "pattern": "^accountToken\\d{4}$"
                },
                "codeSet": {
                  "type": "object",
                  "required": ["codes", "testIds"],
                  "additionalProperties": false,
                  "properties": {
                    "codes": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    },
                    "testIds": {
                      "type": "array",
                      "items": {
                        "type": "string",
                        "pattern": "^OB-\\d{3}-[A-Z]+-\\d{6}$"
                      }
                    }
                  }
                },
                "consentUrl": {
                  "type": "string",
                  "format": "uri"
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "apiSpecification": {
      "type": "object",
      "required": ["name", "url", "version", "schemaVersion", "manifest"],
      "additionalProperties": false,
      "properties": {
        "name": {
          "type": "string"
        },
        "url": {
          "type": "string",
          "format": "uri"
        },
        "version": {
          "type": "string",
          "pattern": "^v\\d+\\.\\d+\\.\\d+$"
        },
        "schemaVersion": {
          "type": "string",
          "format": "uri"
        },
        "manifest": {
          "type": "string",
          "pattern": "^file://.*\\.json$"
        }
      }
    },
    "testCase": {
      "type": "object",
      "required": ["@id", "name", "input", "context", "expect", "apiName", "apiVersion"],
      "additionalProperties": false,
      "properties": {
        "@id": {
          "type": "string",
          "pattern": "^OB-\\d{3}-[A-Z]+-\\d{6}$"
        },
        "name": {
          "type": "string"
        },
        "detail": {
          "type": "string"
        },
        "purpose": {
          "type": "string"
        },
        "refURI": {
          "type": "string",
          "format": "uri"
        },
        "input": {
          "type": "object",
          "required": ["method", "endpoint", "headers"],
          "additionalProperties": false,
          "properties": {
            "method": {
              "type": "string",
              "enum": ["GET", "POST", "PUT", "DELETE", "PATCH"]
            },
            "endpoint": {
              "type": "string",
              "pattern": "^/"
            },
            "headers": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            },
            "queryParameters": {
              "type": "object",
              "additionalProperties": {
                "type": "string"
              }
            },
            "removeheaders": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        },
        "context": {
          "type": "object",
          "additionalProperties": true,
          "properties": {
            "baseurl": {
              "type": "string",
              "format": "uri"
            },
            "accountId": {
              "type": "string"
            },
            "permissions": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "permissions-excluded": {
              "type": "array",
              "items": {
                "type": "string"
              }
            },
            "tokenScope": {
              "type": "string"
            }
          }
        },
        "expect": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "status-code": {
              "type": "integer"
            },
            "schema-validation": {
              "type": "boolean"
            },
            "matches": {
              "type": "array",
              "items": {
                "type": "object",
                "additionalProperties": false,
                "oneOf": [
                  {
                    "required": ["header-present"],
                    "properties": {
                      "header-present": {
                        "type": "string"
                      }
                    }
                  },
                  {
                    "required": ["header", "value"],
                    "properties": {
                      "header": {
                        "type": "string"
                      },
                      "value": {
                        "type": "string"
                      }
                    }
                  },
                  {
                    "required": ["json", "value"],
                    "properties": {
                      "json": {
                        "type": "string"
                      },
                      "value": {
                        "type": "string"
                      }
                    }
                  }
                ]
              }
            },
            "contextPut": {
              "type": "object"
            }
          }
        },
        "expect_one_of": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["status-code", "contextPut"],
            "additionalProperties": false,
            "properties": {
              "status-code": {
                "type": "integer"
              },
              "matches": {
                "type": "array",
                "items": {
                  "type": "object",
                  "required": ["json", "value"],
                  "additionalProperties": false,
                  "properties": {
                    "json": {
                      "type": "string"
                    },
                    "value": {
                      "type": "string"
                    }
                  }
                }
              },
              "contextPut": {
                "type": "object"
              }
            }
          }
        },
        "apiName": {
          "type": "string"
        },
        "apiVersion": {
          "type": "string",
          "pattern": "^v\\d+\\.\\d+\\.\\d+$"
        }
      }
    }
  }
}
```

## Test Runner
### Start Test Run
```
POST /run
```
Initiate a new test run.

### WebSocket Connection for Test Results
```
GET /run/ws
```
WebSocket endpoint for listening to test results in real-time.

**Request Payload**

_EMPTY BODY_

**Response Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["type", "value"],
  "additionalProperties": false,
  "properties": {
    "type": {
      "type": "string",
      "description": "The type of WebSocket event"
    },
    "value": {
      "type": "object",
      "required": ["token_names"],
      "additionalProperties": false,
      "properties": {
        "token_names": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Array of token names that have been acquired"
        }
      }
    }
  },
  "examples": [
    {
      "type": "acquired_all_access_tokens",
      "value": {
        "token_names": ["token1", "token2"]
      }
    }
  ]
}
```

### Stop Test Run
```
DELETE /run
```
Stop an ongoing test run.

## Redirect Handling
### Fragment OK
```
POST /redirect/fragment/ok
```
Handle successful fragment redirects.

**Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "code",
    "scope",
    "id_token",
    "state"
  ],
  "additionalProperties": false,
  "properties": {
    "code": {
      "type": "string",
      "description": "Authorization code"
    },
    "scope": {
      "type": "string",
      "description": "OAuth scope of the request"
    },
    "id_token": {
      "type": "string",
      "description": "OpenID Connect ID Token"
    },
    "state": {
      "type": "string",
      "description": "OAuth state parameter for request verification"
    }
  },
  "examples": [
    {
      "code": "auth_code_123",
      "scope": "openid profile",
      "id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9...",
      "state": "abc123"
    }
  ]
}
```

**Response Payload**

_EMPTY BODY_

### Query OK
```
POST /redirect/query/ok
```
Handle successful query redirects.

**Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "code",
    "scope",
    "id_token",
    "state"
  ],
  "additionalProperties": false,
  "properties": {
    "code": {
      "type": "string",
      "description": "Authorization code"
    },
    "scope": {
      "type": "string",
      "description": "OAuth scope of the request"
    },
    "id_token": {
      "type": "string",
      "description": "OpenID Connect ID Token"
    },
    "state": {
      "type": "string",
      "description": "OAuth state parameter for request verification"
    }
  },
  "examples": [
    {
      "code": "auth_code_123",
      "scope": "openid profile",
      "id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9...",
      "state": "abc123"
    }
  ]
}
```

**Response Payload**

_EMPTY BODY_

### Error
```
POST /redirect/error
```
Handle redirect errors.

*Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "error_description",
    "error",
    "state"
  ],
  "additionalProperties": false,
  "properties": {
    "error_description": {
      "type": "string",
      "description": "Human-readable description of the error"
    },
    "error": {
      "type": "string",
      "description": "OAuth 2.0 error code",
      "enum": [
        "invalid_request",
        "unauthorized_client",
        "access_denied",
        "unsupported_response_type",
        "invalid_scope",
        "server_error",
        "temporarily_unavailable"
      ]
    },
    "state": {
      "type": "string",
      "description": "OAuth state parameter from the original request"
    }
  },
  "examples": [
    {
      "error_description": "The authorization request was invalid",
      "error": "invalid_request",
      "state": "abc123"
    },
    {
      "error_description": "The user denied access to the resource",
      "error": "access_denied",
      "state": "xyz789"
    }
  ]
}
```

**Response Payload**

_EMPTY BODY_

## Export
### Export Data
```
POST /export
```
Endpoint for exporting data.

**Request Payload**

```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "environment",
    "implementer",
    "authorised_by",
    "job_title",
    "products",
    "has_agreed",
    "add_digital_signature"
  ],
  "additionalProperties": false,
  "properties": {
    "environment": {
      "type": "string",
      "description": "Environment used for testing"
    },
    "implementer": {
      "type": "string",
      "description": "Implementer/Brand Name"
    },
    "authorised_by": {
      "type": "string",
      "description": "Authorised by"
    },
    "job_title": {
      "type": "string",
      "description": "Job Title"
    },
    "products": {
      "type": "array",
      "description": "Products tested, e.g., 'Business, Personal, Cards'",
      "items": {
        "type": "string"
      },
      "minItems": 1
    },
    "has_agreed": {
      "type": "boolean",
      "description": "I agree"
    },
    "add_digital_signature": {
      "type": "boolean",
      "description": "Sign this report"
    }
  },
  "examples": [
    {
      "environment": "sandbox",
      "implementer": "Example Bank Ltd",
      "authorised_by": "John Smith",
      "job_title": "Technical Lead",
      "products": ["Business", "Personal", "Cards"],
      "has_agreed": true,
      "add_digital_signature": true
    }
  ]
}
```

**Response Payload**

_ZIP FILE_

## Utility
### Version Check
```
GET /version
```
Check the current version of the API

**Request Payload**
_EMPTY BODY_

**Response Payload**
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "message", "update"],
  "additionalProperties": false,
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^v\\d+\\.\\d+\\.\\d+$",
      "description": "The semantic version of the conformance suite"
    },
    "message": {
      "type": "string",
      "description": "A human-readable message describing the version status"
    },
    "update": {
      "type": "boolean",
      "description": "Indicates whether an update is available"
    }
  },
  "examples": [
    {
      "version": "v1.9.1",
      "message": "Conformance Suite is running the latest version v1.9.1",
      "update": false
    }
  ]
}
```