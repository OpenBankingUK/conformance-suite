#!/usr/bin/env bash
if [[ -z "$DOCKER_USER" ]]; then
  echo -e "\033[91m  ---> cluster: DOCKER_USER is not set \033[0m" 1>&2
  exit 1
fi
if [[ -z "$DOCKER_PASS" ]]; then
  echo -e "\033[91m  ---> cluster: DOCKER_PASS is not set \033[0m" 1>&2
  exit 1
fi

set -ueo pipefail

echo -e "\033[92m  ---> cluster: setting up imagePullSecrets  \033[0m"
kubectl create secret docker-registry openbankingukregistrykey \
  --docker-server="https://index.docker.io/v1/" \
  --docker-username="$DOCKER_USER" \
  --docker-password="$DOCKER_PASS" \
  --docker-email="mohamed.bana@openbanking.org.uk"

kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "openbankingukregistrykey"}]}'
