package outbox

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Save persists a domain event to the outbox within the provided transaction.
// The event will be published to SNS by the Relay once the transaction commits.
func (s *Store) Save(ctx context.Context, tx *sql.Tx, event events.DomainEvent) error {
	const query = `
		INSERT INTO domain_events (id, event_type, aggregate_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.ExecContext(ctx, query,
		event.ID,
		string(event.Type),
		event.AggregateID,
		[]byte(event.Payload),
		event.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("outbox: save event: %w", err)
	}
	return nil
}

// Pending returns up to 100 unpublished events ordered by creation time.
func (s *Store) Pending(ctx context.Context) ([]events.DomainEvent, error) {
	const query = `
		SELECT id, event_type, aggregate_id, payload, occurred_at
		FROM domain_events
		WHERE published_at IS NULL
		ORDER BY created_at
		LIMIT 100`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("outbox: query pending: %w", err)
	}
	defer rows.Close()

	var evts []events.DomainEvent
	for rows.Next() {
		var e events.DomainEvent
		var eventType string
		if err := rows.Scan(&e.ID, &eventType, &e.AggregateID, &e.Payload, &e.OccurredAt); err != nil {
			return nil, fmt.Errorf("outbox: scan event: %w", err)
		}
		e.Type = events.EventType(eventType)
		e.Version = events.SchemaVersion
		evts = append(evts, e)
	}

	return evts, rows.Err()
}

// MarkPublished stamps a domain event as delivered to SNS.
func (s *Store) MarkPublished(ctx context.Context, eventID string) error {
	const query = `UPDATE domain_events SET published_at = NOW() WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, eventID)
	return err
}
