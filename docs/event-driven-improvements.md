# Event-Driven Architecture Improvements

## Overview

This document outlines planned improvements to the event-driven architecture and project structure of the vehicle DMS inventory service. Changes are grouped by category and ordered by priority.

---

## 1. Directory & File Structure Refactor

The project follows **package-by-component** architecture. The current event-related files are misaligned — vehicle-specific event types and handlers live outside the vehicle component.

### Current Structure (event-related)

```
events/
├── event.go          # DomainEvent envelope
├── publisher.go      # Publisher interface
└── vehicle.go        # Vehicle event types  ← misplaced

worker/
├── consumer.go       # Generic router/dispatcher
└── handlers.go       # Vehicle event handlers  ← misplaced
```

### Target Structure

```
events/                              # Shared event contracts only
├── event.go                         # DomainEvent envelope
└── publisher.go                     # Publisher interface + HandlerFunc type

inventory/vehicle/                   # Self-contained vehicle component
├── vehicle.go                       # Domain model + aggregate
├── service.go                       # Use cases
├── handler.go                       # HTTP handlers
├── repository.go                    # Data access
├── request.go                       # Request validation
├── events.go              ← MOVE    # Vehicle event types (from events/vehicle.go)
└── consumer.go            ← MOVE    # Vehicle event handlers (from worker/handlers.go)

infra/sns/
└── publisher.go                     # SNS adapter (unchanged)

internal/outbox/           ← NEW
├── store.go                         # Outbox table DB operations
└── relay.go                         # Background relay process

worker/
└── consumer.go                      # Generic Watermill bootstrap only (slimmed)
```

### Key Principle

Each component exposes its own event handlers via an `EventHandlers()` function. The worker becomes a generic bootstrap that registers handlers from each component:

```go
// worker/consumer.go
for eventType, handler := range vehicle.EventHandlers() {
    consumer.Register(eventType, handler)
}
```

---

## 2. Typed Event Payloads

### Problem

`DomainEvent.Payload` is typed as `any`, which causes `worker/handlers.go` to perform a fragile double-encode (`any → JSON → struct`) to recover typed data. This will break silently if field names change.

### Change

Define concrete payload structs per event type in `inventory/vehicle/events.go`:

```go
type VehicleCreatedPayload struct {
    ID          int
    Name        string
    Description string
}

type VehicleUpdatedPayload struct {
    ID          int
    Name        string
    Description string
}

type VehicleDeletedPayload struct {
    ID int
}
```

Remove `decodeVehiclePayload()` from the consumer entirely.

---

## 3. Transactional Outbox Pattern

### Problem

`service.go:publishEvent()` publishes to SNS after the DB commit but outside the transaction. If SNS is unavailable, the error is suppressed and the event is permanently lost. DB state and the event stream can silently diverge.

### Change

#### 3a. New `domain_events` table (outbox)

```sql
CREATE TABLE domain_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID NOT NULL,
    event_type   VARCHAR(255) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload      JSONB NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON domain_events (published_at) WHERE published_at IS NULL;
```

#### 3b. Service writes event to outbox within the same DB transaction

```go
// service.go
func (s *Service) Create(ctx context.Context, cmd CreateVehicleCommand) (*Vehicle, error) {
    tx, _ := s.db.BeginTx(ctx, nil)
    vehicle, _ := s.repo.CreateTx(ctx, tx, cmd)
    s.outbox.SaveTx(ctx, tx, buildEvent(vehicle))
    tx.Commit()
    return vehicle, nil
}
```

#### 3c. Relay process publishes outbox events to SNS

`internal/outbox/relay.go` runs as a background goroutine (or separate process) that:
1. Polls for rows where `published_at IS NULL`
2. Publishes each event to SNS
3. Marks the row as published (`published_at = NOW()`)

This gives **at-least-once delivery** without coupling the HTTP request to SNS availability.

---

## 4. Context Propagation

### Problem

Service and repository method signatures do not accept `context.Context`, making it impossible to propagate request deadlines, cancellation, or distributed tracing spans.

### Change

Add `ctx context.Context` as the first parameter to all service and repository methods:

```go
// Before
func (s *Service) Create(cmd CreateVehicleCommand) (*Vehicle, error)

// After
func (s *Service) Create(ctx context.Context, cmd CreateVehicleCommand) (*Vehicle, error)
```

---

## 5. Idempotent Event Handlers

### Problem

SQS guarantees at-least-once delivery. The consumer has no deduplication logic, so any handler with side effects will fire multiple times for the same event.

### Change

Track processed event IDs in a `processed_events` table:

```sql
CREATE TABLE processed_events (
    event_id    UUID PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Before executing a handler, check if `event_id` exists. If it does, ack and skip. Insert the ID on successful handling.

---

## 6. Aggregate Emits Its Own Events

### Problem

The application service (`service.go`) constructs and publishes domain events directly. In DDD, a domain event should be raised by the aggregate as a consequence of a state change.

### Change

Have `Vehicle` collect events internally:

```go
// inventory/vehicle/vehicle.go
type Vehicle struct {
    ID          int
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    domainEvents []events.DomainEvent
}

func (v *Vehicle) DomainEvents() []events.DomainEvent { return v.domainEvents }
func (v *Vehicle) ClearEvents()                        { v.domainEvents = nil }
```

The service reads `vehicle.DomainEvents()` after each operation and saves them to the outbox.

---

## 7. Event Type Namespacing

### Problem

Event type strings like `"vehicle.created"` have no domain prefix or schema version, making routing and schema evolution harder as the system grows.

### Change

Adopt the format `<domain>.<aggregate>.<action>.v<version>`:

```go
const (
    VehicleCreated EventType = "inventory.vehicle.created.v1"
    VehicleUpdated EventType = "inventory.vehicle.updated.v1"
    VehicleDeleted EventType = "inventory.vehicle.deleted.v1"
)
```

---

## Implementation Order

| # | Change | Priority | Reason |
|---|--------|----------|--------|
| 1 | Directory & file structure refactor | High | Foundation for all other changes |
| 2 | Typed event payloads | High | Eliminates fragile re-marshalling |
| 3 | Context propagation | High | Required before outbox implementation |
| 4 | Transactional Outbox | High | Prevents silent event loss |
| 5 | Idempotent event handlers | Medium | Safe at-least-once delivery |
| 6 | Aggregate emits events | Medium | Correct DDD ownership |
| 7 | Event type namespacing | Low | Schema evolution readiness |
