BINARY_NAME=commit-gen
GO=go
GOFLAGS=-v
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build clean test install help run release

help:
	@echo "Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Build and install to /usr/local/bin"
	@echo "  make release     - Build cross-platform releases"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run the CLI"
	@echo "  make lint        - Run linter"

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/commit-gen

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	$(GO) clean
	@rm -f $(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GO) test -v ./...

run: build
	@./$(BINARY_NAME)

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && $(GO) install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

mod-tidy:
	@echo "Tidying modules..."
	$(GO) mod tidy

mod-verify:
	@echo "Verifying modules..."
	$(GO) mod verify

all: clean install test

local: clean build 

release:
	@echo "Building cross-platform releases..."
	@./scripts/build-release.sh
	@echo "✓ Release builds available in dist/"
