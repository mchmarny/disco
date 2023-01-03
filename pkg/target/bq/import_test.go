package bq

import (
	"context"
	"errors"
)

type testInserter struct{}

func (i *testInserter) Put(ctx context.Context, items interface{}) error {
	if items == nil {
		return errors.New("items must be specified")
	}
	return nil
}
