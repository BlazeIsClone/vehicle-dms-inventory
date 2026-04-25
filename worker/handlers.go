package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/blazeisclone/vehicle-dms-inventory/events"
	"github.com/blazeisclone/vehicle-dms-inventory/inventory/vehicle"
)

// VehicleEventHandlers returns the handler map for all vehicle domain events.
func VehicleEventHandlers() map[events.EventType]HandlerFunc {
	return map[events.EventType]HandlerFunc{
		events.VehicleCreated: handleVehicleCreated,
		events.VehicleUpdated: handleVehicleUpdated,
		events.VehicleDeleted: handleVehicleDeleted,
	}
}

func handleVehicleCreated(ctx context.Context, event events.DomainEvent) error {
	v, err := decodeVehiclePayload(event)
	if err != nil {
		return err
	}
	log.Printf("worker: [vehicle.created] id=%d name=%q", v.ID, v.Name)
	return nil
}

func handleVehicleUpdated(ctx context.Context, event events.DomainEvent) error {
	v, err := decodeVehiclePayload(event)
	if err != nil {
		return err
	}
	log.Printf("worker: [vehicle.updated] id=%d name=%q", v.ID, v.Name)
	return nil
}

func handleVehicleDeleted(_ context.Context, event events.DomainEvent) error {
	log.Printf("worker: [vehicle.deleted] aggregate_id=%s", event.AggregateID)
	return nil
}

// decodeVehiclePayload re-marshals the generic payload back to JSON then into
// the concrete Vehicle type. This two-step is necessary because encoding/json
// decodes unknown object fields as map[string]any.
func decodeVehiclePayload(event events.DomainEvent) (*vehicle.Vehicle, error) {
	raw, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, fmt.Errorf("re-marshal payload: %w", err)
	}
	var v vehicle.Vehicle
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("decode vehicle payload: %w", err)
	}
	return &v, nil
}
