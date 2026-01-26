.PHONY: build clean test install build-linux build-windows

# Build variables
VERSION ?= dev
GIT_COMMIT ?= unknown
BUILD_DATE ?= unknown
LDFLAGS := -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)

# Default target builds for Linux (for containers)
build: build-linux

build-linux:
	@echo "Building GameKeeper for Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o gamekeeper .

build-windows:
	@echo "Building GameKeeper for Windows..."
	go build -ldflags "$(LDFLAGS)" -o gamekeeper.exe .

install:
	@echo "Installing GameKeeper..."
	@go install -ldflags "$(LDFLAGS)" .

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -f gamekeeper
	@go clean

mod:
	@echo "Tidying go modules..."
	@go mod tidy

deps:
	@echo "Downloading dependencies..."
	@go mod download

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@golangci-lint run

# Docker build (for testing in container)
docker-build:
	@echo "Building Docker image..."
	@docker build -t gamekeeper:$(VERSION) .
