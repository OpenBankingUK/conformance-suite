#!/usr/bin/env bash
# This patches/update the images the deployment are pointing to.
#
# See these guides in particular the section on patching the deployment:
# * https://circleci.com/docs/2.0/deployment-integrations/#google-cloud
set -uexo pipefail

# NB: Value is required to be in a specific format otherwise we get the error below.
# The Deployment "compliance-suite-server-ga" is invalid: spec.template.labels: Invalid value: "2018-08-13T08:35:09+0000": a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')
PATCH="{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"$(date '+%F_%H-%M-%S_%Z')\", \"GIT_REV\":\"${GIT_REV}\"}}}}}"

echo -e "\033[92m  ---> updating ob-api-proxy \033[0m"
# kubectl set image deployment/ob-api-proxy ob-api-proxy=eu.gcr.io/compliance-suite-server/ob-api-proxy:latest
kubectl patch deployment ob-api-proxy -p "${PATCH}"

echo -e "\033[92m  ---> updating compliance-suite-server \033[0m"
# kubectl set image deployment/compliance-suite-server compliance-suite-server=eu.gcr.io/compliance-suite-server/compliance-suite-server:latest
kubectl patch deployment compliance-suite-server -p "${PATCH}"

# we don't really need to update the reverse proxy - it remains static for
# the most part.
echo -e "\033[92m  ---> updating compliance-suite-server-ga-ga \033[0m"
# kubectl set image deployment/compliance-suite-server-ga compliance-suite-server-ga=eu.gcr.io/compliance-suite-server/compliance-suite-server-ga:latest
kubectl patch deployment compliance-suite-server-ga -p "${PATCH}"
