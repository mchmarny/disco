package target

import (
	"context"

	"github.com/mchmarny/disco/pkg/target/bq"
	"github.com/mchmarny/disco/pkg/types"
)

var (
	LicenseImporter       RecordImporter = bq.ImportLicenses
	VulnerabilityImporter RecordImporter = bq.ImportVulnerabilities
)

type RecordImporter func(ctx context.Context, req *types.ImportRequest) error
