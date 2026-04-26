package vehicle

import "github.com/blazeisclone/vehicle-dms-inventory/internal/events"

const (
	VehicleCreated events.EventType = "inventory.vehicle.created.v1"
	VehicleUpdated events.EventType = "inventory.vehicle.updated.v1"
	VehicleDeleted events.EventType = "inventory.vehicle.deleted.v1"
)

type VehicleCreatedPayload struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VehicleUpdatedPayload struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VehicleDeletedPayload struct {
	ID int `json:"id"`
}
