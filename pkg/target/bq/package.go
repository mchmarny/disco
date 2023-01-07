package bq

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	spdx "github.com/spdx/tools-golang/spdx/v2_3"
)

func ImportPackages(ctx context.Context, req *types.ImportRequest) error {
	if req == nil || req.FilePath == "" {
		return errors.Errorf("configured import request is required: %v", req)
	}

	if err := configureTarget(ctx, req); err != nil {
		return errors.Wrap(err, "errors checking target configuration")
	}

	b, err := getFileContent(req.FilePath, req.SBOMFormat)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	var rows []*PackageRow

	switch req.SBOMFormat {
	case types.SBOMFormatSPDX:
		var sbom spdx.Document
		if err = json.Unmarshal(b, &sbom); err != nil {
			return errors.Wrap(err, "failed to unmarshal SPDX SBOM")
		}
		rows = MakeSPDXPackageRows(&sbom)
	case types.SBOMFormatCycloneDX:
		var sbom cyclonedx.BOM
		if err = json.Unmarshal(b, &sbom); err != nil {
			return errors.Wrap(err, "failed to unmarshal CycloneDX SBOM")
		}
		rows = MakeCycloneDXPackageRows(&sbom)
	default:
		return errors.Errorf("unsupported SBOM format: %s", req.SBOMFormat)
	}

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

func getFileContent(path string, format types.SBOMFormat) ([]byte, error) {
	if format == types.SBOMFormatUndefined {
		return nil, errors.Errorf("invalid SBOM format: %v", format)
	}
	if path == "" {
		return nil, errors.Errorf("invalid file path: %s", path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	return b, nil
}
