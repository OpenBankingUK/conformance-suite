ifeq ($(strip $(ENV)),)
ENVFILE = .env
else
ENVFILE = .env.$(ENV)
endif

# Create an empty env file for the environment
# useful for circleci or prod/staging envs where secrets are declared somewhere else
ifeq (,$(wildcard $(ENVFILE)))
define ENVS
VALIDATION_KAFKA_TOPIC=${VALIDATION_KAFKA_TOPIC}\\n
VALIDATION_KAFKA_BROKER=${VALIDATION_KAFKA_BROKER}\\n
KAFKA_HOST=${KAFKA_HOST}\\n
KAFKA_PORT=${KAFKA_PORT}\\n
DATA_DB_HOST=${DATA_DB_HOST}\\n
DATA_DB_NAME=${DATA_DB_NAME}\\n
GOOGLE_OAUTH_CLIENT_ID=${GOOGLE_OAUTH_CLIENT_ID}\\n
GOOGLE_OAUTH_CLIENT_SECRET=${GOOGLE_OAUTH_CLIENT_SECRET}\\n
GUARDIAN_SECRET_KEY=${GUARDIAN_SECRET_KEY}\\n
OB_API_PROXY_URL=${OB_API_PROXY_URL}\\n
ENDPOINT_URL_HOST=${ENDPOINT_URL_HOST}\\n
ENDPOINT_URL_PORT=${ENDPOINT_URL_PORT}
endef
# Create .env for the right environment
$(shell echo $(ENVS) > $(ENVFILE))
endif

# Include .env.ENV and export all variables
include $(ENVFILE)
export

SHELL:=/bin/bash
GIT_REV:=$(shell git rev-parse HEAD)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
PWD:=$(shell pwd)
NPM_INSTALL_LOG_LEVEL:=warn
NPM_INSTALL_ARGS:=--loglevel=$(NPM_INSTALL_LOG_LEVEL) --no-progress

.PHONY: all
all: format lint test

# === Format ===
.PHONY: format_server
format_server:
	@mix format

.PHONY: format
format: format_server
# === Format ===

# === Lint ===
.PHONY: lint_server
lint_server:
	mix credo --all --strict

.PHONY: lint_web
lint_web:
	cd apps/compliance_web/assets && npm run lintjs

.PHONY: lint
lint: lint_server lint_web
# === Lint ===

# === Compile ===
.PHONY: compile_server
compile_server:
	MIX_ENV=prod mix deps.get --only prod
	MIX_ENV=prod mix compile

.PHONY: compile_server_local
compile_server_local:
	mix deps.get
	mix compile

.PHONY: compile
compile: compile_server
# === Compile ===

# === Utils ===
.PHONY: clean
clean:
	@echo -e "\033[92m  ---> clean: removing \033[0m"
	@rm -fr $(shell find . -name "node_modules")
	@rm -fr _build/ deps/ .elixir_ls/
	@rm -fr $(shell find . -name '*.beam')
	@rm -fr $(shell find . -name '.DS_Store')

.PHONY: build_ci
build_ci:
	# curl -o /usr/local/bin/circleci https://circle-downloads.s3.amazonaws.com/releases/build_agent_wrapper/circleci && chmod +x /usr/local/bin/circleci
	@echo -e "\033[92m  ---> circleci: updating \033[0m"
	circleci update
	@echo -e "\033[92m  ---> circleci: validating \033[0m"
	circleci config validate -c .circleci/config.yml
	@echo -e "\033[92m  ---> circleci: running \033[0m"
	@circleci build \
		--branch="$(GIT_BRANCH)" \
		--job "ob-api-proxy" \
		--node-total 4 \
		--repo-url="/fake-remote" \
		--skip-checkout=false \
		--volume="$(PWD)":"/fake-remote"
	@circleci build \
		-e GOOGLE_OAUTH_CLIENT_ID="$${GOOGLE_OAUTH_CLIENT_ID}" \
		-e GOOGLE_OAUTH_CLIENT_SECRET="$${GOOGLE_OAUTH_CLIENT_SECRET}" \
		-e GUARDIAN_SECRET_KEY="$${GUARDIAN_SECRET_KEY}" \
		-e VALIDATION_KAFKA_TOPIC="$${VALIDATION_KAFKA_TOPIC}" \
		-e VALIDATION_KAFKA_BROKER="$${VALIDATION_KAFKA_BROKER}" \
		-e KAFKA_HOST="$${KAFKA_HOST}" \
		-e KAFKA_PORT="$${KAFKA_PORT}" \
		-e DATA_DB_HOST="$${DATA_DB_HOST}" \
		-e DATA_DB_NAME="$${DATA_DB_NAME}" \
		-e OB_API_PROXY_URL="$${OB_API_PROXY_URL}" \
		--branch="$(GIT_BRANCH)" \
		--job "build" \
		--node-total 4 \
		--repo-url="/fake-remote" \
		--skip-checkout=false \
		--volume="$(PWD)":"/fake-remote"
