package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
)

func insert(ctx context.Context, req *types.ImportRequest, items interface{}) error {
	if req == nil || req.ProjectID == "" || req.DatasetID == "" || req.TableID == "" {
		return errors.New("project, dataset and table must be specified")
	}

	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to create bigquery client")
	}
	defer client.Close()

	inserter := client.Dataset(req.DatasetID).Table(req.TableID).Inserter()
	if err := inserter.Put(ctx, items); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}
