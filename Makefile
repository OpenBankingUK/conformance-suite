SHELL:=/bin/bash
# guarantee that go will not reach the network at all (e.g. GOPROXY=off)
export GOPROXY:=off

.PHONY: all
all: help

.PHONY: run
run: init_web ## run binary directly without docker
	@echo -e "\033[92m  ---> Starting web file watcher ... \033[0m"
	cd web && FORCE_COLOR=1 NODE_DISABLE_COLORS=0 yarn build-watch &> $(shell pwd)/web/web.log &
	@echo -e "\033[92m  ---> Starting server ... \033[0m"
	PORT=8080 go run -mod=vendor cmd/server/main.go

.PHONY: run_image
run_image: ## run the docker image
	@echo -e "\033[92m  ---> Running image ... \033[0m"
	docker run \
		--rm \
		-it \
		-p 8080:8080 \
		-v $(shell pwd)/config:/app/config:ro \
		-v $(shell pwd)/swagger:/app/swagger:ro \
		"openbanking/conformance-suite:latest"

.PHONY: build
build: ## build the server binary directly
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -mod=vendor -o server cmd/server/main.go

.PHONY: build_cli
build_cli: ## build the cli binary directly
	@echo -e "\033[92m  ---> Building CLI ... \033[0m"
	go build -mod=vendor -o fcs cmd/cli/main.go

.PHONY: build_image
build_image: ## build the docker image
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	@# We could enable parallel builds for multi-staged builds with `DOCKER_BUILDKIT=1`
	@# See: https://github.com/moby/moby/pull/37151
	@#DOCKER_BUILDKIT=1
	docker build ${DOCKER_BUILD_ARGS} -t "openbanking/conformance-suite:latest" .

.PHONY: init
init: ## initialise
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go get -v ./...

init_web: ./web/node_modules ## install node_modules when not present

./web/node_modules:
	cd web && yarn install

.PHONY: devtools
devtools: ## install dev tools
	@echo -e "\033[92m  ---> Installing golint (golang.org/x/lint/golint) ... \033[0m"
	GOPROXY= go get -u golang.org/x/lint/golint
	@echo -e "\033[92m  ---> Installing gocyclo (github.com/fzipp/gocyclo) ... \033[0m"
	GOPROXY= go get -u github.com/fzipp/gocyclo

.PHONY: lint
lint: ## lint the go code
	@echo -e "\033[92m  ---> Vetting ... \033[0m"
	GOPROXY= go vet $(shell go list ./... | grep -v /vendor/)
	@echo -e "\033[92m  ---> Linting ... \033[0m"
	GOPROXY= golint -min_confidence 1.0 -set_exit_status $(shell go list ./... | grep -v vendor)
	@echo -e "\033[92m  ---> Formatting ... \033[0m"
	@GO_PKGS="$(shell go list -f {{.Dir}} ./...)"; \
	for PKG_DIR in $${GO_PKGS}; do \
		echo -e "\033[92m  ---> Formatting $${PKG_DIR}/*.go ... \033[0m"; \
		gofmt -e -s -w $${PKG_DIR}/*.go; \
	done

.PHONY: cyclomatic
cyclomatic: ## cyclomatic complexity checks
	@echo -e "\033[92m  ---> Checking cyclomatic complexity ... \033[0m"
	gocyclo -over 12 $(shell ls -d */ | grep -v vendor)

.PHONY: clean
clean:
	@echo -e "\033[92m  ---> Cleaning ... \033[0m"
	go clean -i -r -cache -testcache -modcache

.PHONY: test
test: ## run the go tests
	@echo -e "\033[92m  ---> Testing ... \033[0m"
	@# make symbolic link to ./web/public -> ./pkg/server/web/dist so that we can test out that
	@# static files are being served by the Echo web server
	@if [[ ! -d "$$(pwd)/pkg/server/web/dist" ]]; then \
		echo -e "\033[92m  ---> Linking $(shell pwd)/pkg/server/web/dist -> $(shell pwd)/web/public \033[0m"; \
		mkdir -p $(shell pwd)/pkg/server/web; \
		ln -s $(shell pwd)/web/public $(shell pwd)/pkg/server/web/dist; \
	fi
	go test \
		-mod=vendor \
		-v \
		-cover \
		./...

.PHONY: test_coverage
test_coverage: ## run the go tests then open up coverage report
	@echo -e "\033[92m  ---> Testing wth coverage ... \033[0m"
	go test \
		-mod=vendor \
		-v \
		-cover \
		-coverprofile=$(shell pwd)/coverage.out \
		./...
	go tool cover \
		-html=$(shell pwd)/coverage.out


.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
