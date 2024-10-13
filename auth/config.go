package main

import "os"

var (
	hostName     string = getEnv("HOSTNAME", "127.0.0.1")
	hostPort     string = getEnv("PORT", "8887")
	smtpHost     string = getEnv("SMTP_HOST", "mailslurper")
	smtpPort     string = getEnv("SMTP_PORT", "2500")
	fromAddress  string = getEnv("FROM_ADDRESS", "tunnler@up-to.space")
	symmetricKey string = getEnv("SYMMETRICS_KEY", "abcdefghijkl1234567890mnopqrstuv")
)

const (
	authHeaderKey        = "Authorization"
	authHeaderBearerType = "bearer"
)

func getEnv(key, fallback string) string {
	v := os.Getenv(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}
