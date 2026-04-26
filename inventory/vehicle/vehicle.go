package vehicle

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
)

type Vehicle struct {
	ID           int
	Name         string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	domainEvents []events.DomainEvent
}

type CreateVehicleCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateVehicleCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var ErrNotFound = errors.New("vehicle not found")

func (v *Vehicle) DomainEvents() []events.DomainEvent {
	return v.domainEvents
}

func (v *Vehicle) ClearEvents() {
	v.domainEvents = nil
}

func (v *Vehicle) raise(eventType events.EventType, payload any) {
	raw, _ := json.Marshal(payload)
	v.domainEvents = append(v.domainEvents, events.DomainEvent{
		ID:          uuid.NewString(),
		Type:        eventType,
		AggregateID: strconv.Itoa(v.ID),
		Payload:     json.RawMessage(raw),
		OccurredAt:  time.Now().UTC(),
		Version:     events.SchemaVersion,
	})
}
