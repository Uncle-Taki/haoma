# Haoma - Black-Box Carnival Makefile
# Persian god meets Go development

.PHONY: help deps build run test lint clean seed docker dev

# Default target
help: ## Show this help message
	@echo "🎪 Haoma - Black-Box Carnival"
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

deps: ## Install Go dependencies
	@echo "🔮 Gathering mystical dependencies..."
	go mod tidy
	go mod download

build: deps ## Build the carnival
	@echo "🏗️ Building the carnival..."
	go build -o bin/haoma cmd/server/main.go

run: build ## Start the carnival (development)
	@echo "🎪 Opening Haoma's carnival at :8080..."
	./bin/haoma

dev: deps ## Run with hot reload (requires air)
	@echo "🔄 Starting carnival with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

test: deps ## Run all tests
	@echo "🧪 Testing the carnival's mysteries..."
	go test -v ./...

test-coverage: deps ## Run tests with coverage
	@echo "📊 Analyzing test coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: deps ## Lint the codebase
	@echo "🔍 Examining code for impurities..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

fmt: ## Format code
	@echo "✨ Beautifying the carnival's code..."
	go fmt ./...
	goimports -w .

vet: ## Vet the code
	@echo "🔎 Vetting code quality..."
	go vet ./...

seed: build ## Seed database with sample data
	@echo "🌱 Seeding the carnival with wisdom..."
	@mkdir -p cmd/seed && echo "package main\n\nimport (\n\t\"log\"\n\t\"haoma/internal/infrastructure/persistence\"\n\t\"haoma/internal/infrastructure/seeder\"\n)\n\nfunc main() {\n\tdb, err := persistence.NewDatabase()\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\ts := seeder.NewExcelSeeder(db.DB)\n\tif err := s.CreateSampleData(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}" > cmd/seed/main.go
	go run cmd/seed/main.go
	rm -f cmd/seed/main.go

seed-excel: build ## Seed from Excel files (data/SCENARIOS.xlsx, data/questions.xlsx)
	@echo "📊 Loading wisdom from Excel scrolls..."
	@if [ ! -f "data/SCENARIOS.xlsx" ] || [ ! -f "data/questions.xlsx" ]; then \
		echo "⚠️  Excel files not found. Creating sample data instead..."; \
		make seed; \
	else \
		mkdir -p cmd/seed && echo "package main\n\nimport (\n\t\"log\"\n\t\"haoma/internal/infrastructure/persistence\"\n\t\"haoma/internal/infrastructure/seeder\"\n)\n\nfunc main() {\n\tdb, err := persistence.NewDatabase()\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\ts := seeder.NewExcelSeeder(db.DB)\n\tif err := s.SeedFromExcel(\"data/SCENARIOS.xlsx\", \"data/questions.xlsx\"); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}" > cmd/seed/main.go; \
		go run cmd/seed/main.go; \
		rm -f cmd/seed/main.go; \
	fi

swagger: ## Generate Swagger documentation
	@echo "📚 Generating API scrolls..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o api; \
	else \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/server/main.go -o api; \
	fi

clean: ## Clean build artifacts
	@echo "🧹 Cleaning the carnival grounds..."
	rm -rf bin/
	rm -rf api/docs.go api/swagger.json api/swagger.yaml
	rm -f haoma.db
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build Docker image
	@echo "🐳 Containerizing the carnival..."
	docker build -t haoma:latest .

docker-run: docker-build ## Run in Docker
	@echo "🐳 Running carnival in container..."
	docker run -p 8080:8080 --rm haoma:latest

install-tools: ## Install development tools
	@echo "🔧 Installing carnival tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Quality checks
check: lint vet test ## Run all quality checks

# Development workflow
# PostgreSQL operations
pg-start: ## Start PostgreSQL in Docker
	@echo "🐘 Starting PostgreSQL database..."
	docker-compose -f docker-compose.dev.yml up -d postgres

pg-stop: ## Stop PostgreSQL
	@echo "🛑 Stopping PostgreSQL database..."
	docker-compose -f docker-compose.dev.yml stop postgres

pg-logs: ## View PostgreSQL logs
	@echo "📜 PostgreSQL logs:"
	docker-compose -f docker-compose.dev.yml logs -f postgres

pg-shell: ## Connect to PostgreSQL shell
	@echo "🐚 Connecting to PostgreSQL..."
	docker exec -it haoma-postgres psql -U haoma -d haoma

# Full development environment
dev-env: ## Start complete development environment
	@echo "🎪 Starting Haoma's complete carnival environment..."
	docker-compose -f docker-compose.dev.yml up -d

dev-env-stop: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

dev-env-build: ## Build development environment
	@echo "🏗️ Building development environment..."
	docker-compose -f docker-compose.dev.yml up --build -d


dev-setup: deps install-tools pg-start swagger seed ## Complete development setup
	@echo "🎪 Haoma's carnival is ready for development!"
	@echo ""
	@echo "Quick start:"
	@echo "  make dev     # Start with hot reload (ensure pg-start first)"
	@echo "  make run     # Start normally (ensure pg-start first)"
	@echo "  make dev-env # Start full environment (app + db)"
	@echo "  make test    # Run tests"
	@echo ""
	@echo "Database:"
	@echo "  make pg-start    # Start PostgreSQL"
	@echo "  make pg-stop     # Stop PostgreSQL"
	@echo "  make pg-shell    # Connect to DB"
	@echo ""
	@echo "Explore:"
	@echo "  http://localhost:8080/docs    # Swagger UI"
	@echo "  http://localhost:8080/health  # Health check"

.DEFAULT_GOAL := help
