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

var ErrNotFound = errors.New("vehicle not found")
