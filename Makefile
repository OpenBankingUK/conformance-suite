SHELL:=/bin/bash
GOMAXPROCS:=12
PARALLEL:=${GOMAXPROCS}
# guarantee that go will not reach the network at all (e.g. GOPROXY=off)
export GOPROXY:=off

.PHONY: all
all: help

.PHONY: run
run: init_web ## run binary directly without docker
	@echo -e "\033[92m  ---> Starting web file watcher ... \033[0m"
	cd web && FORCE_COLOR=1 NODE_DISABLE_COLORS=0 yarn build-watch &> $(shell pwd)/web/web.log &
	@echo -e "\033[92m  ---> Starting server ... \033[0m"
	PORT=8080 go run main.go

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
build: ## build the binary directly
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build -mod vendor

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

.PHONY: lint
lint: ## lint the go code
	@echo -e "\033[92m  ---> Linting ... \033[0m"
	gofmt -e -s -w .

.PHONY: clean
clean:
	@echo -e "\033[92m  ---> Cleaning ... \033[0m"
	go clean -i -r -cache -testcache -modcache

.PHONY: test
test: ## run the go tests
	@echo -e "\033[92m  ---> Testing ... \033[0m"
	@# make symbolic link to ./web/public -> /lib/server/web/dist so that we can test out that
	@# static files are being served by the Echo web server
	@if [[ ! -d "$$(pwd)/lib/server/web/dist" ]]; then \
		echo -e "\033[92m  ---> Linking $(shell pwd)/lib/server/web/dist -> $(shell pwd)/web/public \033[0m"; \
		mkdir -p $(shell pwd)/lib/server/web; \
		ln -s $(shell pwd)/web/public $(shell pwd)/lib/server/web/dist; \
	fi
	GOMAXPROCS=${GOMAXPROCS} go test \
		-mod vendor \
		-v \
		-cover \
		-parallel ${PARALLEL} \
		./...

.PHONY: test_coverage
test_coverage: ## run the go tests then open up coverage report
	@echo -e "\033[92m  ---> Testing wth coverage ... \033[0m"
	-GOMAXPROCS=${GOMAXPROCS} go test \
		-mod vendor \
		-v \
		-cover \
		-parallel ${PARALLEL} \
		-coverprofile=$(shell pwd)/coverage.out \
		./...
	go tool cover \
		-html=$(shell pwd)/coverage.out


.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
