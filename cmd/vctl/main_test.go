package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	initLogging(zerolog.DebugLevel)
	log.Info().Msg("starting tests")
	code := m.Run()
	os.Exit(code)
}
