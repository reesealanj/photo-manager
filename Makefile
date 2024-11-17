# Variables
BINARY_NAME=photo-manager
SRC=$(shell find . -name '*.go')
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUTPUT_DIR=bin
DEV_BINARY=$(OUTPUT_DIR)/$(BINARY_NAME)
PROD_BINARY_LINUX=$(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64
PROD_BINARY_DARWIN=$(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64
PROD_BINARY_WINDOWS=$(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe
MODE ?= dev

# Default target
.PHONY: all
all: build-dev

# Build for development (current platform)
.PHONY: build-dev
build-dev: $(DEV_BINARY)
$(DEV_BINARY): $(SRC)
	@echo "Building for development..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DEV_BINARY) $(SRC)
	@echo "Build complete: $(DEV_BINARY)"

# Build for production (all platforms)
.PHONY: build-prod
build-prod:
	@echo "Building for production..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(PROD_BINARY_LINUX) $(SRC)
	GOOS=darwin GOARCH=amd64 go build -o $(PROD_BINARY_DARWIN) $(SRC)
	GOOS=windows GOARCH=amd64 go build -o $(PROD_BINARY_WINDOWS) $(SRC)
	@echo "Production builds complete."

# Run the binary (based on MODE)
.PHONY: run
run: build-$(MODE)
ifeq ($(MODE),dev)
	@echo "Running the development binary..."
	@$(DEV_BINARY)
else ifeq ($(MODE),prod) && ($(GOOS),darwin)
	@echo "Running the production binary ($(GOOS))..."
	@$(PROD_BINARY_DARWIN)
else ifeq ($(MODE),prod) && ($(GOOS),linux)
	@echo "Running the production binary ($(GOOS))..."
	@$(PROD_BINARY_LINUX)
else ifeq ($(MODE),prod) && ($(GOOS),windows) 
	@$(PROD_BINARY_WINDOWS)
else
	@echo "Invalid MODE: $(MODE). Use 'dev' or 'prod'."
	@exit 1
endif

# Cleanup target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTPUT_DIR)
	@echo "Cleanup complete."

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Default target (alias for build-dev)"
	@echo "  build-dev   - Build for development (current platform)"
	@echo "  build-prod  - Build for production (all platforms)"
	@echo "  run         - Run the binary (use MODE=dev or MODE=prod)"
	@echo "  clean       - Clean up build artifacts"
	@echo "  help        - Show this help message"
