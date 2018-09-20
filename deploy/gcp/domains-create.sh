#!/usr/bin/env bash
# NB: In the future this can be replaced by using this
# Kubernetes addon, https://github.com/kubernetes-incubator/external-dns.
# For now, we will just create the records manually.
set -ueo pipefail

PROJECT_ID="compliance-suite-server"
SERVICE_NAME="compliance-suite-server-ga"

# get the external ip address assigned to the service, e.g.,
# $ kubectl get svc
# NAME                         TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                      AGE
# compliance-suite-server-ga   LoadBalancer   10.39.243.185   35.232.238.217   80:30014/TCP,443:31476/TCP   3h
#
# $ svc_ip=$(get_external_ip "compliance-suite-server-ga")
# $ echo $svc_ip
# 35.232.238.217
function get_external_ip {
  local service_name=$1
  >&2 echo -e "\033[92m  ---> domains: finding $service_name loadbalancer IP \033[0m"
  kubectl get svc "$service_name" -o=jsonpath='{.status.loadBalancer.ingress.*.ip}'
}

# remove old transaction.yaml and transaction.yml files
function clear_transaction {
  if [[ -f ./transaction.yml ]]; then
      echo -e "\033[93m  ---> domains: removing old transaction ./transaction.yml"
      rm -fr ./transaction.yml
  fi
  if [[ -f ./transaction.yaml ]]; then
      echo -e "\033[93m  ---> domains: removing old transaction ./transaction.yaml"
      rm -fr ./transaction.yaml
  fi
  if [[ -f ../transaction.yml ]]; then
      echo -e "\033[93m  ---> domains: removing old transaction ../transaction.yml"
      rm -fr ../transaction.yml
  fi
  if [[ -f ../transaction.yaml ]]; then
      echo -e "\033[93m  ---> domains: removing old transaction ../transaction.yaml"
      rm -fr ../transaction.yaml
  fi
}

function create_record {
  local to="$1"
  local service_name="$2"
  local zone_name="$3"
  local record_name="$to."
  local service_ip_address
  local service_ip_address_old

  service_ip_address=$(get_external_ip "$service_name")
  service_ip_address_old=$(gcloud --format="value(rrdatas)" dns record-sets list --zone="$zone_name" --name="$record_name" --type="A")

  # remove old record if there is one
  echo -e "\033[93m  ---> domains: removing $record_name=$service_ip_address_old  \033[0m"
  gcloud dns --project=compliance-suite-server record-sets transaction start --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction remove "$service_ip_address_old" --name="$record_name" --ttl=60 --type=A --zone="$zone_name" || (rm -fr ./transaction.yml || rm -fr ./transaction.yaml)
  gcloud dns --project=compliance-suite-server record-sets transaction execute --zone="$zone_name"

  # create new record
  echo -e "\033[92m  ---> domains: creating $record_name=$service_ip_address  \033[0m"
  gcloud dns --project=compliance-suite-server record-sets transaction start --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction add "$service_ip_address" --name="$record_name" --ttl=60 --type=A --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction execute --zone="$zone_name"
}

function create_records {
  local domain="$1"
  local zone_name="compliance-suite-server-$domain-zone"

  echo -e "\033[92m  ---> domains: creating records compliance-suite-server.$domain  \033[0m"

  gcloud dns --project="$PROJECT_ID" managed-zones create $zone_name --description= --dns-name="compliance-suite-server.$domain." || true
  create_record "compliance-suite-server.$domain" "$SERVICE_NAME" "$zone_name"

  # For example:
  # map www.compliance-suite-server.tk to compliance-suite-server.tk
  # map www.compliance-suite-server.ga to compliance-suite-server.ga
  echo -e "\033[93m  ---> domains: removing www.compliance-suite-server.$domain=compliance-suite-server.$domain  \033[0m"
  gcloud dns --project=compliance-suite-server record-sets transaction start --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction remove "compliance-suite-server.$domain." --name="www.compliance-suite-server.$domain." --ttl=60 --type=CNAME --zone="$zone_name" || (rm -fr ./transaction.yml || rm -fr ./transaction.yaml)
  gcloud dns --project=compliance-suite-server record-sets transaction execute --zone="$zone_name"

  echo -e "\033[92m  ---> domains: creating www.compliance-suite-server.$domain=compliance-suite-server.$domain  \033[0m"
  gcloud dns --project=compliance-suite-server record-sets transaction start --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction add "compliance-suite-server.$domain." --name="www.compliance-suite-server.$domain." --ttl=60 --type=CNAME --zone="$zone_name"
  gcloud dns --project=compliance-suite-server record-sets transaction execute --zone="$zone_name"
}

function create_records_all {
  create_records "tk"
  echo
  create_records "ga"
  # echo
  # create_records "ml"
  # echo
  # create_records "cf"
  # echo
  # create_records "gq"
}

echo -e "\033[92m  ---> domains: creating zone, pwd: $(pwd)  \033[0m"

while [[ -z $(get_external_ip "$SERVICE_NAME") ]]
do
  echo -e "\033[93m  ---> domains: waiting for $SERVICE_NAME ...  \033[0m"
  sleep 1
done

clear_transaction
create_records_all
