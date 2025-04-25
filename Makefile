# Variables
CMD_DIR := backend/cmd
VERSION ?= "0.1.0"
LINTER := golangci-lint

# Docker
DOCKER_IMAGE := hyde
DOCKER_CONTAINER := hyde-playground
DOCKER_TAG := v0.1.0

# Source and Output Files/Dirs
CLI_SOURCE := $(CMD_DIR)/cli/main.go 
API_SOURCE := $(CMD_DIR)/api/main.go 

TEST_FILE := test-file.txt
ENCRYPTED_FILE := encrypted.enc

# Password - Allow overriding via environment variable (e.g., make encrypt PASSWORD=example)
PASSWORD ?= "abc123"

# Default target (runs when you just type `make`)
.DEFAULT_GOAL := help

# Phony targets are rules that don't produce an output file with the name of the target.
.PHONY: build-cli lint format encrypt decrypt help 

build-cli-bulk:
	@bash ./scripts/build-cli.bash $(VERSION)

encrypt:
	go run $(CLI_SOURCE) encrypt -i $(TEST_FILE) -o $(ENCRYPTED_FILE) -p $(PASSWORD)
	@rm $(TEST_FILE)

decrypt:
	go run $(CLI_SOURCE) decrypt -i $(ENCRYPTED_FILE) -o $(TEST_FILE) -p $(PASSWORD)
	@rm $(ENCRYPTED_FILE)

dev:
	cd frontend && pnpm run dev &
	go run $(API_SOURCE)

lint:
	$(LINTER) run
	cd frontend && pnpm run lint

lint-fix:
	$(LINTER) run --fix
	cd frontend && pnpm run lint:fix

format:
	$(LINTER) fmt
	cd frontend && pnpm run format

act-release:
	act --var-file .github/.env --workflows .github/workflows/release.yml

encrypt-help:
	go run $(CLI_SOURCE) encrypt --help

decrypt-help:
	go run $(CLI_SOURCE) decrypt --help

help:
	go run $(CLI_SOURCE) --help

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-build-no-cache:
	docker build --no-cache -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker rm -f $(DOCKER_CONTAINER) || true
	docker run -d --name $(DOCKER_CONTAINER) -h $(DOCKER_CONTAINER) -p 3000:3000 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-exec:
	docker exec -it -u hyde-user $(DOCKER_CONTAINER) /bin/bash