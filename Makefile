SHELL:=/bin/bash
GOMAXPROCS:=12
PARALLEL:=${GOMAXPROCS}

.PHONY: all
all: help

.PHONY: run
run: ## run binary directly without docker
	@echo -e "\033[92m  ---> Starting web file watcher ... \033[0m"
	cd web && FORCE_COLOR=1 NODE_DISABLE_COLORS=0 yarn build-watch &> $(shell pwd)/web/web.log &
	@echo -e "\033[92m  ---> Starting server ... \033[0m"
	go run main.go

.PHONY: run_image
run_image: ## run the docker image
	@echo -e "\033[92m  ---> Running image ... \033[0m"
	docker run --rm -it -p 8080:8080 "openbanking/conformance-suite:latest"

.PHONY: build
build: ## build the binary directly
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build

.PHONY: build_image
build_image: ## build the docker image
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	docker build -t "openbanking/conformance-suite:latest" .

.PHONY: init
init: ## initialise
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go get -v ./...

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
	GOMAXPROCS=${GOMAXPROCS} go test \
		-v \
		-cover \
		-parallel ${PARALLEL} \
		./...


.PHONY: test_coverage
test_coverage: ## run the go tests then open up coverage report
	@echo -e "\033[92m  ---> Testing wth coverage ... \033[0m"
	-GOMAXPROCS=${GOMAXPROCS} go test \
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
