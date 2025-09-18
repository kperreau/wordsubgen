# Makefile for wordsubgen

# Variables
BINARY_NAME=wordsubgen
CLI_BINARY=cmd/wordsubgen/wordsubgen
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: build

# Build the CLI binary
.PHONY: build
build: $(CLI_BINARY)

$(CLI_BINARY): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(CLI_BINARY) ./cmd/wordsubgen
	@echo "✓ Built $(CLI_BINARY)"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ Tests completed"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Run CLI with sample input
.PHONY: run-cli
run-cli: build
	@echo "Running CLI with sample input..."
	@./$(CLI_BINARY) --lines "Hello world|This is a test|Word by word fade effect" --out sample.ass
	@echo "✓ Generated sample.ass"

# Run CLI with custom styling
.PHONY: run-cli-custom
run-cli-custom: build
	@echo "Running CLI with custom styling..."
	@./$(CLI_BINARY) --lines "Custom styling|Red text|Larger font" --out custom.ass --color "#FF0000" --fontsize 72 --delay 500
	@echo "✓ Generated custom.ass"

# Run CLI with karaoke mode
.PHONY: run-cli-karaoke
run-cli-karaoke: build
	@echo "Running CLI with karaoke mode..."
	@./$(CLI_BINARY) --lines "Karaoke mode|Word by word|Highlight effect" --out karaoke.ass --karaoke --delay 400
	@echo "✓ Generated karaoke.ass"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(CLI_BINARY)
	@rm -f *.ass
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

# Install gofumpt
.PHONY: install-fmt
install-fmt:
	@echo "Installing gofumpt..."
	@if ! command -v gofumpt >/dev/null 2>&1; then \
		go install mvdan.cc/gofumpt@latest; \
	fi
	@echo "✓ gofumpt installed"

# Format code
.PHONY: fmt
fmt: install-fmt
	@echo "Formatting code with gofumpt..."
	@gofumpt -l -w .
	@echo "✓ Code formatted"

# Install golangci-lint
.PHONY: install-lint
install-lint:
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	fi
	@echo "✓ golangci-lint installed"

# Lint code
.PHONY: lint
lint: install-lint
	@echo "Linting code with golangci-lint..."
	@golangci-lint run
	@echo "✓ Code linted"

# Run all checks
.PHONY: check
check: fmt lint test
	@echo "✓ All checks passed"

# Install the CLI to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@go install $(LDFLAGS) ./cmd/wordsubgen
	@echo "✓ Installed $(BINARY_NAME)"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the CLI binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run-cli        - Run CLI with sample input"
	@echo "  run-cli-custom - Run CLI with custom styling"
	@echo "  run-cli-karaoke- Run CLI with karaoke mode"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  install-fmt    - Install gofumpt"
	@echo "  fmt            - Format code with gofumpt"
	@echo "  install-lint   - Install golangci-lint"
	@echo "  lint           - Lint code with golangci-lint"
	@echo "  check          - Run all checks (fmt, lint, test)"
	@echo "  install        - Install CLI to GOPATH/bin"
	@echo "  help           - Show this help"

# Development workflow
.PHONY: dev
dev: deps fmt lint test build
	@echo "✓ Development build completed"

.PHONY: ci-test
ci-test:
	@./_scripts/ci-test.sh

.PHONY: vet
vet:
	@go vet ./...

# Release build
.PHONY: release
release: clean deps fmt vet lint test
	@echo "Building release binaries..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/wordsubgen
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/wordsubgen
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/wordsubgen
	@echo "✓ Release binaries built in dist/"
