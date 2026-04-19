# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

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

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@bash scripts/watch.sh

migrate:
	@bash scripts/migrate.sh $(action)

.PHONY: all build run test clean watch docker-run docker-down itest