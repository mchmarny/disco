package main

import (
	"github.com/mchmarny/vctl/cmd/vctl/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	name = "vctl"
)

var (
	version = "v0.0.1-default"
)

func main() {
	initLogging(name, version)
	fatalErr(cli.Run(name, version))
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error")
	}
}

func initLogging(name, version string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	zerolog.ErrorFieldName = "err"
	zerolog.CallerFieldName = "caller"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
