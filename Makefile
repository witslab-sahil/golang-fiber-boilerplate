.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the application
	go run cmd/api/main.go

.PHONY: build
build: ## Build the application
	go build -o bin/api cmd/api/main.go

.PHONY: test
test: ## Run tests
	go test -v -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: mod
mod: ## Download and tidy modules
	go mod download
	go mod tidy

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t fiber-boilerplate:latest .

.PHONY: docker-run
docker-run: ## Run docker container
	docker run -p 8080:8080 --env-file .env fiber-boilerplate:latest

.PHONY: docker-up
docker-up: ## Start docker compose services
	docker compose up -d

.PHONY: docker-down
docker-down: ## Stop docker compose services
	docker compose down

.PHONY: dev
dev: ## Start development environment with hot reload
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up

.PHONY: dev-down
dev-down: ## Stop development environment
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

.PHONY: docker-logs
docker-logs: ## View docker compose logs
	docker compose logs -f

.PHONY: docker-test
docker-test: ## Run tests in docker
	docker compose --profile test run --rm test

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	docker compose --profile tools run --rm migrate

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	docker compose --profile tools run --rm migrate -path=/migrations -database postgres://postgres:postgres@db:5433/golang_boilerplate?sslmode=disable down

.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

.PHONY: install-tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: swagger
swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go -o docs

.PHONY: mock
mock: ## Generate mocks
	mockery --all --output=mocks

.PHONY: bench
bench: ## Run benchmarks
	go test -bench=. -benchmem ./...