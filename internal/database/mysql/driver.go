package mysql

import (
	"database/sql"
	"fmt"
	"os"

	"context"

	driver "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func getDBConfig() driver.Config {
	return driver.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		AllowNativePasswords: true,
	}
}

func Init() (*sql.DB, error) {
	dbConfig := getDBConfig()

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open %w", err)
	}

	if err := Ping(db); err != nil {
		return nil, fmt.Errorf("sql.error %w", err)
	}

	return db, nil
}

func Migrate() (*migrate.Migrate, error) {
	dbConfig := getDBConfig()

	return migrate.New("file://internal/database/mysql/migrations", "mysql://"+dbConfig.FormatDSN())

}

func Ping(db *sql.DB) error {
	return db.PingContext(context.Background())
}
