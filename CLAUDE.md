# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build both binaries (api + worker)
make build

# Run API server
make run

# Run SQS worker
make run-worker

# Live reload (uses air)
make watch

# Run all tests
make test

# Run a single test or package
go test ./inventory/vehicle/... -v
go test ./internal/database/... -v -run TestName

# DB integration tests (requires Docker)
make itest

# Start Postgres + LocalStack containers
make docker-run

# Start LocalStack only (SNS/SQS)
make localstack-up

# Run migrations (up or down)
make migrate action=up
make migrate action=down
```

Copy `.env.example` to `.env` before running locally — the defaults point at LocalStack and a local Postgres instance.

## Architecture

The service is split into two independent processes:

**`cmd/api`** — HTTP server that handles vehicle CRUD and runs the outbox relay goroutine.

**`cmd/worker`** — Long-running SQS consumer that processes domain events published by the API.

### Event flow (outbox pattern)

1. The service layer (`inventory/vehicle/service.go`) writes the domain record and saves domain events to `outbox_events` within a single transaction via `outbox.Store.Save`.
2. `outbox.Relay` (started in `cmd/api/main.go`) polls `outbox_events` every 5 s for unpublished rows and publishes them to SNS via `infra/sns.Publisher` (Watermill-wrapped).
3. SQS is subscribed to the SNS topic. `infra/sqs.Consumer` (Watermill router) receives messages and dispatches to registered `events.HandlerFunc` handlers.
4. `outbox.ProcessedStore` (backed by `processed_events` table) prevents duplicate processing — handlers should check this before executing side effects.

### Package layout

- `cmd/api`, `cmd/worker`, `cmd/migrate` — entry points; wire dependencies, no business logic.
- `inventory/vehicle/` — domain package: `Vehicle` aggregate, CRUD service, HTTP handler, Postgres repo, event types, event handlers. All business logic lives here.
- `internal/outbox/` — transactional outbox: `Store` (write/query), `Relay` (poll+publish), `ProcessedStore` (idempotency).
- `internal/events/` — shared types: `DomainEvent`, `Publisher` interface, `HandlerFunc`.
- `internal/database/` — `database.Service` interface, pgx connection, migrations runner.
- `internal/aws/` — AWS config loader, `EnsureResources` (idempotent SNS topic + SQS queue + subscription setup used by the worker on startup).
- `infra/sns/`, `infra/sqs/` — Watermill adapters for SNS publishing and SQS consuming.
- `pkg/api/` — `Path(version, path)` helper that produces `/api/v1/…` prefixed routes.

### Key patterns

**Transactional writes with events**: the `DBTX` interface in `inventory/vehicle/repository.go` accepts both `*sql.DB` and `*sql.Tx`, so write methods participate in the caller-managed transaction that also saves the outbox entry. Never bypass this — domain writes and outbox saves must be atomic.

**Adding a new domain event**: define the `EventType` constant and payload struct in `inventory/vehicle/events.go`, call `vehicle.raise(...)` in the service, add a handler in `inventory/vehicle/consumer.go`, and register it in `EventHandlers()`.

**Adding a new domain entity**: create a new package under `inventory/`, following the same handler → service → repo → events → consumer structure as `inventory/vehicle/`.

### Infrastructure

- PostgreSQL via `pgx/v5` (`database/sql` stdlib interface).
- Migrations in `internal/database/migrations/` using golang-migrate numbered files (`000001_*.up.sql` / `000001_*.down.sql`).
- SNS + SQS via AWS SDK v2; Watermill (`github.com/ThreeDotsLabs/watermill-aws`) wraps both.
- LocalStack (`localhost:4566`) provides SNS/SQS locally. Set `AWS_ENDPOINT_URL=http://localhost:4566` in `.env`.
- DB integration tests use testcontainers to spin up a real Postgres instance.
