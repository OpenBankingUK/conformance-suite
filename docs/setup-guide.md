# Functional Conformance Tool and Ozone Model Bank

This guide will assist you with the technical steps required to setup the Functional Conformance Tool and run your first test. In this guide we will be connecting to and running tests against the Ozone Model Bank.

Please note the following goals of this document:

* Register with Ozone model bank as a Third Party Provider (TPP).
* Generate certificates for transport and signing of requests to Ozone - self signed certs to be used.
* Setup the FCS to run on local machine (via Docker container).
* Execute a test on the FCS using the Payment Service User (PSU) consent flow. The PSU consent flow involves user
interaction during the authentication and consent steps when connecting to Ozone.
* Review the results of the test run.

## Prerequisites

This guide assumes the following tools are installed and functioning correctly. Versions specified used when writing this guide.

* Docker (Client: 18.09.1, Server: 18.09.1 on OSX)
* OpenSSL (LibreSSL 2.6.4 on OSX)
* Google login if using Ozone "self-serve"
* Access to the following hosts from your computer - See Appendix A.

*Note for Windows 10 users - Docker on Windows 10 requires Hyper-V to be installed. Hyper-V is only available
on Pro or Enterprise versions. Please refer to [this guide](https://techcommunity.microsoft.com/t5/ITOps-Talk-Blog/Step-By-Step-Enabling-Hyper-V-for-use-on-Windows-10/ba-p/267945) for more information.*

## Step 1: Register a client with Ozone Bank (Model Bank)

Ozone Bank is an Mock Account Servicing Payment Service Provider (ASPSP), which the FCS will connect to as a TPP.

In order to register a new client you need to use [Dynamic Client Registration Protocol](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1078034771/Dynamic+Client+Registration+-+v3.2). Please follow the steps described in this page to register a new client with Ozone [https://openbanking.atlassian.net/wiki/spaces/DZ/pages/313918598/Integrating+a+TPP+with+Ozone+Model+Banks+Using+Postman+on+Directory+Sandbox](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/313918598/Integrating+a+TPP+with+Ozone+Model+Banks+Using+Postman+on+Directory+Sandbox) 

Once completed, make a note of the certificates and the `CLIENT ID` and `CLIENT SECRET` values.

### Generate transport and signing certificates

Mutually Authenticated TLS (MA-TLS) is required to communicate with Ozone model bank service.

In the interest of expediency, Ozone allows developers to create their own certificates, provided that they set the appropriate fields in accordance with
what is specified on the "Certificate" tab of your registration page on the Ozone site.

_Please note: This is a very
streamlined approach and in normal circumstances, you would typically be required to register with the Open Banking Directory service
prior to engagement with a bank / ASPSP, as per the Open Banking model._

The following three sections will generate various files. As an output, you will require the following files to proceed:

* signing.key (private key)
* signing.pem (certificate)
* transport.key
* transport.pem

#### Precursor: CA Generation

A precursor step of creating a Certificate Authority (CA) is required to sign the transport and signing certs.

`openssl req -new -x509 -days 3650 -keyout ca.key -out ca.pem -nodes`

The following certificate information (DN) values shall be provided, based on what is mentioned on your Ozone registration page
under the Certificates tab:

* C (Country Name) = GB
* O (Organisation Name) = Ozone Financial Technology Limited
* OU (Organisation Unit)
* CN (Common Name)

#### Transport Certificate

* Execute `openssl genrsa -out transport.key 2048`

* Execute `openssl req -new -sha256 -key transport.key -out transport.csr` (Provide same DN information as above, no passphrase)

* Execute `openssl x509 -req -days 3650 -in transport.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out transport.pem`

#### Signing Certificate

* Execute `openssl genrsa -out signing.key 2048`

* Execute `openssl req -new -sha256 -key signing.key -out signing.csr` (Provide same DN information as above, no passphrase)

* Execute `openssl x509 -req -days 3650 -in signing.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out signing.pem`

## Step 2: Add Functional Conformance Suite Server Certificates

The suite runs on https using localhost, you can trust the certificate or add as an exception.

Certificates and be downloaded [https://bitbucket.org/openbankingteam/conformance-suite/src/develop/certs/](here).

## Step 3: Download the Functional Conformance Suite

The following command will download the latest FCS image from docker hub and run it. You may need to login to Docker Hub
at this point by running `docker login`.

### Production

```sh
docker run --rm -it -p 8443:8443 -e LOG_LEVEL=debug -e LOG_TRACER=true -e LOG_HTTP_TRACE=true "openbanking/conformance-suite:v1.6.3"
```

*Please substitute the container version (1.6.3 in this example) with the version of the latest stable release.*

### Production - JWS signatures disabled (for testing purposes only)

```sh
docker run --rm -it -p 8443:8443 -e LOG_LEVEL=debug -e LOG_TRACER=true -e LOG_HTTP_TRACE=true -e DISABLE_JWS=TRUE "openbanking/conformance-suite:v1.6.3"
```

### Non-production run

```sh
docker run --rm -it -p 8443:8443 -e LOG_LEVEL=debug -e LOG_TRACER=true -e LOG_HTTP_TRACE=true -e DISABLE_JWS=TRUE "openbanking/conformance-suite:latest"
```

If all goes well you should be able to launch the FCS UI from you browser via `https://0.0.0.0:8443`

### Dynamic Resource Allocation

Setting the environment variable:

`DYNRES=true`

enables Dynamic Resource Allocation for the Account and Transaction Apis for Account Numbers.
When this flag is active, for each token obtained, an additional call is made to the /accounts endpoint. This call retrieves a list of account numbers associated with the consent.
The first account number in this list is used to as the account resource id for the tests associated with the specific token.

This is a new feature, and as such will rely on feedback from ASPSPs to align with variations in Dynamic Resource Allocation implementations.

### Optional - Docker Content Trust (recommended)

Docker Content Trust *(DCT)* ensures that all content is securely received and verified. Open Banking cryptographically signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flaG. For example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-suite:TAG
    DOCKER_CONTENT_TRUST=1 docker RUN openbanking/conformance-suite:TAG

Alternatively, you can set DCT with an environmental variable. `export DOCKER_CONTENT_TRUST=1`

Once DCT is enabled remote trust is checked on every pull request. If no trust data for a tag is found you will be presented with an error.

    export DOCKER_CONTENT_TRUST=1
    docker pull openbanking/conformance-suite:TAG
    Error: remote trust data does not exist for docker.io/openbanking/conformance-suite: notary.docker.io does not have trust data for docker.io/openbanking/conformance-suite

[More on Docker Content Trust](docs/docker_content_trust.md)

## Step 4: Config & Run the Functional Conformance Suite

Running a test plan on the FCS involves five steps, as follows:

### Start/Load test - To start a new test select the Ozone PSU template.

### Discovery - Review the discovery file and update as required.

Select the Ozone PSU template.

### Configuration
* Client
	* Provide the keys, as created earlier signing and transport (first 4 boxes) 
	* Account ID: `500000000000000000000001`
	* Statement ID: `140000000000000000000001`
	* Transaction _from_/_to_ dates: _pre-populated value ok_
	* Enter the _client ID_ and _secret_ from Ozone Bank, generated during on-boarding in Step 1
	* _x-fapi-financial-id_ = `0015800001041RHAAY` (Ozone Bank)
	* Send _x-fapi-customer-ip-address_ is an optional value. Is used to indicate presence of PSU during interaction.
* Well-Known
	* Token Endpoint: _pre-populated value ok_
	* OAuth 2.0 response_type: `code id_token`
	* Token Endpoint Auth Method: `client_secret_basic`
	* Request object signing algorithm: `PS256`
	* Authorization Endpoint: _pre-populated value ok_
	* Resource Base URL: `https://ob19-rs1.o3bank.co.uk:4501`
	* Issuer: _pre-populated value ok_
	* Redirect URL: `https://0.0.0.0:8443/conformancesuite/callback`
	* Use NON OB Directory: false (not checked)
* Payments
	* Frequency: `EvryDay`
	* Creditor Account
		* SchemeName: `UK.OBIE.IBAN`
		* Identification: `GB29PAPA20000390210002`
		* Name: `Mr Jackson`
		* Instructed Amount Value: `1.00`
		* Instructed Amount Currency: `GBP`
		* Currency of transfer for international payments: `USD`

_Please note: If an item has been pre-populated, that is a generally acceptable default, unless specified above or specific tests
are being defined._

4. Run / Overview

    This screen shows the tests that will be run. Once ready, click "Start PSU Consent" in API Specification section. This should load up Ozone PSU authentication page. Provide mits/mits as login name and password.
    On the account selection page that follows, select at least one account and click Confirm button. On the next page, click Yes button to grant consent and see the authorization code page.
    Go back to the FCS Testcases page to select Pending PSU Consent button at the bottom of the page. The tests should run and go to the "PENDING" status. Once complete the status should move to "PASSED", if everything ran ok. If any of the tests failed, you can click the "FAILED" badge to view more information on the cause of failure.

5. Export Report

    **TBC**

### Review test results

**TBC**

# Appendix A

The following hosts are required to be accessible for the Functional Conformance Suite to function correctly:

| Protocol   | Host | Ports | Comment |
| ---------- | ---- | ----- | ------- |
| TCP, HTTPS | modelobankauth2018.o3bank.co.uk | 4101 | Only required when testing against Ozone Model Bank.
| TCP, HTTPS | modelo2018.o3bank.co.uk | 4201,4501 | Only required when testing against Ozone Model Bank.
| TCP, HTTPS | github.map.fastly.net | 443 | DNS Alias for `raw.githubusercontent.com` - CDN to access OBIE Swagger spec files.
| TCP, HTTPS | api.bitbucket.org | 443 | Used to get version information for Conformance Suite - Update available check.
| TCP, HTTPS | production.cloudflare.docker.com | 443 | Access to Docker repository.
| TCP, HTTPS | registry-1.docker.io | 443 | Access to Docker repository.
| TCP, HTTPS | auth.docker.io | 443 | Authenticating with Docker Hub.
