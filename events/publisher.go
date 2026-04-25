package events

import "context"

// Publisher is the outbound port for domain event publishing.
// Implementations must be safe for concurrent use.
type Publisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}
