package database

import (
	"github.com/UpTo-Space/tunnler/common"
)

var (
	postgresUsername string = common.GetEnv("POSTGRES_USER", "postgres")
	postgresPassword string = common.GetEnv("POSTGRES_PASSWORD", "postgres")
	postgresServer   string = common.GetEnv("POSTGRES_SERVER", "localhost")
	postgresPort     string = common.GetEnv("POSTGRES_PORT", "5432")
	postgresDatabase string = common.GetEnv("POSTGRES_DATABASE", "postgres")
	postgresSslMode  string = common.GetEnv("POSTGRES_SSLMODE", "disable")
)
