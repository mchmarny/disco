package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
)

type ImportRequest struct {
	ProjectID string
	DatasetID string
	TableID   string
	ObjectURI string
}

func ImportJSON(ctx context.Context, req *ImportRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}
	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to create bigquery client")
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference(req.ObjectURI)
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.AutoDetect = true
	gcsRef.AllowJaggedRows = true
	gcsRef.AllowQuotedNewlines = true
	gcsRef.IgnoreUnknownValues = true

	loader := client.Dataset(req.DatasetID).Table(req.TableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteAppend

	job, err := loader.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "error running the loader")
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "error waiting for the job to complete")
	}

	if status.Err() != nil {
		return errors.Wrap(status.Err(), "error loading data")
	}

	return nil
}
