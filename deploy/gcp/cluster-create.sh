#!/usr/bin/env bash
set -ueo pipefail
if [[ -z "$DOCKER_USER" ]]; then
  echo -e "\033[91m  ---> cluster: DOCKER_USER is not set \033[0m" 1>&2
  exit 1
fi
if [[ -z "$DOCKER_PASS" ]]; then
  echo -e "\033[91m  ---> cluster: DOCKER_PASS is not set \033[0m" 1>&2
  exit 1
fi

PROJECT_ID="compliance-suite-server"
# london
COMPUTE_ZONE="europe-west2-a"
CLUSTER_NAME="compliance-suite-server-cluster"
NUM_NODES="4"

gcloud config set project $PROJECT_ID
gcloud config set compute/zone $COMPUTE_ZONE

echo -e "\033[92m  ---> cluster: creating  \033[0m"
gcloud beta container --project "$PROJECT_ID" clusters create "$CLUSTER_NAME" --zone "$COMPUTE_ZONE" --username "admin" --cluster-version "1.10.5-gke.0" --machine-type "g1-small" --image-type "COS" --disk-type "pd-standard" --disk-size "30" --node-labels machine-type=g1-small,machine-size=medium --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "$NUM_NODES" --enable-stackdriver-kubernetes --network "default" --subnetwork "default" --addons HorizontalPodAutoscaling,HttpLoadBalancing,KubernetesDashboard --enable-autoupgrade --enable-autorepair --maintenance-window "23:00"

# WARNING: Currently VPC-native is not the default mode during cluster creation. In the future, this will become the default mode and can be disabled using `--no-enable-ip-alias` flag. Use `--[no-]enable-ip-alias` flag to suppress this warning.
# This will enable the autorepair feature for nodes. Please see
# https://cloud.google.com/kubernetes-engine/docs/node-auto-repair for more
# information on node autorepairs.

# This will enable the autoupgrade feature for nodes. Please see
# https://cloud.google.com/kubernetes-engine/docs/node-management for more
# information on node autoupgrades.

# Creating cluster compliance-suite-server-cluster...done.
# Created [https://container.googleapis.com/v1beta1/projects/compliance-suite-server/zones/us-central1-a/clusters/compliance-suite-server-cluster].
# To inspect the contents of your cluster, go to: https://console.cloud.google.com/kubernetes/workload_/gcloud/us-central1-a/compliance-suite-server-cluster?project=compliance-suite-server
# kubeconfig entry generated for compliance-suite-server-cluster.
# NAME                             LOCATION       MASTER_VERSION  MASTER_IP       MACHINE_TYPE  NODE_VERSION  NUM_NODES  STATUS
# compliance-suite-server-cluster  us-central1-a  1.10.5-gke.0    35.192.101.137  g1-small      1.10.5-gke.0  4          RUNNING

echo -e "\033[92m  ---> cluster: configuring kubectl  \033[0m"
gcloud container clusters get-credentials $CLUSTER_NAME --zone $COMPUTE_ZONE --project $PROJECT_ID
