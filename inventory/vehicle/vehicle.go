package vehicle

import (
	"errors"
	"time"
)

type Vehicle struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
