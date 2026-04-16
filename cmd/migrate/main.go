package main

import (
	"fmt"
	"log"

	mysqlDB "github.com/blazeisclone/vehicle-dms-inventory/internal/database/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Env vars not loaded from file")
	}

	m, err := mysqlDB.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Migrations applied successfully")
}
