.PHONY: help generate generate-watch run build test clean deps dev

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

generate: ## Generate GraphQL code
	@echo "🔄 Generating GraphQL code..."
	@$(shell go env GOPATH)/bin/gqlgen generate
	@echo "✅ GraphQL generation complete!"

generate-watch: ## Watch for changes and auto-generate GraphQL code
	@echo "👀 Starting GraphQL watch mode..."
	@./scripts/watch-gqlgen.sh

ent-generate: ## Generate Ent code
	@echo "🔄 Generating Ent code..."
	@go generate ./ent

build: generate ## Build the application
	@echo "🔨 Building application..."
	@go build -o bin/server ./server/server.go
	@echo "✅ Build complete!"

run: generate ## Run the application
	@echo "🚀 Starting server..."
	@go run ./server/server.go

dev: ## Start development mode (with auto-reload)
	@echo "🛠️  Starting development mode..."
	@air -c .air.toml

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./...

test-coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📈 Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

gql-playground: ## Open GraphQL playground
	@echo "🎮 Opening GraphQL playground..."
	@open http://localhost:8080/graphql

format: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@goimports -w .

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t lingano-go .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 lingano-go

migrate-up: ## Run database migrations up
	@echo "⬆️  Running migrations up..."
	@go run scripts/main_goose.go up

migrate-down: ## Run database migrations down
	@echo "⬇️  Running migrations down..."
	@go run scripts/main_goose.go down

migrate-status: ## Check migration status
	@echo "📋 Checking migration status..."
	@go run scripts/main_goose.go status
