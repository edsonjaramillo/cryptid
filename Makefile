# Variables
CMD_DIR := backend/cmd
LINTER := golangci-lint

# Source and Output Files/Dirs
CLI_SOURCE := $(CMD_DIR)/cli/main.go # Assuming cli.go is the entry point

TEST_FILE := test-file.txt
ENCRYPTED_FILE := encrypted.enc

# Password - Allow overriding via environment variable (e.g., make encrypt PASSWORD=example)
PASSWORD ?= "abc123"

# Default target (runs when you just type `make`)
.DEFAULT_GOAL := help

# Phony targets are rules that don't produce an output file with the name of the target.
.PHONY: prepare lint format encrypt decrypt help 

# Prepare the test file
prepare:
	@[ -f $(TEST_FILE) ] || echo "Hello, World!" > $(TEST_FILE)

# Encrypt the test file
encrypt:
	go run $(CLI_SOURCE) encrypt -i $(TEST_FILE) -o $(ENCRYPTED_FILE) -p $(PASSWORD)
	@rm $(TEST_FILE)

decrypt:
	go run $(CLI_SOURCE) decrypt -i $(ENCRYPTED_FILE) -o $(TEST_FILE) -p $(PASSWORD)
	@rm $(ENCRYPTED_FILE)

lint:
	$(LINTER) run 

format:
	$(LINTER) fmt

help:
	go run $(CLI_SOURCE) --help
