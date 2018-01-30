package main

import (
	"github.com/abhishekkr/gol/golenv"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
)

var (
	TIMEFS_AT = golenv.OverrideIfEnv("LOGR_AT", "127.0.0.1:7999")
)

func main() {
	timefsClient.GrpcClient(TIMEFS_AT)
}
