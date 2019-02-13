# Functional Conformance Suite Setup Guide

This guide will assist you with the technical steps required to setup the Functional Conformance Suite (FCS) and run your first test. In this guide we will be connecting to and running tests against the Ozone Model Bank.

Please note the following goals of this document:
* Register with Ozone model bank as a Third Party Provider (TPP).
* Generate certificates for transport and signing of requests to Ozone - self signed certs to be used.
* Setup the FCS to run on local machine (via Docker container).
* Execute a test on the FCS using the Payment Service User (PSU) consent flow. The PSU consent flow involves user
interaction during the authentication and consent steps when connecting to Ozone.
* Review the results of the test run.

# Prerequisites

This guide assumes the following tools are installed and functioning correctly. Versions specified used when writing this guide.
* Docker (Client: 18.09.1, Server: 18.09.1 on OSX)
* OpenSSL (LibreSSL 2.6.4 on OSX)
* Postman Native (6.7.3 on OSX)

# Steps

## Ozone Bank (Model Bank)

Ozone Bank is an Account Servicing Payment Service Provider (ASPSP), which the FCS will connect to as a TPP. With that in mind,
you must [enrol with Ozone](https://ob2018.o3bank.co.uk:444/pub/home).

* Use Google or LinkedIn as identity provider
* Any organisation name
* Redirect URI should be `https://0.0.0.0:8443/conformancesuite/callback`
* Click "Go"
* Registration should be successful, keep the page open for reference later.

## Generate transport and signing certificates

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

### Precursor: CA Generation

A precursor step of creating a Certificate Authority (CA) is required to sign the transport and signing certs.

`openssl req -new -x509 -days 3650 -keyout ca.key -out ca.pem -nodes`

The following certificate information (DN) values shall be provided, based on what is mentioned on your Ozone registration page
under the Certificates tab:
* C (Country Name)
* O (Organisation Name)
* OU (Organisation Unit) 
* CN (Common Name)

### Transport Certificate

Execute `openssl genrsa -out transport.key 2048`

Execute `openssl req -new -sha256 -key transport.key -out transport.csr` (Provide same DN information as above, no passphrase) 

Execute `openssl x509 -req -days 3650 -in transport.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out transport.pem`

### Signing Certificate

Execute `openssl genrsa -out signing.key 2048`

Execute `openssl req -new -sha256 -key signing.key -out signing.csr` (Provide same DN information as above, no passphrase) 

Execute `openssl x509 -req -days 3650 -in signing.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out signing.pem`

## Run the Functional Conformance Suite

### FCS Server Certificates

* Step 1: Download certificates (localhost)[https://bitbucket.org/openbankingteam/conformance-suite/src/develop/certs/]
* Step 2: Trust the certificate or add as an exception.

**TBC** Is the same as running `make run_parallel` ?

The following command will download the latest FCS image from docker hub and run it. You may need to login to Docker Hub
at this point by running `docker login`. 

`docker run \
        --rm \
        -it \
        -p 8443:8443 \
        "openbanking/conformance-suite:latest"`

If all goes well you should be able to launch the FCS UI from you browser via `https://0.0.0.0:8443`

### Docker Content Trust

Docker ensures that all content is securely received and verified by Open Banking. Docker cryptographically signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flag, for example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-suite:TAG
    DOCKER_CONTENT_TRUST=1 docker RUN openbanking/conformance-suite:TAG

## Execute a test

Running a test plan on the FCS involves five steps, as follows:

1. Start / Load test - Select a template.

2. Discovery - Review the discovery file and update as required.

    You will need to update the discovery file with some values should be retrieved from Ozone. To view these values,
    download the "Postman Environment" file from the "Postman" tab on the Ozone registration page. _I used [jsonlint.com](https://www.jsonlint.com)
    to format the JSON._

    Please keep a note of the following values from the JSON, which are stored as key/value records:
    * obParticipantId
    * oidcClient
    * basicToken
    * redirectUrl

    In the discovery JSON you will need to set the following values in `discoveryModel.customTests.replacementParameters`, based on the above values:
    * client_id = oidcClient
    * fapi_financial_id = obParticipantId 
    * basic_authentication = basicToken
    * redirect_url = redirectUrl

3. Configuration

    Provide the keys, as created earlier. The naming should be self explanatory if you ran the cert generation commands as shown.

4. Run / Overview

    This screen shows the tests that will be run. Once ready, click "Run" at the end of the page. The tests should run and go to the "PENDING" status. Once complete the status should move to "PASSED", if everything ran ok. If any of the tests failed, you can click the "FAILED" badge to view more information on the cause of failure. 

5. Export Report

    **TBC**

## Review test results

**TBC**

# How to get help

**TBC**
