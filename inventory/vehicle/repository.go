package vehicle

import (
	"database/sql"
	"errors"
	"fmt"
)

type VehicleRepo interface {
	Create(vehicle *Vehicle) error
	GetAll() ([]Vehicle, error)
	FindByID(id int) (*Vehicle, error)
	UpdateByID(id int, vehicle *Vehicle) error
	DeleteByID(id int) error
}

type PgSQLVehicleRepo struct {
	db *sql.DB
}

func NewPgSQLVehicleRepo(db *sql.DB) *PgSQLVehicleRepo {
	return &PgSQLVehicleRepo{db: db}
}

func (r *PgSQLVehicleRepo) GetAll() ([]Vehicle, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM vehicles ORDER BY id")
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

func (r *PgSQLVehicleRepo) Create(vehicle *Vehicle) error {
	const q = `INSERT INTO vehicles (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(q, vehicle.Name, vehicle.Description).
		Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
}

func (r *PgSQLVehicleRepo) FindByID(id int) (*Vehicle, error) {
	var v Vehicle
	err := r.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM vehicles WHERE id = $1", id,
	).Scan(&v.ID, &v.Name, &v.Description, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find vehicle: %w", err)
	}
	return &v, nil
}

func (r *PgSQLVehicleRepo) UpdateByID(id int, vehicle *Vehicle) error {
	const q = `UPDATE vehicles SET name=$1, description=$2, updated_at=NOW() WHERE id=$3 RETURNING updated_at`
	err := r.db.QueryRow(q, vehicle.Name, vehicle.Description, id).Scan(&vehicle.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}
	vehicle.ID = id
	return nil
}

func (r *PgSQLVehicleRepo) DeleteByID(id int) error {
	res, err := r.db.Exec("DELETE FROM vehicles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
