[![CircleCI](https://circleci.com/gh/OpenBankingUK/tpp-reference-server.svg?style=svg)](https://circleci.com/gh/OpenBankingUK/tpp-reference-server)

# TPP Reference Server

This application simulates a typical Third Party Provider (TPP) backend server.
Its primary function is to provide Open Banking processes to a client.

The implementation uses
[Node.js](https://nodejs.org/), and [express](https://github.com/expressjs/express).

## Use latest release

Use the latest release [v0.8.0](https://github.com/OpenBankingUK/tpp-reference-server/releases/tag/v0.8.0).

To obtain the latest release:

```sh
git clone https://github.com/OpenBankingUK/tpp-reference-server
git checkout v0.8.0
```

Note: latest `master` branch code is actively under development and may not be stable.

## Table of contents

* [TPP Reference Server](#tpp-reference-server)
  * [Use latest release](#use-latest-release)
  * [Quick start with mocked API](#quick-start-with-mocked-api)
    * [Installation](#installation---for-quick-start-with-mocked-api)
    * [Server setup](#server-setup)
    * [Configure ASPSP Authorisation Servers](#configure-aspsp-authorisation-servers)
    * [Run server](#run-server)
  * [Quick start after Open Banking Directory enrolment](#quick-start-after-open-banking-directory-enrolment)
    * [Installation](#installation---for-after-open-banking-directory-enrolment)
    * [Server setup](#server-setup-1)
    * [Turn on OB Directory access and mTLS](#turn-on-ob-directory-access-and-mtls)
    * [Configure Certs and Keys](#configure-certs-and-keys)
    * [Configure ASPSP Authorisation Servers](#configure-aspsp-authorisation-servers-1)
    * [Run server](#run-server-1)
  * [Use cases](#use-cases)
    * [Authenticating with the server](#authenticating-with-the-server)
    * [List ASPSP Authorisation Servers](#list-aspsp-authorisation-servers)
    * [Basic AISP functionality and consent flow (API v1.1)](#basic-aisp-functionality-and-consent-flow-api-v11)
    * [Basic PISP functionality and consent flow (API v1.1)](#basic-pisp-functionality-and-consent-flow-api-v11)
  * [Explaining environment variables](#explaining-environment-variables)
  * [Deploy to heroku](#deploy-to-heroku)
  * [Testing](#testing)
  * [eslint](#eslint)

## Quick start with mocked API

This assumes you do not have OB Directory access but want to kick the tyres to see what's possible.

Use our [reference mock server](https://github.com/OpenBankingUK/reference-mock-server).
It creates simulated endpoints to showcase what the Read/Write API can provide.

__BEFORE PROCEEDING FURTHER__ install and run the reference mock server [as per these instructions](https://github.com/OpenBankingUK/reference-mock-server).

### Installation - for quick start with mocked API

Below are instructions for installing directly on your local machine.

Alternatively you can [install via Docker containers](./README-DOCKER.md).

#### NodeJS

We assume [NodeJS](https://nodejs.org/en/) ver8.4+ is installed.

On Mac OSX, use instructions here [Installing Node.js Tutorial](https://nodesource.com/blog/installing-nodejs-tutorial-mac-os-x/).

On Linux, use instructions in [How To Install Node.js On Linux](https://www.ostechnix.com/install-node-js-linux/).

On Windows, use instructions provided here [Installing Node.js Tutorial: Windows](https://nodesource.com/blog/installing-nodejs-tutorial-windows/).

#### Redis

On Mac OSX, you can install via [homebrew](https://brew.sh). Then.

```sh
brew install redis
```

On Linux, use instructions in the [Redis Quick Start guide](https://redis.io/topics/quickstart).

On Windows, use instructions provided here [Installing Redis on a Windows Workstation](https://essenceofcode.com/2015/03/18/installing-redis-on-a-windows-workstation/).

Then set the environment variables `REDIS_PORT` and `REDIS_HOST` as per redis instance. Example in [`.env.sample`](https://github.com/OpenBankingUK/tpp-reference-server/blob/master/.env.sample)

#### MongoDB

On Mac OSX, you can install via [homebrew](https://brew.sh). Then

```sh
brew install mongodb
```

On Linux, use instructions in the [Install MongoDB Community Edition on Linux](https://docs.mongodb.com/manual/administration/install-on-linux/).

On Windows, use instructions provided here [Install MongoDB Community Edition on Windows](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows/).

Then set the environment variable `MONGODB_URI` as per your mongodb instance, e.g. `MONGODB_URI=mongodb://localhost:27017/sample-tpp-server`. Example in [`.env.sample`](https://github.com/OpenBankingUK/tpp-reference-server/blob/master/.env.sample)

### Server setup

Go to project root after [cloning the repo](#use-latest-release):

```sh
cd to/cloned/project/root
```

Install npm packages:

```sh
npm install
```

Make a local .env based on our .env.sample:

```sh
cp .env.sample .env
```

### Configure ASPSP Authorisation Servers

__Note:__ NOTHING WORKS IF THIS IS NOT SUCCESSFUL.

Complete the sections below to bootstrap the server with essential ASPSP Authorisation Server data.

#### Adding and Updating ASPSP authorisation servers

Use `updateAuthServersAndOpenIds` to bootstrap or update the list of ASPSP authorisation servers used by the server. These are stored in MongoDB. For each authorisation server the OpenId config is also fetched and stored in the database.

Run:

```sh
DEBUG=debug,log npm run updateAuthServersAndOpenIds
```

#### List available ASPSP authorisation servers

Use the `listAuthServers` script to check that Authorisation Servers are available in the database.

This provides useful info like:
* Are `clientCredentialsPresent` already present for an ASPSP Authorisation Server?
* Is `openIdConfigPresent` already present for an ASPSP Authorisation Server?

Run:

```sh
DEBUG=debug,log npm run listAuthServers
```

Output on terminal is TSV that looks like this:
```
id                 CustomerFriendlyName OrganisationCommonName Authority  OBOrganisationId clientCredentialsPresent openIdConfigPresent
aaaj4NmBD8lQxmLh2O AAA Example Bank     AAA Example PLC        GB:FCA:123 aaax5nTR33811Qy  false                    true
bbbX7tUB4fPIYB0k1m BBB Example Bank     BBB Example PLC        GB:FCA:456 bbbUB4fPIYB0k1m  false                    true
cccbN8iAsMh74sOXhk CCC Example Bank     CCC Example PLC        GB:FCA:789 cccMh74sOXhkQfi  false                    true
```

#### Add Client Credentials for ASPSP Authorisation Servers

There is a script to input and store client credentials against ASPSP Auth Server configuration.

Run:

```sh
DEBUG=debug,log npm run saveCreds authServerId=aaaj4NmBD8lQxmLh2O clientId=spoofClientId clientSecret=spoofClientSecret

DEBUG=debug,log npm run saveCreds authServerId=bbbX7tUB4fPIYB0k1m clientId=spoofClientId clientSecret=spoofClientSecret

DEBUG=debug,log npm run saveCreds authServerId=cccbN8iAsMh74sOXhk clientId=spoofClientId clientSecret=spoofClientSecret
```

### Run server

Run using foreman, this will pick ENVs from the `.env` file [setup earlier](#server-setup):

```sh
npm run foreman
# [OKAY] Loaded ENV .env File as KEY=VALUE Format
# web.1 | log App listening on port 8003 ...
```

Assuming the TPP server is now running, install and run the [TPP Reference Client](https://github.com/OpenBankingUK/tpp-reference-client) to view accounts data or make a single payment.

Alternatively, check the [supported use cases](#use-cases) to issue `CURL` commands and explore features.

## Quick start after Open Banking Directory enrolment

This assumes you have enrolled successfully with OB Directory and have access to all necessary credentials, CERTS and KEYS.

Check [here for more info on how to complete this](https://www.openbanking.org.uk/directory/).

### Installation - for after Open Banking Directory enrolment

#### NodeJS

We assume [NodeJS](https://nodejs.org/en/) ver8.4+ is installed.

On Mac OSX, use instructions here [Installing Node.js Tutorial](https://nodesource.com/blog/installing-nodejs-tutorial-mac-os-x/).

On Linux, use instructions in [How To Install Node.js On Linux](https://www.ostechnix.com/install-node-js-linux/).

On Windows, use instructions provided here [Installing Node.js Tutorial: Windows](https://nodesource.com/blog/installing-nodejs-tutorial-windows/).

#### Redis

On Mac OSX, you can install via [homebrew](https://brew.sh). Then.

```sh
brew install redis
```

On Linux, use instructions in the [Redis Quick Start guide](https://redis.io/topics/quickstart).

On Windows, use instructions provided here [Installing Redis on a Windows Workstation](https://essenceofcode.com/2015/03/18/installing-redis-on-a-windows-workstation/).

Then set the environment variables `REDIS_PORT` and `REDIS_HOST` as per redis instance. Example in [`.env.sample`](https://github.com/OpenBankingUK/tpp-reference-server/blob/master/.env.sample)

#### MongoDB

On Mac OSX, you can install via [homebrew](https://brew.sh). Then

```sh
brew install mongodb
```

On Linux, use instructions in the [Install MongoDB Community Edition on Linux](https://docs.mongodb.com/manual/administration/install-on-linux/).

On Windows, use instructions provided here [Install MongoDB Community Edition on Windows](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-windows/).

Then set the environment variable `MONGODB_URI` as per your mongodb instance, e.g. `MONGODB_URI=mongodb://localhost:27017/sample-tpp-server`. Example in [`.env.sample`](https://github.com/OpenBankingUK/tpp-reference-server/blob/master/.env.sample)

### Server setup

Go to project root after [cloning the repo](#use-latest-release):

```sh
cd to/cloned/project/root
```

Install npm packages:

```sh
npm install
```

Make a local .env based on our .env.sample:

```sh
cp .env.sample .env
```

### Turn on OB Directory access and mTLS

#### OB Directory access

Update the following ENVs in your `.env` file:
* `OB_PROVISIONED=true`
* `OB_DIRECTORY_AUTH_HOST=<enter as per OB Directory instructions>`
* `OB_DIRECTORY_HOST=<enter as per OB Directory instructions>`
* `SIGNING_KID=<enter as per OB Directory instructions>`
* `SOFTWARE_STATEMENT_ID=<enter as per OB Directory instructions>`

#### mTLS access

MTLS is enabled by default. Though it is disabled when uri contains
domain "localhost" or "reference-mock-server".

### Configure Certs and Keys

First ensure:

* You have downloaded the required `Transport` and `Signing` Certs (follow OB
   Directory issued instructions).
* You have access to the `private key` used when generating the `Signing` Cert CSR.
* You have access to the `private key` used when generating the `Transport` Cert CSR.

ENVs to be configured with Certs and Keys require these values as base64 encoded strings. This applies to

* `SIGNING_KEY` - private key used to generate `Signing` cert OB Directory CSR.
* `OB_ISSUING_CA` - Downloaded `OB Issuing chain` cert from OB Directory.
* `TRANSPORT_CERT` - Downloaded `Transport` cert from OB Directory console.
* `TRANSPORT_KEY` - private key used to generate `Transport` cert OB Directory CSR.

#### Encode

Encode using our `base64-cert-or-key` script.

```
npm run base64-cert-or-key <path/to/cert> or <path/to/key>
```

This produces something similar to ...
```
BASE64 ENCODING COMPLETE (Please copy the text below to the required ENV):

LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZxekNDQkpPZ0F3SUJBZ0lFV1d2OG5UQU5CZ2...
```

As per instructions in the output: __copy the base64 encoded string to the relevant ENV in your `.env` file.__

### Configure ASPSP Authorisation Servers

__Note:__ NOTHING WORKS IF THIS IS NOT SUCCESSFUL.

Complete the sections below to bootstrap the server with essential ASPSP Authorisation Server data.

#### Adding and Updating ASPSP authorisation servers

Use `updateAuthServersAndOpenIds` to bootstrap or update the list of ASPSP authorisation servers used by the server. These are stored in MongoDB. For each authorisation server the OpenId config is also fetched and stored in the database.

Run:

```sh
DEBUG=debug,log,error npm run updateAuthServersAndOpenIds
```

#### List available ASPSP authorisation servers

Use the `listAuthServers` script to check that Authorisation Servers are available in the database.

This provides useful info like:
* Are `clientCredentialsPresent` already present for an ASPSP Authorisation Server?
* Is `openIdConfigPresent` already present for an ASPSP Authorisation Server?

Run:

```sh
DEBUG=debug,log,error npm run listAuthServers
```

Output on terminal is TSV that looks like this:
```
id                 CustomerFriendlyName OrganisationCommonName Authority  OBOrganisationId clientCredentialsPresent openIdConfigPresent
aaaj4NmBD8lQxmLh2O AAA Example Bank     AAA Example PLC        GB:FCA:123 aaax5nTR33811Qy  false                    true
bbbX7tUB4fPIYB0k1m BBB Example Bank     BBB Example PLC        GB:FCA:456 bbbUB4fPIYB0k1m  false                    true
cccbN8iAsMh74sOXhk CCC Example Bank     CCC Example PLC        GB:FCA:789 cccMh74sOXhkQfi  false                    true
```

#### Add Client Credentials for ASPSP Authorisation Servers

There is a script to input and store client credentials against ASPSP Auth Server configuration.

Please consult Open Banking documentation on how to acquire a `clientId` and `clientSecret` for an ASPSP Authorisation Server.

Run:

```sh
DEBUG=debug,log npm run saveCreds authServerId=<server_id> clientId=<client_id> clientSecret=<client_secret>
```

### Run server

Run using foreman, this will pick ENVs from the `.env` file [setup earlier](#server-setup-1):

```sh
npm run foreman
# [OKAY] Loaded ENV .env File as KEY=VALUE Format
# web.1 | log App listening on port 8003 ...
```

Assuming the TPP server is now running, install and run the [TPP Reference Client](https://github.com/OpenBankingUK/tpp-reference-client) to view accounts data or make a single payment.

Alternatively, check the [supported use cases](#use-cases) to issue `CURL` commands and explore features.

## Use cases

__Work in progress__ - so far we provide,
* Authenticating with the server.
* List ASPSP Authorisation Servers - actual & simulated based on ENVs.
* Basic AISP functionality and consent flow.
* Basic PISP functionality and consent flow.

### Authenticating with the server

#### Login

```sh
curl -X POST --data 'u=alice&p=wonderland' http://localhost:8003/login
```

This returns a session ID token as a `sid`. Use this for further authorized access.

This is an example.

```sh
{
  "sid": "896beb20-affc-11e7-a5e6-a941c8c37252"
}
```

#### Logout

Please __change__ the `Authorization` header to use the `sid` obtained after login.

This destroys the session established by the `sid` token obtained after login.

```sh
curl -X GET -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' http://localhost:8003/logout
```

### List ASPSP Authorisation Servers

Please __change__ the `Authorization` header to use the `sid` obtained after logging inT.

```sh
curl -X GET -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' http://localhost:8003/account-payment-service-provider-authorisation-servers
```

Here's a sample list of test ASPSPs. This is __NOT__ the raw response from the Open Banking Directory. It has been adapted to simulate what a typical client app would require.

```sh
[
  {
    "id": "aaaj4NmBD8lQxmLh2O",
    "logoUri": "",
    "name": "AAA Example Bank",
  },
  {
    "id": "bbbX7tUB4fPIYB0k1m",
    "logoUri": "",
    "name": "BBB Example Bank",
  },
  {
    "id": "cccbN8iAsMh74sOXhk",
    "logoUri": "",
    "name": "CCC Example Bank",
  }
]
```

### Basic AISP functionality and consent flow (API v1.1)

We support a simple AISP workflow where a PSU authorises a TPP to view account information on their behalf. This showcases the required oAuth consent flow and hits the relevant [proxied APIs](#proxied-api-path).

#### OIDC Authorization Flow

We implement the OIDC Authorization Flow which generates an `access-token` on the PSUs behalf to allow the TPP to access their accounts.

This `access-token`
* Is stored on the ASPSP Authorization Server.
* Post generation marks the PSU's account as `Authorised` on the ASPSP Resource Server.

Here's how to start the flow,

```sh
curl -X POST -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' -H 'x-authorization-server-id: aaaj4NmBD8lQxmLh2O' -H 'Accept: application/json' -d '{"authorisationServerId": "aaaj4NmBD8lQxmLh2O"}' http://localhost:8003/account-request-authorise-consent
```

This will yield a URI required to perform the oAuth flow. Sample below:

```json
{
  "uri": "http://localhost:8001/aaaj4NmBD8lQxmLh2O/authorize?redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Ftpp%2Fauthorized&state=eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiYzEzOWI5M2UtYWU1NC00YzI0LWIzNjEtZWUyODZlOWRjYTBlIiwic2Vzc2lvbklkIjoiNWFjMTEzMzAtZGY1MC0xMWU3LTgzZWEtZDE2NzFhOTU5ZGUwIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMifQ%3D%3D&client_id=spoofClientId&response_type=code&request=eyJhbGciOiJub25lIn0.eyJpc3MiOiJzcG9vZkNsaWVudElkIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJjbGllbnRfaWQiOiJzcG9vZkNsaWVudElkIiwicmVkaXJlY3RfdXJpIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3RwcC9hdXRob3JpemVkIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMiLCJzdGF0ZSI6ImV5SmhkWFJvYjNKcGMyRjBhVzl1VTJWeWRtVnlTV1FpT2lKaFlXRnFORTV0UWtRNGJGRjRiVXhvTWs4aUxDSnBiblJsY21GamRHbHZia2xrSWpvaVl6RXpPV0k1TTJVdFlXVTFOQzAwWXpJMExXSXpOakV0WldVeU9EWmxPV1JqWVRCbElpd2ljMlZ6YzJsdmJrbGtJam9pTldGak1URXpNekF0WkdZMU1DMHhNV1UzTFRnelpXRXRaREUyTnpGaE9UVTVaR1V3SWl3aWMyTnZjR1VpT2lKdmNHVnVhV1FnWVdOamIzVnVkSE1pZlE9PSIsIm5vbmNlIjoiZHVtbXktbm9uY2UiLCJtYXhfYWdlIjo4NjQwMCwiY2xhaW1zIjp7InVzZXJpbmZvIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJ2YWx1ZSI6IjAxMmZhZGY1LWI4NWYtNDZjNS04MTc1LTk4NTRkMzZkZjVmYSIsImVzc2VudGlhbCI6dHJ1ZX19LCJpZF90b2tlbiI6eyJvcGVuYmFua2luZ19pbnRlbnRfaWQiOnsidmFsdWUiOiIwMTJmYWRmNS1iODVmLTQ2YzUtODE3NS05ODU0ZDM2ZGY1ZmEiLCJlc3NlbnRpYWwiOnRydWV9LCJhY3IiOnsiZXNzZW50aWFsIjp0cnVlfX19fQ.&scope=openid%20accounts"
}
```

A `redirect_uri` is included in the querystring. This will be used by the ASPSP server to redirect back to the intended endpoint.

Perform a GET request against the `uri` in the payload to continue the flow,

```sh
curl -X GET --url "http://localhost:8001/aaaj4NmBD8lQxmLh2O/authorize?redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Ftpp%2Fauthorized&state=eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiOGEyYzk1NzItMWYxZi00MDdhLTk1MjYtNWY4MzRlN2ZjMWFjIiwic2Vzc2lvbklkIjoiYzc1OGIwMDAtZGY1Mi0xMWU3LTgzZWEtZDE2NzFhOTU5ZGUwIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMifQ%3D%3D&client_id=spoofClientId&response_type=code&request=eyJhbGciOiJub25lIn0.eyJpc3MiOiJzcG9vZkNsaWVudElkIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJjbGllbnRfaWQiOiJzcG9vZkNsaWVudElkIiwicmVkaXJlY3RfdXJpIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3RwcC9hdXRob3JpemVkIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMiLCJzdGF0ZSI6ImV5SmhkWFJvYjNKcGMyRjBhVzl1VTJWeWRtVnlTV1FpT2lKaFlXRnFORTV0UWtRNGJGRjRiVXhvTWs4aUxDSnBiblJsY21GamRHbHZia2xrSWpvaU9HRXlZemsxTnpJdE1XWXhaaTAwTURkaExUazFNall0TldZNE16UmxOMlpqTVdGaklpd2ljMlZ6YzJsdmJrbGtJam9pWXpjMU9HSXdNREF0WkdZMU1pMHhNV1UzTFRnelpXRXRaREUyTnpGaE9UVTVaR1V3SWl3aWMyTnZjR1VpT2lKdmNHVnVhV1FnWVdOamIzVnVkSE1pZlE9PSIsIm5vbmNlIjoiZHVtbXktbm9uY2UiLCJtYXhfYWdlIjo4NjQwMCwiY2xhaW1zIjp7InVzZXJpbmZvIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJ2YWx1ZSI6ImY2MjRkYWVkLThkYWMtNGExOS1hYmU1LWNlMjgwNWYzNDliOSIsImVzc2VudGlhbCI6dHJ1ZX19LCJpZF90b2tlbiI6eyJvcGVuYmFua2luZ19pbnRlbnRfaWQiOnsidmFsdWUiOiJmNjI0ZGFlZC04ZGFjLTRhMTktYWJlNS1jZTI4MDVmMzQ5YjkiLCJlc3NlbnRpYWwiOnRydWV9LCJhY3IiOnsiZXNzZW50aWFsIjp0cnVlfX19fQ.&scope=openid%20accounts"
```

This results in a redirection url with the correct `code` as per the OIDC Authorization Flow. Also, it includes a `state` query param that includes a base64 encoded string to identify the ASPSP Authorization Server for further queries.

```sh
Found. Redirecting to http://localhost:8080/tpp/authorized?code=spoofAuthorisationCode&state=eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiOGEyYzk1NzItMWYxZi00MDdhLTk1MjYtNWY4MzRlN2ZjMWFjIiwic2Vzc2lvbklkIjoiYzc1OGIwMDAtZGY1Mi0xMWU3LTgzZWEtZDE2NzFhOTU5ZGUwIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMifQ==
```

To conclude the flow and ensure the `access-token` is generated,
* First, parse the `state` query param using the [Node.js `REPL`](https://nodejs.org/api/repl.html).

```sh
node
# Node.js REPL
new Buffer("eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiOGEyYzk1NzItMWYxZi00MDdhLTk1MjYtNWY4MzRlN2ZjMWFjIiwic2Vzc2lvbklkIjoiYzc1OGIwMDAtZGY1Mi0xMWU3LTgzZWEtZDE2NzFhOTU5ZGUwIiwic2NvcGUiOiJvcGVuaWQgYWNjb3VudHMifQ==", 'base64').toString('ascii');
```

This produces
```
'{"authorisationServerId":"aaaj4NmBD8lQxmLh2O","interactionId":"8a2c9572-1f1f-407a-9526-5f834e7fc1ac","sessionId":"c758b000-df52-11e7-83ea-d1671a959de0","scope":"openid accounts"}'
```

* Close the `REPL`:

```sh
# Node.js REPL
.exit
```

* Then, perform a POST request against the TPP Server using `code` and parsed `authorisationServerId` from the `state`.

```sh
curl -X POST -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' -H 'Accept: application/json' -d '{"authorisationServerId": "aaaj4NmBD8lQxmLh2O", "interactionId": "8a2c9572-1f1f-407a-9526-5f834e7fc1ac", "sessionId": "c758b000-df52-11e7-83ea-d1671a959de0", "scope": "accounts", "authorisationCode": "spoofAuthorisationCode"}' http://localhost:8003/tpp/authorized
```

This creates an `access-token` and allows authorized access to the APIs.

#### Proxied API path

To interact with proxied Open Banking Read/Write APIs please use the path `/open-banking/[api_version]` in all requests.

For example `/open-banking/v1.1` gives access to the 1.1 Read write Apis.

#### GET Accounts for a user (Account and Transaction API)

We have a hardcoded demo user `alice` in [mock server](https://github.com/OpenBankingUK/reference-mock-server). To
access demo accounts for this user please setup the following `ENVS` (already configured in [`.env.sample`](https://github.com/OpenBankingUK/tpp-reference-server/blob/master/.env.sample).


There is a second hardcoded user `kate` with password `lookingglass`, for the purpose of testing one user with consent and one without.

```sh
curl -X GET -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' -H 'x-authorization-server-id: aaaj4NmBD8lQxmLh2O' -H 'Accept: application/json'  http://localhost:8003/open-banking/v1.1/accounts
```

[Here's a sample response](https://www.openbanking.org.uk/read-write-apis/account-transaction-api/v1-1-0/#accounts-bulk-response).

#### GET Product associated with an account (Account and Transaction API)

Using the same demo account as above.

```sh
curl -X GET -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' -H 'x-authorization-server-id: aaaj4NmBD8lQxmLh2O' -H 'Accept: application/json'  http://localhost:8003/open-banking/v1.1/accounts/22289/product
```

[Here's a sample response](https://www.openbanking.org.uk/read-write-apis/account-transaction-api/v1-1-0/#product-specific-account-response).

#### GET Balances associated with an account (Account and Transaction API)

Using the same demo account as above.

```sh
curl -X GET -H 'Authorization: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' -H 'x-authorization-server-id: aaaj4NmBD8lQxmLh2O' -H 'Accept: application/json'  http://localhost:8003/open-banking/v1.1/accounts/22289/balances
```

[Here's a sample response](https://www.openbanking.org.uk/read-write-apis/account-transaction-api/v1-1-0/#balances-specific-account-response).


### Basic PISP functionality and consent flow (API v1.1)

We support a simple PISP workflow where the PSU authorises a TPP to initialise a payment from an ASPSP to a third party.  The current use case with v1.1 is a Single Immediate Payment.

There are 5 steps in the Single Immediate Payment flow

1) Request payment initiation (PSU > TPP)
2) Setup single payment initiation (TPP > ASPSP)
3) Authorise consent (PSU > ASPSP)
4) Create payment submission (TPP > ASPSP)
5) Check Payment Status (TPP <> ASPSP)


#### Step 1 Request Payment Initiation

The PSU needs to be logged into the TPP Server (Ref App) and get a session to do anything.
Without the TPP Ref Client We can simulate this with a CURL request like this:

```sh
curl -X POST --data 'u=alice&p=wonderland' http://localhost:8003/login
```

This should yield a payload with a Session ID like this
```sh
{
	"sid": "eefbda80-ec93-11e7-a1f6-0f0979b77e2b"
}
```

We kick off the payment initiation with the Client calling to the TPP Server's
`/payment-authorise-consent` endpoint with a payload like that can be simulated with a CURL command like this

```sh
curl -X POST -H 'Authorization: eefbda80-ec93-11e7-a1f6-0f0979b77e2b' -H 'x-authorization-server-id: aaaj4NmBD8lQxmLh2O' -H 'Accept: application/json' -H "Content-Type: application/json" -d '{"authorisationServerId":"aaaj4NmBD8lQxmLh2O","InstructedAmount":{"Amount":"10.00","Currency":"GBP"},"CreditorAccount":{"SchemeName":"SortCodeAccountNumber","Identification":"11111112345678","Name":"Sam Morse"}}' http://localhost:8003/payment-authorise-consent
```

__Note:__ The Authorisation header is the `sid` value we obtained earlier.

This kicks off a train of events (Step 2) which - if the PSU authorises the payment - gives back
a Redirect URL, the payload of which also contains an embedded redirect URL (see later).

#### Step 2 - Setup single payment initiation

In overview:
1) The TPP Server calls out to the `/token/` endpoint at the ASPSP Auth server to
get an `access-token`
2) The TPP server calls out to the `/payments/` endpoint using the `access-token`, and sends payment details
3) ... The ASPSP creates a *Payment Resource* - with a `PaymentId` - bound to the ClientId (TPP detail stored in Directory)
4) ... the ASPSP sends back a 201 response to the TPP Server with the `PaymentId`
5) The TPP Server stores the `PaymentId` for later use (bound to the PSU Session and / or state - see later)
6) TPP Generates a Signed Request Object with requested claims (including PaymentId)
7) TPP Server generates a Redirct URL (embedded in a JSON Object) with various parameters relevant to the Payment Resource
and sends this to the Client.


From 2 above (`payments` endpoint) - [Here's a sample response](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5786479/Payment+Initiation+API+Specification+-+v1.1.0#PaymentInitiationAPISpecification-v1.1.0-POST/paymentresponse)


From 7 above - Example Redirect URL JSON object:

```sh
{"uri":"http://localhost:8001/aaaj4NmBD8lQxmLh2O/authorize?redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Ftpp%2Fauthorized&state=eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiYTM1Y2QzNGQtYzA3Yi00MWZhLWJjZGQtYjc5YTQ5NGE4NTE4Iiwic2Vzc2lvbklkIjoiZDllZTJmNzAtZWM5NC0xMWU3LWExZjYtMGYwOTc5Yjc3ZTJiIiwic2NvcGUiOiJvcGVuaWQgcGF5bWVudHMifQ%3D%3D&client_id=spoofClientId&response_type=code&request=eyJhbGciOiJub25lIn0.eyJpc3MiOiJzcG9vZkNsaWVudElkIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJjbGllbnRfaWQiOiJzcG9vZkNsaWVudElkIiwicmVkaXJlY3RfdXJpIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3RwcC9hdXRob3JpemVkIiwic2NvcGUiOiJvcGVuaWQgcGF5bWVudHMiLCJzdGF0ZSI6ImV5SmhkWFJvYjNKcGMyRjBhVzl1VTJWeWRtVnlTV1FpT2lKaFlXRnFORTV0UWtRNGJGRjRiVXhvTWs4aUxDSnBiblJsY21GamRHbHZia2xrSWpvaVlUTTFZMlF6TkdRdFl6QTNZaTAwTVdaaExXSmpaR1F0WWpjNVlUUTVOR0U0TlRFNElpd2ljMlZ6YzJsdmJrbGtJam9pWkRsbFpUSm1OekF0WldNNU5DMHhNV1UzTFdFeFpqWXRNR1l3T1RjNVlqYzNaVEppSWl3aWMyTnZjR1VpT2lKdmNHVnVhV1FnY0dGNWJXVnVkSE1pZlE9PSIsIm5vbmNlIjoiZHVtbXktbm9uY2UiLCJtYXhfYWdlIjo4NjQwMCwiY2xhaW1zIjp7InVzZXJpbmZvIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJ2YWx1ZSI6IjVmMGNiYjAxLTQzOTctNDhmZi04MDE3LTQ3OTA4YmU0NWNlYiIsImVzc2VudGlhbCI6dHJ1ZX19LCJpZF90b2tlbiI6eyJvcGVuYmFua2luZ19pbnRlbnRfaWQiOnsidmFsdWUiOiI1ZjBjYmIwMS00Mzk3LTQ4ZmYtODAxNy00NzkwOGJlNDVjZWIiLCJlc3NlbnRpYWwiOnRydWV9LCJhY3IiOnsiZXNzZW50aWFsIjp0cnVlfX19fQ.&scope=openid%20payments"}
```

This URL redirects to the ASPSP for the PSU to log in to their ASPSP
and give consent for the payment to be made.

The embedded Redirect URL mentioned above is the "Software Statement Redirect URL",
and it's the URL that the ASPSP redirects the
PSU client back to once the PSU has granted consent.

#### Step 3 - Authorise Consent

This is in the realm of the ASPSP.  Assuming the PSU gives consent the ASPSP will redirect back to
the embedded redirect URL found in the Software Statement.
In the case of the TPP Reference App this URL is `http://localhost:8080/tpp/authorized`

The TPP Client code picks up the redirected URL which contains two parameters:

`code, state`

Here's an example URL with these params:

```sh
http://localhost:8080/tpp/authorized?code=spoofAuthorisationCode&state=eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiYTM1Y2QzNGQtYzA3Yi00MWZhLWJjZGQtYjc5YTQ5NGE4NTE4Iiwic2Vzc2lvbklkIjoiZDllZTJmNzAtZWM5NC0xMWU3LWExZjYtMGYwOTc5Yjc3ZTJiIiwic2NvcGUiOiJvcGVuaWQgcGF5bWVudHMifQ%3D%3D
```

`code` is the `authorization-code` from the ASPSP Auth Server, which the TPP Server uses
 to call out to the `/token` endpoint again to swap for a new `access-token` which will be used in
 Step 4 - payment submission.


At this point the state payload can be inspected using the Node REPL like this:

```sh
new Buffer("eyJhdXRob3Jpc2F0aW9uU2VydmVySWQiOiJhYWFqNE5tQkQ4bFF4bUxoMk8iLCJpbnRlcmFjdGlvbklkIjoiYTM1Y2QzNGQtYzA3Yi00MWZhLWJjZGQtYjc5YTQ5NGE4NTE4Iiwic2Vzc2lvbklkIjoiZDllZTJmNzAtZWM5NC0xMWU3LWExZjYtMGYwOTc5Yjc3ZTJiIiwic2NvcGUiOiJvcGVuaWQgcGF5bWVudHMifQ==", 'base64').toString('ascii');
```

This results in (example only, and formatted for ease of reading):

```sh
{
	"authorisationServerId": "aaaj4NmBD8lQxmLh2O",
	"interactionId": "a35cd34d-c07b-41fa-bcdd-b79a494a8518",
	"sessionId": "eefbda80-ec93-11e7-a1f6-0f0979b77e2b",
	"scope": "openid payments"
}
```


#### Step 4 - Create Payment Submission

The TPP Server retrieves the `PaymentId` (and various other parameters) it saved earlier.

TPP Calls out to the `/payment-submissions` endpoint using the stored payment initiation details and new access token.

ASPSP Responds with a HTTP Status 201 and a `PaymentSubmissionId`

[Here's a sample response](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5786479/Payment+Initiation+API+Specification+-+v1.1.0#PaymentInitiationAPISpecification-v1.1.0-POST/payment-submissionsResponse)

#### Step 5 - Get payment submission status

Depending on how old the access-toke is we can either use the existing one
OR kick off a new CLient Credentials Grant at the token endpoint to get a new one

The call out to the `/payment-submissions` endpoint with a `GET` and the `PaymentSubmissionId`

[Here's an example response](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5786479/Payment+Initiation+API+Specification+-+v1.1.0#PaymentInitiationAPISpecification-v1.1.0-GET/payment-submissionsrequest)

## Explaining environment variables

* `CLIENT_SCOPES='ASPSPReadAccess TPPReadAccess AuthoritiesReadAccess'` - access scopes for API features.
* `DEBUG=error,log,debug` - enables debugging.
* `MONGODB_URI=mongodb://localhost:27017/sample-tpp-server` - MongoDB configuration.
* `OB_DIRECTORY_AUTH_HOST=xxx` - OB Directory Auth host for issuing tokens.
* `OB_DIRECTORY_HOST=xxx` - OB Directory host.
* `OB_ISSUING_CA=''` - base64 encoded `OB Issuing chain` cert from OB Directory.
* `OB_PROVISIONED=false` - enables / disables enrolled with OB Directory mode.
* `TPP_REF_SERVER_PORT=8003` - port where the app is running.
* `REDIS_HOST=localhost` - Redis configuration.
* `REDIS_PORT=6379` - Redis port.
* `SIGNING_KEY=''` - base64 encoded private key used to generate `Signing` cert OB Directory CSR.
* `SIGNING_KID=xxx` - kid retrieved from OB Directory console.
* `SOFTWARE_STATEMENT_ID=xxx` - softwareStatementId of software statement generated using OB Directory console.
* `TRANSPORT_CERT=''` - base64 encoded downloaded `Transport` cert from OB Directory console.
* `TRANSPORT_KEY=''` - base64 encoded private key used to generate `Transport` cert OB Directory CSR.

## Deploy to heroku

You can follow these [instructions to deploy to heroku](./README-HEROKU.md).

## Validating API responses

You can follow these [instructions to configure optional validation of API responses](./README_VALIDATION.md).

## Testing

Run unit tests with:

```sh
npm run test
```

Run tests continuously on file changes in watch mode via:

```sh
npm run test:watch
```

Manual Testing  
Sending Form Data to login with POstman - use `x-www-form-urlencoded`

## eslint

Run eslint checks with:

```sh
npm run eslint
```
