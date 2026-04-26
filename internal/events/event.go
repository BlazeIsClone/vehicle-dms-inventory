package events

import (
	"encoding/json"
	"time"
)

type EventType string

const SchemaVersion = "1"

// DomainEvent is the canonical event envelope published to SNS and consumed from SQS.
// Fields must remain stable — add optional fields; never remove or rename existing ones.
type DomainEvent struct {
	ID          string          `json:"id"`
	Type        EventType       `json:"type"`
	AggregateID string          `json:"aggregate_id"`
	Payload     json.RawMessage `json:"payload"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Version     string          `json:"version"`
}
