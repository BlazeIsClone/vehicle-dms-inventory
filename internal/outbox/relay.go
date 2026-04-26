package outbox

import (
	"context"
	"log"
	"time"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
)

// Relay polls the outbox table and publishes pending events to SNS.
// It provides at-least-once delivery — events are published then marked delivered.
type Relay struct {
	store    *Store
	pub      events.Publisher
	interval time.Duration
}

func NewRelay(store *Store, pub events.Publisher) *Relay {
	return &Relay{store: store, pub: pub, interval: 5 * time.Second}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.publishPending(ctx)
		}
	}
}

func (r *Relay) publishPending(ctx context.Context) {
	evts, err := r.store.Pending(ctx)
	if err != nil {
		log.Printf("outbox relay: fetch pending: %v", err)
		return
	}

	for _, e := range evts {
		if err := r.pub.Publish(ctx, e); err != nil {
			log.Printf("outbox relay: publish event %s (%s): %v", e.ID, e.Type, err)
			continue
		}
		if err := r.store.MarkPublished(ctx, e.ID); err != nil {
			log.Printf("outbox relay: mark published %s: %v", e.ID, err)
		}
	}
}
