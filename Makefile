# Go Proxy Rotator Makefile

# Variables
BINARY_NAME=go-proxy-rotator
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p dist
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	
	@echo "Build complete! Binaries are in the dist/ directory."

# Run the application
.PHONY: run
run: build
	./${BINARY_NAME}

# Run with development settings
.PHONY: dev
dev:
	@echo "Running in development mode..."
	PORT=3000 DATABASE_PATH=./dev.db LOG_LEVEL=debug go run .

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -rf dist/
	@rm -f *.db

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	go vet ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Create release packages
.PHONY: package
package: build-all
	@echo "Creating release packages..."
	@mkdir -p dist/packages
	
	# Linux AMD64
	@mkdir -p dist/tmp/${BINARY_NAME}-linux-amd64
	@cp dist/${BINARY_NAME}-linux-amd64 dist/tmp/${BINARY_NAME}-linux-amd64/
	@cp README.md sample_proxies.txt dist/tmp/${BINARY_NAME}-linux-amd64/
	@cp -r static dist/tmp/${BINARY_NAME}-linux-amd64/
	@cd dist/tmp && tar -czf ../packages/${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
	
	# Linux ARM64
	@mkdir -p dist/tmp/${BINARY_NAME}-linux-arm64
	@cp dist/${BINARY_NAME}-linux-arm64 dist/tmp/${BINARY_NAME}-linux-arm64/
	@cp README.md sample_proxies.txt dist/tmp/${BINARY_NAME}-linux-arm64/
	@cp -r static dist/tmp/${BINARY_NAME}-linux-arm64/
	@cd dist/tmp && tar -czf ../packages/${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
	
	# Windows AMD64
	@mkdir -p dist/tmp/${BINARY_NAME}-windows-amd64
	@cp dist/${BINARY_NAME}-windows-amd64.exe dist/tmp/${BINARY_NAME}-windows-amd64/
	@cp README.md sample_proxies.txt dist/tmp/${BINARY_NAME}-windows-amd64/
	@cp -r static dist/tmp/${BINARY_NAME}-windows-amd64/
	@cd dist/tmp && zip -r ../packages/${BINARY_NAME}-windows-amd64.zip ${BINARY_NAME}-windows-amd64
	
	# macOS AMD64
	@mkdir -p dist/tmp/${BINARY_NAME}-darwin-amd64
	@cp dist/${BINARY_NAME}-darwin-amd64 dist/tmp/${BINARY_NAME}-darwin-amd64/
	@cp README.md sample_proxies.txt dist/tmp/${BINARY_NAME}-darwin-amd64/
	@cp -r static dist/tmp/${BINARY_NAME}-darwin-amd64/
	@cd dist/tmp && tar -czf ../packages/${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
	
	# macOS ARM64
	@mkdir -p dist/tmp/${BINARY_NAME}-darwin-arm64
	@cp dist/${BINARY_NAME}-darwin-arm64 dist/tmp/${BINARY_NAME}-darwin-arm64/
	@cp README.md sample_proxies.txt dist/tmp/${BINARY_NAME}-darwin-arm64/
	@cp -r static dist/tmp/${BINARY_NAME}-darwin-arm64/
	@cd dist/tmp && tar -czf ../packages/${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
	
	@rm -rf dist/tmp
	@echo "Packages created in dist/packages/"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t ${BINARY_NAME}:${VERSION} .

# Docker run
.PHONY: docker-run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 3000:3000 -v proxy_data:/data ${BINARY_NAME}:${VERSION}

# Swagger/OpenAPI
.PHONY: swagger-validate
swagger-validate:
	@echo "Validating OpenAPI specification..."
	go run tools/swagger-gen.go validate

.PHONY: swagger-json
swagger-json:
	@echo "Converting OpenAPI spec to JSON..."
	go run tools/swagger-gen.go json

.PHONY: swagger-yaml
swagger-yaml:
	@echo "Converting OpenAPI spec to YAML..."
	go run tools/swagger-gen.go yaml

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  build-all       - Build for all platforms"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run in development mode"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  deps            - Install dependencies"
	@echo "  package         - Create release packages"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Build and run Docker container"
	@echo "  swagger-validate - Validate OpenAPI specification"
	@echo "  swagger-json    - Convert OpenAPI spec to JSON"
	@echo "  swagger-yaml    - Convert OpenAPI spec to YAML"
	@echo "  help            - Show this help message"