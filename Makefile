# Makefile for IS CLI on macOS

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=is-cli
MAIN_PATH=.

# Installation directory (should be in your PATH)
INSTALL_PATH=$(HOME)/bin

.PHONY: build run install

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Install the application
install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) has been installed to $(INSTALL_PATH)"
	@echo "Make sure $(INSTALL_PATH) is in your PATH"

# Default target
.DEFAULT_GOAL := run
