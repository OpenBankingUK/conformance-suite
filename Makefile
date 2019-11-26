.DEFAULT_GOAL:=help
SHELL:=/bin/bash
IMAGE_TAG=latest
ENABLE_IMAGE_SIGNING=0

.PHONY: all


.PHONY: help
help: ## Displays this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building & Running:

.PHONY: run
run: WEB_LOG_FILE:=$(shell pwd)/web/web.log
run: init_web ## run binary directly without docker.
	@if [[ -f "${WEB_LOG_FILE}" ]]; then rm "${WEB_LOG_FILE}"; fi
	@./scripts/web &> "${WEB_LOG_FILE}" &
	@./scripts/server

.PHONY: run_parallel
run_parallel: init_web ## run binary directly with logging without docker
	cd scripts && parallel -j2 --linebuffer --verbose --tag ::: ./web ./server

.PHONY: run_image
run_image: ## run the 'latest' docker image.
	@echo -e "\033[92m  ---> Running image ... \033[0m"
	docker run \
		--rm \
		-it \
		-p 8443:8443 \
		"openbanking/conformance-suite:latest"

.PHONY: build
build: ## build the server binary directly.
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -o fcs_server cmd/fcs_server/*.go

.PHONY: build_cli
build_cli: ## build the cli binary directly.
	@echo -e "\033[92m  ---> Building CLI ... \033[0m"
	go build -o fcs cmd/cli/*.go

.PHONY: build_image
build_image: ## build the docker image. Use available args IMAGE_TAG=v1.x.y, ENABLE_IMAGE_SIGNING=1
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	@# We could enable parallel builds for multi-staged builds with `DOCKER_BUILDKIT=1`
	@# See: https://github.com/moby/moby/pull/37151
	@#DOCKER_BUILDKIT=1
	@export DOCKER_CONTENT_TRUST=${ENABLE_IMAGE_SIGNING}
	docker build ${DOCKER_BUILD_ARGS} -t "openbanking/conformance-suite:${IMAGE_TAG}" .

##@ Dependencies:

.PHONY: init
init: ## initialise.
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go mod download

init_web: ./web/node_modules ## install node_modules when not present.

./web/node_modules:
	cd web && yarn install --frozen-lockfile --non-interactive

.PHONY: devtools
devtools: ## install dev tools.
	@echo -e "\033[92m  ---> Installing mockery (github.com/vektra/mockery) ... \033[0m"
	go get github.com/vektra/mockery
	@echo -e "\033[92m  ---> Installing golangci-lint (https://github.com/golangci/golangci-lint) ... \033[0m"
	curl -sfL "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" | sh -s -- -b $(shell go env GOPATH)/bin v1.16.0

##@ Cleanup:

.PHONY: lint
lint: ## lint the go code.
	@echo -e "\033[92m  ---> Checking other qa tools ... \033[0m"
	golangci-lint run --config ./.golangci.yml ./...

.PHONY: qa
qa: test lint ## run all known quality assurance tools

.PHONY: clean
clean:
	@echo -e "\033[92m  ---> Cleaning ... \033[0m"
	go clean -i -r -cache -testcache -modcache

##@ Testing:

.PHONY: test
test: ## run the go tests.
	@echo -e "\033[92m  ---> Testing ... \033[0m"
	@# make symbolic link to ./web/public -> ./pkg/server/web/dist so that we can test out that
	@# static files are being served by the Echo web server
	@if [[ ! -d "$$(pwd)/pkg/server/web/dist" ]]; then \
		echo -e "\033[92m  ---> Linking $(shell pwd)/pkg/server/web/dist -> $(shell pwd)/web/public \033[0m"; \
		mkdir -p $(shell pwd)/pkg/server/web; \
		ln -s $(shell pwd)/web/public $(shell pwd)/pkg/server/web/dist; \
	fi
	go test \
		-cover \
		./...

.PHONY: test_coverage
test_coverage: ## run the go tests then open up coverage report.
	@echo -e "\033[92m  ---> Testing with coverage ... \033[0m"
	go test \
		-v \
		-cover \
		-coverprofile=$(shell pwd)/coverage.out \
		./...
	go tool cover \
		-html=$(shell pwd)/coverage.out
