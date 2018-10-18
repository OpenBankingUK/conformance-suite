# `conformance-suite`

# Conformance-suite

## Overview

```go

                               +  Access Tokens
                               |
                               |
                               v
  +------------+      +-----------------+             +-----------------+
  |            |      |                 |   MATLS     |                 |
  |  Driver    +----->|   Proxy         +------------>|   Ozone Aspsp   |
  |            |      |                 |             |                 |
  +------------+      +--+--------------+             +-----------------+
                         |           ^
                         |           |
                         |           |
                     +---v-----------+----+           +-----------------+
                     |  Request/Response  |           |                 |
                     |                    |+--------->|  Log Reporter   |
                     |    Validation      |           |                 |
                     |                    |           +-----------------+
                     +--------------------+
                              ^
                              |
                      +-------+---------+
                      |                 |
                      |  Swagger Spec   |
                      |                 |
                      +-----------------+
```

### Overview Detail

* A driver application sends http request as part of a test suite to the Proxy
* The Proxy is provided with valid access tokens and certificates for matls
* The Proxy connects to ozone over matls using the provide access tokens
* The driver sends request to the proxy
  * The proxy passes all request and responses to the Validator for validation
  * The validator using the swagger spec as the basis for validation
  * The validator records validation results to the log Reporter
  * The proxy ensures the requests/resposes travel to ozone

### Application code structure

* main.go
  * get application config
    * as a temporary measure reads this from a json file config/config.json
    * also reads transport certs from config directory

```javascript
{
  "softwareStatementId": "1p8bGrHhJRrphjFk0qwNAU",
  "clientScopes": "AuthoritiesReadAccess ASPSPReadAccess TPPReadAccess",
  "keyId": "BFnipP2g4ZaaFySsIaigOUoCP2E",
  "verbose":true,
  "specLocation":"swagger/rw20test.json",
  "bindAddress":":8989",
  "targetHost":"http://localhost:3000",
  "accountAccessToken": "d2d50414-d7a2-48b7-9ddd-aa4bf3cb47ac"
}
```

* start proxy on port 8989
  * pass through to target server on http://localhost:3000
  * you'll need to put somethere there for a response - will be ozone hopfully by Monday
  * If you put something that can gives a 200 response to anything then you should be able to test
  * if you have httpie installed - the following command will push some data through
  * http -v get :8989/open-banking/v2.0/account-requests/myreq x-fapi-financial-id:123 Authorization:123

* Loads the swagger specification
  * currently from swagger/rw20test.json - configured in the config.json file ("specLocation")
* Creates a Proxy Object using the spec, app config, and passes in a LogReporter
* Waits for requests
* Shutdown (ctrl-c typically) does a bit of tidyup

## Other notes

The original code for the swagger-proxy is at
[https://github.com/gchaincl/swagger-proxy](https://github.com/gchaincl/swagger-proxy). I've improved the code and fixed some things - like the abilty to configure TLS - so the original code is not a direct dropin any more. Plus it hasn't changed in a year, so really just a starting point rather than something we need to actively watch. I guess at some stage we may want to offer changes back to it.

Currently not proper package managment (go dep etc...)

Logging is done with Logrus.
I've also used a throwaway console output lib at github.com/x-cray/logrus-prefixed-formatter which enhances logrus output - but this is just a toy. Also how I've setup Logrus is fairly arbitary and open to suggestions for better ways.

### LogReporter

### Stuff which needs fixing

I need to check which http headers are being passed through, not all of them are - specifically "applicaiton/json; charset=UTF-8" appears to struggle. Not a show stopper, just needs investigation.


## Development
### Docker
#### build
```sh
$ docker build -t "openbanking/conformance-suite:latest" .
$ docker run --rm -it -p 8080:8080 "openbanking/conformance-suite:latest"
web_1  | conformance-suite2018/10/10 14:21:15 starting ...
web_1  | conformance-suite2018/10/10 14:21:15 started ...
web_1  |
web_1  |    ____    __
web_1  |   / __/___/ /  ___
web_1  |  / _// __/ _ \/ _ \
web_1  | /___/\__/_//_/\___/ v3.2.1
web_1  | High performance, minimalist Go web framework
web_1  | https://echo.labstack.com
web_1  | ____________________________________O/_______
web_1  |                                     O\
web_1  | ⇨ http server started on [::]:8080
```

If you have `make` installed this becomes:

```sh
$ make build_image
$ make run_image
web_1  | conformance-suite2018/10/10 14:21:15 starting ...
web_1  | conformance-suite2018/10/10 14:21:15 started ...
web_1  |
web_1  |    ____    __
web_1  |   / __/___/ /  ___
web_1  |  / _// __/ _ \/ _ \
web_1  | /___/\__/_//_/\___/ v3.2.1
web_1  | High performance, minimalist Go web framework
web_1  | https://echo.labstack.com
web_1  | ____________________________________O/_______
web_1  |                                     O\
web_1  | ⇨ http server started on [::]:8080
```

#### verify
In another terminal

```sh
$ curl http://localhost:8080/health
OK%
```

### Locally
#### build `web/`
```sh
$ (cd web/ && yarn install && yarn build)
```

#### build `go` server
```sh
$ make init && make build && make test && make run
```
