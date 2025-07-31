# Makefile for ocloud CLI tool

# Application Details
APP_NAME := ocloud
PKG := github.com/rozdolsky33/ocloud
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
OUTPUT_DIR := bin
PLATFORMS := darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 linux/amd64 linux/arm64
GOLANGCI_LINT := golangci-lint

# Determine GOOS and GOARCH for local build
ifeq ($(OS),Windows_NT)
	GOOS_COMPILE := windows
else
	GOOS_COMPILE := $(shell uname -s | tr '[:upper:]' '[:lower:]')
endif

GOARCH_COMPILE := $(shell uname -m | sed -e 's/x86_64/amd64/' -e 's/i[3-6]86/386/')

# Build Flags
LDFLAGS := -X '$(PKG)/buildinfo.Version=$(VERSION)' \
           -X '$(PKG)/buildinfo.CommitHash=$(COMMIT_HASH)' \
           -X '$(PKG)/buildinfo.BuildTime=$(BUILD_TIME)'

# Default target
.DEFAULT_GOAL := help

# Targets
.PHONY: all build run install test fmt fmt-check vet lint clean help generate release compile zip check-env

all: build

# Build the application for the local environment
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) .

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(OUTPUT_DIR)/$(APP_NAME)

# Install the application
install:
	@echo "Installing $(APP_NAME)..."
	@go install -ldflags="$(LDFLAGS)" .

# Run tests
test:
	@echo "Running tests..."
	@go test -v -count=1 ./... -cover

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Check if code is formatted
fmt-check:
	@echo "Checking if code is formatted..."
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not gofmted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Lint code
lint:
	@echo "Linting code..."
	@$(GOLANGCI_LINT) run --no-config ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)

# Release target for multiple platforms
release: clean vet lint compile zip

# Compile binaries for multiple platforms
compile:
	@echo "Compiling binaries for multiple platforms..."
	@mkdir -p $(OUTPUT_DIR)
	$(foreach platform, $(PLATFORMS), \
		$(eval os_arch = $(subst /, ,$(platform))) \
		GOOS=$(word 1,${os_arch}) GOARCH=$(word 2,${os_arch}) \
		go build -v -ldflags="$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(APP_NAME)_$(VERSION)_$(word 1,${os_arch})_$(word 2,${os_arch});)

# Zip compiled binaries
zip:
	@echo "Creating ZIP archives for binaries..."
	@find $(OUTPUT_DIR) -type f ! -name "*.zip" -exec zip -j {}.zip {} \;

# Check environment setup
check-env:
	@echo "Checking environment..."
	@command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed. Aborting."; exit 1; }
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo >&2 "golangci-lint is not installed. Run 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest' to install it."; exit 1; }
	@command -v zip >/dev/null 2>&1 || { echo >&2 "Zip is not installed. Aborting."; exit 1; }

# Help target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all (default)  Builds the project"
	@echo "  build          Compiles the binary into $(OUTPUT_DIR)/$(APP_NAME)"
	@echo "  run            Builds and runs the CLI"
	@echo "  install        Installs the binary to \$(GOBIN) or \$\$(go env GOPATH)/bin"
	@echo "  test           Runs all tests"
	@echo "  fmt            Formats Go source files"
	@echo "  fmt-check      Checks if Go source files are formatted correctly"
	@echo "  vet            Runs go vet on the code"
	@echo "  lint           Runs golangci-lint on the code"
	@echo "  clean          Removes build artifacts"
	@echo "  release        Builds binaries for all supported platforms and creates zip archives"
	@echo "  compile        Compiles binaries for all supported platforms"
	@echo "  zip            Creates zip archives for all binaries"
	@echo "  check-env      Checks if the required tools are installed"
	@echo "  help           Displays this help message"
