.PHONY: help build test test-unit test-integration test-e2e test-all coverage run clean docker-build docker-up docker-down

# Variáveis
DOCKER_COMPOSE := docker-compose
DOCKER := docker
GO := go
BINARY_NAME := server

help:
	@echo "╔════════════════════════════════════════════════════════╗"
	@echo "║     Weather Application - Makefile                     ║"
	@echo "╚════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Targets disponíveis:"
	@echo ""
	@echo "  make build              - Build da aplicação"
	@echo "  make test               - Executar todos os testes"
	@echo "  make test-unit          - Executar apenas testes unitários"
	@echo "  make test-integration   - Executar apenas testes de integração"
	@echo "  make test-e2e           - Executar apenas testes E2E"
	@echo "  make test-all           - Executar todos os testes (unit + integration + e2e)"
	@echo "  make coverage           - Gerar relatório de cobertura"
	@echo "  make run                - Executar a aplicação localmente"
	@echo "  make clean              - Limpeza (remover binary)"
	@echo ""
	@echo "  Docker targets:"
	@echo "  make docker-build       - Build da imagem Docker"
	@echo "  make docker-up          - Iniciar containers"
	@echo "  make docker-down        - Parar containers"
	@echo "  make docker-logs        - Visualizar logs dos containers"
	@echo ""

## Build targets
build:
	@echo "Building application..."
	$(GO) build -o $(BINARY_NAME) main.go
	@echo "✓ Build successful! Binary: $(BINARY_NAME)"

## Test targets
test: test-unit test-integration
	@echo "✓ All tests passed!"

test-unit:
	@echo "Running unit tests..."
	$(GO) test -v -race ./internal/... -coverprofile=coverage-unit.out
	@echo "✓ Unit tests completed!"

test-integration:
	@echo "Running integration tests..."
	$(GO) test -v -race -tags=integration ./... -coverprofile=coverage-integration.out
	@echo "✓ Integration tests completed!"

test-e2e:
	@echo "Running E2E tests..."
	$(GO) test -v -tags=e2e ./...
	@echo "✓ E2E tests completed!"

test-all: test-unit test-integration test-e2e
	@echo "✓ All tests (unit + integration + e2e) passed!"

coverage:
	@echo "Generating coverage report..."
	$(GO) test -v -race ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated at coverage.html"

## Application targets
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out coverage-unit.out coverage-integration.out coverage.html
	@echo "✓ Clean complete!"

## Docker targets
docker-build:
	@echo "Building Docker image..."
	$(DOCKER_COMPOSE) build
	@echo "✓ Docker image built!"

docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "✓ Containers started!"

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down
	@echo "✓ Containers stopped!"

docker-logs:
	$(DOCKER_COMPOSE) logs -f

## Combined targets
docker-test-all: docker-build
	@echo "Running all tests with Docker..."
	$(DOCKER_COMPOSE) run --rm test-unit
	$(DOCKER_COMPOSE) run --rm test-integration
	@echo "✓ All Docker tests completed!"

## Development helpers
dev: docker-build docker-up
	@echo "✓ Development environment ready!"
	@echo "  App running at: http://localhost:8080"

dev-stop: docker-down
	@echo "✓ Development environment stopped!"

dev-logs:
	$(DOCKER_COMPOSE) logs -f app
