package target

import (
	"context"
	"strings"

	"github.com/mchmarny/disco/pkg/target/bq"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
)

const (
	importDefaultLocation              = "US"
	importTargetProtocolParts          = 2
	importTargetProjectAndDatasetParts = 2

	ImportTargetProtocolUndefined ImportTargetProtocol = iota
	ImportTargetProtocolBigQuery

	ImportTargetProtocolUndefinedName = "undefined"
	ImportTargetProtocolBigQueryName  = "bq"
)

var (
	VulnerabilityImporter VulnerabilityImporterFunc = bq.ImportVulnerabilities
	PackageImporter       PackageImporterFunc       = bq.ImportPackages
	LicenseImporter       LicenseImporterFunc       = bq.ImportLicenses
)

type LicenseImporterFunc func(ctx context.Context, req *types.ImportRequest, in ...*types.LicenseReport) error

type PackageImporterFunc func(ctx context.Context, req *types.ImportRequest, in ...*types.PackageReport) error

type VulnerabilityImporterFunc func(ctx context.Context, req *types.ImportRequest, in ...*types.VulnerabilityReport) error

type ImportTargetProtocol int64

func (i ImportTargetProtocol) String() string {
	switch i {
	case ImportTargetProtocolUndefined:
		return ImportTargetProtocolBigQueryName
	case ImportTargetProtocolBigQuery:
		return ImportTargetProtocolBigQueryName
	default:
		return ImportTargetProtocolUndefinedName
	}
}

func parseImportTargetProtocol(v string) ImportTargetProtocol {
	switch v {
	case ImportTargetProtocolBigQueryName:
		return ImportTargetProtocolBigQuery
	default:
		return ImportTargetProtocolUndefined
	}
}

// ParseImportRequest parses import request.
// e.g. bq://cloudy-demos.disco.packages
func ParseImportRequest(req *types.SimpleQuery) (*types.ImportRequest, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}
	p := strings.Split(req.TargetRaw, "://")
	if len(p) != importTargetProtocolParts {
		return nil, errors.Errorf("invalid import target: %s", req.TargetRaw)
	}

	protocol := parseImportTargetProtocol(p[0])
	if protocol != ImportTargetProtocolBigQuery {
		return nil, errors.Errorf("invalid import target protocol: %s", p[0])
	}

	var r *types.ImportRequest
	var err error

	switch protocol {
	case ImportTargetProtocolBigQuery:
		r, err = bq.ParseImportRequest(req.Kind, p[1])
		if err != nil {
			return nil, errors.Wrap(err, "error parsing BigQuery import request")
		}
	default:
		return nil, errors.Errorf("invalid import target protocol %s", protocol)
	}

	return r, nil
}
