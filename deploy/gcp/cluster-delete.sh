#!/usr/bin/env bash
set -ueo pipefail

PROJECT_ID="compliance-suite-server"
# london
COMPUTE_ZONE="europe-west2-a"
CLUSTER_NAME="compliance-suite-server-cluster"

echo -e "\033[92m  ---> cluster: deleting $CLUSTER_NAME  \033[0m"
gcloud container clusters delete "$CLUSTER_NAME" --zone "$COMPUTE_ZONE" --project "$PROJECT_ID"
