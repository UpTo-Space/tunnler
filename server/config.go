package main

import "os"

var (
	hostName string = getEnv("HOSTNAME", "127.0.0.1")
	hostPort string = getEnv("PORT", "8888")
)

func getEnv(key, fallback string) string {
	v := os.Getenv(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}
