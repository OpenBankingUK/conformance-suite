.DEFAULT_GOAL:=help
SHELL:=/bin/bash
GO_PKGS=$(shell go list ./...)
GO_PKGS_FOLDERS=$(shell go list -f '{{.Dir}}/' ./...)

.PHONY: all


.PHONY: help
help: ## Displays this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building & Running:

.PHONY: run
run: WEB_LOG_FILE:=$(shell pwd)/web/web.log
run: init_web ## run binary directly without docker.
	@if [[ -f "${WEB_LOG_FILE}" ]]; then rm "${WEB_LOG_FILE}"; fi
	@./scripts/run_web.sh &> "${WEB_LOG_FILE}" &
	@./scripts/run_server.sh

.PHONY: run_parallel
run_parallel: init_web ## run binary directly with logging without docker
	parallel -j2 --linebuffer --verbose --tag ::: ./scripts/run_web.sh ./scripts/run_server.sh

.PHONY: run_image
run_image: ## run the 'latest' docker image.
	@echo -e "\033[92m  ---> Running image ... \033[0m"
	docker run \
		--rm \
		-it \
		-p 8080:8080 \
		-v $(shell pwd)/config:/app/config:ro \
		-v $(shell pwd)/swagger:/app/swagger:ro \
		"openbanking/conformance-suite:latest"

.PHONY: build
build: ## build the server binary directly.
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -o server cmd/server/*.go

.PHONY: build_cli
build_cli: ## build the cli binary directly.
	@echo -e "\033[92m  ---> Building CLI ... \033[0m"
	go build -o fcs cmd/cli/*.go

.PHONY: build_image
build_image: ## build the docker image.
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	@# We could enable parallel builds for multi-staged builds with `DOCKER_BUILDKIT=1`
	@# See: https://github.com/moby/moby/pull/37151
	@#DOCKER_BUILDKIT=1
	docker build ${DOCKER_BUILD_ARGS} -t "openbanking/conformance-suite:latest" .

##@ Dependencies:

.PHONY: init
init: ## initialise.
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go mod download

init_web: ./web/node_modules ## install node_modules when not present.

./web/node_modules:
	cd web && yarn install

.PHONY: devtools
devtools: ## install dev tools.
	@echo -e "\033[92m  ---> Installing golint (golang.org/x/lint/golint) ... \033[0m"
	go get -u golang.org/x/lint/golint
	@echo -e "\033[92m  ---> Installing gocyclo (github.com/fzipp/gocyclo) ... \033[0m"
	go get -u github.com/fzipp/gocyclo
	@echo -e "\033[92m  ---> Installing mockery (github.com/vektra/mockery) ... \033[0m"
	go get -u github.com/vektra/mockery
	@echo -e "\033[92m  ---> Installing gometalinter (github.com/alecthomas/gometalinter) ... \033[0m"
	curl -L https://git.io/vp6lP | BINDIR="$GOPATH/bin" sh

##@ Cleanup:

.PHONY: lint
lint: ## lint the go code.
	@echo -e "\033[92m  ---> Vetting ... \033[0m"
	go vet ${GO_PKGS}
	@echo -e "\033[92m  ---> Linting ... \033[0m"
	golint -min_confidence 1.0 -set_exit_status ${GO_PKGS}
	@echo -e "\033[92m  ---> Formatting ... \033[0m"
	@for GO_PKG_DIR in ${GO_PKGS_FOLDERS}; do \
		echo -e "\033[92m  ---> Formatting $${GO_PKG_DIR}*.go ... \033[0m"; \
		gofmt -e -s -w $${GO_PKG_DIR}*.go; \
	done

.PHONY: cyclomatic
cyclomatic: ## cyclomatic complexity checks.
	@echo -e "\033[92m  ---> Checking cyclomatic complexity ... \033[0m"
	gocyclo -over 12 ${GO_PKGS_FOLDERS}

.PHONY: cyclomatic
metalinter: ## other qa tools (linter).
	@echo -e "\033[92m  ---> Checking other qa tools ... \033[0m"
	gometalinter --disable-all --enable=structcheck --enable=megacheck --enable=misspell -enable=vetshadow --enable=goconst --enable=nakedret --enable=deadcode --enable=unparam --deadline=30s --line-length=1747 --enable=lll ${GO_PKGS_FOLDERS}

.PHONY: qa
qa: test lint cyclomatic metalinter ## run all known quality assurance tools

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
		-v \
		-cover \
		./...

.PHONY: test_coverage
test_coverage: ## run the go tests then open up coverage report.
	@echo -e "\033[92m  ---> Testing wth coverage ... \033[0m"
	go test \
		-v \
		-cover \
		-coverprofile=$(shell pwd)/coverage.out \
		./...
	go tool cover \
		-html=$(shell pwd)/coverage.out
