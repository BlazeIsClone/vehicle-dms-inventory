package vehicle

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/blazeisclone/vehicle-dms-inventory/events"
)

type Service struct {
	repo VehicleRepo
	pub  events.Publisher
}

func NewVehicleSvc(repo VehicleRepo, pub events.Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

func (svc *Service) Create(cmd CreateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	if err := svc.repo.Create(vehicle); err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}

	svc.publishEvent(events.VehicleCreated, vehicle.ID, vehicle)
	return vehicle, nil
}

func (svc *Service) GetAll() ([]Vehicle, error) {
	return svc.repo.GetAll()
}

func (svc *Service) FindByID(id int) (*Vehicle, error) {
	return svc.repo.FindByID(id)
}

func (svc *Service) Update(id int, cmd UpdateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	if err := svc.repo.UpdateByID(id, vehicle); err != nil {
		return nil, fmt.Errorf("update vehicle: %w", err)
	}

	svc.publishEvent(events.VehicleUpdated, id, vehicle)
	return vehicle, nil
}

func (svc *Service) Delete(id int) error {
	if err := svc.repo.DeleteByID(id); err != nil {
		return err
	}

	svc.publishEvent(events.VehicleDeleted, id, map[string]int{"id": id})
	return nil
}

// publishEvent is best-effort: publish failures are logged but do not fail the
// primary operation. The DB write is the source of truth. context.Background()
// is used so HTTP request cancellation does not abort an in-flight publish that
// follows a successful DB write.
func (svc *Service) publishEvent(eventType events.EventType, aggregateID int, payload any) {
	event := events.DomainEvent{
		ID:          uuid.NewString(),
		Type:        eventType,
		AggregateID: strconv.Itoa(aggregateID),
		Payload:     payload,
		OccurredAt:  time.Now().UTC(),
		Version:     events.SchemaVersion,
	}

	if err := svc.pub.Publish(context.Background(), event); err != nil {
		log.Printf("vehicle service: publish %s event: %v", eventType, err)
	}
}
