package common

import (
	"os"
	"strconv"
)

const (
	AuthHeaderKey        = "Authorization"
	AuthHeaderBearerType = "bearer"
)

func GetEnv(key, fallback string) string {
	v := os.Getenv(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}

func GetEnvBool(key string, fallback bool) bool {
	v := GetEnv(key, strconv.FormatBool(fallback))

	if len(v) == 0 {
		return fallback
	}

	b, err := strconv.ParseBool(v)

	if err != nil {
		panic(err)
	}

	return b
}
