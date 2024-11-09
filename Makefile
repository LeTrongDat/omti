# Project details
BINARY_NAME := omti
BUILD_DIR := build
INSTALL_DIR := $(GOPATH)/bin

# If GOPATH is not set, use default Go install path
ifeq ($(GOPATH),)
    INSTALL_DIR := $(HOME)/go/bin
endif

# Default target: build the binary
.PHONY: all
all: build

# Build with optimizations, saving output in build directory
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) with optimizations..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)

# Install the binary from the build directory to the Go bin path
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully to $(INSTALL_DIR)"

# Clean up the build directory
.PHONY: clean
clean:
	@echo "Cleaning up build files..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned up."

# Uninstall the binary from the Go bin path
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully from $(INSTALL_DIR)"
