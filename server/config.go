package main

import "os"

var (
	postgresUsername string = getEnv("POSTGRES_USER", "postgres")
	postgresPassword string = getEnv("POSTGRES_PASSWORD", "postgres")
	postgresServer   string = getEnv("POSTGRES_SERVER", "localhost")
	postgresPort     string = getEnv("POSTGRES_PORT", "5432")
	postgresDatabase string = getEnv("POSTGRES_DATABASE", "postgres")
	postgresSslMode  string = getEnv("POSTGRES_SSLMODE", "disable")
	hostName         string = getEnv("HOSTNAME", "127.0.0.1")
	hostPort         string = getEnv("PORT", "8888")
)

func getEnv(key, fallback string) string {
	v := os.Getenv(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}
