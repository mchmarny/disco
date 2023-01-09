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

func DiscoverPackages(ctx context.Context, in *types.PackageQuery, ir *types.ImportRequest) error {
	if in == nil {
		return errors.New("nil input")
	}

	log.Debug().Msgf("discovering packages with: %s", in)
	printProjectScope(in.ProjectID, "packages")

	if err := scanPackages(ctx, &in.SimpleQuery, makePackageFilter(in), ir); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func makePackageFilter(in *types.PackageQuery) types.ItemFilter {
	in.NamePart = strings.ToLower(strings.TrimSpace(in.NamePart))
	return func(v interface{}) bool {
		lic := v.(*types.Package)

		if in.NamePart != "" {
			exclude := !strings.Contains(strings.ToLower(lic.Package), in.NamePart)
			log.Debug().Msgf("filter on package (want: %s, got: %s, filter out: %t",
				in.NamePart, lic.Package, exclude)
			return exclude
		}

		return false
	}
}

func makePackageHandler(ctx context.Context, in *types.SimpleQuery, results []*types.PackageReport, filter types.ItemFilter) itemHandler {
	return func(dir, uri string) error {
		scannerResultPath := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting packages for %s (file: %s)", uri, scannerResultPath)

		rez, err := scanner.GetPackages(uri, scannerResultPath, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting packages for %s", uri)
		}
		log.Info().Msgf("found %d packages in %s", len(rez.Packages), uri)
		if len(rez.Packages) > 0 {
			results = append(results, rez)
		}
		if in.Bucket != "" {
			objName := fmt.Sprintf("%s/%s/packages-%d.json",
				time.Now().UTC().Format("2006/01/02"),
				types.ParseImageShaFromDigestWithoutPrefix(uri),
				time.Now().UTC().Unix())
			if err := object.Save(ctx, in.Bucket, objName, scannerResultPath); err != nil {
				return errors.Wrapf(err, "error saving %s to %s", scannerResultPath, in.Bucket)
			}
		}
		return nil
	}
}

func scanPackages(ctx context.Context, in *types.SimpleQuery, filter types.ItemFilter, ir *types.ImportRequest) error {
	results := make([]*types.PackageReport, 0)

	h := makePackageHandler(ctx, in, results, filter)

	if err := handleImages(ctx, in, h); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	report := types.NewItemReport(in, results...)

	if ir != nil {
		if err := target.PackageImporter(ctx, ir, report.Items...); err != nil {
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

func MeterPackage(ctx context.Context, reportPath string) ([]*metric.Record, error) {
	var report types.ItemReport[types.PackageReport]
	if err := types.UnmarshalFromFile(reportPath, &report); err != nil {
		return nil, errors.Wrapf(err, "error parsing report file: %s", reportPath)
	}

	if len(report.Items) == 0 {
		log.Debug().Msgf("no packages found in %s", reportPath)
		return nil, nil
	}

	imageCounter := 0
	packageCounter := 0

	for _, item := range report.Items {
		imageCounter++
		packageCounter += len(item.Packages)
	}

	list := make([]*metric.Record, 0)
	list = append(list, &metric.Record{
		MetricType:  metric.MakeMetricType("disco/package/image"),
		MetricValue: int64(imageCounter),
	})
	// packages have too high (and unpredictable) cardinality for labels
	list = append(list, &metric.Record{
		MetricType:  metric.MakeMetricType("disco/package/count"),
		MetricValue: int64(packageCounter),
	})

	return list, nil
}
