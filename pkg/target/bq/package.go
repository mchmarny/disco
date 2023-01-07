package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ImportPackages(ctx context.Context, req *types.ImportRequest, in ...*types.PackageReport) error {
	if err := configureTarget(ctx, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	rows := MakePackageRows(in...)
	if err := insert(ctx, req, rows); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	log.Info().Msgf("inserted %d rows into %s.%s.%s", len(rows), req.ProjectID, req.DatasetID, req.TableID)

	return nil
}

type PackageRow struct {
	BatchID        int64
	Image          string
	Sha            string
	Format         string
	Provider       string
	Package        string
	PackageVersion string
	Source         string
	License        string
	Updated        string
}

func (i *PackageRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"batch_id": i.BatchID,
		"image":    i.Image,
		"sha":      i.Sha,
		"format":   i.Format,
		"provider": i.Provider,
		"package":  i.Package,
		"version":  i.PackageVersion,
		"source":   i.Source,
		"license":  i.License,
		"updated":  i.Updated,
	}, "", nil
}
