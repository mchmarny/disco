package vctl

import (
	"context"
	"os"

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

type ImagesQuery struct {
	SimpleQuery
	OnlyDigest bool
}

// DiscoverImages discovers all deployed images in the project.
func DiscoverImages(ctx context.Context, in *ImagesQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	printProjectScope(in.ProjectID)

	images, err := getDeployedImages(ctx, in.ProjectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	log.Info().Msgf("Found the following %d images:", len(images))

	if in.OnlyDigest {
		for _, img := range images {
			os.Stdout.WriteString(img.Image.URI())
			os.Stdout.WriteString("\n")
		}
		return nil
	}

	list := make([]*ImageReport, 0)
	for _, img := range images {
		list = append(list, &ImageReport{
			Region:  img.Region.ID,
			Project: img.Project.ID,
			Service: img.Service.Metadata.Name,
			Image:   img.Image.URL(),
		})
	}

	if err := writeOutput(in.OutputPath, in.OutputFmt, list); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}

func getDeployedImages(ctx context.Context, projectID string) ([]*RunningImage, error) {
	if projectID != "" {
		log.Debug().Msgf("discovering images for project: %s", projectID)
	}

	projects, err := project.GetProjects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting projects")
	}

	list := make([]*RunningImage, 0)

	for _, p := range projects {
		if !isQualifiedProject(ctx, p, projectID) {
			continue
		}

		reg, err := region.GetRegions(ctx, p.Number)
		if err != nil {
			log.Error().Err(err).Msgf("error getting regions for project: %s (#%s)", p.ID, p.Number)
			continue
		}
		log.Info().Msgf("found %d regions where Cloud Run is supported", len(reg))

		for _, r := range reg {
			svcs, err := service.GetServices(ctx, p.Number, r.ID)
			if err != nil {
				log.Error().Err(err).Msgf("error getting services for project: %s in region %s", p.Number, r.ID)
				continue
			}

			log.Debug().Msgf("found %d services in: %s/%s", len(svcs), p.ID, r.ID)
			for _, s := range svcs {
				log.Info().Msgf("processing service: %s (Project: %s, Region: %s)", s.Metadata.Name, p.ID, r.ID)

				for _, c := range s.Spec.Template.Spec.Containers {
					f, err := registry.GetImageInfo(ctx, c.Image)
					if err != nil {
						log.Error().Err(err).Msgf("error getting manifest for: %s", c.Image)
						continue
					}

					if f.Deployed != f.Digest {
						log.Info().Msgf("resolved %s to %s", f.Deployed, f.Digest)
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
	log.Debug().Msgf("qualifying project: %s (%s - %s)", p.ID, p.Number, p.State)

	if filterID != "" && p.ID != filterID {
		log.Debug().Msgf("skipping: %s (filter: %s)", p.ID, filterID)
		return false
	}

	if p.State != project.ProjectStateActive {
		log.Debug().Msgf("skipping: %s (inactive)", p.ID)
		return false
	}

	on, err := usage.IsAPIEnabled(ctx, p.Number, usage.CloudRunAPI)
	if err != nil {
		log.Error().Err(err).Msgf("error checking Cloud Run API: %s", p.ID)
		return false
	}

	if !on {
		log.Debug().Msgf("skipping: %s (API not enabled)", p.ID)
		return false
	}

	return true
}
