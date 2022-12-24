package vctl

import (
	"context"

	"github.com/mchmarny/vctl/pkg/analysis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type VulnsQuery struct {
	*SimpleQuery
	CVE string
}

func DiscoverVulns(ctx context.Context, in *VulnsQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	printProjectScope(in.ProjectID)

	images, err := getDeployedImages(ctx, in.ProjectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	// get unique list of iamges
	m := make(map[string]*RunningImage)
	for _, img := range images {
		u := img.Image.URL()
		if _, ok := m[u]; !ok {
			m[u] = img
		}
	}

	log.Info().Msgf("Found %d unique images.", len(m))

	list := make([]*analysis.Occurrence, 0)

	for k, img := range m {
		oc, err := analysis.GetImageVulnerabilities(ctx, img.Project.ID, k)
		if err != nil {
			log.Error().Err(err).Msgf("Error getting vulnerabilities for: %s.", k)
			continue
		}
		if oc == nil {
			log.Debug().Msgf("No vulnerabilities found for: %s.", k)
			continue
		}
		for _, o := range oc {
			log.Info().Msgf("%-14s (%s) in %s (Project: %s, Region: %s).", o.Vulnerability.ShortDescription, o.Vulnerability.Severity, img.Service.Metadata.Name, img.Project.ID, img.Region.ID)

			// TODO: add filter by CVE
			list = append(list, o)
		}
	}

	if err := writeOutput(in.OutputPath, list); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}
