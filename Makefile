.PHONY: build run clean docker-build docker-run docker-stop test help

# Variables
BINARY_NAME=proxypal
DOCKER_IMAGE=proxypal-nvidia:latest
CONFIG_FILE=config.yaml

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/proxypal
	@chmod +x $(BINARY_NAME)
	@echo "Build complete! Binary: ./$(BINARY_NAME)"

run: ## Run the application
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Error: $(CONFIG_FILE) not found. Copy config.example.yaml to config.yaml and add your API keys."; \
		exit 1; \
	fi
	@./$(BINARY_NAME)

dev: build run ## Build and run the application

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run with Docker Compose
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Error: $(CONFIG_FILE) not found. Copy config.example.yaml to config.yaml and add your API keys."; \
		exit 1; \
	fi
	@docker-compose up -d
	@echo "ProxyPal is running in Docker"
	@echo "Check logs with: docker-compose logs -f"

docker-stop: ## Stop Docker Compose
	@docker-compose down
	@echo "ProxyPal stopped"

docker-logs: ## Show Docker logs
	@docker-compose logs -f

test: ## Run tests
	@go test -v ./...

init: ## Initialize config file from example
	@if [ -f $(CONFIG_FILE) ]; then \
		echo "$(CONFIG_FILE) already exists!"; \
	else \
		cp config.example.yaml $(CONFIG_FILE); \
		echo "Created $(CONFIG_FILE) from example. Please edit it and add your API keys."; \
	fi

deps: ## Download dependencies
	@go mod download
	@go mod tidy

check: ## Check for issues
	@echo "Running go vet..."
	@go vet ./...
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Check complete!"
