# The name of your application's binary
BINARY_NAME=smokescreen

# Use the latest Git tag as the version (or commit hash)
VERSION := $(shell git describe --tags --always)

# Define the ldflags for injecting version information
LDFLAGS := -ldflags="-X 'starless.dev/smokescreen/cmd.Version=${VERSION}'"

all: mod fmt lint clean build

# Build the Go application
build:
	@echo "==> Building application..."
	@echo "Version: ${VERSION}"
	go build ${LDFLAGS} -o $(BINARY_NAME) .

# Tidy and verify module dependencies
mod:
	@echo "==> Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify

# Format the source code
fmt:
	@echo "==> Formatting source code..."
	go fmt ./...

# Run the linter
lint: mod
	@echo "==> Running linter..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

# Remove the previously built binary
clean:
	@echo "==> Cleaning up..."
	rm -f $(BINARY_NAME)