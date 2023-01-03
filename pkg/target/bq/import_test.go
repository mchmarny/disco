package bq

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

type testInserter struct{}

func (i *testInserter) Put(ctx context.Context, items interface{}) error {
	if items == nil {
		return errors.New("items must be specified")
	}
	log.Info().Msgf("Inserting %v", items)
	return nil
}
