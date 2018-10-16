SHELL:=/bin/bash
GOMAXPROCS:=12
PARALLEL:=${GOMAXPROCS}

.PHONY: all
all: init lint test build

.PHONY: run
run:
	@echo -e "\033[92m  ---> Starting web file watcher ... \033[0m"
	cd web && FORCE_COLOR=1 NODE_DISABLE_COLORS=0 yarn build-watch &> $(shell pwd)/web/web.log &
	@echo -e "\033[92m  ---> Starting server ... \033[0m"
	go run main.go

.PHONY: run_image
run_image:
	@echo -e "\033[92m  ---> Running image ... \033[0m"
	docker run --rm -it -p 8080:8080 "openbanking/conformance-suite:latest"

.PHONY: build
build:
	@echo -e "\033[92m  ---> Building ... \033[0m"
	go build

.PHONY: build_image
build_image:
	@echo -e "\033[92m  ---> Building image ... \033[0m"
	docker build -t "openbanking/conformance-suite:latest" .

.PHONY: init
init:
	@echo -e "\033[92m  ---> Initialising ... \033[0m"
	go get -d -v ./...

.PHONY: lint
lint:
	@echo -e "\033[92m  ---> Linting ... \033[0m"
	gofmt -e -s -w .

.PHONY: test
test:
	@echo -e "\033[92m  ---> Testing ... \033[0m"
	GOMAXPROCS=${GOMAXPROCS} go test \
		-v \
		-cover \
		-parallel ${PARALLEL} \
		./...

.PHONY: test_coverage
test_coverage:
	@echo -e "\033[92m  ---> Testing wth coverage ... \033[0m"
	-GOMAXPROCS=${GOMAXPROCS} go test \
		-v \
		-cover \
		-parallel ${PARALLEL} \
		-coverprofile=$(shell pwd)/coverage.out \
		./...
	go tool cover \
		-html=$(shell pwd)/coverage.out
