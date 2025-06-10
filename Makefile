# Makefile for ocloud CLI tool

# Configuration
GO        := go
BIN       := ocloud
BUILD_DIR := bin

# Default target
.DEFAULT_GOAL := help

# Targets
.PHONY: all build run install test fmt vet clean help

all: build

build:
	@echo "Building $(BIN)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BIN) .

run: build
	@echo "Running $(BIN)..."
	@$(BUILD_DIR)/$(BIN)

install:
	@echo "Installing $(BIN)..."
	@$(GO) install .

test:
	@echo "Running tests..."
	@$(GO) test ./...

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet:
	@echo "Vet code..."
	@$(GO) vet ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all (default)  Builds the project"
	@echo "  build          Compiles the binary into $(BUILD_DIR)/$(BIN)"
	@echo "  run            Builds and runs the CLI"
	@echo "  install        Installs the binary to \$(GOBIN) or \$\$(go env GOPATH)/bin"
	@echo "  test           Runs all tests"
	@echo "  fmt            Formats Go source files"
	@echo "  vet            Runs go vet on the code"
	@echo "  clean          Removes build artifacts"
	@echo "  help           Displays this help message"