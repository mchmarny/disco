package vctl

import (
	"context"
	"encoding/json"
	"os"

	"github.com/mchmarny/vctl/pkg/project"
	"github.com/mchmarny/vctl/pkg/region"
	"github.com/mchmarny/vctl/pkg/service"
	"github.com/mchmarny/vctl/pkg/usage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DiscoReport struct {
	Items []*DiscoItem `json:"items"`
}

type DiscoItem struct {
	Project  *project.Project   `json:"project"`
	Services []*service.Service `json:"services"`
}

func DiscoverVulns(ctx context.Context, projectID string) error {
	log.Debug().Msgf("discovering vulnerabilities for project: %s", projectID)
	projects, err := project.GetProjects(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting projects")
	}

	// TODO: for now just spool everything
	report := &DiscoReport{
		Items: make([]*DiscoItem, 0),
	}

	for _, p := range projects {
		log.Debug().Msgf("getting project details: %+v", p)

		if projectID != "" && p.ProjectID != projectID {
			log.Debug().Msgf("skipping project (not selected): %s", p.ProjectID)
			continue
		}

		if p.LifecycleState != project.ProjectStateActive {
			log.Debug().Msgf("skipping project (inactive): %s", p.ProjectID)
			continue
		}
		on, err := usage.IsAPIEnabled(ctx, p.ProjectNumber, usage.CloudRunAPI)
		if err != nil {
			return errors.Wrap(err, "error checking if Cloud Run API is enabled")
		}
		if !on {
			log.Debug().Msgf("skipping project (API not enabled): %s", p.ProjectID)
			continue
		}

		if err := discoServices(ctx, p, report); err != nil {
			return errors.Wrapf(err, "error discovering services: %s", p.ProjectID)
		}
	}

	// TODO: process the report

	// TODO: convert report to vuln report

	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		return errors.Wrap(err, "error encoding report")
	}

	return nil
}

func discoServices(ctx context.Context, project *project.Project, report *DiscoReport) error {
	if project == nil {
		return errors.New("project is nil")
	}
	if report == nil {
		return errors.New("report is nil")
	}

	reg, err := region.GetRegions(ctx, project.ProjectNumber)
	if err != nil {
		return errors.Wrapf(err, "error getting regions for project: %s (#%s)",
			project.ProjectID, project.ProjectNumber)
	}
	log.Debug().Msgf("found %d regions", len(reg))

	item := &DiscoItem{
		Project:  project,
		Services: make([]*service.Service, 0),
	}

	for _, r := range reg {
		log.Debug().Msgf("region: %s (%s)", r.ID, r.ID)
		svcs, err := service.GetServices(ctx, project.ProjectNumber, r.ID)
		if err != nil {
			return errors.Wrapf(err, "error getting services for project: %s in region %s",
				project.ProjectID, r.ID)
		}
		log.Debug().Msgf("found %d services in %s/%s", len(svcs), project.ProjectID, r.ID)
		item.Services = append(item.Services, svcs...)
	}

	report.Items = append(report.Items, item)
	return nil
}
