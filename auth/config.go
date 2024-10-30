package main

import (
	"github.com/UpTo-Space/tunnler/common"
)

var (
	listenHostName string = common.GetEnv("LISTENHOSTNAME", "127.0.0.1")
	hostName       string = common.GetEnv("HOSTNAME", "127.0.0.1")
	hostPort       string = common.GetEnv("PORT", "8887")
	smtpHost       string = common.GetEnv("SMTP_HOST", "127.0.0.1")
	smtpPort       string = common.GetEnv("SMTP_PORT", "1025")
	smtpIdentity   string = common.GetEnv("SMTP_IDENTITY", "")
	smtpUser       string = common.GetEnv("SMTP_USER", "")
	smtpPassword   string = common.GetEnv("SMTP_PASSWORD", "")
	fromAddress    string = common.GetEnv("FROM_ADDRESS", "tunnler@up-to.space")
	symmetricKey   string = common.GetEnv("SYMMETRICS_KEY", "abcdefghijkl1234567890mnopqrstuv")
	devEnv         bool   = common.GetEnvBool("DEV", true)

	EMail struct {
		ActivationSubject string
	} = struct{ ActivationSubject string }{
		ActivationSubject: common.GetEnv("EMAIL_ACTIVATION_SUBJECT", "Tunnler Activation"),
	}
)

const (
	authHeaderKey        = "Authorization"
	authHeaderBearerType = "bearer"
)
