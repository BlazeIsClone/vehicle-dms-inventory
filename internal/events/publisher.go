package events

import "context"

// Publisher is the outbound port for domain event publishing.
// Implementations must be safe for concurrent use.
type Publisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// HandlerFunc processes a single domain event. Return a non-nil error to nack
// the message; it becomes visible again after the queue's visibility timeout.
type HandlerFunc func(ctx context.Context, event DomainEvent) error
