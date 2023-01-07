package disco

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/metric"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func DiscoverVulns(ctx context.Context, in *types.VulnsQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	log.Debug().Msgf("discovering vulnerabilities with: %s", in)
	printProjectScope(in.ProjectID, "vulnerabilities")

	f := func(v interface{}) bool {
		vul := v.(*types.Vulnerability)

		if in.CVE != "" {
			exclude := !strings.EqualFold(in.CVE, vul.ID)
			log.Debug().Msgf("filter on cve (want: %s, got: %s, filter out: %t",
				in.CVE, vul.ID, exclude)
			return exclude
		}

		if in.MinVulnSev != types.VulnSevUndefined {
			vs := types.ParseMinVulnSeverityOrDefault(vul.Severity)
			exclude := !(vs >= in.MinVulnSev)
			log.Debug().Msgf("filter on severity (want: %s, got: %s, filter out: %t",
				in.MinVulnSev, vul.Severity, exclude)
			return exclude
		}
		return false
	}

	if err := scanVulnerabilities(ctx, in, f); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func scanVulnerabilities(ctx context.Context, in *types.VulnsQuery, filter types.ItemFilter) error {
	results := make([]*types.VulnerabilityReport, 0)
	vSpacer := "vulnerabilities"
	if in.CVE != "" {
		vSpacer = fmt.Sprintf("%ss", in.CVE)
	}
	h := func(dir, uri string) error {
		p := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting vulnerabilities for %s (file: %s)", uri, p)

		rez, err := scanner.GetVulnerabilities(uri, p, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting vulnerabilities for %s", uri)
		}
		log.Info().Msgf("found %d %s in %s", len(rez.Vulnerabilities), vSpacer, uri)
		if len(rez.Vulnerabilities) > 0 {
			results = append(results, rez)
		}
		return nil
	}

	if err := handleImages(ctx, &in.SimpleQuery, h); err != nil {
		return errors.Wrap(err, "error handling images")
	}

	report := types.NewItemReport(&in.SimpleQuery, results...)

	if err := writeOutput(in.OutputPath, in.OutputFmt, report); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}

func MeterVulns(ctx context.Context, reportPath string) ([]*metric.Record, error) {
	var report types.ItemReport[types.VulnerabilityReport]
	if err := types.UnmarshalFromFile(reportPath, &report); err != nil {
		return nil, errors.Wrapf(err, "error parsing report file: %s", reportPath)
	}

	if len(report.Items) == 0 {
		log.Debug().Msgf("no vulnerabilities found in %s", reportPath)
		return nil, nil
	}

	imageCounter := 0
	severityCounter := make(map[string]int64)

	for _, item := range report.Items {
		imageCounter++
		for _, vuln := range item.Vulnerabilities {
			severityCounter[vuln.Severity]++
		}
	}

	list := make([]*metric.Record, 0)
	list = append(list, &metric.Record{
		MetricType:  metric.MakeMetricType("disco/vulnerability/image"),
		MetricValue: int64(imageCounter),
	})
	for k, v := range severityCounter {
		list = append(list, &metric.Record{
			MetricType:  metric.MakeMetricType("disco/vulnerability/severity"),
			MetricValue: v,
			Labels: map[string]string{
				"level": metric.MakeMetricLabelSafe(k),
			},
		})
	}

	return list, nil
}
