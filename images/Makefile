SHELL=/bin/bash

.PHONY: all
all: build push

.PHONY: build
build:
	@for d in `find ./ -type d`; do \
		if [[ -f "$$d/docker-compose.yml" && -f "$$d/Dockerfile" ]]; then \
			echo -e "\033[92m" "Building $$d..." "\033[0m"; \
			(cd "$$d" && docker-compose build); \
			if [[ $$? == 0 ]]; then \
				echo -e "\033[92m" "OK $$d..." "\033[0m"; \
			else \
				echo -e "\033[91m" "FAIL $$d..." "\033[0m"; \
				exit 1; \
			fi; \
		fi; \
	done

.PHONY: push
push:
	@for d in `find ./ -type d`; do \
		if [[ -f "$$d/docker-compose.yml" && -f "$$d/Dockerfile" ]]; then \
			echo -e "\033[92m" "Pushing $$d..." "\033[0m"; \
			(cd "$$d" && docker-compose push); \
			if [[ $$? == 0 ]]; then \
				echo -e "\033[92m" "OK $$d..." "\033[0m"; \
			else \
				echo -e "\033[91m" "FAIL $$d..." "\033[0m"; \
				exit 1; \
			fi; \
		fi; \
	done
