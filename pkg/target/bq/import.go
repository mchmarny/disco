package bq

import (
	"context"
	"strings"

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
	inserter.SkipInvalidRows = true
	if err := inserter.Put(ctx, items); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

const (
	importDefaultLocation     = "US"
	importTargetProtocolParts = 2
	importTargetProjectPart   = 1
	importTargetDatasetPart   = 2
	importTargetTablePart     = 3
)

// ParseImportRequest parses import request.
// e.g. bq://cloudy-demos.disco.packages
func ParseImportRequest(k types.DiscoKind, v string) (*types.ImportRequest, error) {
	t := &types.ImportRequest{
		Location: importDefaultLocation,
	}

	switch k {
	case types.KindLicense:
		t.TableID = types.TableKindLicenseName
		t.TableKind = types.TableKindLicense
	case types.KindVulnerability:
		t.TableID = types.TableKindVulnerabilityName
		t.TableKind = types.TableKindVulnerability
	case types.KindPackage:
		t.TableID = types.TableKindPackageName
		t.TableKind = types.TableKindPackage
	default:
		return nil, errors.Errorf("invalid table kind: %s", k)
	}

	parts := strings.Split(v, ".")
	if len(parts) < importTargetProjectPart || len(parts) > importTargetTablePart {
		return nil, errors.Errorf("invalid import target: %s", v)
	}

	t.ProjectID = parts[0]

	if len(parts) == importTargetTablePart {
		t.DatasetID = parts[1]
		t.TableID = parts[2]
	} else if len(parts) == importTargetDatasetPart {
		t.DatasetID = parts[1]
		t.TableID = getTableFromKind(k)
	} else {
		t.DatasetID = DatasetNameDefault
		t.TableID = getTableFromKind(k)
	}

	if t.ProjectID == "" || t.DatasetID == "" || t.TableID == "" {
		return nil, errors.Errorf("invalid import target: %s", v)
	}

	return t, nil
}

func getTableFromKind(k types.DiscoKind) string {
	switch k {
	case types.KindLicense:
		return types.TableKindLicenseName
	case types.KindVulnerability:
		return types.TableKindVulnerabilityName
	case types.KindPackage:
		return types.TableKindPackageName
	default:
		return ""
	}
}
