package main

import (
	"fmt"
	"os"

	"github.com/mchmarny/vctl/cmd/vctl/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	name           = "vctl"
	logLevelEnvVar = "debug"
)

var (
	version = "v0.0.1-default"
)

func main() {
	logLevel := zerolog.InfoLevel
	levStr := os.Getenv(logLevelEnvVar)
	if levStr != "" {
		lev, err := zerolog.ParseLevel(levStr)
		if err != nil {
			fmt.Printf("invalid log level: %s, using default: %v", levStr, logLevel)
		}
		logLevel = lev
	}
	initLogging(logLevel, version)
	fatalErr(cli.Run(name, version))
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error")
	}
}

func initLogging(level zerolog.Level, version string) {
	log.Logger = log.With().Caller().Str("ver", version).Logger()
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "lev"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)
}
