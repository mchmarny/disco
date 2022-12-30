package disco

import (
	"context"
	"path"
	"strings"

	"github.com/google/uuid"
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
	printProjectScope(in.ProjectID)

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

	if err := scanVulnerabilities(ctx, &in.SimpleQuery, f); err != nil {
		return errors.Wrap(err, "error scanning")
	}

	return nil
}

func scanVulnerabilities(ctx context.Context, in *types.SimpleQuery, filter types.ItemFilter) error {
	results := make([]*types.VulnerabilityReport, 0)
	h := func(dir, uri string) error {
		p := path.Join(dir, uuid.NewString())
		log.Debug().Msgf("getting %s for %s (file: %s)", types.KindVulnerabilityName, uri, p)

		rez, err := scanner.GetVulnerabilities(uri, p, filter)
		if err != nil {
			return errors.Wrapf(err, "error getting licenses for %s", uri)
		}
		log.Info().Msgf("found %d licenses in %s", len(rez.Vulnerabilities), uri)
		if len(rez.Vulnerabilities) > 0 {
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
