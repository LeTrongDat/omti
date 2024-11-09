# Makefile for the omti CLI tool

# Project details
BINARY_NAME := omti
INSTALL_DIR := $(GOPATH)/bin

# If GOPATH is not set, use default Go install path
ifeq ($(GOPATH),)
    INSTALL_DIR := $(HOME)/go/bin
endif

# Default target: build the binary
.PHONY: all
all: build

# Build with optimizations
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) with optimizations..."
	go build -ldflags="-s -w" -o $(BINARY_NAME)

# Install the binary in $GOPATH/bin or $HOME/go/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully to $(INSTALL_DIR)"

# Clean up binary
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@echo "Cleaned up."

# Uninstall the binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully from $(INSTALL_DIR)"

