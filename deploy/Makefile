SHELL:=/bin/bash
# BROKERS:=$(shell echo "$$(kubectl get svc | grep kafka | awk '{ print $$4; }'):9092")
BROKERS:="kafka.compliance-suite-server.tk:9092"
TOPIC:="kafka-test-topic"
GIT_REV:=$(shell git rev-parse --short HEAD)

.PHONY: convert
convert:
	@echo -e "\033[92m  ---> Kubernetes: converting from Docker Compose to Kubernetes  \033[0m"
	rm -fr ./cluster/out
	mkdir -p ./cluster/out
	cd .. && ENVFILE=./.env.prod COMPOSE_BUILD_TAG=$(GIT_REV) kompose convert --verbose -f ./docker-compose.yml -o deploy/cluster/out/

.PHONY: apply
apply:
	@echo -e "\033[92m  ---> Kubernetes: applying  \033[0m"

	kubectl apply -f ./cluster/out/ob-api-proxy-env-ob-api-proxy-prod-configmap.yaml
	kubectl apply -f ./cluster/out/ob-api-proxy-deployment.yaml
	kubectl apply -f ./cluster/out/ob-api-proxy-service.yaml

	kubectl apply -f ./cluster/out/compliance-suite-server-env-prod-configmap.yaml
	kubectl apply -f ./cluster/out/compliance-suite-server-deployment.yaml
	kubectl apply -f ./cluster/out/compliance-suite-server-service.yaml

.PHONY: apply_services
apply_services:
	@echo -e "\033[92m  ---> Kubernetes: applying services  \033[0m"

	kubectl apply -f ./cluster/out/zookeeper-deployment.yaml
	kubectl apply -f ./cluster/out/zookeeper-service.yaml

	kubectl apply -f ./cluster/out/kafka-deployment.yaml
	kubectl apply -f ./cluster/out/kafka-service.yaml

	kubectl apply -f ./cluster/out/mongo-deployment.yaml
	kubectl apply -f ./cluster/out/mongo-service.yaml

	kubectl apply -f ./cluster/out/redis-deployment.yaml
	kubectl apply -f ./cluster/out/redis-service.yaml

# dont delete volumes or service load balancer rules
.PHONY: delete
delete:
	@echo -e "\033[92m  ---> Kubernetes: deleting \033[0m"

	kubectl delete -f ./cluster/out/compliance-suite-server-env-prod-configmap.yaml
	kubectl delete -f ./cluster/out/compliance-suite-server-deployment.yaml
	kubectl delete -f ./cluster/out/compliance-suite-server-service.yaml

	kubectl delete -f ./cluster/out/ob-api-proxy-env-ob-api-proxy-prod-configmap.yaml
	kubectl delete -f ./cluster/out/ob-api-proxy-deployment.yaml
	kubectl delete -f ./cluster/out/ob-api-proxy-service.yaml

.PHONY: delete_services
delete_services:
	@echo -e "\033[92m  ---> Kubernetes: deleting services \033[0m"

	kubectl delete -f ./cluster/out/mongo-deployment.yaml
	kubectl delete -f ./cluster/out/mongo-service.yaml

	kubectl delete -f ./cluster/out/redis-deployment.yaml
	kubectl delete -f ./cluster/out/redis-service.yaml

	kubectl delete -f ./cluster/out/kafka-deployment.yaml
	kubectl delete -f ./cluster/out/kafka-service.yaml

	kubectl delete -f ./cluster/out/zookeeper-deployment.yaml
	kubectl delete -f ./cluster/out/zookeeper-service.yaml

# dont replace volumes or service load balancer rules
.PHONY: replace
replace:
	@echo -e "\033[92m  ---> Kubernetes: replacing \033[0m"

	kubectl replace --force -f ./cluster/out/compliance-suite-server-env-prod-configmap.yaml
	kubectl replace --force -f ./cluster/out/compliance-suite-server-deployment.yaml
	kubectl replace --force -f ./cluster/out/compliance-suite-server-service.yaml

	kubectl replace --force -f ./cluster/out/ob-api-proxy-env-ob-api-proxy-prod-configmap.yaml
	kubectl replace --force -f ./cluster/out/ob-api-proxy-deployment.yaml
	kubectl replace --force -f ./cluster/out/ob-api-proxy-service.yaml

.PHONY: replace_services
replace_services:
	@echo -e "\033[92m  ---> Kubernetes: replacing services \033[0m"

	kubectl replace --force -f ./cluster/out/kafka-deployment.yaml
	kubectl replace --force -f ./cluster/out/kafka-service.yaml

	kubectl replace --force -f ./cluster/out/mongo-deployment.yaml
	kubectl replace --force -f ./cluster/out/mongo-service.yaml

	kubectl replace --force -f ./cluster/out/redis-deployment.yaml
	kubectl replace --force -f ./cluster/out/redis-service.yaml

	kubectl replace --force -f ./cluster/out/zookeeper-deployment.yaml
	kubectl replace --force -f ./cluster/out/zookeeper-service.yaml

.PHONY: cluster_create
cluster_create:
	@echo -e "\033[92m  ---> Kubernetes: creating cluster  \033[0m"
	@cd gcp && ./cluster-create.sh
	@cd gcp && ./cluster-secrets.sh

.PHONY: cluster_delete
cluster_delete:
	@echo -e "\033[92m  ---> Kubernetes: deleting cluster  \033[0m"
	@cd gcp && ./cluster-delete.sh

# we only need to create the volumes and ingress rule once
# then we wait whilst the ingress rule has been setup so
# that we can create the domain name mappings
#
# Note:
# * That `compliance-suite-server-ga-deployment` doesn't get deleted or replaced, it remains static.
# * That `compliance-suite-server-ga-service` doesn't get deleted or replaced, it remains static.
.PHONY: initial_apply
initial_apply:
	@echo -e "\033[92m  ---> Kubernetes: initial_apply \033[0m"

	@echo -e "\033[92m  ---> Kubernetes: initial_apply creating volumes \033[0m"
	kubectl apply -f ./cluster/static/caddy-certs-persistentvolumeclaim.yaml

	@echo -e "\033[92m  ---> Kubernetes: initial_apply creating compliance-suite-server-ga \033[0m"
	kubectl apply -f ./cluster/out/compliance-suite-server-ga-service.yaml

	@cd gcp && ./domains-create.sh
	@echo -e "\033[92m  ---> Kubernetes: initial_apply waiting for compliance-suite-server.ga domain changes to propogate \033[0m"
	@sleep 120

	@echo -e "\033[92m  ---> Kubernetes: creating Caddy server for compliance-suite-server.ga \033[0m"
	kubectl apply -f ./cluster/out/compliance-suite-server-ga-deployment.yaml

	@make apply

# === Utils ===
.PHONY: kafka_tail
kafka_tail:
	@echo -e "\033[92m  ---> Kafka: tailing topic  \033[0m"
	kafka-console-consumer --brokers=$(BROKERS) --topic=$(TOPIC) --verbose

.PHONY: kafka_echo
kafka_echo:
	@echo -e "\033[92m  ---> Kafka: echoing 'shrinky-dinks' into topic  \033[0m"
	echo "shrinky-dinks" | kafka-console-producer --brokers=$(BROKERS) --topic=$(TOPIC)
# === Utils ===
