package vctl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mchmarny/vctl/pkg/analysis"
	"github.com/mchmarny/vctl/pkg/project"
	"github.com/mchmarny/vctl/pkg/region"
	"github.com/mchmarny/vctl/pkg/registry"
	"github.com/mchmarny/vctl/pkg/service"
	"github.com/mchmarny/vctl/pkg/usage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type RunningImage struct {
	Image   *registry.ImageInfo
	Service *service.Service
	Project *project.Project
	Region  *region.Region
}

type ImageReport struct {
	Region  string `json:"region"`
	Project string `json:"project"`
	Service string `json:"service"`
	Image   string `json:"image"`
}

func printProjectScope(projectID string) {
	if projectID != "" {
		log.Info().Msgf("Scanning project: '%s'.", projectID)
	} else {
		log.Info().Msgf("Scanning all projects accessible to current user.")
	}
}

// DiscoverImages discovers all deployed images in the project.
func DiscoverImages(ctx context.Context, projectID string) error {
	printProjectScope(projectID)

	images, err := getDeployedImages(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	log.Info().Msgf("Found the following %d images:", len(images))

	list := make([]*ImageReport, 0)
	for _, img := range images {
		list = append(list, &ImageReport{
			Region:  img.Region.ID,
			Project: img.Project.ID,
			Service: img.Service.Metadata.Name,
			Image:   img.Image.URL(),
		})
	}

	fmt.Println()

	if err := json.NewEncoder(os.Stdout).Encode(list); err != nil {
		return errors.Wrap(err, "error encoding report")
	}

	return nil
}

func DiscoverVulns(ctx context.Context, projectID, cve string) error {
	printProjectScope(projectID)

	images, err := getDeployedImages(ctx, projectID)
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

	fmt.Println()

	if err := json.NewEncoder(os.Stdout).Encode(list); err != nil {
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
				log.Info().Msgf("Processing service: %s (Project: %s, Region: %s).", s.Metadata.Name, p.ID, r.ID)

				for _, c := range s.Spec.Template.Spec.Containers {
					f, err := registry.GetImageInfo(ctx, c.Image)
					if err != nil {
						log.Error().Err(err).Msgf("Error getting manifest for: %s.", c.Image)
						continue
					}

					if f.Deployed != f.Digest {
						log.Info().Msgf("Resolved %s to %s.", f.Deployed, f.Digest)
					}
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
