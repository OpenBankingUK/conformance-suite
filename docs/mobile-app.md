# Mobile Application support

FCS allows testing applications that allows mobile interactions only and cannot redirect to localhost without tunneling
software.

To allow this OBIE provides a proxy that allows storing consent callback parameters.

## Configuration

Sample configuration is provided. This configuration uses Ozone model Bank through a mobile browser.

Key elements of the discovery file are:

```json5
{
  "discoveryModel": {
    "name": "ob-v3.1-ozone",
    "description": "O3 Mobile PSU consent flow. An Open Banking UK discovery template for v3.1 of Accounts and Payments with pre-populated model Bank (Ozone) data.",
    "discoveryVersion": "v0.4.0",

    //Required for mobile app support
    "tokenAcquisition": "mobile",
    
    //Required for mobile app support
    "consentCallbackUrl": "https://fcs-callback-proxy.openbanking.rocks",
    "discoveryItems": [
      {
        //...
      }
    ]
  }
}

```

