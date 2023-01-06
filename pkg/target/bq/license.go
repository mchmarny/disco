package bq

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ImportLicenses(ctx context.Context, req *types.ImportRequest) error {
	if req == nil || req.FilePath == "" {
		return errors.Errorf("configured import request is required: %v", req)
	}

	f := func(v interface{}) bool {
		return false
	}

	var report *types.LicenseReport
	var err error

	switch req.LicenseFormat {
	case types.LicenseReportFormatTrivy:
		report, err = scanner.ParseLicense(req.FilePath, f)
		if err != nil {
			return errors.Wrapf(err, "failed to parse report from %s", req.FilePath)
		}
	case types.LicenseReportFormatDisco:
		if err := types.UnmarshalFromFile(req.FilePath, &report); err != nil {
			return errors.Wrapf(err, "failed to parse report from %s", req.FilePath)
		}
	default:
		return errors.Errorf("unsupported vulnerability report format: %s", req.VulnFormat)
	}

	if err := configureTarget(ctx, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	rows := MakeLicenseRows(report)
	if err := insert(ctx, req, rows); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	log.Info().Msgf("inserted %d rows into %s.%s.%s", len(rows), req.ProjectID, req.DatasetID, req.TableID)

	return nil
}

func MakeLicenseRows(in *types.LicenseReport) []*LicenseRow {
	list := make([]*LicenseRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)
	batchID := time.Now().UTC().Unix()

	for _, l := range in.Licenses {
		log.Info().Msgf("adding license %s from %s", l.Name, l.Source)
		list = append(list, &LicenseRow{
			BatchID: batchID,
			Image:   types.ParseImageNameFromDigest(in.Image),
			Sha:     types.ParseImageShaFromDigest(in.Image),
			Name:    l.Name,
			Package: l.Source,
			Updated: updated,
		})
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