# === Utils ===

# === Serve ===
# example:
# $ DEBUG_ENVS=true ENV=serve_web make serve_web
.PHONY: serve_web
serve_web:
	./web-serve.sh

.PHONY: serve_web_docker
serve_web_docker:
	docker-compose up -d --build --scale compliance-suite-server-ga=0

.PHONY: start_server
start_server:
	PORT=4000 MIX_ENV=prod mix phx.server

.PHONY: start_server_local
start_server_local:
	iex -S mix phx.server
# === Serve ===

# === Init ===
.PHONY: init_server
init_server:
	# MIX_ENV=prod mix ecto.drop
	# MIX_ENV=prod mix ecto.create
	MIX_ENV=prod mix ecto.migrate

.PHONY: init_server_local
init_server_local:
	# mix ecto.drop
	# mix ecto.create
	mix ecto.migrate
# === Init ===

# === Build ===
.PHONY: build_web
build_web:
	cd apps/compliance_web/assets && npm install $(NPM_INSTALL_ARGS)
	cd apps/compliance_web/assets && npm run build
	cd apps/compliance_web && MIX_ENV=prod mix phx.digest

.PHONY: build_web_local
build_web_local:
	cd apps/compliance_web/assets && npm install $(NPM_INSTALL_ARGS)
	cd apps/compliance_web/assets && npm run build
	cd apps/compliance_web && mix phx.digest
# === Build ===

# === Docker ===
.PHONY: build_images_compliance
build_images_compliance:
	@echo -e "\033[92m  ---> image: building compliance-suite-server \033[0m"
	@docker-compose build --no-cache

.PHONY: build_image_reference_mock_server
build_image_reference_mock_server:
	@cd ../reference-mock-server && docker-compose build

.PHONY: build_images
build_images:
	@echo -e "\033[92m  ---> images: building \033[0m"
	@make build_images_compliance
	@make build_image_reference_mock_server
	@make images_tag

.PHONY: images_tag
images_tag:
	@echo -e "\033[92m  ---> images: tagging \033[0m"
	docker tag eu.gcr.io/compliance-suite-server/ob-api-proxy:latest eu.gcr.io/compliance-suite-server/ob-api-proxy:$(GIT_REV)
	docker tag eu.gcr.io/compliance-suite-server/compliance-suite-server:latest eu.gcr.io/compliance-suite-server/compliance-suite-server:$(GIT_REV)
	docker tag eu.gcr.io/compliance-suite-server/compliance-suite-server-ga:latest eu.gcr.io/compliance-suite-server/compliance-suite-server-ga:$(GIT_REV)

	@docker images

.PHONY: stop_images
stop_images:
	@echo -e "\033[92m  ---> stopping dependencies \033[0m"
	@docker-compose down --volumes --remove-orphans
# === Docker ===

# === Deployment/Kubernetes ===
.PHONY: deploy_details
deploy_details:
	@echo -e "\033[92m  ---> deploy: details \033[0m"
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD)/deploy/gcloud_scripts:/gcloud_scripts:ro \
		-e GCLOUD_SERVICE_KEY='$(GCLOUD_SERVICE_KEY)' \
		google/cloud-sdk:latest \
		/bin/bash -c '/gcloud_scripts/gcloud_authenticate.sh && /gcloud_scripts/deploy_details.sh'

.PHONY: deploy_images_push
deploy_images_push:
	@echo -e "\033[92m  ---> deploy: push images \033[0m"
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD)/deploy/gcloud_scripts:/gcloud_scripts:ro \
		-e GCLOUD_SERVICE_KEY='$(GCLOUD_SERVICE_KEY)' \
		google/cloud-sdk:latest \
		/bin/bash -c '/gcloud_scripts/gcloud_authenticate.sh && /gcloud_scripts/deploy_images_push.sh'

