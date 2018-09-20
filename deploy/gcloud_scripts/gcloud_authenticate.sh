#!/usr/bin/env bash
# This authenticates gcloud and docker within the container.
#
# See these guides:
# * https://cloud.google.com/sdk/docs/authorizing#authorizing_with_a_service_account
# * https://circleci.com/docs/2.0/google-auth/
set -uexo pipefail

echo "${GCLOUD_SERVICE_KEY}" > /tmp/gcloud-service-key.json

gcloud auth activate-service-account --key-file=/tmp/gcloud-service-key.json
gcloud --quiet config set project "compliance-suite-server"
gcloud --quiet config set compute/zone "europe-west2-a"
gcloud --quiet container clusters get-credentials "compliance-suite-server-cluster" --project="compliance-suite-server"
gcloud --quiet auth configure-docker
