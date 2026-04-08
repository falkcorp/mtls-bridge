BINARY := mtls-bridge
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build build-all test coverage lint clean help

help:
	@echo "Targets:"
	@echo "  make build       - Build for current platform"
	@echo "  make build-all   - Cross-compile for all platforms"
	@echo "  make test        - Run tests"
	@echo "  make coverage    - Run tests with coverage"
	@echo "  make lint        - Run go vet"
	@echo "  make clean       - Remove build artifacts"

build:
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/mtls-bridge

build-all:
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-darwin-amd64 ./cmd/mtls-bridge
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-darwin-arm64 ./cmd/mtls-bridge
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-linux-amd64 ./cmd/mtls-bridge
	@GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-linux-arm64 ./cmd/mtls-bridge
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BINARY)-windows-amd64.exe ./cmd/mtls-bridge

test:
	@go test ./... -v -count=1

coverage:
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out

lint:
	@go vet ./...

clean:
	@rm -f $(BINARY) $(BINARY)-* coverage.out
