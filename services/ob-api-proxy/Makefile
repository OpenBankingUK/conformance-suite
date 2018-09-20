SHELL=/bin/bash

GIT_REV=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: test
all: test


.PHONY: deps_up
deps_up:
	@echo -e "\033[92m  ---> Starting Dependencies ...  \033[0m"
	docker-compose up -d --force-recreate --remove-orphans mongo redis reference-mock-server

.PHONY: deps_down
deps_down:
	@echo -e "\033[92m  ---> Stopping Dependencies ...  \033[0m"
	docker-compose down --volumes --remove-orphans


.PHONY: init_aspsps
init_aspsps:
	@echo -e "\033[92m  ---> Initialise list of ASPSPs ...  \033[0m"
	npm run updateAuthServersAndOpenIds
	npm run saveCreds authServerId=aaaj4NmBD8lQxmLh2O clientId=spoofClientId clientSecret=spoofClientSecret
	npm run saveCreds authServerId=bbbX7tUB4fPIYB0k1m clientId=spoofClientId clientSecret=spoofClientSecret
	npm run saveCreds authServerId=cccbN8iAsMh74sOXhk clientId=spoofClientId clientSecret=spoofClientSecret

.PHONY: test
test:
	make deps_up
	@echo -e "\033[92m  ---> Linting ... \033[0m"
	npm run eslint
	@echo -e "\033[92m  ---> Testing ...  \033[0m"
	npm run test
	make deps_down

.PHONY: image_build
image_build:
	@echo -e "\033[92m  ---> image: building openbankinguk/tpp-reference-server:latest \033[0m"
	docker-compose build
	@echo -e "\033[92m  ---> image: building openbankinguk/tpp-reference-server:$(GIT_REV) \033[0m"
	COMPOSE_BUILD_TAG=$(GIT_REV) docker-compose build

.PHONY: image_push
image_push:
	@echo -e "\033[92m  ---> image: pushing openbankinguk/tpp-reference-server:latest \033[0m"
	docker-compose push
	@echo -e "\033[92m  ---> image: pushing openbankinguk/tpp-reference-server:$(GIT_REV) \033[0m"
	COMPOSE_BUILD_TAG=$(GIT_REV) docker-compose push
