package vehicle

import (
	"database/sql"
	"fmt"
)

type MySQLVehicleRepository struct {
	db *sql.DB
}

func NewMySQLVehicleRepository(db *sql.DB) *MySQLVehicleRepository {
	return &MySQLVehicleRepository{db: db}
}

func (r *MySQLVehicleRepository) GetAll() (*[]Vehicle, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM vehicles"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle

	for rows.Next() {
		var vehicle Vehicle
		var createdAt, updatedAt []uint8

		err := rows.Scan(&vehicle.ID, &vehicle.Name, &vehicle.Description, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}
	if err = rows.Err(); err != nil {
		return &vehicles, err
	}

	return &vehicles, nil
}

func (r *MySQLVehicleRepository) Create(vehicle *Vehicle) error {
	query := "INSERT INTO vehicles (name, description) VALUES (?, ?)"
	result, err := r.db.Exec(query, vehicle.Name, vehicle.Description)
	if err != nil {
		return fmt.Errorf("error creating vehicle: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	vehicle.ID = int(id)
	return nil
}

func (r *MySQLVehicleRepository) FindByID(id int) (*Vehicle, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM vehicles WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var vehicle Vehicle
	var createdAt, updatedAt []uint8
	err := row.Scan(&vehicle.ID, &vehicle.Name, &vehicle.Description, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding vehicle: %w", err)
	}

	return &vehicle, nil
}

func (r *MySQLVehicleRepository) UpdateByID(id int, vehicle *Vehicle) error {
	query := "UPDATE vehicles SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, vehicle.Name, vehicle.Description, id)
	if err != nil {
		return fmt.Errorf("error updating vehicle: %w", err)
	}
	return nil
}

func (r *MySQLVehicleRepository) DeleteByID(id int) error {
	query := "DELETE FROM vehicles WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting vehicle: %w", err)
	}
	return nil
}
