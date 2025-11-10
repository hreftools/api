# Binary name
BINARY_NAME=api

# Build the application
build:
	PORT=3000 go build -o ${BINARY_NAME} cmd/api/main.go

# Run the built binary (production-like)
run: build
	./${BINARY_NAME}

# Development mode with live reload
dev:
	air

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Install development tools
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Default target (what runs when you just type 'make')
.PHONY: build run dev clean test test-coverage install-tools
