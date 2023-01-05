package bq

import (
	"context"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

var (
	licenseSchema = bigquery.Schema{
		{Name: "batch_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "image", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "package", Type: bigquery.StringFieldType},
		{Name: "updated", Type: bigquery.TimestampFieldType, Required: true},
	}

	vulnerabilitySchema = bigquery.Schema{
		{Name: "batch_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "image", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType},
		{Name: "cve", Type: "STRING", Required: true},
		{Name: "severity", Type: "STRING"},
		{Name: "package", Type: bigquery.StringFieldType},
		{Name: "version", Type: "STRING"},
		{Name: "title", Type: "STRING"},
		{Name: "description", Type: "STRING"},
		{Name: "url", Type: "STRING"},
		{Name: "updated", Type: bigquery.TimestampFieldType, Required: true},
	}

	packageSchema = bigquery.Schema{
		{Name: "batch_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "image", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType},
		{Name: "format", Type: bigquery.StringFieldType, Required: true},
		{Name: "provider", Type: bigquery.StringFieldType, Required: true},
		{Name: "originator", Type: bigquery.StringFieldType},
		{Name: "package", Type: bigquery.StringFieldType, Required: true},
		{Name: "version", Type: bigquery.StringFieldType, Required: true},
		{Name: "source", Type: bigquery.StringFieldType},
		{Name: "license", Type: bigquery.StringFieldType},
		{Name: "updated", Type: bigquery.TimestampFieldType, Required: true},
	}
)

func configureTarget(ctx context.Context, req *types.ImportRequest) error {
	if req == nil {
		return errors.New("nil request")
	}

	log.Debug().
		Str("kind", req.TableKind.String()).
		Str("project", req.ProjectID).
		Str("dataset", req.DatasetID).
		Str("table", req.TableID).
		Str("file", req.FilePath).
		Msg("configuring target")

	exists, err := datasetExists(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to check if dataset exists")
	}

	if !exists {
		if err := createDataset(ctx, req); err != nil {
			return errors.Wrap(err, "failed to create dataset")
		}
	}

	exists, err = tableExists(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to check if table exists")
	}

	if !exists {
		var tableSchema bigquery.Schema
		switch req.TableKind {
		case types.TableKindLicense:
			tableSchema = licenseSchema
		case types.TableKindVulnerability:
			tableSchema = vulnerabilitySchema
		case types.TableKindPackage:
			tableSchema = packageSchema
		default:
			return errors.Errorf("unknown table kind: %d", req.TableKind)
		}

		if err := createTable(ctx, req, tableSchema); err != nil {
			return errors.Wrapf(err, "failed to create table %s with ID: %s", req.TableKind, req.TableID)
		}
	}

	return nil
}

func createTable(ctx context.Context, req *types.ImportRequest, schema bigquery.Schema) error {
	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return errors.Wrapf(err, "failed to create bigquery client for project %s", req.ProjectID)
	}
	defer client.Close()

	metaData := &bigquery.TableMetadata{
		Schema: schema,
	}

	tableRef := client.Dataset(req.DatasetID).Table(req.TableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return errors.Wrapf(err, "failed to create table %s", req.TableID)
	}
	return nil
}

func createDataset(ctx context.Context, req *types.ImportRequest) error {
	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return errors.Wrapf(err, "failed to create bigquery client for project %s", req.ProjectID)
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{Location: req.Location}
	if err := client.Dataset(req.DatasetID).Create(ctx, meta); err != nil {
		return errors.Wrapf(err, "failed to create dataset %s", req.DatasetID)
	}
	return nil
}

func datasetExists(ctx context.Context, req *types.ImportRequest) (bool, error) {
	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create bigquery client for project %s", req.ProjectID)
	}
	defer client.Close()

	it := client.Datasets(ctx)
	for {
		dataset, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return false, errors.Wrapf(err, "failed to list datasets for project %s", req.ProjectID)
		}

		if strings.EqualFold(dataset.DatasetID, req.DatasetID) {
			return true, nil
		}
	}
	return false, nil
}

func tableExists(ctx context.Context, req *types.ImportRequest) (bool, error) {
	client, err := bigquery.NewClient(ctx, req.ProjectID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create bigquery client for project %s", req.ProjectID)
	}
	defer client.Close()

	it := client.Dataset(req.DatasetID).Tables(ctx)
	for {
		table, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return false, errors.Wrapf(err, "failed to list datasets for project %s", req.ProjectID)
		}

		if strings.EqualFold(table.TableID, req.TableID) {
			return true, nil
		}
	}
	return false, nil
}
