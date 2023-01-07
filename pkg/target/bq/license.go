package bq

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ImportLicenses(ctx context.Context, req *types.ImportRequest, in ...*types.LicenseReport) error {
	if err := configureTarget(ctx, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	rows := MakeLicenseRows(in...)
	if err := insert(ctx, req, rows); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	log.Info().Msgf("inserted %d records into %s.%s.%s", len(rows), req.ProjectID, req.DatasetID, req.TableID)

	return nil
}

func MakeLicenseRows(in ...*types.LicenseReport) []*LicenseRow {
	list := make([]*LicenseRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)
	batchID := time.Now().UTC().Unix()

	for _, r := range in {
		for _, l := range r.Licenses {
			log.Debug().Msgf("adding license %s from %s", l.Name, l.Source)
			list = append(list, &LicenseRow{
				BatchID: batchID,
				Image:   types.ParseImageNameFromDigest(r.Image),
				Sha:     types.ParseImageShaFromDigest(r.Image),
				Name:    l.Name,
				Package: l.Source,
				Updated: updated,
			})
		}
	}

	return list
}

type LicenseRow struct {
	BatchID int64
	Image   string
	Sha     string
	Name    string
	Package string
	Updated string
}

func (i *LicenseRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"batch_id": i.BatchID,
		"image":    i.Image,
		"sha":      i.Sha,
		"name":     i.Name,
		"package":  i.Package,
		"updated":  i.Updated,
	}, "", nil
}
