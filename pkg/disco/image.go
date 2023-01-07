package disco

import (
	"context"
	"os"

	"github.com/mchmarny/disco/pkg/source"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DiscoverImages discovers all deployed images in the project.
func DiscoverImages(ctx context.Context, in *types.ImagesQuery) error {
	if in == nil {
		return errors.New("nil input")
	}
	log.Debug().Msgf("discovering images with: %s", in)
	printProjectScope(in.ProjectID, "images")

	images, err := source.ImageProvider(ctx, in)
	if err != nil {
		return errors.Wrap(err, "error getting images")
	}

	if in.URIOnly {
		if err := writeList(in.OutputPath, images); err != nil {
			return errors.Wrap(err, "error writing output")
		}
		return nil
	}

	report := types.NewItemReport(&in.SimpleQuery, images...)

	if err := writeOutput(in.OutputPath, in.OutputFmt, report); err != nil {
		return errors.Wrap(err, "error writing output")
	}

	return nil
}

func discoverImageURIs(ctx context.Context, in *types.ImagesQuery) ([]string, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	images, err := source.ImageProvider(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "error getting images")
	}

	list := make([]string, 0, len(images))
	for _, img := range images {
		list = append(list, img.URI)
	}
	return list, nil
}

func writeList(path string, images []*types.ImageItem) error {
	if path == "" {
		for _, img := range images {
			os.Stdout.WriteString(img.URI)
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
		if _, err := f.WriteString(img.URI + "\n"); err != nil {
			return errors.Wrapf(err, "error writing to file: %s", path)
		}
	}
	return nil
}
