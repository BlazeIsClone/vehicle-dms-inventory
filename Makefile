# Build the application
all: build-all test

build:
	@echo "Building API..."
	@go build -o api cmd/api/main.go

build-worker:
	@echo "Building worker..."
	@go build -o worker_bin cmd/worker/main.go

build-all: build build-worker

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
	@rm -f api worker_bin

# Live Reload
watch:
	@bash scripts/watch.sh

migrate:
	@bash scripts/migrate.sh $(action)

.PHONY: all build build-worker build-all run run-worker test clean watch docker-run docker-down itest localstack-up localstack-down