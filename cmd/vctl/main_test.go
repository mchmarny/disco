package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	initLogging(name, version)
	log.Info().Msg("starting tests")
	code := m.Run()
	os.Exit(code)
}
