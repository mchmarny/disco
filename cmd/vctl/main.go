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
	initLogging(zerolog.InfoLevel)
	fatalErr(cli.Run(name, version))
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error")
	}
}

func initLogging(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)
}
