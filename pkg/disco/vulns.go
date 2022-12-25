package disco

import (
	"context"
	"fmt"

	"github.com/mchmarny/disco/pkg/analysis"
	"github.com/mchmarny/disco/pkg/project"
	"github.com/mchmarny/disco/pkg/scanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type VulnsQuery struct {
	SimpleQuery
	CVE   string
	CAAPI bool
}

func (q *VulnsQuery) String() string {
	return fmt.Sprintf("ProjectID:%s, CVE:%s, Output:%s, Format:%s, CA-API:%t",
		q.ProjectID, q.CVE, q.OutputPath, q.OutputFmt, q.CAAPI)
}

func DiscoverVulns(ctx context.Context, in *VulnsQuery) error {
	if in == nil {
		return errors.New("nil input")
	}

	if !in.CAAPI {
		return DiscoverVulnsLocally(ctx, in)
	}

	log.Debug().Msgf("Discovering vulnerabilities with: %s", in)
	printProjectScope(in.ProjectID)

	var list []*analysis.Occurrence
	var err error

	// for single project
	if in.ProjectID != "" {
		if in.CVE == "" {
			list, err = discoverImageVulns(ctx, in.ProjectID)
		} else {
			list, err = discoverProjectCVEs(ctx, in.ProjectID, in.CVE)
		}
		if err != nil {
			return errors.Wrapf(err, "error discovering vulnerabilities for project: %s", in.ProjectID)
		}
		if err := writeOutput(in.OutputPath, in.OutputFmt, list); err != nil {
			return errors.Wrap(err, "error writing output")
		}
		return nil
	}

	// for all projects
	projects, err := project.GetProjects(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting projects")
	}

	var subList []*analysis.Occurrence

	for _, p := range projects {
		if in.CVE == "" {
			subList, err = discoverImageVulns(ctx, p.ID)
		} else {
			subList, err = discoverProjectCVEs(ctx, p.ID, in.CVE)
		}
		if err != nil {
			return errors.Wrapf(err, "error discovering vulnerabilities for project: %s", p.ID)
		}
		list = append(list, subList...)
	}

	if err := writeOutput(in.OutputPath, in.OutputFmt, list); err != nil {
		return errors.Wrap(err, "error writing output")
	}
	return nil
}

func discoverProjectCVEs(ctx context.Context, projectID, cveID string) ([]*analysis.Occurrence, error) {
	if cveID == "" || projectID == "" {
		return nil, errors.New("projectID and cveID required")
	}

	list, err := analysis.GetCVEVulnerabilities(ctx, projectID, cveID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting vulnerabilities for: %s in: %s", cveID, projectID)
	}

	for _, o := range list {
		log.Info().Msgf("%s [%s] in %s", o.Vulnerability.ShortDescription, o.Vulnerability.Severity, o.URI)
	}

	return list, nil
}

func discoverImageVulns(ctx context.Context, projectID string) ([]*analysis.Occurrence, error) {
	if projectID == "" {
		return nil, errors.New("projectID required")
	}

	images, err := getDeployedImages(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting images")
	}

	if len(images) == 0 {
		return make([]*analysis.Occurrence, 0), nil
	}

	// get unique list of iamges
	m := make(map[string]*RunningImage)
	for _, img := range images {
		u := img.Image.URL()
		if _, ok := m[u]; !ok {
			m[u] = img
		}
	}

	log.Info().Msgf("found %d unique images", len(m))

	list := make([]*analysis.Occurrence, 0)

	for k, img := range m {
		oc, err := analysis.GetImageVulnerabilities(ctx, img.Project.ID, k)
		if err != nil {
			log.Error().Err(err).Msgf("error getting vulnerabilities for: %s", k)
			continue
		}
		if oc == nil {
			log.Debug().Msgf("no vulnerabilities found for: %s", k)
			continue
		}
		for _, o := range oc {
			log.Info().Msgf("%-14s - %s in %s (Project: %s, Region: %s)", o.Vulnerability.ShortDescription, o.Vulnerability.Severity, img.Service.Metadata.Name, img.Project.ID, img.Region.ID)

			list = append(list, o)
		}
	}

	return list, nil
}

func DiscoverVulnsLocally(ctx context.Context, in *VulnsQuery) error {
	vulnFilter := func(v string) bool {
		if in.CVE == "" {
			return true
		}
		return v == in.CVE
	}

	if err := scan(ctx, scanner.VulnerabilityScanner, &in.SimpleQuery, vulnFilter); err != nil {
		return errors.Wrap(err, "error scanning for vulnerabilities")
	}

	return nil
}
