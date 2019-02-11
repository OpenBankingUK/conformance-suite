# Functional Conformance Suite Setup Guide

This guide assists with the technical steps to set-up the Functional Conformance Suite.

## Installing localhost certificates

* Step 1: Download certificates (localhost)[https://bitbucket.org/openbankingteam/conformance-suite/src/develop/certs/]
* Step 2: Trust the certificate or add as an exception.

## Running the UI (Docker)

The fastest way to get up and running is by running the docker image.

    docker run \
            --rm \
            -it \
            -p 8443:8443 \
            "openbanking/conformance-suite:latest"

Functional Conformance Suite should now be available on your localhost @ https://localhost

## Docker Content Trust

Docker ensures that all content is securely received and verified by Open Banking. Docker cryptographically signs the images upon completion of a satisfactory image check, so that implementers can verify and trust certified content.

To verify the content has not been tampered with you can you the `DOCKER_CONTENT_TRUST` flag, for example:

    DOCKER_CONTENT_TRUST=1 docker pull openbanking/conformance-suite:TAG
    DOCKER_CONTENT_TRUST=1 docker RUN openbanking/conformance-suite:TAG
