package vehicle

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
)

// EventHandlers returns the handler map for all vehicle domain events.
func EventHandlers() map[events.EventType]events.HandlerFunc {
	return map[events.EventType]events.HandlerFunc{
		VehicleCreated: handleVehicleCreated,
		VehicleUpdated: handleVehicleUpdated,
		VehicleDeleted: handleVehicleDeleted,
	}
}

func handleVehicleCreated(_ context.Context, event events.DomainEvent) error {
	var payload VehicleCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode vehicle.created payload: %w", err)
	}
	log.Printf("[inventory.vehicle.created] id=%d name=%q", payload.ID, payload.Name)
	return nil
}

func handleVehicleUpdated(_ context.Context, event events.DomainEvent) error {
	var payload VehicleUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode vehicle.updated payload: %w", err)
	}
	log.Printf("[inventory.vehicle.updated] id=%d name=%q", payload.ID, payload.Name)
	return nil
}

func handleVehicleDeleted(_ context.Context, event events.DomainEvent) error {
	var payload VehicleDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode vehicle.deleted payload: %w", err)
	}
	log.Printf("[inventory.vehicle.deleted] id=%d", payload.ID)
	return nil
}
