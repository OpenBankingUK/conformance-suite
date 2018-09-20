#!/usr/bin/env bash
set -uexo pipefail

echo -e "\033[92m  ---> pushing ob-api-proxy \033[0m"
docker push eu.gcr.io/compliance-suite-server/ob-api-proxy:latest
echo -e "\033[92m  ---> pushing compliance-suite-server \033[0m"
docker push eu.gcr.io/compliance-suite-server/compliance-suite-server:latest
echo -e "\033[92m  ---> pushing compliance-suite-server-ga \033[0m"
docker push eu.gcr.io/compliance-suite-server/compliance-suite-server-ga:latest
