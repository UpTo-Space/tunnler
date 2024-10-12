package database

import "os"

var (
	postgresUsername string = getEnv("POSTGRES_USER", "postgres")
	postgresPassword string = getEnv("POSTGRES_PASSWORD", "postgres")
	postgresServer   string = getEnv("POSTGRES_SERVER", "localhost")
	postgresPort     string = getEnv("POSTGRES_PORT", "5432")
	postgresDatabase string = getEnv("POSTGRES_DATABASE", "postgres")
	postgresSslMode  string = getEnv("POSTGRES_SSLMODE", "disable")
)

func getEnv(key, fallback string) string {
	v := os.Getenv(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}
