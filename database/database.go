package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	connectionString string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", postgresUsername, postgresPassword, postgresServer, postgresPort, postgresDatabase, postgresSslMode)
)

type DatabaseConnection struct {
	Database *sql.DB
	logf     func(f string, v ...interface{})
}

func NewDatabaseConnection() (*DatabaseConnection, error) {
	connection := &DatabaseConnection{
		logf: log.Printf,
	}

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		connection.logf("Failed to connect to database: %v", err)
		return nil, err
	}

	connection.Database = db

	if _, err := os.Stat("./migrations"); os.IsNotExist(err) {
		connection.logf("no migrations found, skipping migrations")
		return connection, nil
	}

	if err := connection.CheckMigrations(); err != nil {
		connection.logf("Error checking migrations: %v", err)
		return nil, err
	}

	return connection, nil
}

func (db *DatabaseConnection) CheckMigrations() error {
	driver, err := postgres.WithInstance(db.Database, &postgres.Config{})
	if err != nil {
		db.logf("Could not start db driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		db.logf("Could not start migrations: %v", err)
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		db.logf("Migration failed: %v", err)
		return err
	}

	return nil
}
