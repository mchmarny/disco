package main

import (
	"os"

	"github.com/mchmarny/disco/cmd/disco/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logLevelEnvVar = "debug"
)

var (
	version = "v0.0.1-default"
	commit  = "none"
)

func main() {
	logLevel := zerolog.InfoLevel
	levStr := os.Getenv(logLevelEnvVar)
	if levStr == "true" {
		logLevel = zerolog.DebugLevel
	}
	initLogging(logLevel)
	err := cli.Execute(version, commit)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func initLogging(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	})
}
