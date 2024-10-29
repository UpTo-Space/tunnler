package main

import (
	"os"
	"strconv"
)

var (
	listenHostName string = getEnv("LISTENHOSTNAME", "127.0.0.1")
	hostName       string = getEnv("HOSTNAME", "127.0.0.1")
	hostPort       string = getEnv("PORT", "8887")
	smtpHost       string = getEnv("SMTP_HOST", "127.0.0.1")
	smtpPort       string = getEnv("SMTP_PORT", "1025")
	smtpIdentity   string = getEnv("SMTP_IDENTITY", "")
	smtpUser       string = getEnv("SMTP_USER", "")
	smtpPassword   string = getEnv("SMTP_PASSWORD", "")
	fromAddress    string = getEnv("FROM_ADDRESS", "tunnler@up-to.space")
	symmetricKey   string = getEnv("SYMMETRICS_KEY", "abcdefghijkl1234567890mnopqrstuv")
	devEnv         bool   = getEnvBool("DEV", true)

	EMail struct {
		ActivationSubject string
	} = struct{ ActivationSubject string }{
		ActivationSubject: getEnv("EMAIL_ACTIVATION_SUBJECT", "Tunnler Activation"),
	}
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

func getEnvBool(key string, fallback bool) bool {
	v := getEnv(key, strconv.FormatBool(fallback))

	if len(v) == 0 {
		return fallback
	}

	b, err := strconv.ParseBool(v)

	if err != nil {
		panic(err)
	}

	return b
}
