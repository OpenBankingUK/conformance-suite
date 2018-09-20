SHELL:=/bin/bash
GIT_REV:=$(shell git rev-parse --short HEAD)

.PHONY: all
all: image_build image_push

.PHONY: image_build
image_build:
	@echo -e "\033[92m  ---> image: building \033[0m"
	docker-compose build

.PHONY: image_push
image_push:
	@echo -e "\033[92m  ---> image: pushing \033[0m"
	docker-compose push
