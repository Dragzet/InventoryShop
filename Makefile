# Makefile for InventoryShop
# Основные команды: сборка образов, запуск/остановка compose, логи, запуск CLI, тесты

DOCKER_COMPOSE ?= docker compose
CLIENT = $(DOCKER_COMPOSE) run --rm client

.PHONY: all build build-all build-inventory build-orders build-frontend build-client up down restart logs logs-inventory logs-orders logs-frontend logs-client ps client test test-all test-inventory test-orders fmt clean

all: build-all

# Build targets
build-all:
	$(DOCKER_COMPOSE) build

build-inventory:
	docker build -t inventory-service:latest ./services/inventory

build-orders:
	docker build -t orders-service:latest ./services/orders

build-frontend:
	docker build -t inventory-frontend:latest ./project

build-client:
	docker build -t inventory-client:latest ./cmd/client

# Compose lifecycle
up:
	$(DOCKER_COMPOSE) up -d --build

down:
	$(DOCKER_COMPOSE) down

restart: down up

ps:
	$(DOCKER_COMPOSE) ps

# Logs
logs:
	$(DOCKER_COMPOSE) logs -f

logs-inventory:
	$(DOCKER_COMPOSE) logs -f inventory

logs-orders:
	$(DOCKER_COMPOSE) logs -f orders

logs-frontend:
	$(DOCKER_COMPOSE) logs -f frontend

logs-client:
	$(DOCKER_COMPOSE) logs -f client

# Run client CLI (пример: make client CMD="create-item \"Hat\" 10 9.99")
client:
	@sh -c 'if [ -z "$(CMD)" ]; then echo "usage: make client CMD=\"create-order 1:2,2:1\""; exit 1; fi'
	$(CLIENT) sh -c 'exec /client $(CMD)'

# Tests
test-inventory:
	cd services/inventory && go test ./... -v

test-orders:
	cd services/orders && go test ./... -v

test: test-inventory test-orders

# Misc
fmt:
	go fmt ./...

# Cleanup images/volumes created by compose
clean:
	$(DOCKER_COMPOSE) down --rmi local --volumes --remove-orphans
