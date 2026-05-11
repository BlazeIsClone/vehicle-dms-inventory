package database

import (
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate() (*migrate.Migrate, error) {
	db := New()

	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{
		SchemaName: os.Getenv("BLUEPRINT_DB_SCHEMA"),
	})
	if err != nil {
		return nil, fmt.Errorf("create postgres driver: %w", err)
	}

	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, os.Getenv("BLUEPRINT_DB_DATABASE"), driver)
	if err != nil {
		return nil, fmt.Errorf("create migrate instance: %w", err)
	}

	return m, nil
}
