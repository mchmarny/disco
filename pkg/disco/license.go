package disco

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/metric"
	"github.com/mchmarny/disco/pkg/object"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/target"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DiscoverLicenses(ctx context.Context, in *types.LicenseQuery, ir *types.ImportRequest) error {
	if in == nil {
		return errors.New("nil input")
	}

	log.Debug().Msgf("discovering licenses with: %s", in)
	printProjectScope(in.ProjectID, "licenses")

	if err := scanLicenses(ctx, &in.SimpleQuery, makeLicenseFilter(in), ir); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func makeLicenseFilter(in *types.LicenseQuery) types.ItemFilter {
	return func(v interface{}) bool {
		lic := v.(*types.License)

		if in.TypeFilter != "" {
			exclude := !strings.HasPrefix(
				strings.ToLower(lic.Name),
				strings.ToLower(in.TypeFilter))
			log.Debug().Msgf("filter on license (want: %s, got: %s, filter out: %t",
				in.TypeFilter, lic.Name, exclude)
			return exclude
		}

		return false
	}
}

func scanLicenses(ctx context.Context, in *types.SimpleQuery, filter types.ItemFilter, ir *types.ImportRequest) error {
	results := make([]*types.LicenseReport, 0)

	h := func(dir, uri string) error {
		scannerResultPath := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting licenses for %s (file: %s)", uri, scannerResultPath)

		rez, err := scanner.GetLicenses(uri, scannerResultPath, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting licenses for %s", uri)
		}
		log.Info().Msgf("found %d licenses in %s", len(rez.Licenses), uri)
		if len(rez.Licenses) > 0 {
			results = append(results, rez)
		}
		if in.Bucket != "" {
			objName := fmt.Sprintf("%s/%s/licenses-%d.json",
				time.Now().UTC().Format("2006/01/02"),
				types.ParseImageShaFromDigestWithoutPrefix(uri),
				time.Now().UTC().Unix())
			if err := object.Save(ctx, in.Bucket, objName, scannerResultPath); err != nil {
				return errors.Wrapf(err, "error saving %s to %s", scannerResultPath, in.Bucket)
			}
		}
		return nil
	}

	if err := handleImages(ctx, in, h); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	report := types.NewItemReport(in, results...)

	if ir != nil {
		if err := target.LicenseImporter(ctx, ir, report.Items...); err != nil {
			return errors.Wrap(err, "error importing")
		}
	}

	if in.Quiet && in.OutputPath == "" {
		return nil
	}

	if err := writeOutput(in.OutputPath, in.OutputFmt, report); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}

func MeterLicense(ctx context.Context, reportPath string) ([]*metric.Record, error) {
	var report types.ItemReport[types.LicenseReport]
	if err := types.UnmarshalFromFile(reportPath, &report); err != nil {
		return nil, errors.Wrapf(err, "error parsing report file: %s", reportPath)
	}

	if len(report.Items) == 0 {
		log.Debug().Msgf("no licenses found in %s", reportPath)
		return nil, nil
	}

	imageCounter := 0
	licenseCounter := 0

	for _, item := range report.Items {
		imageCounter++
		licenseCounter += len(item.Licenses)
	}

	list := make([]*metric.Record, 0)
	list = append(list, &metric.Record{
		MetricType:  metric.MakeMetricType("disco/license/image"),
		MetricValue: int64(imageCounter),
	})
	// licenses have too high (and unpredictable) cardinality for labels
	list = append(list, &metric.Record{
		MetricType:  metric.MakeMetricType("disco/license/count"),
		MetricValue: int64(licenseCounter),
	})

	return list, nil
}