.PHONY: deploy_images_update
deploy_images_update:
	@echo -e "\033[92m  ---> deploy: updating images \033[0m"
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD)/deploy/gcloud_scripts:/gcloud_scripts:ro \
		-e GCLOUD_SERVICE_KEY='$(GCLOUD_SERVICE_KEY)' \
		-e GIT_REV='$(GIT_REV)' \
		google/cloud-sdk:latest \
		/bin/bash -c '/gcloud_scripts/gcloud_authenticate.sh && /gcloud_scripts/deploy_images_update.sh'
# === Deployment/Kubernetes ===

# === Test ===
.PHONY: test_server_local
test_server_local:
	docker-compose up -V -d mongo && \
		mix test --trace --color && \
		docker-compose down -v --remove-orphans

.PHONY: test_server
test_server:
	@echo -e "\033[92m  ---> Running Server Tests \033[0m"
	mix deps.get --only test
	mix test --trace --color

.PHONY: test_web
test_web:
	@echo -e "\033[92m  ---> Running Unit Tests \033[0m"
	cd apps/compliance_web/assets && npm install $(NPM_INSTALL_ARGS)
	cd apps/compliance_web/assets && npm run test

.PHONY: test
test: test_server test_web

.PHONY: test_integration
test_integration:
	@echo -e "\033[92m  ---> Running end-2-end Tests \033[0m"
	cd ./e2e/ && npm install $(NPM_INSTALL_ARGS)
	cd ./e2e/ && npm run headless

.PHONY: test_integration_local
test_integration_local:
	@echo -e "\033[92m  ---> Running end-2-end Tests \033[0m"
	@cd ./e2e/ && npm ci && npm run local

.PHONY: run_e2e_tests
run_e2e_tests: NETWORK=$(shell docker network ls --filter name=compliance --format "{{.Name}}")
run_e2e_tests:
	@echo -e "\033[92m  ---> Running run_e2e_tests \033[0m"
	@docker run \
		--network="$(NETWORK)" \
		--rm \
		--name="tests_integration" \
		--volume $(PWD)/e2e:/e2e \
		--volume ~/.cache:/root/.cache \
		--volume ~/.npm:/root/.npm \
		"cypress/base:8" \
		/bin/bash -c 'cd /e2e && npm install && npm run headless'

.PHONY: run_web_tests
run_web_tests:
	@echo -e "\033[92m  ---> Running run_web_tests \033[0m"
	docker run \
		--rm \
		--name="tests_web" \
		--volume $(PWD)/Makefile:/root/app/Makefile:ro \
		--volume $(PWD)/apps/compliance_web/assets/:/root/app/apps/compliance_web/assets/:ro \
		"cypress/base:8" \
			/bin/bash -c 'make test_web'

.PHONY: run_server_tests
run_server_tests:
	@echo -e "\033[92m  ---> Running run_server_tests \033[0m"
	docker-compose run \
		--rm \
		--name="tests_server" \
		compliance-suite-server \
			/bin/bash -c 'make test_server'

# This allows you to debug X11 applications in this cause it gives you a gui of cypress.
# how to debug failing test: https://sourabhbajaj.com/blog/2017/02/07/gui-applications-docker-mac/
.PHONY: run_tests_debug
run_tests_debug: IP=$(shell ifconfig en0 | grep inet | awk '$$1=="inet" {print $$2}')
run_tests_debug:
	@echo -e "\033[92m  ---> Running Tests \033[0m"
	# xhost + $${IP}
	docker run \
		--network="compliance-suite-server_openbanking_network" \
		--rm \
		--name "tests_e2e_debug" \
		-e DISPLAY=$(IP):0 \
		--volume /tmp/.X11-unix:/tmp/.X11-unix \
		--volume $(PWD):/root/app/:ro \
		"cypress/base:8" \
			/bin/bash -c 'cd ./e2e/ && npm install && npm run open'
# === Test ===
#
release:
	@PORT=4000 MIX_ENV=prod mix release

serve:
	@/opt/app/bin/compliance_suite_server
