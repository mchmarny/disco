package disco

import (
	"context"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/metric"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DiscoverLicenses(ctx context.Context, in *types.LicenseQuery) error {
	if in == nil {
		return errors.New("nil input")
	}

	f := func(v interface{}) bool {
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

	log.Debug().Msgf("discovering licenses with: %s", in)
	printProjectScope(in.ProjectID, "licenses")

	if err := scanLicenses(ctx, &in.SimpleQuery, f); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func scanLicenses(ctx context.Context, in *types.SimpleQuery, filter types.ItemFilter) error {
	results := make([]*types.LicenseReport, 0)
	h := func(dir, uri string) error {
		p := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting %s for %s (file: %s)", types.KindLicenseName, uri, p)

		rez, err := scanner.GetLicenses(uri, p, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting licenses for %s", uri)
		}
		log.Info().Msgf("found %d licenses in %s", len(rez.Licenses), uri)
		if len(rez.Licenses) > 0 {
			results = append(results, rez)
		}
		return nil
	}

	if err := handleImages(ctx, in, h); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	report := types.NewItemReport(in, results...)

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
