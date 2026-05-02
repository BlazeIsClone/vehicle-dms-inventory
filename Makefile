SHELL=/bin/bash

all: build test

bin:
	@mkdir -p bin

build: bin
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/worker cmd/worker/main.go

run:
	@go run cmd/api/main.go

run-worker:
	@go run cmd/worker/main.go

docker-run:
	@bash scripts/docker-run.sh

docker-down:
	@bash scripts/docker-down.sh

test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

localstack-up:
	@docker compose up localstack -d

localstack-down:
	@docker compose stop localstack

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Live Reload
watch:
	@bash scripts/watch.sh

migrate:
	@bash scripts/migrate.sh $(action)

.PHONY: all bin build run run-worker test clean watch docker-run docker-down itest localstack-up localstack-down
