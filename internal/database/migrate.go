package database

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func Migrate() (*migrate.Migrate, error) {
	db := New()

	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{
		SchemaName: os.Getenv("BLUEPRINT_DB_SCHEMA"),
	})
	if err != nil {
		return nil, fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		os.Getenv("BLUEPRINT_DB_DATABASE"),
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("create migrate instance: %w", err)
	}

	return m, nil
}
