package main

import (
	"os"

	disco "github.com/mchmarny/disco/cmd/disco/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logLevelEnvVar = "debug"
)

var (
	version = "v0.0.1-default"
	commit  = "none"
	date    = "unknown"
)

func main() {
	initLogging()
	err := disco.Execute(version, commit, date, os.Args)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func initLogging() {
	level := zerolog.InfoLevel
	levStr := os.Getenv(logLevelEnvVar)
	if levStr == "true" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	out := zerolog.ConsoleWriter{
		Out: os.Stderr,
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}

	log.Logger = zerolog.New(out)
}
