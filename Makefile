# Makefile for Income & Expense Tracker

.PHONY: help build run clean test docker-build docker-run docker-stop deps install

# Variables
BINARY_NAME=myexpress-tracker
MAIN_PATH=./cmd/server
DOCKER_IMAGE=expense-tracker:latest

# Default target
help:
	@echo "Income & Expense Tracker - Available Commands:"
	@echo ""
	@echo "  make build        - Build the application binary"
	@echo "  make run          - Build and run the application"
	@echo "  make clean        - Remove binary and database"
	@echo "  make deps         - Download Go dependencies"
	@echo "  make test         - Run tests"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make install      - Install binary to system"
	@echo ""

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Build complete!"

# Build and run
run: build
	@echo "Starting server on http://localhost:8080..."
	@./$(BINARY_NAME).exe

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME).exe
	@rm -rf data/
	@echo "Clean complete!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built!"

# Run with Docker Compose
docker-run:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d
	@echo "Application running! Visit http://localhost:8080"

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped!"

# Install binary (Linux/Mac)
install: build
	@echo "Installing to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Development mode with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
	fi
