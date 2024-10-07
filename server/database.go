package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	connectionString string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", postgresUsername, postgresPassword, postgresServer, postgresPort, postgresDatabase, postgresSslMode)
	database         *sql.DB
)

func initializeDatabase() {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database = db

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start db driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
}

func addListener(clientIp, localPort, remotePort, subDomain string) error {
	_, err := database.Exec(`
	INSERT INTO connections (client_ip, local_port, remote_port, subdomain) 
	VALUES ($1, $2, $3, $4)	
	`, clientIp, localPort, remotePort, subDomain)

	return err
}
