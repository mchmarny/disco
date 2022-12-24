package vctl

import (
	"context"
	"encoding/json"
	"os"

	"github.com/mchmarny/vctl/pkg/project"
	"github.com/mchmarny/vctl/pkg/region"
	"github.com/mchmarny/vctl/pkg/registry"
	"github.com/mchmarny/vctl/pkg/service"
	"github.com/mchmarny/vctl/pkg/usage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DiscoReport struct {
	Images []*RunningImage
}

type RunningImage struct {
	Image   *registry.ImageInfo
	Service *service.Service
	Project *project.Project
	Region  *region.Region
}

func DiscoverVulns(ctx context.Context, projectID string) error {
	if projectID != "" {
		log.Info().Msgf("Discovering vulnerabilities for project: '%s'.", projectID)
	} else {
		log.Info().Msgf("Discovering vulnerabilities for all projects accessible to current user.")
	}

	images, err := getDeployedImages(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	report := &DiscoReport{
		Images: images,
	}

	// TODO: convert report to vuln report

	log.Info().Msgf("Done, found %d images.", len(images))

	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		return errors.Wrap(err, "error encoding report")
	}

	return nil
}

func getDeployedImages(ctx context.Context, projectID string) ([]*RunningImage, error) {
	if projectID != "" {
		log.Debug().Msgf("discovering images for project: %s.", projectID)
	}

	projects, err := project.GetProjects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting projects.")
	}

	list := make([]*RunningImage, 0)

	for _, p := range projects {
		log.Debug().Msgf("Quering project: %s (%s - %s).", p.ID, p.Number, p.State)

		if !isQualifiedProject(ctx, p, projectID) {
			continue
		}

		reg, err := region.GetRegions(ctx, p.Number)
		if err != nil {
			log.Error().Err(err).Msgf("Error getting regions for project: %s (#%s).", p.ID, p.Number)
			continue
		}
		log.Info().Msgf("Found %d regions where Cloud Run is supported.", len(reg))

		for _, r := range reg {
			svcs, err := service.GetServices(ctx, p.Number, r.ID)
			if err != nil {
				log.Error().Err(err).Msgf("error getting services for project: %s in region %s.", p.Number, r.ID)
				continue
			}

			log.Debug().Msgf("Found %d services in: %s/%s.", len(svcs), p.ID, r.ID)
			for _, s := range svcs {
				log.Info().Msgf("Processing service: %s.", s.Metadata.Name)

				for _, c := range s.Spec.Template.Spec.Containers {
					f, err := registry.GetImageInfo(ctx, c.Image)
					if err != nil {
						log.Error().Err(err).Msgf("Error getting manifest for: %s.", c.Image)
						continue
					}

					log.Info().Msgf("Found container digest %s.", f.Digest)
					list = append(list, &RunningImage{
						Project: p,
						Region:  r,
						Service: s,
						Image:   f,
					})
				}
			}
		}
	}

	return list, nil
}

func isQualifiedProject(ctx context.Context, p *project.Project, filterID string) bool {
	if filterID != "" && p.ID != filterID {
		log.Debug().Msgf("Skipping project: %s (filter: %s).", p.ID, filterID)
		return false
	}

	if p.State != project.ProjectStateActive {
		log.Debug().Msgf("Skipping project: %s (inactive).", p.ID)
		return false
	}

	on, err := usage.IsAPIEnabled(ctx, p.Number, usage.CloudRunAPI)
	if err != nil {
		log.Error().Err(err).Msgf("Error checking if Cloud Run API is enabled for project: %s.", p.ID)
		return false
	}

	if !on {
		log.Debug().Msgf("Skipping project: %s (API not enabled).", p.ID)
		return false
	}

	return true
}
