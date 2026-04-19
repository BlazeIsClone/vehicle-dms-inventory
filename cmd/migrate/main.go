package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	m, err := database.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		if err := m.Down(); errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to rollback")
			return
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migrations rolled back successfully")

	default:
		log.Fatalf("Unknown command %q. Use 'up' or 'down'", command)
	}
}
