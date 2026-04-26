package vehicle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// DBTX is satisfied by both *sql.DB and *sql.Tx, allowing write methods to
// participate in caller-managed transactions without exposing *sql.Tx directly.
type DBTX interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type VehicleRepo interface {
	Create(ctx context.Context, tx DBTX, vehicle *Vehicle) error
	GetAll(ctx context.Context) ([]Vehicle, error)
	FindByID(ctx context.Context, id int) (*Vehicle, error)
	UpdateByID(ctx context.Context, tx DBTX, id int, vehicle *Vehicle) error
	DeleteByID(ctx context.Context, tx DBTX, id int) error
}

type PgSQLVehicleRepo struct {
	db *sql.DB
}

func NewPgSQLVehicleRepo(db *sql.DB) *PgSQLVehicleRepo {
	return &PgSQLVehicleRepo{db: db}
}

func (repo *PgSQLVehicleRepo) GetAll(ctx context.Context) ([]Vehicle, error) {
	const query = `SELECT id, name, description, created_at, updated_at FROM vehicles ORDER BY id`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(&v.ID, &v.Name, &v.Description, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, rows.Err()
}

func (repo *PgSQLVehicleRepo) Create(ctx context.Context, tx DBTX, vehicle *Vehicle) error {
	const query = `INSERT INTO vehicles (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	return tx.QueryRowContext(ctx, query, vehicle.Name, vehicle.Description).
		Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
}

func (repo *PgSQLVehicleRepo) FindByID(ctx context.Context, id int) (*Vehicle, error) {
	const query = `SELECT id, name, description, created_at, updated_at FROM vehicles WHERE id = $1`
	var vehicle Vehicle

	err := repo.db.QueryRowContext(ctx, query, id).
		Scan(&vehicle.ID, &vehicle.Name, &vehicle.Description, &vehicle.CreatedAt, &vehicle.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find vehicle: %w", err)
	}

	return &vehicle, nil
}

func (repo *PgSQLVehicleRepo) UpdateByID(ctx context.Context, tx DBTX, id int, vehicle *Vehicle) error {
	const query = `UPDATE vehicles SET name=$1, description=$2, updated_at=NOW() WHERE id=$3 RETURNING updated_at`

	err := tx.QueryRowContext(ctx, query, vehicle.Name, vehicle.Description, id).Scan(&vehicle.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}

	vehicle.ID = id
	return nil
}

func (repo *PgSQLVehicleRepo) DeleteByID(ctx context.Context, tx DBTX, id int) error {
	const query = `DELETE FROM vehicles WHERE id = $1`

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}
