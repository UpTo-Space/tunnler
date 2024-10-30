package main

import (
	"github.com/UpTo-Space/tunnler/common"
)

var (
	hostName string = common.GetEnv("HOSTNAME", "127.0.0.1")
	hostPort string = common.GetEnv("PORT", "8888")
)
