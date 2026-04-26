package vehicle

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/outbox"
)

type Service struct {
	db     *sql.DB
	repo   VehicleRepo
	outbox *outbox.Store
}

func NewVehicleSvc(db *sql.DB, repo VehicleRepo, store *outbox.Store) *Service {
	return &Service{db: db, repo: repo, outbox: store}
}

func (svc *Service) Create(ctx context.Context, cmd CreateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create vehicle: begin tx: %w", err)
	}

	if err := svc.repo.Create(ctx, tx, vehicle); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create vehicle: %w", err)
	}

	vehicle.raise(VehicleCreated, VehicleCreatedPayload{
		ID:          vehicle.ID,
		Name:        vehicle.Name,
		Description: vehicle.Description,
	})

	if err := svc.flushEvents(ctx, tx, vehicle); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create vehicle: commit: %w", err)
	}

	return vehicle, nil
}

func (svc *Service) GetAll(ctx context.Context) ([]Vehicle, error) {
	return svc.repo.GetAll(ctx)
}

func (svc *Service) FindByID(ctx context.Context, id int) (*Vehicle, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *Service) Update(ctx context.Context, id int, cmd UpdateVehicleCommand) (*Vehicle, error) {
	vehicle := &Vehicle{Name: cmd.Name, Description: cmd.Description}

	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("update vehicle: begin tx: %w", err)
	}

	if err := svc.repo.UpdateByID(ctx, tx, id, vehicle); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update vehicle: %w", err)
	}

	vehicle.raise(VehicleUpdated, VehicleUpdatedPayload{
		ID:          vehicle.ID,
		Name:        vehicle.Name,
		Description: vehicle.Description,
	})

	if err := svc.flushEvents(ctx, tx, vehicle); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("update vehicle: commit: %w", err)
	}

	return vehicle, nil
}

func (svc *Service) Delete(ctx context.Context, id int) error {
	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("delete vehicle: begin tx: %w", err)
	}

	if err := svc.repo.DeleteByID(ctx, tx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete vehicle: %w", err)
	}

	vehicle := &Vehicle{ID: id}
	vehicle.raise(VehicleDeleted, VehicleDeletedPayload{ID: id})

	if err := svc.flushEvents(ctx, tx, vehicle); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("delete vehicle: commit: %w", err)
	}

	return nil
}

// flushEvents saves all pending domain events from the aggregate to the outbox
// within the active transaction, then clears them from the aggregate.
func (svc *Service) flushEvents(ctx context.Context, tx *sql.Tx, vehicle *Vehicle) error {
	for _, e := range vehicle.DomainEvents() {
		if err := svc.outbox.Save(ctx, tx, e); err != nil {
			return fmt.Errorf("save domain event: %w", err)
		}
	}
	vehicle.ClearEvents()
	return nil
}
