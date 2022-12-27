package disco

import (
	"context"
	"os"

	"github.com/mchmarny/disco/pkg/gcp"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type RunningImage struct {
	Image    *gcp.ImageInfo
	Service  *gcp.Service
	Project  *gcp.Project
	Location *gcp.Location
}

// DiscoverImages discovers all deployed images in the project.
func DiscoverImages(ctx context.Context, in *types.ImagesQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	log.Debug().Msgf("discovering images with: %s", in)
	printProjectScope(in.ProjectID)

	images, err := getDeployedImages(ctx, in.ProjectID)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	log.Info().Msgf("found %d images", len(images))

	if in.URIOnly {
		if err := writeList(in.OutputPath, images); err != nil {
			return errors.Wrap(err, "error writing output")
		}
		return nil
	}

	list := make([]*types.ImageReport, 0)
	for _, img := range images {
		list = append(list, &types.ImageReport{
			Location: img.Location.ID,
			Project:  img.Project.ID,
			Service:  img.Service.Metadata.Name,
			Image:    img.Image.URI(),
		})
	}

	if err := writeOutput(in.OutputPath, in.OutputFmt, list); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}

func getDeployedImageURIs(ctx context.Context, projectID string) ([]string, error) {
	images, err := getDeployedImages(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting images")
	}
	imageURIs := make([]string, 0)
	for _, img := range images {
		imageURIs = append(imageURIs, img.Image.URI())
	}
	return imageURIs, nil
}

func getDeployedImages(ctx context.Context, projectID string) ([]*RunningImage, error) {
	if projectID != "" {
		log.Debug().Msgf("discovering images for project: %s", projectID)
	}

	projects, err := getProjectsFunc(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error getting projects")
	}
	log.Debug().Msgf("found %d projects", len(projects))

	list := make([]*RunningImage, 0)

	for _, p := range projects {
		if !isQualifiedProject(ctx, p, projectID) {
			continue
		}

		reg, err := getLocationsFunc(ctx, p.Number)
		if err != nil {
			log.Error().Err(err).Msgf("error getting regions for project: %s (#%s)", p.ID, p.Number)
			continue
		}
		log.Info().Msgf("found %d regions where Cloud Run is supported, processing...", len(reg))

		for _, r := range reg {
			svcs, err := getServicesFunc(ctx, p.Number, r.ID)
			if err != nil {
				log.Error().Err(err).Msgf("error getting services for project: %s in region %s", p.Number, r.ID)
				continue
			}

			log.Debug().Msgf("found %d services in: %s/%s", len(svcs), p.ID, r.ID)
			for _, s := range svcs {
				log.Info().Msgf("processing service: %s (project: %s, region: %s)", s.Metadata.Name, p.ID, r.ID)

				for _, c := range s.Spec.Template.Spec.Containers {
					f, err := getImageInfoFunc(ctx, c.Image)
					if err != nil {
						log.Error().Err(err).Msgf("error getting manifest for: %s", c.Image)
						continue
					}

					if f.Deployed != f.Digest {
						log.Info().Msgf("resolved %s to %s", f.Deployed, f.Digest)
					}
					list = append(list, &RunningImage{
						Project:  p,
						Location: r,
						Service:  s,
						Image:    f,
					})
				}
			}
		}
	}

	return list, nil
}

func isQualifiedProject(ctx context.Context, p *gcp.Project, filterID string) bool {
	log.Debug().Msgf("qualifying project: %s (%s - %s)", p.ID, p.Number, p.State)

	if filterID != "" && p.ID != filterID {
		log.Debug().Msgf("skipping: %s (filter: %s)", p.ID, filterID)
		return false
	}

	if p.State != gcp.ProjectStateActive {
		log.Debug().Msgf("skipping: %s (inactive)", p.ID)
		return false
	}

	on, err := isAPIEnabledFunc(ctx, p.Number, gcp.CloudRunAPI)
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

func writeList(path string, images []*RunningImage) error {
	if path == "" {
		for _, img := range images {
			os.Stdout.WriteString(img.Image.URI())
			os.Stdout.WriteString("\n")
		}
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "error creating file: %s", path)
	}
	defer f.Close()

	for _, img := range images {
		if _, err := f.WriteString(img.Image.URI() + "\n"); err != nil {
			return errors.Wrapf(err, "error writing to file: %s", path)
		}
	}
	return nil
}
