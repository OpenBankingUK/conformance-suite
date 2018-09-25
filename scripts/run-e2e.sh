#!/usr/bin/env bash
set -ueo pipefail
docker network ls

NETWORK=$(docker network ls --filter name=compliance --format "{{.Name}}")
if [[ -z "$NETWORK" ]]; then
  # could not find network named *compliance* so re-try with openbanking instead
  NETWORK=$(docker network ls --filter name=openbanking --format "{{.Name}}")
  if [[ -z "$NETWORK" ]]; then
    NETWORKS=$(docker network ls)
    echo -e "\\033[91m ---> e2e tests, failed to find NETWORK in NETWORKS=${NETWORKS} ... \\033[0m"
    exit 1
  fi
fi

echo -e "\\033[92m ---> e2e tests, NETWORK=${NETWORK} ... \\033[0m"
docker run \
  --network="${NETWORK}" \
  --rm \
  --volume "${BITBUCKET_CLONE_DIR}/e2e":/e2e \
  --volume "${BITBUCKET_CLONE_DIR}/.cache/cypress/.cache":/root/.cache \
  --volume "${BITBUCKET_CLONE_DIR}/.cache/cypress/.npm":/root/.npm \
  cypress/base:8 \
    /bin/bash -c 'cd /e2e && CI=true npm install && npm run headless'
