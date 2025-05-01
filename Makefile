# Variables
CMD_DIR := backend/cmd
VERSION ?= "0.1.0"
LINTER := golangci-lint

# Docker
DOCKER_IMAGE := hyde
DOCKER_CONTAINER := hyde-playground

# Source and Output Files/Dirs
CLI_SOURCE := $(CMD_DIR)/cli/main.go
API_SOURCE := $(CMD_DIR)/api/main.go

TEST_FILE := test-file.txt
ENCRYPTED_FILE := test-file.txt.enc

# Password - Allow overriding via environment variable (e.g., make encrypt PASSWORD=example)
PASSWORD ?= "abc123"

# Default target 
.DEFAULT_GOAL := help

# Phony targets are rules that don't produce an output file with the name of the target.
.PHONY: help build-cli build-cli-bulk encrypt decrypt dev go-lint go-lint-fix web-lint web-lint-fix lint test format go-format web-format act-release docker-build docker-build-no-cache docker-run docker-exec

# ====================== Build targets =====================

build-cli-bulk:
	@bash ./scripts/build-cli.bash $(VERSION)

build-cli:
	go build -o ./hyde $(CLI_SOURCE)

# ===================== Development targets =====================

dev:
	cd frontend && pnpm run dev &
	go run $(API_SOURCE)

api-dev:
	go run $(API_SOURCE)

web-dev:
	cd frontend && pnpm run dev

encrypt:
	go run $(CLI_SOURCE) encrypt $(TEST_FILE) --password $(PASSWORD) --delete

decrypt:
	go run $(CLI_SOURCE) decrypt $(ENCRYPTED_FILE) --password $(PASSWORD) --delete

# ===================== Linting targets =====================

go-lint:
	$(LINTER) run

go-lint-fix:
	$(LINTER) run --fix

web-lint:
	cd frontend && pnpm lint 

web-lint-fix:
	cd frontend && pnpm lint --fix

lint:
	$(LINTER) run
	cd frontend && pnpm lint

# ===================== Testing targets =====================

test:
	go test ./...

# ===================== Formatters =====================

format:
	$(LINTER) fmt
	cd frontend && pnpm run format

go-format:
	$(LINTER) fmt

web-format:
	cd frontend && pnpm format

# ===================== Local action for testing =====================

act-release:
	act --var-file .github/.env --workflows .github/workflows/release.yml

# ===================== Help targets =====================

encrypt-help:
	go run $(CLI_SOURCE) encrypt --help

decrypt-help:
	go run $(CLI_SOURCE) decrypt --help

help:
	go run $(CLI_SOURCE) --help

# ===================== Docker targets =====================

docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

docker-build-no-cache:
	docker build --no-cache -t $(DOCKER_IMAGE):$(VERSION) .

docker-run:
	docker rm -f $(DOCKER_CONTAINER) || true
	docker run -d --name $(DOCKER_CONTAINER) -h $(DOCKER_CONTAINER) -p 3000:3000 -p 8080:8080 $(DOCKER_IMAGE):$(VERSION)

docker-exec:
	docker exec -it -u hyde-user $(DOCKER_CONTAINER) /bin/bash