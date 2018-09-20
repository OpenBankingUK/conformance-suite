#!/usr/bin/env bash
# Prints out the date and git revision for each of the deployments
set -uexo pipefail
OUTPUT_PATH='{"date="}{.spec.template.metadata.labels.date}{"\n"}{"GIT_REV="}{.spec.template.metadata.labels.GIT_REV}{"\n"}'

kubectl get deployment/zookeeper --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/kafka --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/redis --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/mongo --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/ob-api-proxy --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/compliance-suite-server --output=jsonpath="${OUTPUT_PATH}"
kubectl get deployment/compliance-suite-server-ga --output=jsonpath="${OUTPUT_PATH}"

kubectl get pods
