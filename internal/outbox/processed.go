package outbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ProcessedStore tracks consumed event IDs to guard against duplicate delivery.
// SQS provides at-least-once semantics, so any handler with side effects must
// check this store before executing.
type ProcessedStore struct {
	db *sql.DB
}

func NewProcessedStore(db *sql.DB) *ProcessedStore {
	return &ProcessedStore{db: db}
}

func (s *ProcessedStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	const query = `SELECT 1 FROM processed_events WHERE event_id = $1`
	var exists int
	err := s.db.QueryRowContext(ctx, query, eventID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check processed event: %w", err)
	}
	return true, nil
}

func (s *ProcessedStore) MarkProcessed(ctx context.Context, eventID string) error {
	const query = `INSERT INTO processed_events (event_id) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := s.db.ExecContext(ctx, query, eventID)
	return err
}
