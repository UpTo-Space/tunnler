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

var (
	AddUserQuery string = `
		INSERT INTO users (username, password_hash, email, activation_code) 
		VALUES ('$1', crypt('$2', gen_salt('bf')), '$3', CAST(1000000000 + floor(random() * 9000000000) AS bigint));`

	ChangePasswordQuery string = `
		UPDATE users 
		SET password_hash = crypt('$1', gen_salt('bf')) 
		WHERE id = $2;`

	AttemptLoginQuery string = `
		SELECT (password_hash = crypt('$1', password_hash) AND activated)
    	AS password_match 
		FROM users 
		WHERE username = '$2';`

	GetUserNameQuery string = `
		SELECT username
		FROM users
		WHERE id = $1;`

	SetActivatedQuery string = `
		UPDATE users
		set activated = TRUE
		WHERE username = '$1'`

	CheckActivationCodeQuery string = `
		SELECT (activation_code = $1)
		AS code_match
		FROM users
		WHERE username = '$1'`
)

func initializeDatabase() {
	log.Printf("Setting up Postgres Connection...\n")
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
	log.Printf("Postgre Connection established...\n")
}

func registerUser(username, password, email string) error {
	_, err := database.Exec(AddUserQuery, username, password, email)

	return err
}

func changePassword(id, newPassword string) error {
	_, err := database.Exec(ChangePasswordQuery, newPassword, id)

	return err
}

func attemptLogin(username, password string) (bool, error) {
	var result string
	err := database.QueryRow(AttemptLoginQuery, password, username).Scan(&result)

	if err != nil {
		return false, err
	}

	return result == "t", nil
}

func getUserName(id string) (string, error) {
	var result string
	err := database.QueryRow(GetUserNameQuery).Scan(&result)

	if err != nil {
		return "", err
	}

	return result, nil
}

func activateUser(username string) error {
	_, err := database.Exec(SetActivatedQuery, username)

	return err
}

func checkActivationCode(activationCode int, username string) (bool, error) {
	var result string
	err := database.QueryRow(CheckActivationCodeQuery, activationCode, username).Scan(&result)

	if err != nil {
		return false, err
	}

	return result == "t", nil
}
